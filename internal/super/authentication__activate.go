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

// activate handles supervisor requests for account activation and
// calls the correct function depending on the requested task
func (s *Supervisor) activate(q *msg.Request) {
	result := msg.FromRequest(q)
	result.Super.Verdict = 403

	// start response delay timer
	timer := time.NewTimer(1 * time.Second)

	// start assembly of auditlog entry
	logEntry := singleton.auditLog.
		WithField(`RequestID`, q.ID.String()).
		WithField(`KexID`, q.Super.Encrypted.KexID).
		WithField(`IPAddr`, q.RemoteAddr)

	// account activations are master instance functions
	if s.readonly {
		result.ReadOnly()
		goto returnImmediate
	}

	// filter requests with invalid task
	switch q.Super.Task {
	case msg.SubjectUser:
	default:
		result.UnknownRequest(q)
		goto returnImmediate
	}

	// update auditlog entry
	logEntry = logEntry.WithField(`Type`,
		fmt.Sprintf("%s/%s:%s", q.Section, q.Action, q.Super.Task))

	// select correct taskhandler
	switch q.Super.Task {
	case msg.SubjectUser:
		s.activateUser(q, &result, logEntry)
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
