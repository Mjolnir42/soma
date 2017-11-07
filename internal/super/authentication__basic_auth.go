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
	"github.com/mjolnir42/soma/lib/auth"
)

// IMPORTANT!
//
// differentiated error returns are for logging purposes only. Failed
// Authentication is returned to the client as 401/Unauthorized.

// authenticateBasicAuth performs BasicAuth authentication
func (s *Supervisor) authenticateBasicAuth(q *msg.Request, mr *msg.Result) {

	// basic auth always fails for root if root is disabled
	if q.Super.BasicAuth.User == msg.SubjectRoot && s.rootDisabled {
		mr.Forbidden(fmt.Errorf(
			`Attempted authentication on disabled root account`))
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return
	}

	// basic auth always fails for root if root is restricted and
	// the request comes from an unrestricted endpoint. Note: there
	// are currently no restricted endpoints (https over unix socket)
	if q.Super.BasicAuth.User == msg.SubjectRoot && s.rootRestricted && !q.Super.RestrictedEndpoint {
		mr.Forbidden(fmt.Errorf(
			`Attempted root authentication on unrestricted endpoint`))
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return
	}

	// unknown or incorrect provided hmac authentication token will
	// fail here, since it will neither be found within the in-memory
	// map or the database
	tk := s.tokens.read(q.Super.BasicAuth.Token)
	if tk == nil && !s.readonly {
		// rw instance knows every token, therefor the token
		// missing from the credential cache is a critical fault
		mr.NotFound(fmt.Errorf(
			`Unknown Token: not found in in-memory TokenMap`))
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return
	} else if tk == nil {
		// load the missing token from the database into the cache
		if !s.fetchTokenFromDB(q.Super.BasicAuth.Token) {
			mr.NotFound(fmt.Errorf(
				`Unknown Token: not found in pgSQL database`))
			mr.Super.Audit.
				WithField(`Code`, mr.Code).
				Warningln(mr.Error)
			return
		}
		tk = s.tokens.read(q.Super.BasicAuth.Token)
	}

	// the token hmac does not contain tk.validFrom,
	// only tk.expiresAt, so test both here
	if time.Now().UTC().Before(tk.validFrom.UTC()) ||
		time.Now().UTC().After(tk.expiresAt.UTC()) {
		mr.Unauthorized(fmt.Errorf(
			`Authentication failed, token invalid: ` +
				`expired or not valid yet`))
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return
	}

	// recalculate hmac and compare with the the binary version of the
	// provided token, this checks that the known token that was
	// provided is used by the correct user and from the source IP
	// address it was issued to.
	if !auth.VerifyExtracted(q.Super.BasicAuth.User, q.RemoteAddr,
		tk.binToken, s.key, s.seed, tk.binExpiresAt, tk.salt) {
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(`Authentication failed`)
		return
	}

	// check if the tokens for the user have been revoked
	if revokedAt, revoked := s.tokens.isExpired(
		q.Super.BasicAuth.User,
	); revoked {
		// token was issued prior to the revocation
		if tk.validFrom.UTC().Before(revokedAt.UTC()) {
			// issue invalidate request for this specific token
			returnChannel := make(chan msg.Result)
			request := msg.Request{
				ID:         q.ID,
				Section:    msg.SectionSupervisor,
				Action:     msg.ActionToken,
				Reply:      returnChannel,
				RemoteAddr: q.RemoteAddr,
				AuthUser:   q.Super.BasicAuth.User,
				Super: &msg.Supervisor{
					Task:      msg.TaskInvalidate,
					AuthToken: q.Super.BasicAuth.Token,
				},
			}
			s.Input <- request
			<-returnChannel

			mr.Unauthorized(fmt.Errorf(
				`Authentication failed, tokens for user are revoked`))
			mr.Super.Audit.
				WithField(`Code`, mr.Code).
				Warningln(mr.Error)
			return
		}
	}

	// the provided hmac token was valid
	mr.OK()
	mr.Super.Verdict = 200
	mr.Super.Audit.
		WithField(`Code`, mr.Code).
		WithField(`Verdict`, mr.Super.Verdict).
		Infoln(`Authentication OK`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
