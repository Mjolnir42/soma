/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"fmt"

	"github.com/mjolnir42/soma/internal/config"
	"github.com/mjolnir42/soma/internal/msg"
)

// cfg is a copy of the config.Config given to Supervisor so
// exported functions can access it
var cfg *config.Config

// IsAuthorized returns if the request is permitted
func IsAuthorized(q *msg.Request) bool {
	// instance is configured as wild-west instance
	if cfg.OpenInstance {
		return true
	}

	// assembly of the auditlog entry
	audit := singleton.auditLog.
		WithField(`RequestID`, q.ID.String()).
		WithField(`IPAddr`, q.RemoteAddr).
		WithField(`UserName`, q.AuthUser).
		WithField(`Section`, q.Section).
		WithField(`Action`, q.Action).
		WithField(`Code`, 403).
		WithField(`Verdict`, 403).
		WithField(`Request`, fmt.Sprintf("%s::%s", q.Section, q.Action)).
		WithField(`Supervisor`, fmt.Sprintf("%s::%s=%s", msg.SectionSupervisor, msg.ActionAuthorize, msg.ActionAuthorize))

	// the original request is wrapped because the http handler
	// is reading from q.Reply
	returnChannel := make(chan msg.Result)
	singleton.Input <- msg.Request{
		ID:      q.ID,
		Section: msg.SectionSupervisor,
		Action:  msg.ActionAuthorize,
		Reply:   returnChannel,
		Super: &msg.Supervisor{
			Authorize: q,
			Audit:     audit,
		},
	}
	result := <-returnChannel

	switch result.Super.Verdict {
	case 200:
		// the request is authorized
		result.Super.Audit.WithField(`Code`, result.Code).
			WithField(`Verdict`, result.Super.Verdict).
			Infoln(`Request authorized`)
		return true
	default:
		result.Super.Audit.WithField(`Code`, result.Code).
			WithField(`Verdict`, result.Super.Verdict).
			Warningln(`Request forbidden`)
	}

	return false
}

// authorize forwards the request to the permission cache for
// assessment
func (s *Supervisor) authorize(q *msg.Request) {
	q.Reply <- s.permCache.IsAuthorized(q)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
