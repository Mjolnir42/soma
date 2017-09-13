/*-
 * Copyright (c) 2016, Jörg Pernfuß
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

	"github.com/mjolnir42/scrypth64"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/auth"
	uuid "github.com/satori/go.uuid"
)

func (s *Supervisor) activateUser(q *msg.Request) {
	result := msg.FromRequest(q)
	result.Super = &msg.Supervisor{}

	var (
		timer                               *time.Timer
		plain                               []byte
		err                                 error
		kex                                 *auth.Kex
		validFrom, expiresAt, credExpiresAt time.Time
		token                               auth.Token
		userID                              string
		userUUID                            uuid.UUID
		ok, active                          bool
		mcf                                 scrypth64.Mcf
		tx                                  *sql.Tx
	)
	data := q.Super.Encrypted.Data

	if s.readonly {
		result.Conflict(fmt.Errorf(`Readonly instance`))
		goto conflict
	}

	// start response timer
	timer = time.NewTimer(1 * time.Second)
	defer timer.Stop()

	// -> get kex
	if kex = s.kex.read(q.Super.Encrypted.KexID); kex == nil {
		//    --> reply 404 if not found
		result.NotFound(fmt.Errorf(`Key exchange not found`))
		goto dispatch
	}
	// -> check kex.SameSource
	if !kex.IsSameSourceExtractedString(q.RemoteAddr) {
		//    --> reply 404 if !SameSource
		result.NotFound(fmt.Errorf(`Key exchange not found`))
		goto dispatch
	}
	// -> delete kex from s.kex (kex is now used)
	s.kex.remove(q.Super.Encrypted.KexID)
	// -> rdata = kex.DecodeAndDecrypt(data)
	if err = kex.DecodeAndDecrypt(&data, &plain); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> json.Unmarshal(rdata, &token)
	if err = json.Unmarshal(plain, &token); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// request has been decrypted, log it
	s.reqLog.Printf(msg.LogStrSRq, q.Section, q.Action, token.UserName, q.RemoteAddr)

	// -> check token.UserName != `root`
	if token.UserName == `root` {
		//    --> reply 401
		result.Unauthorized(fmt.Errorf(`Cannot activate root`))
		goto dispatch
	}

	// check we have the user
	if err = s.stmtFindUserID.QueryRow(token.UserName).Scan(&userID); err == sql.ErrNoRows {
		result.Unauthorized(fmt.Errorf("Unknown user: %s", token.UserName))
		goto dispatch
	} else if err != nil {
		result.ServerError(err)
	}
	userUUID, _ = uuid.FromString(userID)

	// check the user is not already active
	if err = s.stmtCheckUserActive.QueryRow(userID).Scan(&active); err == sql.ErrNoRows {
		result.Unauthorized(fmt.Errorf("Unknown user: %s", token.UserName))
		goto dispatch
	}
	if active {
		result.Conflict(fmt.Errorf("User %s (%s) is already active", token.UserName, userID))
		goto dispatch
	}

	// no account ownership verification in open mode
	if !s.conf.OpenInstance {
		switch s.activation {
		case `ldap`:
			if ok, err = validateLdapCredentials(token.UserName, token.Token); err != nil {
				result.ServerError(err)
				goto dispatch
			} else if !ok {
				result.Unauthorized(fmt.Errorf(`Invalid LDAP credentials`))
				goto dispatch
			}
			// fail activation if local password is the same as the
			// upstream password
			if token.Token == token.Password {
				result.Unauthorized(fmt.Errorf("User %s denied: matching local/upstream passwords", token.UserName))
				goto dispatch
			}
		case `token`: // TODO
			result.ServerError(fmt.Errorf(`Not implemented`))
			goto dispatch
		default:
			result.ServerError(fmt.Errorf("Unknown activation: %s",
				s.conf.Auth.Activation))
			goto dispatch
		}
	}
	// OK: validation success

	// -> scrypth64.Digest(Password, nil)
	if mcf, err = scrypth64.Digest(token.Password, nil); err != nil {
		result.Unauthorized(nil)
		goto dispatch
	}
	// -> generate token
	token.SetIPAddressExtractedString(q.RemoteAddr)
	if err = token.Generate(mcf, s.key, s.seed); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	validFrom, _ = time.Parse(msg.RFC3339Milli, token.ValidFrom)
	expiresAt, _ = time.Parse(msg.RFC3339Milli, token.ExpiresAt)
	credExpiresAt = validFrom.Add(time.Duration(s.credExpiry) * time.Hour * 24).UTC()

	// -> open transaction
	if tx, err = s.conn.Begin(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	defer tx.Rollback()
	// -> DB Insert: password data
	if _, err = tx.Exec(
		stmt.SetUserCredential,
		userUUID,
		mcf.String(),
		validFrom.UTC(),
		credExpiresAt.UTC(),
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> DB Update: activate user
	if _, err = tx.Exec(
		stmt.ActivateUser,
		userUUID,
	); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> DB Insert: token data
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
	// -> s.credentials Update
	s.credentials.insert(token.UserName, userUUID, validFrom.UTC(),
		credExpiresAt.UTC(), mcf)
	// -> s.tokens Update
	if err = s.tokens.insert(token.Token, token.ValidFrom, token.ExpiresAt,
		token.Salt); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// commit transaction
	if err = tx.Commit(); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> sdata = kex.EncryptAndEncode(&token)
	plain = []byte{}
	data = []byte{}
	if plain, err = json.Marshal(token); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	if err = kex.EncryptAndEncode(&plain, &data); err != nil {
		result.ServerError(err)
		goto dispatch
	}
	// -> send sdata reply
	result.Super.Verdict = 200
	result.Super.Encrypted.Data = data
	result.OK()

dispatch:
	<-timer.C

conflict:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
