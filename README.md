# Diameter Base Protocol

Package [diameter](http://godoc.org/github.com/fiorix/go-diameter) is an
implementation of a stack and
Diameter Base Protocol [rfc6733](http://tools.ietf.org/html/rfc6733)
for the [Go programming language](http://golang.org).

## Status

It can parse dictionaries and read and write messages, but still needs a lot
of work.

## Features

This package implements an API based on [net/http](http://golang.org/pkg/net/http/)
aiming for simplicity of diameter message handling.

- Comprehensive XML dictionary format
- Embedded base protocol dictionary (one less file to carry around)
- Human readable AVP representation
- Dictionary AVPs are can be overloaded
- TLS support (untested, but is there)
