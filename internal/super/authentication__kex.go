/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"fmt"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
)

func (s *Supervisor) kexInit(q *msg.Request) {
	result := msg.FromRequest(q)
	result.Action = `reply`
	kex := q.Super.Kex
	var err error

	// kexInit is a master instance function
	if s.readonly {
		result.Conflict(fmt.Errorf(`Readonly instance`))
		goto dispatch
	}

	// check the client submitted IV for fishyness
	err = s.checkIV(kex.InitializationVector)
	for err != nil {
		if err = kex.GenerateNewVector(); err != nil {
			continue
		}
		err = s.checkIV(kex.InitializationVector)
	}

	// record the kex submission time
	kex.SetTimeUTC()

	// record the client ip address
	kex.SetIPAddressExtractedString(q.RemoteAddr)

	// generate a request ID
	kex.GenerateNewRequestID()

	// set the client submitted public key as peer key
	kex.SetPeerKey(kex.PublicKey())

	// generate our own keypair
	kex.GenerateNewKeypair()

	// save kex
	s.kex.insert(kex)

	// send out reply
	result.Super = &msg.Supervisor{
		Verdict: 200,
		Kex: auth.Kex{
			Public:               kex.Public,
			InitializationVector: kex.InitializationVector,
			Request:              kex.Request,
		},
	}
	result.OK()

dispatch:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
