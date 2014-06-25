// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Sn🔍 op agent can sit in between two diameter peers and snoop all messages
// in real time, printing them to the console.
//
// It's a simple 1:1 proxy.

package main

import (
	"flag"
	"log"
	"runtime"
	"strings"
	"sync"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/diamdict"
)

type Bridge struct {
	Client chan *diam.Message // Remote client connecting to this server
	Server chan *diam.Message // Upstream connection (real server)
}

var (
	UpstreamAddr string
	LiveMu       sync.RWMutex
	Live         = make(map[string]*Bridge) // ip:bridge
)

func main() {
	local := flag.String("local", ":3868", "set local addr")
	remote := flag.String("remote", "", "set remote addr")
	files := flag.String("dict", "", "comma separated list of dictionaries")
	flag.Parse()
	UpstreamAddr = *remote
	log.Println("Diameter sn🔍 op agent")
	if len(*remote) == 0 {
		log.Fatal("Missing argument -remote")
	}
	if *local == *remote {
		log.Fatal("Local and remote address are the same. Duh?")
	}
	// Load dictionary files onto the default (base protocol) diamdict.
	if *files != "" {
		for _, f := range strings.Split(*files, ",") {
			log.Println("Loading dictionary", f)
			if err := diamdict.Default.LoadFile(f); err != nil {
				log.Fatal(err)
			}
		}
	}
	// Use all CPUs.
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Prepare the server.
	diam.HandleFunc("ALL", func(c diam.Conn, m *diam.Message) {
		// Forward incoming messages to the upstream server.
		if b := GetBridge(c); b != nil {
			b.Server <- m
		} else {
			// Upstream server unavailable, bye.
			c.Close()
		}
	})
	// Start the server using default handler and dict.
	log.Printf("Starting server on %s", *local)
	diam.ListenAndServe(*local, nil, nil)
}

// GetBridge returns the Bridge object for a given client, if it exists.
// Otherwise GetBridge connects to the upstream server and set up the
// bridge with the client, returning the newly created Bridge object.
func GetBridge(c diam.Conn) *Bridge {
	LiveMu.RLock()
	if b, ok := Live[c.RemoteAddr().String()]; ok {
		LiveMu.RUnlock()
		return b
	}
	LiveMu.RUnlock()
	LiveMu.Lock()
	defer LiveMu.Unlock()
	b := &Bridge{
		Client: make(chan *diam.Message),
		Server: make(chan *diam.Message),
	}
	Live[c.RemoteAddr().String()] = b
	// Prepare for the upstream connection.
	mux := diam.NewServeMux()
	mux.HandleFunc("ALL", func(c diam.Conn, m *diam.Message) {
		// Forward incoming messages to the client.
		b.Client <- m
	})
	// Connect to upstream server.
	s, err := diam.Dial(UpstreamAddr, mux, nil)
	if err != nil {
		return nil
	}
	log.Printf("Creating bridge from %s to %s",
		c.RemoteAddr().String(), s.RemoteAddr().String())
	go Pump(c, s, b.Client, b.Server)
	go Pump(s, c, b.Server, b.Client)
	return b
}

// Pump messages from one side to the other.
func Pump(src, dst diam.Conn, srcChan, dstChan chan *diam.Message) {
	for {
		select {
		case m := <-srcChan:
			if m == nil {
				src.Close()
				return
			}
			log.Printf(
				"Message from %s to %s\n%s",
				src.RemoteAddr().String(),
				dst.RemoteAddr().String(),
				m,
			)
			if _, err := m.WriteTo(src); err != nil {
				src.Close() // triggers the case below
			}
		case <-src.(diam.CloseNotifier).CloseNotify():
			LiveMu.Lock()
			defer LiveMu.Unlock()
			if _, ok := Live[src.RemoteAddr().String()]; ok {
				delete(Live, src.RemoteAddr().String())
				log.Printf(
					"Destroying bridge from %s to %s",
					src.RemoteAddr().String(),
					dst.RemoteAddr().String(),
				)
			} else {
				delete(Live, dst.RemoteAddr().String())
				log.Printf(
					"Destroying bridge from %s to %s",
					dst.RemoteAddr().String(),
					src.RemoteAddr().String(),
				)
			}
			src.Close()
			dstChan <- nil
			return
		}
	}
}
