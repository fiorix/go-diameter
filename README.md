# Diameter Base Protocol

Package [diameter](http://godoc.org/github.com/fiorix/go-diameter) is an
implementation of a stack and
Diameter Base Protocol [rfc6733](http://tools.ietf.org/html/rfc6733)
for the [Go programming language](http://golang.org).

### Status

It can parse dictionaries and read and write messages, but still needs a lot
of work.  See the [TODO](./TODO) list for details.

[![Build Status](https://secure.travis-ci.org/fiorix/go-diameter.png)](http://travis-ci.org/fiorix/go-diameter)

## Modules

- [base](./base): Diameter Base Protocol
- [dict](./dict): XML Dictionary Parser
- [stack](./stack): Server stack for handling multiple clients and requests.

## Features

- Comprehensive XML dictionary format
- Embedded base protocol dictionary (one less file to carry around)
- Human readable AVP representation
- TLS support (untested, but is there)
- Stack based on [net/http](http://golang.org/pkg/net/http/) for simplicity
