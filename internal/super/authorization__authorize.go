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
		},
	}
	result := <-returnChannel

	// the request is authorized
	if result.Super.Verdict == 200 {
		// TODO auditlog
		return true
	}

	// the request is not authorized
	// XXX should be auditlog
	singleton.errLog.Printf(
		"Section=%s, Action=%s, InternalCode=%d, Error=%s",
		msg.SectionSupervisor,
		msg.ActionAuthorize,
		result.Super.Verdict,
		fmt.Sprintf("Forbidden: %s, %s, %s/%s",
			q.AuthUser,
			q.RemoteAddr,
			q.Section,
			q.Action),
	)
	return false
}

// authorize forwards the request to the permission cache for
// assessment
func (s *Supervisor) authorize(q *msg.Request) {
	q.Reply <- s.permCache.IsAuthorized(q)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
