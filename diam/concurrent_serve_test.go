// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// cerPayload builds a minimal CER payload usable by the test mux.
func cerPayload(tb testing.TB) []byte {
	tb.Helper()
	msg := NewRequest(257, 0, dict.Default)
	msg.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	msg.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("test"))
	msg.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	msg.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(0))
	msg.NewAVP(avp.ProductName, 0, 0, datatype.UTF8String("test"))
	b, err := msg.Serialize()
	if err != nil {
		tb.Fatalf("serialize: %v", err)
	}
	return b
}

// TestConcurrentServeThroughput verifies that enabling MaxConcurrentHandlers
// actually dispatches handlers concurrently. With a 1ms handler and 1000
// messages, sequential dispatch would take >= 1s; concurrent should be
// dramatically faster.
func TestConcurrentServeThroughput(t *testing.T) {
	var received int64
	const total = 1000

	mux := NewServeMux()
	mux.HandleFunc("ALL", func(c Conn, m *Message) {
		time.Sleep(1 * time.Millisecond)
		atomic.AddInt64(&received, 1)
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	srv := &Server{Handler: mux, MaxConcurrentHandlers: 256}
	go srv.Serve(ln)

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	payload := cerPayload(t)
	start := time.Now()
	for i := 0; i < total; i++ {
		if _, err := conn.Write(payload); err != nil {
			t.Fatal(err)
		}
	}
	for time.Since(start) < 5*time.Second {
		if atomic.LoadInt64(&received) >= total {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	elapsed := time.Since(start)
	got := atomic.LoadInt64(&received)
	t.Logf("Sent %d, received %d in %v (%.0f msg/s)", total, got, elapsed, float64(got)/elapsed.Seconds())
	if got < total {
		t.Errorf("Lost %d messages", total-got)
	}
	if elapsed > 500*time.Millisecond {
		t.Errorf("Too slow (%v); handlers may still be sequential", elapsed)
	}
}

// TestConcurrentServePanicRecovery verifies that a panic inside a
// concurrently-dispatched handler does not crash the server and that
// subsequent messages are still processed.
func TestConcurrentServePanicRecovery(t *testing.T) {
	var received int64

	mux := NewServeMux()
	mux.HandleFunc("ALL", func(c Conn, m *Message) {
		if atomic.AddInt64(&received, 1) == 1 {
			panic("boom")
		}
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	srv := &Server{Handler: mux, MaxConcurrentHandlers: 4}
	go srv.Serve(ln)

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	payload := cerPayload(t)
	for i := 0; i < 3; i++ {
		if _, err := conn.Write(payload); err != nil {
			t.Fatal(err)
		}
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt64(&received) >= 3 {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	if got := atomic.LoadInt64(&received); got != 3 {
		t.Errorf("want 3 received, got %d (server may have crashed on panic)", got)
	}
}

// TestConcurrentServeBounded verifies that MaxConcurrentHandlers caps the
// number of simultaneously-running handlers.
func TestConcurrentServeBounded(t *testing.T) {
	const limit = 4
	const total = 20

	var inflight, maxInflight int64
	release := make(chan struct{})

	mux := NewServeMux()
	mux.HandleFunc("ALL", func(c Conn, m *Message) {
		cur := atomic.AddInt64(&inflight, 1)
		for {
			m := atomic.LoadInt64(&maxInflight)
			if cur <= m || atomic.CompareAndSwapInt64(&maxInflight, m, cur) {
				break
			}
		}
		<-release
		atomic.AddInt64(&inflight, -1)
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	srv := &Server{Handler: mux, MaxConcurrentHandlers: limit}
	go srv.Serve(ln)

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	payload := cerPayload(t)
	for i := 0; i < total; i++ {
		if _, err := conn.Write(payload); err != nil {
			t.Fatal(err)
		}
	}

	// Give the server time to saturate the semaphore.
	time.Sleep(200 * time.Millisecond)
	peak := atomic.LoadInt64(&maxInflight)
	close(release)
	if peak > int64(limit) {
		t.Errorf("in-flight handlers %d exceeded limit %d", peak, limit)
	}
	if peak == 0 {
		t.Errorf("no handlers ran")
	}
}

// The benchmarks below quantify the throughput advantage of concurrent
// handler dispatch over sequential dispatch under two regimes:
//
//   - 1ms handler latency: realistic for handlers that do any I/O (DB
//     lookup, remote call, structured logging). Sequential throughput is
//     bounded by 1/latency (~1000 msg/s). Concurrent dispatch parallelizes
//     those waits and lifts the ceiling to MaxConcurrentHandlers/latency.
//
//   - No-op handler: measures raw dispatch overhead only. There is no
//     latency to hide, so the two modes converge and concurrent mode pays
//     a small goroutine+semaphore tax. This is the worst case for
//     concurrent dispatch and is why it remains opt-in.
//
// runDispatchBenchmark drives the real server/read path with a handler of
// the given latency, dispatching `total` messages and reports ns/op (per
// message).
func runDispatchBenchmark(b *testing.B, maxConcurrent int, handlerLatency time.Duration, total int) {
	b.Helper()
	for i := 0; i < b.N; i++ {
		var received int64
		done := make(chan struct{})

		mux := NewServeMux()
		mux.HandleFunc("ALL", func(c Conn, m *Message) {
			if handlerLatency > 0 {
				time.Sleep(handlerLatency)
			}
			if atomic.AddInt64(&received, 1) == int64(total) {
				close(done)
			}
		})

		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			b.Fatal(err)
		}
		srv := &Server{Handler: mux, MaxConcurrentHandlers: maxConcurrent}
		go srv.Serve(ln)

		conn, err := net.Dial("tcp", ln.Addr().String())
		if err != nil {
			ln.Close()
			b.Fatal(err)
		}
		payload := cerPayload(b)

		b.StartTimer()
		for j := 0; j < total; j++ {
			if _, err := conn.Write(payload); err != nil {
				b.Fatal(err)
			}
		}
		<-done
		b.StopTimer()

		conn.Close()
		ln.Close()
	}
	b.ReportMetric(float64(total), "msgs/op")
}

// BenchmarkDispatchSequential1ms models a realistic ~1ms handler with the
// default (sequential) dispatch. Expected throughput: ~1000 msg/s.
func BenchmarkDispatchSequential1ms(b *testing.B) {
	b.StopTimer()
	runDispatchBenchmark(b, 0, 1*time.Millisecond, 200)
}

// BenchmarkDispatchConcurrent1ms models the same handler with concurrent
// dispatch. Expected throughput: orders of magnitude higher.
func BenchmarkDispatchConcurrent1ms(b *testing.B) {
	b.StopTimer()
	runDispatchBenchmark(b, 256, 1*time.Millisecond, 200)
}

// BenchmarkDispatchSequentialNoLatency measures pure dispatch overhead with
// a no-op handler under sequential mode.
func BenchmarkDispatchSequentialNoLatency(b *testing.B) {
	b.StopTimer()
	runDispatchBenchmark(b, 0, 0, 1000)
}

// BenchmarkDispatchConcurrentNoLatency measures pure dispatch overhead with
// a no-op handler under concurrent mode (shows the overhead of the
// semaphore + goroutine path).
func BenchmarkDispatchConcurrentNoLatency(b *testing.B) {
	b.StopTimer()
	runDispatchBenchmark(b, 256, 0, 1000)
}
