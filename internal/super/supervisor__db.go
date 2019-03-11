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

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	uuid "github.com/satori/go.uuid"
)

func (s *Supervisor) txExpireCred(tx *sql.Tx, at time.Time, user uuid.UUID, mr *msg.Result) bool {
	if _, err := tx.Exec(
		stmt.InvalidateUserCredential,
		at,
		user,
	); err != nil {
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return true
	}
	return false
}

func (s *Supervisor) txInsertCred(tx *sql.Tx, user uuid.UUID, subject, mcf string, validFrom, expireAt time.Time, mr *msg.Result) bool {
	var err error

	switch subject {
	case msg.SubjectAdmin:
		_, err = tx.Exec(
			stmt.SetAdminCredential,
			user,
			mcf,
			validFrom,
			expireAt,
		)
	case msg.SubjectUser:
		_, err = tx.Exec(
			stmt.SetUserCredential,
			user,
			mcf,
			validFrom,
			expireAt,
		)
	case msg.SubjectRoot:
		_, err = tx.Exec(
			stmt.SetRootCredentials,
			user,
			mcf,
			validFrom,
		)
	}

	if err != nil {
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return true
	}
	return false
}

func (s *Supervisor) txInsertToken(tx *sql.Tx, token, salt string, validFrom, expireAt time.Time, mr *msg.Result) bool {
	if _, err := tx.Exec(
		stmt.InsertToken,
		token,
		salt,
		validFrom,
		expireAt,
	); err != nil {
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
