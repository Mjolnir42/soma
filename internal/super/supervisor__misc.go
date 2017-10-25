/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
)

func (s *Supervisor) fetchTokenFromDB(token string) bool {
	var (
		err                       error
		salt, strValid, strExpire string
		validF, validU            time.Time
	)

	err = s.stmtTokenSelect.QueryRow(token).Scan(&salt, &validF, &validU)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		// XXX log error
		return false
	}

	strValid = validF.UTC().Format(msg.RFC3339Milli)
	strExpire = validU.UTC().Format(msg.RFC3339Milli)

	if err = s.tokens.insert(token, strValid, strExpire, salt); err == nil {
		return true
	}
	return false
}

func (s *Supervisor) fetchRootToken() (string, error) {
	var (
		err   error
		token string
	)

	err = s.conn.QueryRow(stmt.SelectRootToken).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

// the nonces used for encryption are implemented as
// a counter on top of the agreed upon IV. The first
// nonce used is IV+1.
// Check that the IV is not 0, this is likely to indicate
// a bad client. An IV of -1 would be worse, resulting in
// an initial nonce of 0 which can always lead to crypto
// swamps. Why are safe from that, since the Nonce calculation
// always takes the Abs value of the IV, stripping the sign.
func (s *Supervisor) checkIV(iv string) error {
	var (
		err       error
		bIV       []byte
		iIV, zero *big.Int
	)
	zero = big.NewInt(0)

	if bIV, err = hex.DecodeString(iv); err != nil {
		return err
	}

	iIV = big.NewInt(0)
	iIV.SetBytes(bIV)
	iIV.Abs(iIV)
	if iIV.Cmp(zero) == 0 {
		return fmt.Errorf(`Invalid Initialization vector`)
	}
	return nil
}

// checkUser verifies that a user exists and is either active or
// inactive depending on the target state
func (s *Supervisor) checkUser(name string, mr *msg.Result, audit *logrus.Entry, target bool) (string, error) {
	var (
		userID string
		active bool
		err    error
	)

	// verify the user to activate exists
	if err = s.stmtFindUserID.QueryRow(
		name,
	).Scan(
		&userID,
	); err == sql.ErrNoRows {
		str := fmt.Sprintf("Unknown user: %s", name)
		mr.NotFound(fmt.Errorf(str), mr.Section)
		audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	} else if err != nil {
		mr.ServerError(err, mr.Section)
		audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	}

	// query user active status
	if err = s.stmtCheckUserActive.QueryRow(
		userID,
	).Scan(
		&active,
	); err == sql.ErrNoRows {
		str := fmt.Sprintf("Unknown user: %s", name)
		mr.NotFound(fmt.Errorf(str), mr.Section)
		audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	} else if err != nil {
		mr.ServerError(err, mr.Section)
		audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	}

	// verify the user is not already in target state
	if active != target {
		str := fmt.Sprintf("User %s (%s) is active: %t",
			name, userID, target)

		mr.BadRequest(fmt.Errorf(str), mr.Section)
		audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	}

	return userID, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
