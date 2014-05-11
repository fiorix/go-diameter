// Copyright 2013-2014 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

/*
Package go-diameter provides support for the diameter protocol for Go.
See RFC 6733.

go-diameter is minimalist implementation of the diameter base protocol,
organized in sub-packages with specific functionality:

 * diam: the main package, provides the capability of encoding and
 decoding messages, and a client and server API similar to net/http.
 * diamdict: a dictionary parser that can combine dictionaries.
 * diamtype: the diameter data types used in messaging.

If you're looking to go right into code, see the examples subdirectory for
applications like clients and servers.


Diameter Applications

All diameter applications require at least the following:

 * A dictionary with the application id, its commands and message formats
 * A program that implements the application, driven by the dictionary

The diamdict sub-package supports the base application (id 0) and
the credit control application (id 4), as per RFC 4006. Each application
has its own commands and messages, and their pre-defined AVPs.

AVPs are of specific types, which are the diameter types like UTF8String,
Unsigned32 and so on. Fortunately, those types map 1:1 with Go types,
which makes things easier for us. However, the diameter types have
specific properties like padding for certain strings, which have to
be taken care of. The sub-package diamtype handles it all.

At last, the main package is used to build clients and servers using
an API very similar to the one of net/http. To initiate the client or
server, you'll have to pass a dictionary. Messages sent and received
are encoded and decoded using the dictionary automatically.

The API of clients and servers require that you assign handlers for
certain messages, similar to how you route HTTP endpoints. In the
handlers, you'll receive messages already decoded, but their data
is always encoded by the types of diamtype.

To reply to those messages, or to create new ones any time, you'll
have to use the diameter types.
*/
package diam
