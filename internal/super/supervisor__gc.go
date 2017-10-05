/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

func (s *Supervisor) gc() {
	// TODO purge expired KEX in s.kex
	// TODO purge expired tokens in s.tokens
	// TODO purge expired credentials in s.credentials
}

func (s *Supervisor) gcMark() {
}

func (s *Supervisor) gcSweep() {
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
