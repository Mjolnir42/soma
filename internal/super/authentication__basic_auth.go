/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"fmt"
	"time"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
)

func (s *Supervisor) authenticate(q *msg.Request) {
	var tok *token
	result := msg.FromRequest(q)
	result.Code = 401
	result.Super.Verdict = 401

	// assembly of the auditlog entry
	audit := s.auditLog.
		WithField(`RequestID`, q.ID.String()).
		WithField(`IPAddr`, q.RemoteAddr).
		WithField(`UserName`, q.Super.BasicAuth.User).
		WithField(`Section`, q.Section).
		WithField(`Action`, q.Action).
		WithField(`Code`, result.Code).
		WithField(`Verdict`, result.Super.Verdict).
		WithField(`RequestType`, fmt.Sprintf("%s/%s", q.Section, q.Action)).
		WithField(`Supervisor`, fmt.Sprintf("%s/%s:%s", msg.SectionSupervisor, msg.ActionAuthenticate, `basicAuth`))

	// basic auth always fails for root if root is disabled
	if q.Super.BasicAuth.User == `root` && s.rootDisabled {
		err := fmt.Errorf(`Attempted authentication on disabled root account`)
		result.Forbidden(err)
		audit.WithField(`Code`, result.Code).Warningln(err)
		goto unauthorized
	}

	// basic auth always fails for root if root is restricted and
	// the request comes from an unrestricted endpoint. Note: there
	// are currently no restricted endpoints (https over unix socket)
	if q.Super.BasicAuth.User == `root` && s.rootRestricted && !q.Super.RestrictedEndpoint {
		err := fmt.Errorf(`Attempted root authentication on unrestricted endpoint`)
		result.Forbidden(err)
		audit.WithField(`Code`, result.Code).Warningln(err)
		goto unauthorized
	}

	tok = s.tokens.read(q.Super.BasicAuth.Token)
	if tok == nil && !s.readonly {
		// rw instance knows every token
		// 404 is only logged, BasicAuth always replies 401 on error
		err := fmt.Errorf(`Unknown Token: not found in in-memory TokenMap`)
		result.NotFound(err)
		audit.WithField(`Code`, result.Code).Warningln(err)
		goto unauthorized
	} else if tok == nil {
		if !s.fetchTokenFromDB(q.Super.BasicAuth.Token) {
			// 404 is only logged, BasicAuth always replies 401 on error
			err := fmt.Errorf(`Unknown Token: not found in pgSQL database`)
			result.NotFound(err)
			audit.WithField(`Code`, result.Code).Warningln(err)
			goto unauthorized
		}
		tok = s.tokens.read(q.Super.BasicAuth.Token)
	}
	if time.Now().UTC().Before(tok.validFrom.UTC()) ||
		time.Now().UTC().After(tok.expiresAt.UTC()) {
		err := fmt.Errorf(`Authentication failed, token invalid: expired or not valid yet`)
		result.Unauthorized(err)
		audit.WithField(`Code`, result.Code).Warningln(err)
		goto unauthorized
	}

	if auth.VerifyExtracted(q.Super.BasicAuth.User, q.RemoteAddr, tok.binToken, s.key,
		s.seed, tok.binExpiresAt, tok.salt) {
		// valid token
		result.Super.Verdict = 200
		result.OK()
		audit.WithField(`Code`, result.Code).
			WithField(`Verdict`, result.Super.Verdict).
			Infoln(`Authentication OK`)
	} else {
		audit.WithField(`Code`, result.Code).Warningln(`Authentication failed`)
	}

unauthorized:
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
