// Copyright 2013 Alexandre Fiori
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"log"
	"time"

	diam "github.com/fiorix/go-diameter/stack"
)

const ADDR = ":3868"

func main() {
	// base.Dict.LoadFile("custom-dict.xml")
	go func() {
		time.Sleep(1 * time.Second)
		log.Println("Connecting to", ADDR)
		c, err := diam.Dial(ADDR, nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		SendCER(c)
	}()
	log.Println("Starting server on", ADDR)
	diam.ListenAndServe(ADDR, nil, nil)
}
