/*-
 * Copyright (c) 2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"database/sql"
	"time"

	"github.com/mjolnir42/soma/internal/msg"
	uuid "github.com/satori/go.uuid"
)

// passwordRevoke handles revoking a user's credentials. It is called
// internally when users are removed and not an E2E encrypted part of
// the client API
func (s *Supervisor) passwordRevoke(q *msg.Request, mr *msg.Result) {
	var (
		err      error
		userID   string
		userUUID uuid.UUID
		tx       *sql.Tx
	)

	// check the requesting user exists and is active, this is only for
	// updating the auditlog
	if userID, err = s.checkUser(q.AuthUser, mr, true); err != nil {
		return
	}

	// check the user to revoke exists and is active
	switch {
	case q.Super.RevokeForName != ``:
		q.Super.RevokeForID, err = s.checkUser(
			q.Super.RevokeForName,
			mr,
			true,
		)
	case q.Super.RevokeForID != ``:
		q.Super.RevokeForName, err = s.checkUserByID(
			q.Super.RevokeForID,
			mr,
			true,
		)
	}
	if err != nil {
		return
	}

	// update auditlog
	mr.Super.Audit = mr.Super.Audit.
		WithField(`UserName`, q.AuthUser).
		WithField(`UserID`, userID).
		WithField(`RevokedUserName`, q.Super.RevokeForName).
		WithField(`RevokedUserID`, q.Super.RevokeForID)
	userUUID, _ = uuid.FromString(q.Super.RevokeForID)

	// credentials are revoked as of 1 second ago
	revocationTime := time.Now().UTC().Add(time.Second * -1)

	// Open transaction to revoke credentials
	if tx, err = s.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// Revoke existing credentials
	if s.txExpireCred(tx, revocationTime, userUUID, mr) {
		tx.Rollback()
		return
	}
	s.credentials.revoke(q.Super.RevokeForName)

	// commit transaction
	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	mr.Super.Audit.Infoln(`Successfully revoked password`)
	mr.OK()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
