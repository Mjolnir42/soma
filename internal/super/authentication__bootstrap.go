/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mjolnir42/scrypth64"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/auth"
	"github.com/satori/go.uuid"
)

func (s *Supervisor) bootstrap(q *msg.Request) {
	result := msg.FromRequest(q)
	result.Super.Verdict = 401
	q.Log(s.reqLog)

	kexID := q.Super.Encrypted.KexID
	var (
		kex                  *auth.Kex
		err                  error
		token                *auth.Token
		rootToken            string
		mcf                  scrypth64.Mcf
		tx                   *sql.Tx
		validFrom, expiresAt time.Time
		timer                *time.Timer
		ok                   bool
	)

	// bootstrapRoot is a master instance function
	if s.readonly {
		result.ReadOnly()
		goto returnImmediate
	}

	// start response timer
	timer = time.NewTimer(1 * time.Second)
	defer timer.Stop()

	// check if root is not already active
	if s.credentials.read(msg.SubjectRoot) != nil {
		result.Forbidden(fmt.Errorf(`Root account is already active`))
		// delete kex, this is done
		s.kex.remove(kexID)
		goto returnImmediate
	}

	// decrypt e2e encrypted request
	if token, kex, ok = s.decrypt(q, &result); !ok {
		return
	}

	if token.UserName != msg.SubjectRoot {
		result.Forbidden(nil)
		goto dispatch
	}

	if s.rootRestricted && !q.Super.RestrictedEndpoint {
		result.ServerError(fmt.Errorf(
			`Root bootstrap requested on unrestricted endpoint`))
		goto dispatch
	}

	if rootToken, err = s.fetchRootToken(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if token.Token != rootToken || len(token.Password) == 0 {
		result.Forbidden(nil)
		goto dispatch
	}

	// generate password hash that the server will store
	if mcf, err = scrypth64.Digest(token.Password, nil); err != nil {
		result.ServerError(nil)
		goto dispatch
	}

	// generate a token for the user that can be used for further
	// requests
	token.SetIPAddressExtractedString(q.RemoteAddr)
	if err = token.Generate(mcf, s.key, s.seed); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	validFrom, _ = time.Parse(msg.RFC3339Milli, token.ValidFrom)
	expiresAt, _ = time.Parse(msg.RFC3339Milli, token.ExpiresAt)

	if tx, err = s.conn.Begin(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	defer tx.Rollback()

	// insert hashed root password
	if _, err = tx.Exec(
		stmt.SetRootCredentials,
		uuid.Nil,
		mcf.String(),
		validFrom.UTC(),
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}

	// insert generated token
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

	// update credential store
	s.credentials.insert(
		msg.SubjectRoot,
		uuid.Nil,
		validFrom.UTC(),
		msg.PosTimeInf.UTC(),
		mcf,
	)

	// update token store
	if err = s.tokens.insert(
		token.Token,
		token.ValidFrom,
		token.ExpiresAt,
		token.Salt,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}

	if err = tx.Commit(); err != nil {
		result.ServerError(err)
		goto dispatch
	}

	// send encrypted reply
	if err = s.encrypt(kex, token, &result); err != nil {
		result.ServerError(err, q.Section)
		goto dispatch
	}
	// XXX BUG: write out auditlog entry

dispatch:
	<-timer.C

returnImmediate:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
