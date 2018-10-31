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

// IMPORTANT!
//
// all errors returned from encrypted supervisor methods are
// returned to the client as 403/Forbidden. Provided error details
// are used only for serverside logging.

// password handles supervisor requests that modify an account's
// stored password credentials and calls the correct function
// depending on the requested task
func (s *Supervisor) password(q *msg.Request) {
	result := msg.FromRequest(q)
	// default result is for the request to fail
	result.Code = 403
	result.Super.Verdict = 403

	// start response delay timer
	timer := time.NewTimer(1 * time.Second)

	// start assembly of auditlog entry
	result.Super.Audit = s.auditLog.
		WithField(`RequestID`, q.ID.String()).
		WithField(`KexID`, q.Super.Encrypted.KexID).
		WithField(`IPAddr`, q.RemoteAddr).
		WithField(`UserName`, `AnonymousCoward`).
		WithField(`UserID`, `ffffffff-ffff-ffff-ffff-ffffffffffff`).
		WithField(`Section`, q.Section).
		WithField(`Action`, q.Action).
		WithField(`Code`, result.Code).
		WithField(`Verdict`, result.Super.Verdict).
		WithField(`Request`, fmt.Sprintf(
			"%s::%s", q.Section, q.Action)).
		WithField(`Supervisor`, fmt.Sprintf(
			"%s::%s=%s", q.Section, q.Action, q.Super.Task))

	// password changes are master instance functions
	if s.readonly {
		result.ReadOnly()
		result.Super.Audit.
			WithField(`Code`, result.Code).
			Warningln(result.Error)
		goto returnImmediate
	}

	// filter requests with invalid task
	switch q.Super.Task {
	case msg.TaskReset:
	case msg.TaskChange:
	case msg.TaskRevoke:
	default:
		result.UnknownTask(q)
		result.Super.Audit.
			WithField(`Code`, result.Code).
			Warningln(result.Error)
		goto returnImmediate
	}

	// select correct taskhandler
	switch q.Super.Task {
	case msg.TaskReset, msg.TaskChange:
		s.passwordModify(q, &result)
	case msg.TaskRevoke:
		// TaskRevoke is internal and not connected to the user
		// interface, no need to add the processing delay
		s.passwordRevoke(q, &result)
		goto returnImmediate
	}

	// wait for delay timer to trigger
	<-timer.C

returnImmediate:
	// cleanup delay timer
	if timer.Stop() {
		<-timer.C
	}
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
