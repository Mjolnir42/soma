/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
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

// token handles supervisor requests for token and calls the correct
// function depending on the requested task
func (s *Supervisor) token(q *msg.Request) {
	result := msg.FromRequest(q)
	result.Super.Verdict = 403

	// start response delay timer
	timer := time.NewTimer(1 * time.Second)

	// tokenRequest/tokenInvalidate are master instance functions
	if s.readonly {
		result.ReadOnly()
		goto returnImmediate
	}

	// select correct taskhandler
	switch q.Super.Task {
	case msg.TaskRequest:
		s.tokenRequest(q, &result)
	case msg.TaskInvalidateAll:
		s.tokenInvalidateAll(q, &result)
	case msg.TaskInvalidate:
		s.tokenInvalidate(q, &result)
	default:
		result.UnknownRequest(q)
		goto returnImmediate
	}

	// wait for delay timer to trigger
	<-timer.C

returnImmediate:
	// cleanup delay timer
	if !timer.Stop() {
		<-timer.C
	}
	q.Reply <- result
}

// tokenRequest handles requests for new tokens to be issued
func (s *Supervisor) tokenRequest(q *msg.Request, mr *msg.Result) {
	var (
		cred                 *credential
		err                  error
		kex                  *auth.Kex
		token                *auth.Token
		ok                   bool
		tx                   *sql.Tx
		validFrom, expiresAt time.Time
	)

	// start assembly of auditlog entry
	logEntry := singleton.auditLog.
		WithField(`Type`, fmt.Sprintf("%s/%s:%s", q.Section, q.Action, q.Super.Task)).
		WithField(`RequestID`, q.ID.String()).
		WithField(`KexID`, q.Super.Encrypted.KexID).
		WithField(`IPAddr`, q.RemoteAddr)

	// decrypt e2e encrypted request
	if token, ok = s.decrypt(q, mr, logEntry); !ok {
		return
	}

	// check if root is available
	if token.UserName == `root` && s.rootRestricted && !q.Super.RestrictedEndpoint {
		mr.ServerError(
			fmt.Errorf(`Restricted-mode root token requested on unrestricted endpoint`))
		logEntry.WithField(`Code`, mr.Super.Verdict).Warningln(`Restricted-mode root token requested on unrestricted endpoint`)
		return
	}

	// check if the user exists
	if cred = s.credentials.read(token.UserName); cred == nil {
		mr.Forbidden(fmt.Errorf("Unknown user: %s", token.UserName))
		logEntry.WithField(`Code`, mr.Super.Verdict).Warningln(fmt.Errorf("Unknown user: %s", token.UserName))
		return
	}
	logEntry = logEntry.WithField(`User`, token.UserName)

	// check if the user is active
	if !cred.isActive {
		mr.Forbidden(fmt.Errorf("Inactive user: %s", token.UserName))
		logEntry.WithField(`Code`, mr.Super.Verdict).Warningln(`User is inactive`)
		return
	}

	// check if the token has either expired or not become valid yet
	if time.Now().UTC().Before(cred.validFrom.UTC()) ||
		time.Now().UTC().After(cred.expiresAt.UTC()) {
		mr.Forbidden(fmt.Errorf(`Token expired`))
		logEntry.WithField(`Code`, mr.Super.Verdict).Warningln(`Token expired`)
		return
	}

	// generate token if the provided credentials are valid
	token.SetIPAddressExtractedString(q.RemoteAddr)
	if err = token.Generate(cred.cryptMCF, s.key, s.seed); err != nil {
		mr.ServerError(err)
		logEntry.WithField(`Code`, mr.Super.Verdict).Warningln(err)
		return
	}

	// persist generated token into database
	validFrom, _ = time.Parse(msg.RFC3339Milli, token.ValidFrom)
	expiresAt, _ = time.Parse(msg.RFC3339Milli, token.ExpiresAt)
	if tx, err = s.conn.Begin(); err != nil {
		mr.ServerError(err)
		logEntry.WithField(`Code`, mr.Super.Verdict).Warningln(`Database error`)
		return
	}
	defer tx.Rollback()

	if _, err = tx.Exec(
		stmt.InsertToken,
		token.Token,
		token.Salt,
		validFrom.UTC(),
		expiresAt.UTC(),
	); err != nil {
		mr.ServerError(err)
		logEntry.WithField(`Code`, mr.Super.Verdict).Warningln(`Database error`)
		return
	}

	// store token in inmemory token map while db transaction is still
	// open
	if err = s.tokens.insert(token.Token, token.ValidFrom, token.ExpiresAt,
		token.Salt); err != nil {
		mr.ServerError(err)
		return
	}
	if err = tx.Commit(); err != nil {
		mr.ServerError(err)
		logEntry.WithField(`Code`, mr.Super.Verdict).Warningln(err)
		return
	}

	// encrypt generated token for client transmission
	plain := []byte{}
	data := []byte{}
	if plain, err = json.Marshal(token); err != nil {
		mr.ServerError(err)
		logEntry.WithField(`Code`, mr.Super.Verdict).Warningln(err)
		return
	}
	if err = kex.EncryptAndEncode(&plain, &data); err != nil {
		mr.ServerError(err)
		logEntry.WithField(`Code`, mr.Super.Verdict).Warningln(err)
		return
	}

	// prepare result for client transmission
	mr.Super = msg.Supervisor{
		Verdict: 200,
		Encrypted: struct {
			KexID string
			Data  []byte
		}{
			Data: data,
		},
	}
	logEntry.WithField(`Code`, mr.Super.Verdict).Infoln(`Successfully issued token`)
	mr.OK()
}

// tokenInvalidateAll invalidates all tokens
func (s *Supervisor) tokenInvalidateAll(q *msg.Request, mr *msg.Result) {
	// XXX TODO
}

// tokenInvalidate marks all tokens of a user as invalidate-on-use
func (s *Supervisor) tokenInvalidate(q *msg.Request, mr *msg.Result) {
	// XXX TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
