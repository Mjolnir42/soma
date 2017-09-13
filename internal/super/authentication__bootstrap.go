/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mjolnir42/scrypth64"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/auth"
	"github.com/satori/go.uuid"
)

func (s *Supervisor) bootstrapRoot(q *msg.Request) {
	result := msg.FromRequest(q)
	kexID := q.Super.Encrypted.KexID
	data := q.Super.Encrypted.Data
	var kex *auth.Kex
	var err error
	var plain []byte
	var token auth.Token
	var rootToken string
	var mcf scrypth64.Mcf
	var tx *sql.Tx
	var validFrom, expiresAt time.Time
	var timer *time.Timer

	// bootstrapRoot is a master instance function
	if s.readonly {
		result.Conflict(fmt.Errorf(`Readonly instance`))
		goto conflict
	}

	// start response timer
	timer = time.NewTimer(1 * time.Second)
	defer timer.Stop()

	// -> check if root is not already active
	if s.credentials.read(`root`) != nil {
		result.BadRequest(fmt.Errorf(`Root account is already active`))
		//    --> delete kex
		s.kex.remove(kexID)
		goto dispatch
	}
	// -> get kex
	if kex = s.kex.read(kexID); kex == nil {
		//    --> reply 404 if not found
		result.NotFound(fmt.Errorf(`Key exchange not found`))
		goto dispatch
	}
	// -> check kex.SameSource
	if !kex.IsSameSourceExtractedString(q.RemoteAddr) {
		//    --> reply 404 if !SameSource
		result.NotFound(fmt.Errorf(`Key exchange not found`))
		goto dispatch
	}
	// -> delete kex from s.kex (kex is now used)
	s.kex.remove(kexID)
	// -> rdata = kex.DecodeAndDecrypt(data)
	if err = kex.DecodeAndDecrypt(&data, &plain); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> json.Unmarshal(rdata, &token)
	if err = json.Unmarshal(plain, &token); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> check token.UserName == `root`
	if token.UserName != `root` {
		//    --> reply 401
		result.Unauthorized(nil)
		goto dispatch
	}
	if token.UserName == `root` && s.rootRestricted && !q.Super.RestrictedEndpoint {
		result.ServerError(
			fmt.Errorf(`Root bootstrap requested on unrestricted endpoint`))
		goto dispatch
	}
	// -> check token.Token is correct bearer token
	if rootToken, err = s.fetchRootToken(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if token.Token != rootToken || len(token.Password) == 0 {
		//    --> reply 401
		result.Unauthorized(nil)
		goto dispatch
	}
	// -> scrypth64.Digest(Password, nil)
	if mcf, err = scrypth64.Digest(token.Password, nil); err != nil {
		result.Unauthorized(nil)
		goto dispatch
	}
	// -> generate token
	token.SetIPAddressExtractedString(q.RemoteAddr)
	if err = token.Generate(mcf, s.key, s.seed); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	validFrom, _ = time.Parse(msg.RFC3339Milli, token.ValidFrom)
	expiresAt, _ = time.Parse(msg.RFC3339Milli, token.ExpiresAt)

	// -> DB Insert: root password data
	if tx, err = s.conn.Begin(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	defer tx.Rollback()
	if _, err = tx.Exec(
		stmt.SetRootCredentials,
		uuid.Nil,
		mcf.String(),
		validFrom.UTC(),
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> DB Insert: token data
	if _, err = tx.Exec(
		stmt.InsertToken,
		token.Token,
		token.Salt,
		validFrom.UTC(),
		expiresAt.UTC(),
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> s.credentials Update
	s.credentials.insert(`root`, uuid.Nil, validFrom.UTC(),
		msg.PosTimeInf.UTC(), mcf)
	// -> s.tokens Update
	if err = s.tokens.insert(token.Token, token.ValidFrom, token.ExpiresAt,
		token.Salt); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if err = tx.Commit(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> sdata = kex.EncryptAndEncode(&token)
	plain = []byte{}
	data = []byte{}
	if plain, err = json.Marshal(token); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if err = kex.EncryptAndEncode(&plain, &data); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> send sdata reply
	result.Super = &msg.Supervisor{
		Verdict: 200,
		Encrypted: struct {
			KexID string
			Data  []byte
		}{
			Data: data,
		},
	}
	result.OK()

dispatch:
	<-timer.C

conflict:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
