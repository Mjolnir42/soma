/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"sync"

	"github.com/mjolnir42/soma/lib/auth"
)

// kexMap is the internal storage format for key exchanges
type kexMap struct {
	// kexid(uuid.string) -> auth.Kex
	KMap  map[string]auth.Kex
	gcMap map[string]bool
	mutex sync.RWMutex
}

// newKexMap returns a new kexMap
func newKexMap() *kexMap {
	m := kexMap{}
	m.KMap = make(map[string]auth.Kex)
	return &m
}

// Map manipulation

// read returns a copy of the requested kex
//
// the nonce information would normally mean returning
// a copy is problematic, but since these keys are only
// used for one client/server exchange, they are never
// put back
func (k *kexMap) read(kexRequest string) *auth.Kex {
	k.rlock()
	defer k.runlock()
	if kex, ok := k.KMap[kexRequest]; ok {
		return &kex
	}
	return nil
}

// insert adds a new key exchange to the kexMap
func (k *kexMap) insert(kex auth.Kex) {
	k.lock()
	defer k.unlock()

	k.KMap[kex.Request.String()] = kex
}

// remove deletes a key exchange from the kexMap
func (k *kexMap) remove(kexRequest string) {
	k.lock()
	defer k.unlock()

	delete(k.KMap, kexRequest)
	delete(k.gcMap, kexRequest)
}

// Garbage collection bulk functions with external locking

// removeUnlocked deletes a key exchange from the kexMap without acquiring the
// mutex lock. Locking must be done externally.
func (k *kexMap) removeUnlocked(kexRequest string) {
	delete(k.KMap, kexRequest)
	delete(k.gcMap, kexRequest)
}

// iterateUnlocked returns all current key exchanges in a channel
// without acquiring the mutex lock. Locking must be done externally.
func (k *kexMap) interateUnlocked() chan auth.Kex {
	ret := make(chan auth.Kex, len(k.KMap)+1)

	for id := range k.KMap {
		ret <- k.KMap[id]
	}

	close(ret)
	return ret
}

// markUnlocked sets the garbage collection mark on an expired key echange
// without acquiring the mutex lock. Locking must be done externally.
func (k *kexMap) markUnlocked(kexID string) {
	k.gcMap[kexID] = true
}

// sweepUnlocked deletes all key exchanges marked for garbage collection
// without acquiring the mutex lock. Locking must be done externally.
func (k *kexMap) sweepUnlocked() {
	for id := range k.gcMap {
		delete(k.KMap, id)
		delete(k.gcMap, id)
	}
}

// Locking

// lock acquires the writelock on kexMap k
func (k *kexMap) lock() {
	k.mutex.Lock()
}

// rlock acquires the readlock on kexMap k
func (k *kexMap) rlock() {
	k.mutex.RLock()
}

// unlock releases the writelock on kexMap k
func (k *kexMap) unlock() {
	k.mutex.Unlock()
}

// runlock releases the readlock on kexMap k
func (k *kexMap) runlock() {
	k.mutex.RUnlock()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
