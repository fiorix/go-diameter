// Copyright 2013-2020 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

//go:build linux && !386

package diam

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"unsafe"

	"github.com/ishidawataru/sctp"
)

// setSCTPNoDelay must set SCTP_NODELAY even when the goroutine stack grows and
// moves during the call.
//
// The parent checks that NewSCTPConn sets SCTP_NODELAY on a loopback socket
// (skipping immediately when kernel SCTP is unavailable), then re-execs the
// test binary with GODEBUG=efence=1 so the runtime unmaps old stacks after
// moving them. The child sweeps stack depths in fresh goroutines over
// setSCTPNoDelay directly because NewSCTPConn calls SubscribeEvents first,
// whose deeper frame grows the stack before setSCTPNoDelay runs.
//
// The socket pair is made with plain syscalls so the test does not depend on
// option calls inside the sctp package. The aim check sweeps sctp's uintptr
// Setsockopt the same way and requires at least one failure so the test cannot
// pass without moving stacks.

const noDelayChildEnv = "GODIAM_TEST_NODELAY_STACK_CHILD"

// atDepth calls f under n extra stack frames.
//
//go:noinline
func atDepth(n int, f func()) {
	var pad [64]byte
	if n == 0 {
		f()
		return
	}
	pad[n%len(pad)] = byte(n)
	atDepth(n-1, f)
	_ = pad
}

func TestNewSCTPConnSetsNoDelayWhenTheStackMoves(t *testing.T) {
	if os.Getenv(noDelayChildEnv) == "" {
		fd := sctpLoopbackClient(t)
		if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_SCTP, sctp.SCTP_NODELAY, 0); err != nil {
			t.Fatalf("clear SCTP_NODELAY: %v", err)
		}
		NewSCTPConn(sctp.NewSCTPConn(fd, nil))
		if got := readSCTPNoDelay(t, fd); got != 1 {
			t.Fatalf("NewSCTPConn left SCTP_NODELAY=%d, want 1", got)
		}

		godebug := "efence=1"
		if v := os.Getenv("GODEBUG"); v != "" {
			godebug = v + "," + godebug
		}
		cmd := exec.Command(os.Args[0],
			"-test.run=^TestNewSCTPConnSetsNoDelayWhenTheStackMoves$",
			"-test.v")
		cmd.Env = append(os.Environ(), noDelayChildEnv+"=1", "GODEBUG="+godebug)
		out, err := cmd.CombinedOutput()
		t.Logf("child (GODEBUG=%s):\n%s", godebug, out)
		if err != nil {
			t.Fatalf("child failed: %v", err)
		}
		return
	}

	fd := sctpLoopbackClient(t)
	conn := sctp.NewSCTPConn(fd, nil)
	const depths = 1600

	unset := func(set func()) int {
		failed := 0
		for d := 0; d < depths; d++ {
			if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_SCTP, sctp.SCTP_NODELAY, 0); err != nil {
				t.Fatalf("clear SCTP_NODELAY: %v", err)
			}
			done := make(chan struct{})
			go func(d int) {
				atDepth(d, set)
				close(done)
			}(d)
			<-done
			if readSCTPNoDelay(t, fd) != 1 {
				failed++
			}
		}
		return failed
	}

	fixed := unset(func() { setSCTPNoDelay(conn) })
	canary := unset(func() {
		v := 1
		conn.Setsockopt(sctp.SCTP_NODELAY, uintptr(unsafe.Pointer(&v)), unsafe.Sizeof(v))
	})
	t.Logf("SCTP_NODELAY left unset over %d depths: setSCTPNoDelay %d; uintptr Setsockopt (aim) %d", depths, fixed, canary)
	if fixed != 0 {
		t.Errorf("setSCTPNoDelay left SCTP_NODELAY unset at %d of %d stack depths", fixed, depths)
	}
	if canary == 0 {
		t.Fatal("the uintptr Setsockopt never failed: these depths did not move the stack during the call, so this test proves nothing")
	}
}

func readSCTPNoDelay(t *testing.T, fd int) int {
	t.Helper()
	v, err := syscall.GetsockoptInt(fd, syscall.IPPROTO_SCTP, sctp.SCTP_NODELAY)
	if err != nil {
		t.Fatalf("read SCTP_NODELAY: %v", err)
	}
	return v
}

// sctpLoopbackClient returns the client socket of an SCTP association over
// 127.0.0.1, made with plain syscalls.
func sctpLoopbackClient(t *testing.T) int {
	t.Helper()
	ln, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM|syscall.SOCK_CLOEXEC, syscall.IPPROTO_SCTP)
	if err != nil {
		t.Skipf("SCTP is not available here: %v", err)
	}
	t.Cleanup(func() { syscall.Close(ln) })
	if err := syscall.Bind(ln, &syscall.SockaddrInet4{Addr: [4]byte{127, 0, 0, 1}}); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Listen(ln, 1); err != nil {
		t.Fatal(err)
	}
	sa, err := syscall.Getsockname(ln)
	if err != nil {
		t.Fatal(err)
	}
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM|syscall.SOCK_CLOEXEC, syscall.IPPROTO_SCTP)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { syscall.Close(fd) })
	if err := syscall.Connect(fd, sa); err != nil {
		t.Fatal(err)
	}
	return fd
}
