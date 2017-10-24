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

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
)

// IMPORTANT!
//
// all errors returned from encrypted supervisor methods are
// returned to the client as 403/Forbidden. Provided error details
// are used only for serverside logging.

func (s *Supervisor) decrypt(q *msg.Request, mr *msg.Result, audit *logrus.Entry) (*auth.Token, bool) {
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
		audit.WithField(`Code`, mr.Super.Verdict).
			Warningln(str)
		return nil, false
	}

	// check KeyExchange is used by the same source that negotiated it
	if !kex.IsSameSourceExtractedString(q.RemoteAddr) {
		str := `KexID referenced from wrong source system`

		mr.BadRequest(fmt.Errorf(str), q.Section)
		audit.WithField(`Code`, mr.Super.Verdict).
			Errorln(str)
		return nil, false
	}

	// KeyExchanges are single-use and this KexID now has been used,
	// remove it.
	s.kex.remove(q.Super.Encrypted.KexID)

	// attempt decrypting the request data
	if err = kex.DecodeAndDecrypt(&q.Super.Encrypted.Data,
		&plain); err != nil {
		mr.ServerError(err)
		audit.WithField(`Code`, mr.Super.Verdict).
			Warningln(err)
		return nil, false
	}

	// unmarshal the decrypted request data into a auth.Token protocol datastructure
	if err = json.Unmarshal(plain, token); err != nil {
		mr.ServerError(err)
		audit.WithField(`Code`, mr.Super.Verdict).
			Warningln(err)
		return nil, false
	}

	return token, true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
