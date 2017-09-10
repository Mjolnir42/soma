/*-
Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"fmt"
	"time"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
)

func (s *Supervisor) validateBasicAuth(q *msg.Request) {
	var tok *svToken
	result := msg.FromRequest(q)

	// basic auth always fails for root if root is disabled
	if q.Super.BasicAuth.User == `root` && s.rootDisabled {
		result.ServerError(
			fmt.Errorf(`Attempted authentication on disabled root account`))
		goto unauthorized
	}

	// basic auth always fails for root if root is restricted and
	// the request comes from an unrestricted endpoint. Note: there
	// are currently no restricted endpoints (https over unix socket)
	if q.Super.BasicAuth.User == `root` && s.rootRestricted && !q.Super.RestrictedEndpoint {
		result.ServerError(
			fmt.Errorf(`Attempted root authentication on unrestricted endpoint`))
		goto unauthorized
	}

	tok = s.tokens.read(q.Super.BasicAuth.Token)
	if tok == nil && !s.readonly {
		// rw instance knows every token
		result.ServerError(fmt.Errorf(`Unknown Token (TokenMap)`))
		goto unauthorized
	} else if tok == nil {
		if !s.fetchTokenFromDB(q.Super.BasicAuth.Token) {
			result.ServerError(fmt.Errorf(`Unknown Token (pgSQL)`))
			goto unauthorized
		}
		tok = s.tokens.read(q.Super.BasicAuth.Token)
	}
	if time.Now().UTC().Before(tok.validFrom.UTC()) ||
		time.Now().UTC().After(tok.expiresAt.UTC()) {
		result.Unauthorized(fmt.Errorf(`Token expired`))
		goto unauthorized
	}

	if auth.VerifyExtracted(q.Super.BasicAuth.User, q.RemoteAddr, tok.binToken, s.key,
		s.seed, tok.binExpiresAt, tok.salt) {
		// valid token
		result.Super = &msg.Supervisor{Verdict: 200}
		result.OK()
		q.Reply <- result
	}

unauthorized:
	result.Super = &msg.Supervisor{Verdict: 401}
	q.Reply <- result
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix