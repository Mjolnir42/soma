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

	"github.com/mjolnir42/scrypth64"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/auth"
	uuid "github.com/satori/go.uuid"
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
		s.errLog.WithField(`Function`, `fetchTokenFromDB`).Errorln(err)
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
// inactive depending on the target state. It returns the userID of the
// user.
func (s *Supervisor) checkUser(name string, mr *msg.Result, target bool) (string, error) {
	var (
		userID string
		active bool
		err    error
		cred   *credential
	)

	// verify the user to activate exists
	if err = s.stmtFindUserID.QueryRow(
		name,
	).Scan(
		&userID,
	); err == sql.ErrNoRows {
		str := fmt.Sprintf("Unknown user: %s", name)
		mr.NotFound(fmt.Errorf(str), mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	} else if err != nil {
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
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
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	} else if err != nil {
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	}

	// verify the user is not already in target state
	if active != target {
		str := fmt.Sprintf("User %s (%s) is active: %t",
			name, userID, target)
		mr.BadRequest(fmt.Errorf(str), mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	}

	// if the user should already be active, it should also
	// be in the credentials cache and active there as well
	if target {
		if cred = s.credentials.read(name); cred == nil {
			str := fmt.Sprintf("Active user not found in credentials"+
				" cache: %s (%s)", name, userID)
			mr.ServerError(fmt.Errorf(str), mr.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Errorln(mr.Error)
			return ``, mr.Error
		}

		if !cred.isActive {
			str := fmt.Sprintf("Active user inactive in credentials"+
				" cache: %s (%s)", name, userID)
			mr.ServerError(fmt.Errorf(str), mr.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Errorln(mr.Error)
			return ``, mr.Error
		}
	}

	return userID, nil
}

// checkUserByID verifies that a user exists and is either active or
// inactive depending on the target state. It returns the name of the
// user.
func (s *Supervisor) checkUserByID(userID string, mr *msg.Result, target bool) (string, error) {
	var (
		name   string
		active bool
		err    error
		cred   *credential
	)

	// verify the user to activate exists
	if err = s.stmtFindUserName.QueryRow(
		userID,
	).Scan(
		&name,
	); err == sql.ErrNoRows {
		str := fmt.Sprintf("Unknown user ID: %s", userID)
		mr.NotFound(fmt.Errorf(str), mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	} else if err != nil {
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
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
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	} else if err != nil {
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	}

	// verify the user is not already in target state
	if active != target {
		str := fmt.Sprintf("User %s (%s) is active: %t",
			name, userID, target)
		mr.BadRequest(fmt.Errorf(str), mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	}

	// if the user should already be active, it should also
	// be in the credentials cache and active there as well
	if target {
		if cred = s.credentials.read(name); cred == nil {
			str := fmt.Sprintf("Active user not found in credentials"+
				" cache: %s (%s)", name, userID)
			mr.ServerError(fmt.Errorf(str), mr.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Errorln(mr.Error)
			return ``, mr.Error
		}

		if !cred.isActive {
			str := fmt.Sprintf("Active user inactive in credentials"+
				" cache: %s (%s)", name, userID)
			mr.ServerError(fmt.Errorf(str), mr.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Errorln(mr.Error)
			return ``, mr.Error
		}
	}

	return name, nil
}

// saveToken persists a newly generated token
func (s *Supervisor) saveToken(tx *sql.Tx, token *auth.Token, mr *msg.Result) bool {
	validFrom, _ := time.Parse(msg.RFC3339Milli, token.ValidFrom)
	expiresAt, _ := time.Parse(msg.RFC3339Milli, token.ExpiresAt)

	// Insert issued token
	if s.txInsertToken(
		tx,
		token.Token,
		token.Salt,
		validFrom.UTC(),
		expiresAt.UTC(),
		mr,
	) {
		return false
	}

	if err := s.tokens.insert(
		token.Token,
		token.ValidFrom,
		token.ExpiresAt,
		token.Salt,
	); err != nil {
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return false
	}

	return true
}

// saveCred persists newly generated credentials
func (s *Supervisor) saveCred(tx *sql.Tx, user, subject string, userUUID uuid.UUID, mcf scrypth64.Mcf, validFrom, expiresAt time.Time, mr *msg.Result) bool {
	// Insert new credentials
	if s.txInsertCred(
		tx,
		userUUID,
		subject,
		mcf.String(),
		validFrom,
		expiresAt,
		mr,
	) {
		return false
	}
	s.credentials.insert(
		user,
		userUUID,
		validFrom,
		expiresAt,
		mcf,
	)
	return true
}

// checkAdmin verifies that a admin exists and is either active or
// inactive depending on the target state. It returns the adminID of the
// admin.
func (s *Supervisor) checkAdmin(name string, mr *msg.Result, target bool) (string, error) {
	var (
		adminID string
		active  bool
		err     error
		cred    *credential
	)

	// verify the user to activate exists
	if err = s.stmtFindAdminID.QueryRow(
		name,
	).Scan(
		&adminID,
	); err == sql.ErrNoRows {
		str := fmt.Sprintf("Unknown user: %s", name)
		mr.NotFound(fmt.Errorf(str), mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	} else if err != nil {
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	}

	// query user active status
	if err = s.stmtCheckAdminActive.QueryRow(
		adminID,
	).Scan(
		&active,
	); err == sql.ErrNoRows {
		str := fmt.Sprintf("Unknown user: %s", name)
		mr.NotFound(fmt.Errorf(str), mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	} else if err != nil {
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	}

	// verify the user is not already in target state
	if active != target {
		str := fmt.Sprintf("User %s (%s) is active: %t",
			name, adminID, target)
		mr.BadRequest(fmt.Errorf(str), mr.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return ``, mr.Error
	}

	// if the user should already be active, it should also
	// be in the credentials cache and active there as well
	if target {
		if cred = s.credentials.read(name); cred == nil {
			str := fmt.Sprintf("Active user not found in credentials"+
				" cache: %s (%s)", name, adminID)
			mr.ServerError(fmt.Errorf(str), mr.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Errorln(mr.Error)
			return ``, mr.Error
		}

		if !cred.isActive {
			str := fmt.Sprintf("Active user inactive in credentials"+
				" cache: %s (%s)", name, adminID)
			mr.ServerError(fmt.Errorf(str), mr.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Errorln(mr.Error)
			return ``, mr.Error
		}
	}

	return adminID, nil
}

// saveAdminCred persists newly generated credentials
func (s *Supervisor) saveAdminCred(tx *sql.Tx, user, subject string, userUUID uuid.UUID, mcf scrypth64.Mcf, validFrom, expiresAt time.Time, mr *msg.Result) bool {
	// Insert new credentials
	if s.txInsertCred(
		tx,
		userUUID,
		subject,
		mcf.String(),
		validFrom,
		expiresAt,
		mr,
	) {
		return false
	}
	s.credentials.insert(
		user,
		userUUID,
		validFrom,
		expiresAt,
		mcf,
	)
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
