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
	"github.com/mjolnir42/soma/lib/auth"
	"github.com/satori/go.uuid"
)

// activateRoot handles requests to unlock the root account
func (s *Supervisor) activateRoot(q *msg.Request, mr *msg.Result) {

	var (
		kex                  *auth.Kex
		err                  error
		token                *auth.Token
		mcf                  scrypth64.Mcf
		tx                   *sql.Tx
		validFrom, expiresAt time.Time
		ok                   bool
	)

	// decrypt e2e encrypted request
	if token, kex, ok = s.decrypt(q, mr); !ok {
		return
	}

	// check that the encrypted username is actually root
	if token.UserName != msg.SubjectRoot {
		mr.BadRequest(fmt.Errorf(`Root activation attempted with`+
			" account: %s", token.UserName), q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return
	}

	// update auditlog entry
	mr.Super.Audit = mr.Super.Audit.
		WithField(`UserName`, token.UserName).
		WithField(`UserID`, uuid.Nil.String())

	// check if the request came in on a valid endpoint
	if s.rootRestricted && !q.Super.RestrictedEndpoint {
		mr.ServerError(fmt.Errorf(
			`Root bootstrap requested on unrestricted endpoint`),
			q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return
	}

	// check if root is not already active, which means there are
	// credentials stored for it
	if s.credentials.read(msg.SubjectRoot) != nil {
		mr.Forbidden(fmt.Errorf(`Root account is already active`))
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return
	}

	// verify provided RootToken is correct
	if !s.authenticateRootToken(token, mr) {
		return
	}
	// OK: validation success

	// calculate the scrypt KDF hash using scrypth64.DefaultParams()
	if mcf, err = scrypth64.Digest(token.Password, nil); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// generate a token for root. This checks the provided credentials
	// which always always succeeds since mcf was just computed from
	// token.Password, but causes a second scrypt computation delay
	token.SetIPAddressExtractedString(q.RemoteAddr)
	if err = token.Generate(mcf, s.key, s.seed); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// prepare data required for storing the user activation
	validFrom, _ = time.Parse(msg.RFC3339Milli, token.ValidFrom)
	expiresAt, _ = time.Parse(msg.RFC3339Milli, token.ExpiresAt)

	// open multi statement transaction
	if tx, err = s.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).
			Warningln(err)
		return
	}

	// Insert new credentials
	if s.txInsertCred(
		tx,
		uuid.Nil,
		msg.SubjectRoot,
		mcf.String(),
		validFrom.UTC(),
		msg.PosTimeInf.UTC(),
		mr,
	) {
		tx.Rollback()
		return
	}
	s.credentials.insert(
		msg.SubjectRoot,
		uuid.Nil,
		validFrom.UTC(),
		msg.PosTimeInf.UTC(),
		mcf,
	)

	// Insert issued token
	if s.txInsertToken(
		tx,
		token.Token,
		token.Salt,
		validFrom.UTC(),
		expiresAt.UTC(),
		mr,
	) {
		tx.Rollback()
		return
	}
	if err = s.tokens.insert(
		token.Token,
		token.ValidFrom,
		token.ExpiresAt,
		token.Salt,
	); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		tx.Rollback()
		return
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// encrypt e2e encrypted result and store it in mr
	if err = s.encrypt(kex, token, mr); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	mr.Super.Audit.Infoln(`Successfully activated root account`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
