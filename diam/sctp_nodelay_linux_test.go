// Copyright 2013-2020 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

//go:build linux && !386

package diam

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"
	"unsafe"

	"github.com/ishidawataru/sctp"
)

// NewSCTPConn must set SCTP_NODELAY even when the goroutine stack moves
// during the call.
//
// The test re-runs itself in a child "go test" that compiles the sctp package
// with -d=maymorestack=runtime.mayMoreStackMove, the Go toolchain's own mode
// for finding stack-move bugs, which moves the goroutine stack at the entry
// of every function in that package, and runs with GODEBUG=efence=1, which
// unmaps a stack once it has been moved. A pointer into the old stack that
// survives a call into the sctp package then gets EFAULT every time instead
// of quietly reading freed memory.
//
// The socket pair is made with plain syscalls, not with the sctp package, so
// the test does not depend on which version of that package is in use. Its
// aim check calls sctp.SCTPConn.Setsockopt the way NewSCTPConn used to, and
// requires that to fail: otherwise the child is not moving stacks and proves
// nothing.

const noDelayChildEnv = "GODIAM_TEST_NODELAY_STACK_CHILD"

func TestNewSCTPConnSetsNoDelayWhenTheStackMoves(t *testing.T) {
	if os.Getenv(noDelayChildEnv) == "" {
		goTool, err := exec.LookPath("go")
		if err != nil {
			goTool = filepath.Join(runtime.GOROOT(), "bin", "go")
			if _, err := os.Stat(goTool); err != nil {
				t.Skipf("no go command to build the stack-moving child: %v", err)
			}
		}
		godebug := "efence=1"
		if v := os.Getenv("GODEBUG"); v != "" {
			godebug = v + "," + godebug
		}
		cmd := exec.Command(goTool, "test", "-count=1", "-v",
			"-run=^TestNewSCTPConnSetsNoDelayWhenTheStackMoves$",
			"-gcflags=github.com/ishidawataru/sctp=-d=maymorestack=runtime.mayMoreStackMove",
			".")
		cmd.Env = append(os.Environ(), noDelayChildEnv+"=1", "GODEBUG="+godebug)
		out, err := cmd.CombinedOutput()
		t.Logf("child (GODEBUG=%s, sctp built with maymorestack=runtime.mayMoreStackMove):\n%s", godebug, out)
		if err != nil {
			t.Fatalf("child failed: %v", err)
		}
		return
	}

	fd := sctpLoopbackClient(t)
	conn := sctp.NewSCTPConn(fd, nil)
	noDelay := func() int {
		v, err := syscall.GetsockoptInt(fd, syscall.IPPROTO_SCTP, sctp.SCTP_NODELAY)
		if err != nil {
			t.Fatalf("read SCTP_NODELAY: %v", err)
		}
		return v
	}
	unset := func(set func()) int {
		failed := 0
		for i := 0; i < 100; i++ {
			if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_SCTP, sctp.SCTP_NODELAY, 0); err != nil {
				t.Fatalf("clear SCTP_NODELAY: %v", err)
			}
			set()
			if noDelay() != 1 {
				failed++
			}
		}
		return failed
	}

	fixed := unset(func() { NewSCTPConn(conn) })
	canary := unset(func() {
		v := 1
		conn.Setsockopt(sctp.SCTP_NODELAY, uintptr(unsafe.Pointer(&v)), unsafe.Sizeof(v))
	})
	t.Logf("SCTP_NODELAY left unset in 100 calls: NewSCTPConn %d; uintptr Setsockopt (aim) %d", fixed, canary)
	if fixed != 0 {
		t.Errorf("NewSCTPConn left SCTP_NODELAY unset %d of 100 times when the stack moved", fixed)
	}
	if canary == 0 {
		t.Fatal("the uintptr Setsockopt never failed: the stack is not being moved, so this test proves nothing")
	}
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
