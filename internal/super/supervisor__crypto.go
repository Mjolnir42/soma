/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"encoding/json"
	"fmt"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
)

// IMPORTANT!
//
// all errors returned from encrypted supervisor methods are
// returned to the client as 403/Forbidden. Provided error details
// are used only for serverside logging.

// decrypt returns the decrypted auth.Token embedded in msg.Request
func (s *Supervisor) decrypt(q *msg.Request, mr *msg.Result) (*auth.Token, *auth.Kex, bool) {
	var (
		err   error
		kex   *auth.Kex
		plain []byte
		token *auth.Token
	)

	// lookup requested KeyExchange by provided KeyExchangeID
	if kex = s.kex.read(q.Super.Encrypted.KexID); kex == nil {
		str := `Key exchange not found`
		mr.NotFound(fmt.Errorf(str), q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return nil, nil, false
	}

	// check KeyExchange is used by the same source that negotiated it
	if !kex.IsSameSourceExtractedString(q.RemoteAddr) {
		str := `KexID referenced from wrong source system`
		mr.BadRequest(fmt.Errorf(str), q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Errorln(mr.Error)
		return nil, nil, false
	}

	// KeyExchanges are single-use and this KexID now has been used,
	// remove it.
	s.kex.remove(q.Super.Encrypted.KexID)

	// attempt decrypting the request data
	if err = kex.DecodeAndDecrypt(
		&q.Super.Encrypted.Data,
		&plain,
	); err != nil {
		mr.ServerError(err)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return nil, nil, false
	}

	// unmarshal the decrypted request data into a auth.Token protocol datastructure
	if err = json.Unmarshal(plain, token); err != nil {
		mr.ServerError(err)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return nil, nil, false
	}

	return token, kex, true
}

// encrypt embeds the encrypted token into mr
func (s *Supervisor) encrypt(kex *auth.Kex, token *auth.Token, mr *msg.Result) error {
	var plain, data []byte
	var err error

	// serialize token
	if plain, err = json.Marshal(token); err != nil {
		return err
	}

	// encrypt serialized token
	if err = kex.EncryptAndEncode(&plain, &data); err != nil {
		return err
	}

	// update result for dispatching encrypted result data
	mr.Super.Verdict = 200
	mr.Super.Encrypted.Data = data
	mr.OK()
	mr.Super.Audit = mr.Super.Audit.WithField(`Verdict`, mr.Super.Verdict).
		WithField(`Code`, mr.Code)
	return nil
}

// kexInit handles key exchange requests
func (s *Supervisor) kexInit(q *msg.Request) {
	result := msg.FromRequest(q)
	result.Super.Verdict = 401
	q.Log(s.reqLog)

	kex := q.Super.Kex
	var err error
	var attemptIV int

	// kexInit is a master instance function
	if s.readonly {
		result.ReadOnly()
		goto dispatch
	}

	// generate new initialization vector
	err = kex.GenerateNewVector()
	for err != nil {
		attemptIV++
		if attemptIV > 5 {
			result.ServerError(err, q.Section)
			goto dispatch
		}
		err = kex.GenerateNewVector()
	}

	// check no bad IV was generated, with bad being defined
	// as likely to use 0 as counter
	err = s.checkIV(kex.InitializationVector)
	for err != nil {
		attemptIV++
		if attemptIV > 5 {
			result.ServerError(err, q.Section)
			goto dispatch
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
	result.Super = msg.Supervisor{
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
