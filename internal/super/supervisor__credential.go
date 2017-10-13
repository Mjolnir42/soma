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

	"github.com/mjolnir42/scrypth64"
	"github.com/satori/go.uuid"
)

// credential is the internal storage format for user credentials
type credential struct {
	id          uuid.UUID
	name        string
	validFrom   time.Time
	expiresAt   time.Time
	cryptMCF    scrypth64.Mcf
	resetActive bool
	isActive    bool
	gcMark      bool
}

// isExpired returns if credentials are expired
func (c *credential) isExpired() bool {
	return time.Now().UTC().After(c.expiresAt.UTC())
}

// credentialMap is a read/write locked map of credentials
type credentialMap struct {
	// username -> credential
	CMap  map[string]credential
	mutex sync.RWMutex
}

// newCredentialMap returns a new credentialMap
func newCredentialMap() *credentialMap {
	m := credentialMap{}
	m.CMap = make(map[string]credential)
	return &m
}

// Map manipulation

// read returns a copy of the requested user's credentials
func (c *credentialMap) read(user string) *credential {
	c.rlock()
	defer c.runlock()
	if cred, ok := c.CMap[user]; ok {
		return &cred
	}
	return nil
}

// insert adds a new active credential to the credentialMap
func (c *credentialMap) insert(user string, uid uuid.UUID, valid, expires time.Time, mcf scrypth64.Mcf) {
	c.lock()
	defer c.unlock()
	c.CMap[user] = credential{
		id:          uid,
		name:        user,
		validFrom:   valid,
		expiresAt:   expires,
		cryptMCF:    mcf,
		resetActive: false,
		isActive:    true,
	}
}

// restore loads a credential into the credentialMap, allowing different
// activation states
func (c *credentialMap) restore(user string, uid uuid.UUID, valid, expires time.Time, mcf scrypth64.Mcf, reset, active bool) {
	c.lock()
	defer c.unlock()
	c.CMap[user] = credential{
		id:          uid,
		validFrom:   valid,
		expiresAt:   expires,
		cryptMCF:    mcf,
		resetActive: reset,
		isActive:    active,
	}
}

// revoke expires the credentials for a user
func (c *credentialMap) revoke(user string) {
	c.lock()
	defer c.unlock()

	cred := c.CMap[user]
	cred.expiresAt = time.Now().UTC().Add(-5 * time.Second)
	c.CMap[user] = cred
}

// Garbage collection bulk functions with external locking

// removeUnlocked deletes a user's credentials from credentialMap without
// acquiring the mutex lock. Locking must be done externally.
func (c *credentialMap) removeUnlocked(user string) {
	delete(c.CMap, user)
}

// iterateUnlocked returns all current credentials in a channel without
// acquiring the mutex lock. Locking must be done externally.
func (c *credentialMap) iterateUnlocked() chan credential {
	ret := make(chan credential, len(c.CMap)+1)

	for user := range c.CMap {
		ret <- c.CMap[user]
	}

	return ret
}

// markUnlocked sets the garbage collection mark on an expired
// credential without acquiring the mutex lock. Locking must be done externally.
func (c *credentialMap) markUnlocked(user string) {
	cred := c.CMap[user]
	cred.gcMark = true
	c.CMap[user] = cred
}

// sweepUnlocked deletes all credentials marked for garbage collection
// without acquiring the mutex lock. Locking must be done externally.
func (c *credentialMap) sweepUnlocked() {
	for user := range c.CMap {
		if c.CMap[user].gcMark {
			c.removeUnlocked(user)
		}
	}
}

// Locking

// lock acquires the writelock on credentialMap c
func (c *credentialMap) lock() {
	c.mutex.Lock()
}

// rlock acquires the readlock on credentialMap c
func (c *credentialMap) rlock() {
	c.mutex.RLock()
}

// unlock releases the writelock on credentialMap c
func (c *credentialMap) unlock() {
	c.mutex.Unlock()
}

// runlock releases the readlock on credentialMap c
func (c *credentialMap) runlock() {
	c.mutex.RUnlock()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
