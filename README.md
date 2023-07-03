# Diameter Base Protocol

[![Go Report Card](https://goreportcard.com/badge/github.com/fiorix/go-diameter)](https://goreportcard.com/report/github.com/fiorix/go-diameter)
[![Test Status](https://github.com/fiorix/go-diameter/actions/workflows/test.yaml/badge.svg)](https://github.com/fiorix/go-diameter/actions/workflows/test.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/fiorix/go-diameter/v4.svg)](https://pkg.go.dev/github.com/fiorix/go-diameter/v4)
[![Latest](https://img.shields.io/github/v/tag/fiorix/go-diameter.svg?sort=semver&style=flat-square&label=latest)](https://github.com/fiorix/go-diameter/tags)

Package [go-diameter](https://pkg.go.dev/github.com/fiorix/go-diameter/v4) is an
implementation of the
Diameter Base Protocol [RFC 6733](http://tools.ietf.org/html/rfc6733)
and a stack for the [Go programming language](https://go.dev/).

### Status

The current implementation is solid and works fine for general purpose
clients and servers. It can send and receive messages efficiently as
well as build and parse AVPs based on dictionaries.

See the API documentation at https://pkg.go.dev/github.com/fiorix/go-diameter/v4

## Features

- Comprehensive XML dictionary format
- Embedded dictionaries:
  	* Base Protocol [RFC 6733](https://tools.ietf.org/html/rfc6733)
  	* Credit Control [RFC 4006](http://tools.ietf.org/html/rfc4006)
  	* Network Access Server [RFC 7155](http://tools.ietf.org/html/rfc7155)
  	* 3GPP specific AVPs from [TS 32.299 version 12.7.0](http://www.etsi.org/deliver/etsi_ts/132200_132299/132299/12.07.00_60/ts_132299v120700p.pdf)
  	* 3GPP S6a specific commands and AVPs from
  	  	[RFC 5516](https://tools.ietf.org/html/rfc5516) and
  	  	[TS 129 272](http://www.etsi.org/deliver/etsi_ts/129200_129299/129272/10.09.00_60/ts_129272v100900p.pdf)
- Human readable AVP representation (for debugging)
- TLS, IPv4 and IPv6 support for both clients and servers
- Stack based on [net/http](https://pkg.go.dev/net/http) for simplicity
- Ships with sample client, server, snoop agent and benchmark tool
- [State machines](http://tools.ietf.org/html/rfc6733#section-5.6) for CER/CEA and DWR/DWA for clients and servers
- TCP and SCTP support. SCTP support relies on kernel SCTP implementation and external github.com/ishidawataru/sctp
  package and is currently tested and enabled for Go 1.8+ and x86 Linux
  
## Getting started

The easiest way to get started is by trying out the client and server example programs.

With Go 1.11 and newer (preferred), you can start the client and server already:

```
export GO111MODULE=on
go run github.com/fiorix/go-diameter/v4/examples/server
go run github.com/fiorix/go-diameter/v4/examples/client -hello
```

Without modules, use standard procedure:

```
go get github.com/fiorix/go-diameter/examples/...
go run github.com/fiorix/go-diameter/examples/server
go run github.com/fiorix/go-diameter/examples/client -hello
```

Source code is your best friend. Check out other examples and test cases.

## Performance

Clients and servers written with the go-diameter package can be quite
performant if done well. Besides Go benchmarks, the package ships with
a simple benchmark tool to help testing servers and identifying bottlenecks.

In the examples directory, the server has a pprof (http server) that
allows the `go pprof` tool to profile the server in real time. The client
can perform benchmarks using the `-bench` command line flag.

For better performance, avoid logging diameter messages. Although logging
is very useful for debugging purposes, it kills performance due to a number
of conversions to make messages look pretty. If you run benchmarks on the
example server, make sure to use the `-s` (silent) command line switch.

TLS degrades performance a bit, as well as reflection (Unmarshal). Those are
important trade offs you might have to consider.

Besides this, the source code (and sub-packages) have function benchmarks
that can help you understand what's fast and isn't. You will see that
parsing messages is much slower than writing them, for example. This is
because in order to parse messages it makes numerous dictionary lookups
for AVP types, to be able to decode them. Encoding messages require less
lookups and is generally simpler, thus faster.

## Contribute

In case you want to add new AVPs, please add them to `diam/dict/testdata` xml
files. Then regenerate the go models using `./autogen.sh` you will find at 
`diam` folder. This will modify files at `diam/dict` to include your changes.

Before submitting PR, please run `make test` to test your changes. Or do it 
manually:

```
	go test ./...
```

You also have the option to run the test using a Linux VM through Docker (this
is not mandatory). To do so, run `make test_docker`. Runing test on Linux  can 
be useful in case you add `sctp` tests. Note you will need to install
docker and docker-compose.
