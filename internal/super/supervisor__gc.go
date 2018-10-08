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
)

// gc runs garbage collection on various supervisor data structures
func (s *Supervisor) gc() {
	s.appLog.Debug(`Supervisor.GC locking kex map`)
	// lock key-exchange map
	s.kex.lock()
	defer s.kex.unlock()

	// lock token map
	s.appLog.Debug(`Supervisor.GC locking token map`)
	s.tokens.lock()
	defer s.tokens.unlock()

	// lock credentials map
	s.appLog.Debug(`Supervisor.GC locking credential map`)
	s.credentials.lock()
	defer s.credentials.unlock()

	// sweep records marked for garbage collection during the last gc
	// run
	s.appLog.Debug(`Supervisor.GC sweeping records marked for deletion`)
	s.gcSweep()

	// mark records for garbage collection during the next run
	s.appLog.Debug(`Supervisor.GC marking records for cleanup next cycle`)
	s.gcMarkForNext()

	// remove old token expire entries
	s.appLog.Debug(`Supervisor.GC cleanup expired tokens`)
	s.tokens.cleanExpireUnlocked()
	s.appLog.Debug(`Supervisor.GC FINISHED`)
}

// gcMarkForNext marks data structures for deletion during the next
// garbage collection cycle.
func (s *Supervisor) gcMarkForNext() {
	wg := sync.WaitGroup{}
	wg.Add(3)

	// key exchanges
	go func() {
		s.appLog.Debug(`Supervisor.GC: s.gcMarkKex()::start`)
		defer wg.Done()
		s.gcMarkKex()
		s.appLog.Debug(`Supervisor.GC: s.gcMarkKex()::end`)
	}()

	// tokens
	go func() {
		s.appLog.Debug(`Supervisor.GC: s.gcMarkTokens()::start`)
		defer wg.Done()
		s.gcMarkTokens()
		s.appLog.Debug(`Supervisor.GC: s.gcMarkTokens()::end`)
	}()

	// credentials
	go func() {
		s.appLog.Debug(`Supervisor.GC: s.gcMarkCredentials()::start`)
		defer wg.Done()
		s.gcMarkCredentials()
		s.appLog.Debug(`Supervisor.GC: s.gcMarkCredentials()::end`)
	}()
	wg.Wait()
}

// gcMarkKex iterates over stored key exchanges and marks expired ones
// for garbage collection
func (s *Supervisor) gcMarkKex() {
	for kex := range s.kex.interateUnlocked() {
		if kex.IsExpired() {
			s.kex.markUnlocked(kex.Request.String())
		}
	}
}

// gcMarkTokens  iterates over stored tokens and marks expired ones for
// garbage collection
func (s *Supervisor) gcMarkTokens() {
	for token := range s.tokens.iterateUnlocked() {
		if token.isExpired() {
			s.tokens.markUnlocked(hex.EncodeToString(token.binToken))
		}
	}
}

func (s *Supervisor) gcMarkCredentials() {
	for credential := range s.credentials.iterateUnlocked() {
		if credential.isExpired() {
			s.credentials.markUnlocked(credential.name)
		}
	}
}

// gcSweep removes data marked for garbage collection
func (s *Supervisor) gcSweep() {
	wg := sync.WaitGroup{}
	wg.Add(3)

	// sweep key exchanges marked for garbage collection
	go func() {
		s.appLog.Debug(`Supervisor.GC: s.kex.sweepUnlocked()::start`)
		defer wg.Done()
		s.kex.sweepUnlocked()
		s.appLog.Debug(`Supervisor.GC: s.kex.sweepUnlocked()::end`)
	}()

	// sweep tokens marked for garbage collection
	go func() {
		s.appLog.Debug(`Supervisor.GC: s.tokens.sweepUnlocked()::start`)
		defer wg.Done()
		s.tokens.sweepUnlocked()
		s.appLog.Debug(`Supervisor.GC: s.tokens.sweepUnlocked()::end`)
	}()

	// sweep credentials marked for garbage collection
	go func() {
		s.appLog.Debug(`Supervisor.GC: s.credentials.sweepUnlocked()::start`)
		defer wg.Done()
		s.credentials.sweepUnlocked()
		s.appLog.Debug(`Supervisor.GC: s.credentials.sweepUnlocked()::end`)
	}()
	s.appLog.Debug(`Supervisor.GC: s.gcSweep()::waiting`)
	wg.Wait()
	s.appLog.Debug(`Supervisor.GC: s.gcSweep()::done`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
