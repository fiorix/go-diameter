# Diameter Base Protocol

Package [go-diameter](http://godoc.org/github.com/fiorix/go-diameter) is an
implementation of the
Diameter Base Protocol [rfc6733](http://tools.ietf.org/html/rfc6733)
and a stack for the [Go programming language](http://golang.org).

### Status

- v0.5

It can currently send and receive messages and build and parse AVPs based on
dictionaries. API is subject to changes. See the [TODO](./TODO) list for
details.

[![Build Status](https://secure.travis-ci.org/fiorix/go-diameter.png)](http://travis-ci.org/fiorix/go-diameter)

## Features

- Comprehensive XML dictionary format
- Embedded base protocol dictionary (one less file to carry around)
- Human readable AVP representation
- TLS support
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

See the [API reference](http://godoc.org/github.com/fiorix/go-diameter) for
details.
