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

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/auth"
)

// TODO: timer, delete all expired key exchanges
func (s *Supervisor) pruneKex() {
	s.kex.lock()
	defer s.kex.unlock()
	for kexID, kex := range s.kex.KMap {
		if kex.IsExpired() {
			delete(s.kex.KMap, kexID)
		}
	}
}

func (s *Supervisor) newTokenMap() *svTokenMap {
	m := svTokenMap{}
	m.TMap = make(map[string]svToken)
	return &m
}

func (s *Supervisor) newCredentialMap() *svCredMap {
	m := svCredMap{}
	m.CMap = make(map[string]svCredential)
	return &m
}

func (s *Supervisor) newKexMap() *svKexMap {
	m := svKexMap{}
	m.KMap = make(map[string]auth.Kex)
	return &m
}

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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
