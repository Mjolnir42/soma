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
	"github.com/mjolnir42/soma/lib/auth"
)

// authenticateRootToken verifies the provided token to unlock
// the root account
func (s *Supervisor) authenticateRootToken(token *auth.Token, mr *msg.Result) bool {
	var rootToken string
	var err error

	// read rootToken from database
	if rootToken, err = s.fetchRootToken(); err != nil {
		mr.ServerError(err)
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return false
	}

	if token.Token != rootToken || len(token.Password) == 0 {
		mr.Forbidden(fmt.Errorf(`Root activation failed`))
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return false
	}
	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
