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

// IMPORTANT!
//
// differentiated error returns are for logging purposes only. Failed
// Authentication is returned to the client as 401/Unauthorized.

// authenticate handles supervisor requests for authentication
func (s *Supervisor) authenticate(q *msg.Request) {
	result := msg.FromRequest(q)
	// default result is for the request to fail
	result.Code = 401
	result.Super.Verdict = 401

	// assembly of the auditlog entry
	result.Super.Audit = s.auditLog.
		WithField(`RequestID`, q.ID.String()).
		WithField(`IPAddr`, q.RemoteAddr).
		WithField(`UserName`, q.Super.BasicAuth.User).
		WithField(`Section`, q.Section).
		WithField(`Action`, q.Action).
		WithField(`Code`, result.Code).
		WithField(`Verdict`, result.Super.Verdict).
		WithField(`RequestType`, fmt.Sprintf(
			"%s/%s", q.Section, q.Action)).
		WithField(`Supervisor`, fmt.Sprintf(
			"%s/%s:%s", q.Section, q.Action, q.Super.Task))

	// filter requests with invalid task
	switch q.Super.Task {
	case msg.TaskBasicAuth:
	default:
		result.UnknownTask(q)
		result.Super.Audit.
			WithField(`Code`, result.Code).
			Warningln(result.Error)
		goto unauthorized
	}

	// select correct taskhandler
	switch q.Super.Task {
	case msg.TaskBasicAuth:
		s.authenticateBasicAuth(q, &result)
	}

unauthorized:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
