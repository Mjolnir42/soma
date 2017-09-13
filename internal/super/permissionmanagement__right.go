/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"fmt"

	"github.com/mjolnir42/soma/internal/msg"
)

func (s *Supervisor) right(q *msg.Request) {
	result := msg.FromRequest(q)
	q.Log(s.reqLog)

	if q.Grant.RecipientType != msg.SubjectUser {
		result.NotImplemented(
			fmt.Errorf("Rights for recipient type"+
				" %s are currently not implemented",
				q.Grant.RecipientType))
		goto dispatch
	}

	switch q.Action {
	case msg.ActionGrant, msg.ActionRevoke:
		s.rightWrite(q, &result)
	case msg.ActionSearch:
		s.rightRead(q, &result)
	default:
		result.UnknownRequest(q)
	}

dispatch:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
