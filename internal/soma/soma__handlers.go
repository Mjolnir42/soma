/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import "github.com/mjolnir42/soma/internal/super"

// getSupervisor returns the supervisor from the handlermap
func (s *Soma) getSupervisor() *super.Supervisor {
	return s.handlerMap.Get(`supervisor`).(*super.Supervisor)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
