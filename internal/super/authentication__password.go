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

	"github.com/mjolnir42/scrypth64"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
)

// authenticatePassword verifies the provided password in the
// token.Token field
func (s *Supervisor) authenticatePassword(token *auth.Token, mr *msg.Result) bool {
	// read current credentials
	cred := s.credentials.read(token.UserName)

	// check current credentials are still valid
	if time.Now().UTC().Before(cred.validFrom.UTC()) ||
		time.Now().UTC().After(cred.expiresAt.UTC()) {
		mr.Forbidden(fmt.Errorf("Expired credentials: %s",
			token.UserName))
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return false
	}

	// verify the provided password is correct
	if ok, err := scrypth64.Verify(
		token.Token,
		cred.cryptMCF,
	); err != nil {
		mr.ServerError(err, mr.Section)
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return false
	} else if !ok {
		mr.Forbidden(fmt.Errorf(`Credentials do not match`))
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
