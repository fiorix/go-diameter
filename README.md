# Diameter Base Protocol

Package [go-diameter](http://godoc.org/github.com/fiorix/go-diameter) is an
implementation of the
Diameter Base Protocol [RFC 6733](http://tools.ietf.org/html/rfc6733)
and a stack for the [Go programming language](http://golang.org).


### Status

It can currently send and receive messages and build and parse AVPs based on
dictionaries. API is subject to changes. See the [TODO](./TODO) list for
details.

See the API documentation at http://godoc.org/github.com/fiorix/go-diameter


[![Build Status](https://secure.travis-ci.org/fiorix/go-diameter.png)](http://travis-ci.org/fiorix/go-diameter)


## Features

- Comprehensive XML dictionary format
- Embedded dictionaries (base protocol and credit control [RFC 4006](http://tools.ietf.org/html/rfc4006))
- Human readable AVP representation
- TLS support for both clients and servers
- Stack based on [net/http](http://golang.org/pkg/net/http/) for simplicity
- Ships with sample client, server, snoop agent and benchmark tool

For now, [state machines](http://tools.ietf.org/html/rfc6733#section-4) are
not part of this implementation on purpose.


## Install

Make sure Go is installed, and both $GOPATH and $GOROOT are set.

Install:

	go get github.com/fiorix/go-diameter/diam

Check out the examples:

	cd $GOPATH/src/github.com/fiorix/go-diameter/examples

See the test cases for more specific examples.


## Performance

Clients and servers written with the go-diameter package can be quite
performant if done well. There's a simple benchmark tool to help testing
servers and identifying bottlenecks.

The results below are from two Intel i5 quad-core with 4GB ram on a 1Gbps
network, with 4 concurrent clients hammering the example server:

2014/05/11 10:32:26 1000000 messages in 10.449877643s seconds, 95694 msg/s
2014/05/11 10:32:26 1000000 messages in 10.604832166s seconds, 94296 msg/s
2014/05/11 10:32:26 1000000 messages in 10.699197777s seconds, 93464 msg/s
2014/05/11 10:32:27 1000000 messages in 10.727248401s seconds, 93220 msg/s

For better performance, avoid printing diameter messages to the log.
Although they're very useful for debugging purposes, they kill performance
due to a number of convertions to make them pretty. If you run benchmarks
on the example server, make sure to use the `-q` (quiet) command line switch.

TLS degrades performance a bit, as well as reflection (Unmarshal). Those are
important tradeoffs you might have to consider.

Besides this, the source code (and sub-packages) has function benchmarks
that can help you understand what's fast and what's not. You'll see that
parsing messages is much slower than writing them, for example. This is
because in order to parse messages it makes numerous dictionary lookups
for AVP types, to be able to decode them. Encoding messages require less
lookups and is generally simpler, thus faster.
