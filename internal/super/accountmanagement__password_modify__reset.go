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

// IMPORTANT!
//
// all errors returned from encrypted supervisor methods are
// returned to the client as 403/Forbidden. Provided error details
// are used only for serverside logging.

// passwordReset performs the required verification for a password
// reset
func (s *Supervisor) passwordReset(token *auth.Token, mr *msg.Result) bool {

	// decrypt e2e encrypted request
	// token.UserName is the username
	// token.Password is the _NEW_ password that should be set
	// token.Token    is the ldap password or mailtoken

	// validate external credentials
	switch s.activation {
	case `ldap`:
		if ok, err := validateLdapCredentials(token.UserName, token.Token); err != nil {
			mr.ServerError(err, mr.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
			return false
		} else if !ok {
			mr.Forbidden(fmt.Errorf(`Invalid LDAP credentials`))
			mr.Super.Audit.
				WithField(`Code`, mr.Code).
				Warningln(mr.Error)
			return false
		}
	case `token`:
		mr.NotImplemented(fmt.Errorf(`Mail-Token not supported yet`),
			mr.Section)
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return false
	default:
		mr.ServerError(fmt.Errorf("Unknown activation method: %s",
			s.conf.Auth.Activation), mr.Section)
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return false
	}

	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
