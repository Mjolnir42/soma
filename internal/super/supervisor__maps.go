/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"sync"
	"time"

	"github.com/mjolnir42/soma/lib/auth"

	"github.com/mjolnir42/scrypth64"
	"github.com/satori/go.uuid"
)

//
//
// supervisor internal storage format for credentials
type svCredential struct {
	id          uuid.UUID
	name        string
	validFrom   time.Time
	expiresAt   time.Time
	cryptMCF    scrypth64.Mcf
	resetActive bool
	isActive    bool
	gcMark      bool
}

// isExpired returns if a token is expired
func (c *svCredential) isExpired() bool {
	return time.Now().UTC().After(c.expiresAt.UTC())
}

type svCredMap struct {
	// username -> svCredential
	CMap  map[string]svCredential
	mutex sync.RWMutex
}

func (c *svCredMap) read(user string) *svCredential {
	c.rlock()
	defer c.runlock()
	if cred, ok := c.CMap[user]; ok {
		return &cred
	}
	return nil
}

func (c *svCredMap) insert(user string, uid uuid.UUID, valid, expires time.Time, mcf scrypth64.Mcf) {
	c.lock()
	defer c.unlock()
	c.CMap[user] = svCredential{
		id:          uid,
		name:        user,
		validFrom:   valid,
		expiresAt:   expires,
		cryptMCF:    mcf,
		resetActive: false,
		isActive:    true,
	}
}

func (c *svCredMap) restore(user string, uid uuid.UUID, valid, expires time.Time, mcf scrypth64.Mcf, reset, active bool) {
	c.lock()
	defer c.unlock()
	c.CMap[user] = svCredential{
		id:          uid,
		validFrom:   valid,
		expiresAt:   expires,
		cryptMCF:    mcf,
		resetActive: reset,
		isActive:    active,
	}
}

func (c *svCredMap) revoke(user string) {
	c.lock()
	defer c.unlock()
	delete(c.CMap, user)
}

// revokeUnlocked deletes a user's credentials without acquiring the
// mutex lock
func (c *svCredMap) revokeUnlocked(user string) {
	delete(c.CMap, user)
}

// iterateUnlocked returns all current credentials in a channel without
// acquiring the mutex lock
func (c *svCredMap) iterateUnlocked() chan svCredential {
	ret := make(chan svCredential, len(c.CMap)+1)

	for user := range c.CMap {
		ret <- c.CMap[user]
	}

	return ret
}

// markUnlocked sets the garbage collection mark on an expired
// credential without acquiring the mutex lock
func (c *svCredMap) markUnlocked(user string) {
	cred := c.CMap[user]
	cred.gcMark = true
	c.CMap[user] = cred
}

// sweepUnlocked deletes all credentials marked for garbage collection
// without acquiring the mutex lock
func (c *svCredMap) sweepUnlocked() {
	for user := range c.CMap {
		if c.CMap[user].gcMark {
			c.revokeUnlocked(user)
		}
	}
}

// set writelock
func (c *svCredMap) lock() {
	c.mutex.Lock()
}

// set readlock
func (c *svCredMap) rlock() {
	c.mutex.RLock()
}

// release writelock
func (c *svCredMap) unlock() {
	c.mutex.Unlock()
}

// release readlock
func (c *svCredMap) runlock() {
	c.mutex.RUnlock()
}

//
//
// read/write locked map of key exchanges
type svKexMap struct {
	// kexid(uuid.string) -> auth.Kex
	KMap  map[string]auth.Kex
	gcMap map[string]bool
	mutex sync.RWMutex
}

// the nonce information would normally mean returning
// a copy is problematic, but since these keys are only
// used for one client/server exchange, they are never
// put back
func (k *svKexMap) read(kexRequest string) *auth.Kex {
	k.rlock()
	defer k.runlock()
	if kex, ok := k.KMap[kexRequest]; ok {
		return &kex
	}
	return nil
}

func (k *svKexMap) insert(kex auth.Kex) {
	k.lock()
	defer k.unlock()

	k.KMap[kex.Request.String()] = kex
}

func (k *svKexMap) remove(kexRequest string) {
	k.lock()
	defer k.unlock()

	delete(k.KMap, kexRequest)
	delete(k.gcMap, kexRequest)
}

// removeUnlocked
func (k *svKexMap) removeUnlocked(kexRequest string) {
	delete(k.KMap, kexRequest)
	delete(k.gcMap, kexRequest)
}

// sweepUnlocked deletes all key exchanges marked for garbage collection
// without acquiring the mutex lock
func (k *svKexMap) sweepUnlocked() {
	for id := range k.gcMap {
		delete(k.KMap, id)
		delete(k.gcMap, id)
	}
}

// set writelock
func (k *svKexMap) lock() {
	k.mutex.Lock()
}

// set readlock
func (k *svKexMap) rlock() {
	k.mutex.RLock()
}

// release writelock
func (k *svKexMap) unlock() {
	k.mutex.Unlock()
}

// release readlock
func (k *svKexMap) runlock() {
	k.mutex.RUnlock()
}

// markUnlocked sets the garbage collection mark on an expired key echange
// without acquiring the mutex lock
func (k *svKexMap) markUnlocked(kexID string) {
	k.gcMap[kexID] = true
}

// iterateUnlocked returns all current key exchanges in a channel
// without acquiring the mutex lock
func (k *svKexMap) interateUnlocked() chan auth.Kex {
	ret := make(chan auth.Kex, len(k.KMap)+1)

	for id := range k.KMap {
		ret <- k.KMap[id]
	}

	return ret
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
