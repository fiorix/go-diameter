// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package peer

import (
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/sm/command"
	"golang.org/x/net/context"
)

type key int

const metadataKey key = 0

// Metadata contains information about a diameter peer, acquired
// during the CER/CEA handshake.
type Metadata struct {
	OriginHost   datatype.DiameterIdentity
	OriginRealm  datatype.DiameterIdentity
	Applications []uint32 // Acct or Auth IDs supported by the peer.
}

// FromCER creates a Metadata object from data in the CER.
func FromCER(cer *command.CER) *Metadata {
	return &Metadata{
		OriginHost:   cer.OriginHost,
		OriginRealm:  cer.OriginRealm,
		Applications: cer.Applications(),
	}
}

// NewContext returns a new Context that carries a Metadata object.
func NewContext(ctx context.Context, metadata *Metadata) context.Context {
	return context.WithValue(ctx, metadataKey, metadata)
}

// FromContext extracts a Metadata object from the context.
func FromContext(ctx context.Context) (*Metadata, bool) {
	meta, ok := ctx.Value(metadataKey).(*Metadata)
	return meta, ok
}
