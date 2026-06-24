// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpeer

import (
	"context"

	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm/smparser"
)

type key int

const metadataKey key = 0

// Metadata contains information about a diameter peer, acquired
// during the CER/CEA handshake.
type Metadata struct {
	OriginHost   datatype.DiameterIdentity
	OriginRealm  datatype.DiameterIdentity
	Applications []uint32 // Acct or Auth IDs supported by the peer.

	// Cer and Cea retain the full parsed handshake message so callers can
	// read capability AVPs (e.g. Vendor-Id, Product-Name, Supported-Vendor-Id)
	// that are not promoted to the fields above. Exactly one is set, matching
	// which constructor was used: Cer for inbound peers, Cea for outbound.
	Cer *smparser.CER
	Cea *smparser.CEA
}

// FromCER creates a Metadata object from data in the CER.
func FromCER(cer *smparser.CER) *Metadata {
	return &Metadata{
		OriginHost:   cer.OriginHost,
		OriginRealm:  cer.OriginRealm,
		Applications: cer.Applications(),
		Cer:          cer,
	}
}

// FromCEA creates a Metadata object from data in the CEA.
func FromCEA(cea *smparser.CEA) *Metadata {
	return &Metadata{
		OriginHost:   cea.OriginHost,
		OriginRealm:  cea.OriginRealm,
		Applications: cea.Applications(),
		Cea:          cea,
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
