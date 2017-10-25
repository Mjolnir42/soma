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

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/auth"
)

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
		userID               string
	)

	// decrypt e2e encrypted request
	if token, kex, ok = s.decrypt(q, mr); !ok {
		return
	}

	// update auditlog entry
	mr.Super.Audit = mr.Super.Audit.WithField(`UserName`, token.UserName)

	// check if root is available
	if token.UserName == `root` && s.rootRestricted && !q.Super.RestrictedEndpoint {
		str := `Restricted-mode root token requested on ` +
			`unrestricted endpoint`
		mr.BadRequest(fmt.Errorf(str), q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return
	}

	// check the user exists and is active
	if userID, err = s.checkUser(token.UserName, mr, true); err != nil {
		return
	}
	// update auditlog entry
	mr.Super.Audit = mr.Super.Audit.WithField(`UserID`, userID)

	// fetch user credentials, checked to exist by checkUser()
	cred = s.credentials.read(token.UserName)

	// check if the cred has either expired or not become valid yet
	if time.Now().UTC().Before(cred.validFrom.UTC()) ||
		time.Now().UTC().After(cred.expiresAt.UTC()) {
		mr.Forbidden(fmt.Errorf(`Cred expired`), q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(`Token expired`)
		return
	}

	// generate token if the provided credentials are valid
	token.SetIPAddressExtractedString(q.RemoteAddr)
	if err = token.Generate(cred.cryptMCF, s.key, s.seed); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// persist generated token into database
	validFrom, _ = time.Parse(msg.RFC3339Milli, token.ValidFrom)
	expiresAt, _ = time.Parse(msg.RFC3339Milli, token.ExpiresAt)
	if tx, err = s.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
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
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// store token in inmemory token map while db transaction is still
	// open
	if err = s.tokens.insert(token.Token, token.ValidFrom, token.ExpiresAt,
		token.Salt); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// encrypt generated token for client transmission
	if err = s.encrypt(kex, token, mr); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	mr.Super.Audit.Infoln(`Successfully issued token`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
