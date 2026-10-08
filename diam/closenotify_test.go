// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"net"
	"sync"
	"testing"
)

func newCloseNotifyConn(t *testing.T) (*conn, net.Conn) {
	t.Helper()
	local, peer := net.Pipe()
	t.Cleanup(func() {
		local.Close()
		peer.Close()
	})
	c, err := new(Server).newConn(local)
	if err != nil {
		t.Fatal(err)
	}
	return c, peer
}

func TestCloseNotifyBeforeClose(t *testing.T) {
	c, peer := newCloseNotifyConn(t)
	closed := c.writer.CloseNotify()
	select {
	case <-closed:
		t.Fatal("CloseNotify fired before the peer closed")
	default:
	}
	if got := c.writer.CloseNotify(); got != closed {
		t.Fatal("CloseNotify returned a different channel before close")
	}

	peer.Close()
	c.serve()
	select {
	case <-closed:
	default:
		t.Fatal("CloseNotify did not fire after the connection finished serving")
	}
	if got := c.writer.CloseNotify(); got != closed {
		t.Fatal("CloseNotify returned a different channel after close")
	}
}

func TestCloseNotifyAfterClose(t *testing.T) {
	c, peer := newCloseNotifyConn(t)
	peer.Close()
	// Running serve synchronously ensures its deferred notification has
	// completed before the first CloseNotify call, without arming it early.
	c.serve()
	closed := c.writer.CloseNotify()
	select {
	case <-closed:
	default:
		t.Fatal("CloseNotify called after the connection finished serving did not fire")
	}
	if got := c.writer.CloseNotify(); got != closed {
		t.Fatal("CloseNotify returned a different channel after close")
	}
}

func TestCloseNotifyConcurrent(t *testing.T) {
	c, peer := newCloseNotifyConn(t)
	peer.Close()
	const callers = 32
	channels := make(chan (<-chan struct{}), callers)
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		c.serve()
	}()
	for i := 0; i < callers; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			channels <- c.writer.CloseNotify()
		}()
		go func() {
			defer wg.Done()
			<-start
			// Both the reader and serve can report the same disconnect.
			c.notifyClientGone()
		}()
	}
	close(start)
	wg.Wait()
	close(channels)

	want := c.writer.CloseNotify()
	for closed := range channels {
		if closed != want {
			t.Error("concurrent CloseNotify calls returned different channels")
		}
		select {
		case <-closed:
		default:
			t.Error("concurrent CloseNotify caller missed the close")
		}
	}
}
