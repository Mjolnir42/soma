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

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/auth"
)

func (s *Supervisor) issueToken(q *msg.Request) {
	result := msg.FromRequest(q)
	result.Super.Verdict = 401
	var (
		cred                 *svCredential
		err                  error
		kex                  *auth.Kex
		plain                []byte
		timer                *time.Timer
		token                auth.Token
		tx                   *sql.Tx
		validFrom, expiresAt time.Time
	)
	data := q.Super.Encrypted.Data

	// issue_token is a master instance function
	if s.readonly {
		result.ReadOnly()
		goto returnImmediate
	}
	// start response timer
	timer = time.NewTimer(1 * time.Second)
	defer timer.Stop()

	// -> get kex
	if kex = s.kex.read(q.Super.Encrypted.KexID); kex == nil {
		result.NotFound(fmt.Errorf(`Key exchange not found`))
		goto dispatch
	}
	// check kex.SameSource
	if !kex.IsSameSourceExtractedString(q.RemoteAddr) {
		result.NotFound(fmt.Errorf(`Key exchange not found`))
		goto dispatch
	}
	// delete kex from s.kex (kex is now used)
	s.kex.remove(q.Super.Encrypted.KexID)
	// decrypt request
	if err = kex.DecodeAndDecrypt(&data, &plain); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> json.Unmarshal(rdata, &token)
	if err = json.Unmarshal(plain, &token); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if token.UserName == `root` && s.rootRestricted && !q.Super.RestrictedEndpoint {
		result.ServerError(
			fmt.Errorf(`Root token requested on unrestricted endpoint`))
		goto dispatch
	}

	s.reqLog.Printf(msg.LogStrSRq, q.Section, q.Action, token.UserName, q.RemoteAddr)

	if cred = s.credentials.read(token.UserName); cred == nil {
		result.Unauthorized(fmt.Errorf("Unknown user: %s", token.UserName))
		goto dispatch
	}
	if !cred.isActive {
		result.Unauthorized(fmt.Errorf("Inactive user: %s", token.UserName))
		goto dispatch
	}
	if time.Now().UTC().Before(cred.validFrom.UTC()) ||
		time.Now().UTC().After(cred.expiresAt.UTC()) {
		result.Unauthorized(fmt.Errorf("Expired: %s", token.UserName))
		goto dispatch
	}
	// generate token if the provided credentials are valid
	token.SetIPAddressExtractedString(q.RemoteAddr)
	if err = token.Generate(cred.cryptMCF, s.key, s.seed); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	validFrom, _ = time.Parse(msg.RFC3339Milli, token.ValidFrom)
	expiresAt, _ = time.Parse(msg.RFC3339Milli, token.ExpiresAt)
	// -> DB Insert: token data
	if tx, err = s.conn.Begin(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	defer tx.Rollback()
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
	result.Super = msg.Supervisor{
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

returnImmediate:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
