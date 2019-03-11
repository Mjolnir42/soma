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
	uuid "github.com/satori/go.uuid"
)

// IMPORTANT!
//
// all errors returned from encrypted supervisor methods are
// returned to the client as 403/Forbidden. Provided error details
// are used only for serverside logging.

// activateAdmin handles requests to activate inactive admin accounts
func (s *Supervisor) activateAdmin(q *msg.Request, mr *msg.Result) {

	var (
		err                      error
		kex                      *auth.Kex
		validFrom, credExpiresAt time.Time
		token                    *auth.Token
		adminID                  string
		adminUUID                uuid.UUID
		ok                       bool
		mcf                      scrypth64.Mcf
		tx                       *sql.Tx
	)

	// decrypt e2e encrypted request
	if token, kex, ok = s.decrypt(q, mr); !ok {
		return
	}

	// update auditlog entry
	mr.Super.Audit = mr.Super.Audit.WithField(`UserName`, token.UserName)

	// root can not be activated via the user handler
	if token.UserName == msg.SubjectRoot {
		str := `Invalid user activation: root`

		mr.BadRequest(fmt.Errorf(str), q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(str)
		return
	}

	// check the admin exists and is not active
	if adminID, err = s.checkAdmin(token.UserName, mr, false); err != nil {
		return
	}

	// update auditlog entry
	mr.Super.Audit = mr.Super.Audit.WithField(`UserID`, adminID)
	adminUUID, _ = uuid.FromString(adminID)

	// no account ownership verification in open mode
	if !s.conf.OpenInstance {
		if !s.authenticatePassword(token, mr) {
			return
		}
	}
	// OK: validation success

	// calculate the scrypt KDF hash using scrypth64.DefaultParams()
	if mcf, err = scrypth64.Digest(token.Password, nil); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).
			Warningln(err)
		return
	}

	// TODO: refactor token generation
	// generate a token for the user. This checks the provided credentials
	// which always succeeds since mcf was just computed from token.Password,
	// but causes a second scrypt computation delay
	token.SetIPAddressExtractedString(q.RemoteAddr)
	if err = token.Generate(mcf, s.key, s.seed); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).
			Warningln(err)
		return
	}

	// prepare data required for storing the user activation
	validFrom, _ = time.Parse(msg.RFC3339Milli, token.ValidFrom)
	credExpiresAt = validFrom.Add(time.Duration(s.credExpiry) * time.Hour * 24).UTC()

	// open multi statement transaction
	if tx, err = s.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).
			Warningln(err)
		return
	}

	// Insert new credentials
	if !s.saveAdminCred(tx, token.UserName, msg.SubjectAdmin, adminUUID,
		mcf, validFrom.UTC(), credExpiresAt, mr) {
		tx.Rollback()
		return
	}

	// Insert issued token
	if !s.saveToken(tx, token, mr) {
		tx.Rollback()
		return
	}

	// activate user account
	if _, err = tx.Exec(
		stmt.ActivateAdminUser,
		adminUUID,
	); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).
			Warningln(err)
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
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	mr.Super.Audit.Infoln(`Successfully activated admin`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
