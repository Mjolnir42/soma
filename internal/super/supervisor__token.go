/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"encoding/hex"
	"sync"
	"time"

	"github.com/mjolnir42/soma/internal/msg"
)

// token is the internal storage format for authentication tokens
type token struct {
	validFrom    time.Time
	expiresAt    time.Time
	binToken     []byte
	binExpiresAt []byte
	salt         []byte
	gcMark       bool
}

// isExpired returns if a token is expired
func (t *token) isExpired() bool {
	return time.Now().UTC().After(t.expiresAt.UTC())
}

// tokenMap is a read/write locked map of tokens
type tokenMap struct {
	// token(hex.string) -> token
	TMap map[string]token
	// username -> revocation time
	Expire map[string]time.Time
	mutex  sync.RWMutex
}

// newTokenMap returns a new tokenMap
func newTokenMap() *tokenMap {
	m := tokenMap{}
	m.TMap = make(map[string]token)
	m.Expire = make(map[string]time.Time)
	return &m
}

// Map manipulation

// read returns a copy of the requested token
func (t *tokenMap) read(token string) *token {
	t.rlock()
	defer t.runlock()
	if tok, ok := t.TMap[token]; ok {
		return &tok
	}
	return nil
}

// insert adds a new token to the tokenMap
func (t *tokenMap) insert(sToken, valid, expires, salt string) error {
	var (
		err                     error
		valTime, expTime        time.Time
		bExpTime, bSalt, bToken []byte
	)
	// convert input data into the different formats required to
	// perform later actions without conversions
	if valTime, err = time.Parse(msg.RFC3339Milli, valid); err != nil {
		return err
	}
	if expTime, err = time.Parse(msg.RFC3339Milli, expires); err != nil {
		return err
	}
	if bExpTime, err = expTime.MarshalBinary(); err != nil {
		return err
	}
	if bToken, err = hex.DecodeString(sToken); err != nil {
		return err
	}
	if bSalt, err = hex.DecodeString(salt); err != nil {
		return err
	}
	// whiteout unstable subsecond timestamp part with "random" value
	copy(bExpTime[9:], []byte{0xde, 0xad, 0xca, 0xfe})
	// acquire write lock
	t.lock()
	defer t.unlock()

	// insert token
	t.TMap[sToken] = token{
		validFrom:    valTime,
		expiresAt:    expTime,
		binToken:     bToken,
		binExpiresAt: bExpTime,
		salt:         bSalt,
	}
	return nil
}

// remove deletes a token from the token map
func (t *tokenMap) remove(token string) {
	// acquire write lock
	t.lock()
	defer t.unlock()

	delete(t.TMap, token)
}

// expireAccount marks the account as having all tokens issued
// until now expired
func (t *tokenMap) expireAccount(user string) {
	// acquire write lock
	t.lock()
	defer t.unlock()

	t.Expire[user] = time.Now().UTC()
}

// isExpired returns if and when the tokens for this account have
// been expired
func (t *tokenMap) isExpired(user string) (time.Time, bool) {
	t.rlock()
	defer t.runlock()
	if t, ok := t.Expire[user]; ok {
		return t, ok
	}
	return time.Time{}, false
}

// Garbage collection bulk functions with external locking

// removeUnlock deletes a token from the token map without acquiring the
// mutex lock. Locking must be done externally.
func (t *tokenMap) removeUnlocked(token string) {
	delete(t.TMap, token)
}

// iterateUnlocked returns all current tokens in a channel without
// acquiring the mutex lock. Locking must be done externally.
func (t *tokenMap) iterateUnlocked() chan token {
	ret := make(chan token, len(t.TMap)+1)
	defer close(ret)

	for id := range t.TMap {
		ret <- t.TMap[id]
	}

	return ret
}

// iterateStringUnlocked returns all current tokens that expire after
// revokeAt as strings in a channel without acquiring the mutex lock.
// Locking must be done externally.
func (t *tokenMap) iterateStringUnlocked(revokeAt time.Time) chan string {
	ret := make(chan string, len(t.TMap)+1)
	defer close(ret)

	for id := range t.TMap {
		// no need to return tokens that have already expired but have not yet
		// been garbage collected
		if revokeAt.Before(t.TMap[id].expiresAt) {
			ret <- id
		}
	}

	return ret
}

// markUnlocked sets the garbage collection mark on an expired tokens
// without acquiring the mutex lock. Locking must be done externally.
func (t *tokenMap) markUnlocked(token string) {
	tok := t.TMap[token]
	tok.gcMark = true
	t.TMap[token] = tok
}

// sweepUnlocked deletes all tokens marked for garbage collection
// without acquiring the mutex lock. Locking must be done externally.
func (t *tokenMap) sweepUnlocked() {
	for id := range t.TMap {
		if t.TMap[id].gcMark {
			t.removeUnlocked(id)
		}
	}
}

// cleanExpireUnlocked deletes all account revocation markers that are
// older than the expiry time the tokens, since all tokens affected by
// the revocation are no longer valid anyway. It does not acquire the
// mutex lock. Locking must be done externally.
func (t *tokenMap) cleanExpireUnlocked() {
	now := time.Now().UTC()
	offset := time.Duration(singleton.conf.Auth.TokenExpirySeconds) *
		time.Second

	for user, expiredAt := range t.Expire {
		if now.After(expiredAt.UTC().Add(offset)) {
			delete(t.Expire, user)
		}
	}
}

// Locking

// lock acquires the writelock on tokenMap t
func (t *tokenMap) lock() {
	t.mutex.Lock()
}

// rlock acquires the readlock on tokenMap t
func (t *tokenMap) rlock() {
	t.mutex.RLock()
}

// unlock releases the writelock on tokenMap t
func (t *tokenMap) unlock() {
	t.mutex.Unlock()
}

// runlock releases the readlock on tokenMap t
func (t *tokenMap) runlock() {
	t.mutex.RUnlock()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
