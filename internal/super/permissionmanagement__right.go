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

func (s *Supervisor) right(q *msg.Request) {
	result := msg.FromRequest(q)

	// start assembly of auditlog entry
	result.Super.Audit = s.auditLog.
		WithField(`RequestID`, q.ID.String()).
		WithField(`IPAddr`, q.RemoteAddr).
		WithField(`UserName`, q.AuthUser).
		WithField(`Section`, q.Section).
		WithField(`Action`, q.Action)

	switch q.Action {
	case msg.ActionGrant, msg.ActionRevoke:
		s.rightWrite(q, &result)
	case msg.ActionList, msg.ActionSearch:
		s.rightRead(q, &result)
	default:
		result.UnknownRequest(q)
		result.Super.Audit.
			WithField(`Code`, result.Code).
			Warningln(result.Error)
	}

	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
