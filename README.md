# Diameter Base Protocol

Package [go-diameter](http://godoc.org/github.com/fiorix/go-diameter) is an
implementation of the
Diameter Base Protocol [RFC 6733](http://tools.ietf.org/html/rfc6733)
and a stack for the [Go programming language](http://golang.org).

[![GoDoc](https://godoc.org/github.com/fiorix/go-diameter?status.svg)](https://godoc.org/github.com/fiorix/go-diameter)

### Status

The current implementation is solid and works fine for general purpose
clients and servers. It can send and receive messages efficiently as
well as build and parse AVPs based on dictionaries.

See the API documentation at http://godoc.org/github.com/fiorix/go-diameter

[![Build Status](https://secure.travis-ci.org/fiorix/go-diameter.png)](http://travis-ci.org/fiorix/go-diameter)

## Features

- Comprehensive XML dictionary format
- Embedded dictionaries (base protocol and credit control [RFC 4006](http://tools.ietf.org/html/rfc4006))
- Human readable AVP representation (for debugging)
- TLS support for both clients and servers
- Stack based on [net/http](http://golang.org/pkg/net/http/) for simplicity
- Ships with sample client, server, snoop agent and benchmark tool
- [State machines](http://tools.ietf.org/html/rfc6733#section-4) for CER/CEA and DWR/DWA for clients and servers

## Install

go-diameter requires at least Go 1.4.

Make sure Go is installed, and both GOPATH and GOROOT are set.

Install:

	go get github.com/fiorix/go-diameter/diam

Check out the examples:

	cd $GOPATH/src/github.com/fiorix/go-diameter/examples

See the test cases for more specific examples.


## Performance

Clients and servers written with the go-diameter package can be quite
performant if done well. Besides Go benchmarks, the package ships with
a simple benchmark tool to help testing servers and identifying bottlenecks.

The results below are from two Intel i5 quad-core with 4GB ram on a 1Gbps
network, with 4 concurrent clients hammering the example server:

	2014/10/14 17:19:35 200000 messages (request+answer) in 1.044636932s seconds, 191454 msg/s
	2014/10/14 17:19:35 200000 messages (request+answer) in 1.051021203s seconds, 190291 msg/s
	2014/10/14 17:19:35 200000 messages (request+answer) in 1.050285029s seconds, 190424 msg/s
	2014/10/14 17:19:35 200000 messages (request+answer) in 1.070140824s seconds, 186891 msg/s
	2014/10/14 17:19:35 Total of 800000 messages in 1.076188492s: 743364 msg/s

For better performance, avoid printing diameter messages to the log.
Although they're very useful for debugging purposes, they kill performance
due to a number of conversions to make them pretty. If you run benchmarks
on the example server, make sure to use the `-q` (quiet) command line switch.

TLS degrades performance a bit, as well as reflection (Unmarshal). Those are
important trade offs you might have to consider.

Besides this, the source code (and sub-packages) have function benchmarks
that can help you understand what's fast and isn't. You will see that
parsing messages is much slower than writing them, for example. This is
because in order to parse messages it makes numerous dictionary lookups
for AVP types, to be able to decode them. Encoding messages require less
lookups and is generally simpler, thus faster.
