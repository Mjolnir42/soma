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
