# Diameter Base Protocol

Package [diameter](http://godoc.org/github.com/fiorix/go-diameter) is an
implementation of the
[Diameter Base Protocol (rfc3588)](http://tools.ietf.org/html/rfc3588)
for the [Go programming language](http://golang.org).

## Status

It currently has a dictionary loader and incomplete (yet functional)
message header and AVP parser. Can't write messages.

## Features

- Simple XML dictionary format
- Embedded base protocol dictionary (one less file to carry around)
- (almost) Human readable AVP representation
- Dictionary AVPs are can be overloaded
