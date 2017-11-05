/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"database/sql"
	"time"

	"github.com/mjolnir42/scrypth64"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
	uuid "github.com/satori/go.uuid"
)

// IMPORTANT!
//
// all errors returned from encrypted supervisor methods are
// returned to the client as 403/Forbidden. Provided error details
// are used only for serverside logging.

// passwordModify handles requests to perform modifications to an
// account's password
func (s *Supervisor) passwordModify(q *msg.Request, mr *msg.Result) {

	var (
		err                                   error
		kex                                   *auth.Kex
		token                                 *auth.Token
		tx                                    *sql.Tx
		validFrom, expiresAt                  time.Time
		newCredExpiresAt, oldCredDeactivateAt time.Time
		userID                                string
		userUUID                              uuid.UUID
		mcf                                   scrypth64.Mcf
		ok                                    bool
	)

	// decrypt e2e encrypted request
	// token.UserName is the username
	// token.Password is the _NEW_ password that should be set
	if token, kex, ok = s.decrypt(q, mr); !ok {
		return
	}

	// update auditlog entry
	mr.Super.Audit = mr.Super.Audit.WithField(`UserName`, token.UserName)

	// check the user exists and is active
	if userID, err = s.checkUser(token.UserName, mr, true); err != nil {
		return
	}

	// update auditlog entry
	mr.Super.Audit = mr.Super.Audit.WithField(`UserID`, userID)
	userUUID, _ = uuid.FromString(userID)

	switch q.Super.Task {
	case msg.TaskReset:
		if !s.passwordReset(token, mr) {
			return
		}
	case msg.TaskChange:
		if !s.passwordChange(token, mr) {
			return
		}
	}

	// calculate scrypt KDF for new password
	if mcf, err = scrypth64.Digest(token.Password, nil); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// generate a new basic auth token
	token.SetIPAddressExtractedString(q.RemoteAddr)
	if err = token.Generate(mcf, s.key, s.seed); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	validFrom, _ = time.Parse(msg.RFC3339Milli, token.ValidFrom)
	expiresAt, _ = time.Parse(msg.RFC3339Milli, token.ExpiresAt)

	// old, changed credentials expire 1 second before the new
	// token was issued
	oldCredDeactivateAt = validFrom.Add(time.Second * -1).UTC()
	// new credentials expire based on server setting
	newCredExpiresAt = validFrom.Add(
		time.Duration(s.credExpiry) * time.Hour * 24).UTC()

	// Open transaction to update credentials
	if tx, err = s.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// Invalidate existing credentials
	if s.txExpireCred(tx, oldCredDeactivateAt, userUUID, mr) {
		tx.Rollback()
		return
	}
	s.credentials.revoke(token.UserName)

	// Insert new credentials
	if s.txInsertCred(tx, userUUID, msg.SubjectUser, mcf.String(), validFrom.UTC(), newCredExpiresAt, mr) {
		tx.Rollback()
		return
	}
	s.credentials.insert(token.UserName,
		userUUID,
		validFrom.UTC(),
		newCredExpiresAt,
		mcf,
	)

	// Insert issued token
	if s.txInsertToken(tx, token.Token, token.Salt, validFrom.UTC(), expiresAt.UTC(), mr) {
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
	mr.Super.Audit.Infoln(`Successfully updated password`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
