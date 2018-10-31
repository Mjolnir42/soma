/*-
 * Copyright (c) 2017, Jörg Pernfuß
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
)

// tokenInvalidate revokes the currently used access token. Since
// the rest endpoint for this function performs basic auth, the
// token must exist
func (s *Supervisor) tokenInvalidate(q *msg.Request, mr *msg.Result) {
	var (
		userID string
		err    error
		res    sql.Result
		cnt    int64
	)

	// check the user exists and is active, this is for updating
	// the auditlog only
	if userID, err = s.checkUser(q.AuthUser, mr, true); err != nil {
		return
	}

	// update auditlog entry
	mr.Super.Audit = mr.Super.Audit.
		WithField(`UserName`, q.AuthUser).
		WithField(`UserID`, userID).
		WithField(`KexID`, `none`)

	// revocation time for the token
	expiredAt := time.Now().UTC()

	// update token in database
	if res, err = s.conn.Exec(
		stmt.ExpireToken,
		expiredAt,
		q.Super.AuthToken,
	); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return
	}

	// token row has unique constraint
	if cnt, err = res.RowsAffected(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return
	}
	switch cnt {
	case 1:
	default:
		// the token that was used to authenticate this request was
		// not found in the database, the authentication system is
		// corrupt. HARD crash.
		s.errLog.Errorln(`Supervisor corrupted, emergency crash. Check supervisor audit log`)
		mr.Super.Audit.Fatalf("Supervisor corruption detected! "+
			"Token %s used to authenticate this "+
			"Request was found in the database %d times!",
			q.Super.AuthToken,
			cnt,
		)
	}

	// remove the token from the in-memory map. the r/w master instance
	// has the authoritative copy of all tokens in memory and does
	// not load them from the database at runtime
	s.tokens.remove(q.Super.AuthToken)

	mr.Super.Verdict = 200
	mr.OK()
	mr.Super.Audit.
		WithField(`Verdict`, mr.Super.Verdict).
		WithField(`Code`, mr.Code).
		Infoln(`Successfully revoked token`)
}

// tokenInvalidateAccount marks all tokens of a user as
// invalidate-on-use
func (s *Supervisor) tokenInvalidateAccount(q *msg.Request, mr *msg.Result) {
	var (
		userID, victimID string
		err              error
		res              sql.Result
		cnt              int64
	)

	// check the user exists and is active, this is for updating
	// the auditlog only
	if userID, err = s.checkUser(q.AuthUser, mr, true); err != nil {
		return
	}

	// check the user to revoke exists and is active
	switch {
	case q.Super.RevokeForName != ``:
		victimID, err = s.checkUser(
			q.Super.RevokeForName,
			mr,
			true,
		)
	}
	if err != nil {
		return
	}

	// update auditlog entry
	mr.Super.Audit = mr.Super.Audit.
		WithField(`UserName`, q.AuthUser).
		WithField(`UserID`, userID).
		WithField(`KexID`, `none`)

	// insert revocation into database
	if res, err = s.conn.Exec(
		stmt.RevokeTokensForUser,
		victimID,
		time.Now().UTC(),
	); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return
	}

	// check how many rows were inserted
	if cnt, err = res.RowsAffected(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return
	}
	switch cnt {
	case 1:
	default:
		mr.ServerError(fmt.Errorf(
			"Revocation inserted %d rows, expected 1", cnt),
			q.Section,
		)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return
	}

	// add revocation to in-memory map
	s.tokens.expireAccount(q.Super.RevokeForName)

	mr.Super.Verdict = 200
	mr.OK()
	mr.Super.Audit.
		WithField(`Verdict`, mr.Super.Verdict).
		WithField(`Code`, mr.Code).
		Infof(
			"Successfully marked tokens for %s as revoked",
			q.Super.RevokeTokensFor,
		)
}

// tokenInvalidateGlobal invalidates all tokens
func (s *Supervisor) tokenInvalidateGlobal(q *msg.Request, mr *msg.Result) {
	var (
		userID string
		err    error
		res    sql.Result
		cnt    int64
	)

	// check the user exists and is active, this is for updating
	// the auditlog only
	if userID, err = s.checkUser(q.AuthUser, mr, true); err != nil {
		return
	}

	// update auditlog entry
	mr.Super.Audit = mr.Super.Audit.
		WithField(`UserName`, q.AuthUser).
		WithField(`UserID`, userID).
		WithField(`KexID`, `none`)

	// revocation time for the tokens
	revokeAt := time.Now().UTC()

	// lock the token map externally once for the entire job
	s.tokens.lock()
	defer s.tokens.unlock()

	for tokenID := range s.tokens.iterateStringUnlocked(revokeAt) {
		// update token in database
		if res, err = s.conn.Exec(
			stmt.ExpireToken,
			revokeAt,
			tokenID,
		); err != nil {
			mr.ServerError(err, q.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
			return
		}

		// token row has unique constraint
		if cnt, err = res.RowsAffected(); err != nil {
			mr.ServerError(err, q.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
			return
		}
		switch cnt {
		case 1:
		default:
			// the token to invalidate was not found in the database,
			// the authentication system is corrupt. HARD crash.
			s.errLog.Errorln(`Supervisor corrupted, emergency crash. Check supervisor audit log`)
			mr.Super.Audit.Fatalf("Supervisor corruption detected! "+
				"Token %s to be invalidated was found in the database %d times!",
				tokenID,
				cnt,
			)
		}

		// remove the token from the in-memory map. the r/w master instance
		// has the authoritative copy of all tokens in memory and does
		// not load them from the database at runtime
		s.tokens.removeUnlocked(tokenID)
	}

	mr.Super.Verdict = 200
	mr.OK()
	mr.Super.Audit.
		WithField(`Verdict`, mr.Super.Verdict).
		WithField(`Code`, mr.Code).
		Infoln(`Successfully revoked token`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
