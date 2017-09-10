/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"github.com/mjolnir42/soma/internal/msg"
)

func (s *Supervisor) section(q *msg.Request) {
	result := msg.FromRequest(q)
	q.Log(s.reqLog)

	switch q.Action {
	case msg.ActionList, msg.ActionShow, msg.ActionSearch:
		s.sectionRead(q, &result)
	case msg.ActionAdd, msg.ActionRemove:
		s.sectionWrite(q, &result)
	default:
		result.UnknownRequest(q)
	}

	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
