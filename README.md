# Diameter Base Protocol

Package [diameter](http://godoc.org/github.com/fiorix/go-diameter) is an
implementation of the
Diameter Base Protocol [rfc6733](http://tools.ietf.org/html/rfc6733)
and a stack for the [Go programming language](http://golang.org).

### Status

- v0.3

It can currently send and receive messages, build and parse AVPs based on
dictionaries. API is subject to changes. See the [TODO](./TODO) list for
details.

[![Build Status](https://secure.travis-ci.org/fiorix/go-diameter.png)](http://travis-ci.org/fiorix/go-diameter)

## Modules

- [base](./base): Diameter Base Protocol
- [dict](./dict): XML Dictionary Parser
- [stack](./stack): Server stack for handling multiple clients and requests.

## Features

- Comprehensive XML dictionary format
- Embedded base protocol dictionary (one less file to carry around)
- Human readable AVP representation
- TLS support
- Stack based on [net/http](http://golang.org/pkg/net/http/) for simplicity
