/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"fmt"
	"time"

	"github.com/mjolnir42/soma/internal/msg"
)

// token handles supervisor requests for token and calls the correct
// function depending on the requested task
func (s *Supervisor) token(q *msg.Request) {
	result := msg.FromRequest(q)
	// default result is for the request to fail
	result.Code = 403
	result.Super.Verdict = 403

	// start response delay timer
	timer := time.NewTimer(1 * time.Second)

	// start assembly of auditlog entry
	result.Super.Audit = singleton.auditLog.
		WithField(`RequestID`, q.ID.String()).
		WithField(`KexID`, q.Super.Encrypted.KexID).
		WithField(`IPAddr`, q.RemoteAddr).
		WithField(`UserName`, `AnonymousCoward`).
		WithField(`UserID`, `ffffffff-ffff-ffff-ffff-ffffffffffff`).
		WithField(`Code`, result.Code).
		WithField(`Verdict`, result.Super.Verdict).
		WithField(`Section`, q.Section).
		WithField(`Action`, q.Action).
		WithField(`RequestType`, fmt.Sprintf("%s/%s", q.Section, q.Action)).
		WithField(`Supervisor`, fmt.Sprintf("%s/%s:%s", q.Section, q.Action, q.Super.Task))

	// tokenRequest/tokenInvalidate are master instance functions
	if s.readonly {
		result.ReadOnly()
		result.Super.Audit.WithField(`Code`, result.Code).Warningln(result.Error)
		goto returnImmediate
	}

	// filter requests with invalid task
	switch q.Super.Task {
	case msg.TaskRequest:
	case msg.TaskInvalidateGlobal:
	case msg.TaskInvalidateAccount:
	case msg.TaskInvalidate:
	default:
		result.UnknownTask(q)
		result.Super.Audit.WithField(`Code`, result.Code).Warningln(result.Error)
		goto returnImmediate
	}

	// select correct taskhandler
	switch q.Super.Task {
	case msg.TaskRequest:
		s.tokenRequest(q, &result)
	case msg.TaskInvalidateGlobal:
		s.tokenInvalidateGlobal(q, &result)
	case msg.TaskInvalidateAccount:
		s.tokenInvalidateAccount(q, &result)
	case msg.TaskInvalidate:
		s.tokenInvalidate(q, &result)
	}

	// wait for delay timer to trigger
	<-timer.C

returnImmediate:
	// cleanup delay timer
	if !timer.Stop() {
		<-timer.C
	}
	q.Reply <- result
}

// tokenInvalidateGlobal invalidates all tokens
func (s *Supervisor) tokenInvalidateGlobal(q *msg.Request, mr *msg.Result) {
	// XXX TODO
}

// tokenInvalidateAccount marks all tokens of a user as invalidate-on-use
func (s *Supervisor) tokenInvalidateAccount(q *msg.Request, mr *msg.Result) {
	// XXX TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
