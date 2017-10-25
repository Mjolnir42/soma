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

	"github.com/Sirupsen/logrus"
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

// activateUser handles requests to activate inactive user accounts
func (s *Supervisor) activateUser(q *msg.Request, mr *msg.Result, audit *logrus.Entry) {

	var (
		err                                 error
		kex                                 *auth.Kex
		validFrom, expiresAt, credExpiresAt time.Time
		token                               *auth.Token
		userID                              string
		userUUID                            uuid.UUID
		ok, active                          bool
		mcf                                 scrypth64.Mcf
		tx                                  *sql.Tx
	)

	if s.readonly {
		mr.ReadOnly()
		return
	}

	// decrypt e2e encrypted request
	if token, kex, ok = s.decrypt(q, mr, audit); !ok {
		return
	}

	// update auditlog entry
	audit = audit.WithField(`UserName`, token.UserName)

	// TODO: refactor user verification?
	// root can not be activated via the user handler
	if token.UserName == msg.SubjectRoot {
		str := `Invalid user activation: root`

		mr.BadRequest(fmt.Errorf(str), q.Section)
		audit.WithField(`Code`, mr.Code).
			Warningln(str)
		return
	}

	// verify the user to activate exists
	if err = s.stmtFindUserID.QueryRow(
		token.UserName,
	).Scan(
		&userID,
	); err == sql.ErrNoRows {
		str := fmt.Sprintf("Unknown user: %s", token.UserName)

		mr.NotFound(fmt.Errorf(str), q.Section)
		audit.WithField(`Code`, mr.Code).
			Warningln(str)
		return
	} else if err != nil {
		mr.ServerError(err)
		audit.WithField(`Code`, mr.Code).
			Warningln(err)
		return
	}

	// update auditlog entry
	audit = audit.WithField(`UserID`, userID)
	userUUID, _ = uuid.FromString(userID)

	// verify the user is not already active
	if err = s.stmtCheckUserActive.QueryRow(
		userID,
	).Scan(
		&active,
	); err == sql.ErrNoRows {
		str := fmt.Sprintf("Unknown user: %s", token.UserName)

		mr.NotFound(fmt.Errorf(str), q.Section)
		audit.WithField(`Code`, mr.Code).
			Warningln(str)
		return
	} else if err != nil {
		mr.ServerError(err)
		audit.WithField(`Code`, mr.Code).
			Warningln(err)
		return
	}
	if active {
		str := fmt.Sprintf("User %s (%s) is already active",
			token.UserName, userID)

		mr.BadRequest(fmt.Errorf(str), q.Section)
		audit.WithField(`Code`, mr.Code).
			Warningln(str)
		return
	}

	// TODO: refactor ownership verification
	// no account ownership verification in open mode
	if !s.conf.OpenInstance {
		switch s.activation {
		case `ldap`:
			if ok, err = validateLdapCredentials(
				token.UserName,
				token.Token,
			); err != nil {
				mr.ServerError(err)
				audit.WithField(`Code`, mr.Code).
					Warningln(err)
				return
			} else if !ok {
				str := `Invalid LDAP credentials`

				mr.Unauthorized(fmt.Errorf(str), q.Section)
				audit.WithField(`Code`, mr.Code).
					Warningln(str)
				return
			}
			// fail activation if local password is the same as the
			// upstream password. This error _IS_ sent to the user!
			if token.Token == token.Password {
				str := fmt.Sprintf(
					"User %s denied: matching local/upstream passwords",
					token.UserName)

				mr.Conflict(fmt.Errorf(str), q.Section)
				audit.WithField(`Code`, mr.Code).
					Warningln(str)
				return
			}
		case `token`: // TODO
			str := `Mailtoken activation is not implemented`

			mr.NotImplemented(fmt.Errorf(str), q.Section)
			audit.WithField(`Code`, mr.Code).
				Errorln(str)
			return
		default:
			str := fmt.Sprintf("Unknown activation: %s",
				s.conf.Auth.Activation)

			mr.ServerError(fmt.Errorf(str), q.Section)
			audit.WithField(`Code`, mr.Code).
				Errorln(str)
			return
		}
	}
	// OK: validation success

	// calculate the scrypt KDF hash using scrypth64.DefaultParams()
	if mcf, err = scrypth64.Digest(token.Password, nil); err != nil {
		mr.ServerError(err, q.Section)
		audit.WithField(`Code`, mr.Code).
			Warningln(err)
		return
	}

	// TODO: refactor token generation
	// generate a token for the user. This checks the provided credentials
	// which always always succeeds since mcf was just computed from token.Password,
	// but causes a second scrypt computation delay
	token.SetIPAddressExtractedString(q.RemoteAddr)
	if err = token.Generate(mcf, s.key, s.seed); err != nil {
		mr.ServerError(err, q.Section)
		audit.WithField(`Code`, mr.Code).
			Warningln(err)
		return
	}

	// prepare data required for storing the user activation
	validFrom, _ = time.Parse(msg.RFC3339Milli, token.ValidFrom)
	expiresAt, _ = time.Parse(msg.RFC3339Milli, token.ExpiresAt)
	credExpiresAt = validFrom.Add(time.Duration(s.credExpiry) * time.Hour * 24).UTC()

	// open multi statement transaction
	if tx, err = s.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		audit.WithField(`Code`, mr.Code).
			Warningln(err)
		return
	}
	defer tx.Rollback()

	// persist accepted credentials for the user
	if _, err = tx.Exec(
		stmt.SetUserCredential,
		userUUID,
		mcf.String(),
		validFrom.UTC(),
		credExpiresAt.UTC(),
	); err != nil {
		mr.ServerError(err, q.Section)
		audit.WithField(`Code`, mr.Code).
			Warningln(err)
		return
	}

	// activate user account
	if _, err = tx.Exec(
		stmt.ActivateUser,
		userUUID,
	); err != nil {
		mr.ServerError(err, q.Section)
		audit.WithField(`Code`, mr.Code).
			Warningln(err)
		return
	}

	// persist generated token
	if _, err = tx.Exec(
		stmt.InsertToken,
		token.Token,
		token.Salt,
		validFrom.UTC(),
		expiresAt.UTC(),
	); err != nil {
		mr.ServerError(err, q.Section)
		audit.WithField(`Code`, mr.Code).
			Warningln(err)
		return
	}

	// update supervisor private in-memory credentials store
	s.credentials.insert(token.UserName, userUUID, validFrom.UTC(),
		credExpiresAt.UTC(), mcf)

	// update supervisor private in-memory token store
	if err = s.tokens.insert(token.Token, token.ValidFrom, token.ExpiresAt,
		token.Salt); err != nil {
		mr.ServerError(err, q.Section)
		audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// encrypt e2e encrypted result and store it in mr
	if err = s.encrypt(kex, token, mr, audit); err != nil {
		mr.ServerError(err, mr.Section)
		audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	audit.Infoln(`Successfully activated user`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
