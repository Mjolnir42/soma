package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"fmt"
	"strings"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
)

// IMPORTANT!
//
// all errors returned from encrypted supervisor methods are
// returned to the client as 403/Forbidden. Provided error details
// are used only for serverside logging.

// passwordChange performs the required verification for a password
// change
func (s *Supervisor) passwordChange(token *auth.Token, mr *msg.Result) bool {

	// token.UserName is the username
	// token.Password is the _NEW_ password that should be set
	// token.Token    is the old password

	// validate provided credentials
	if !s.authenticatePassword(token, mr) {
		return false
	}

	// check if the new password is the same, prefixed or suffixed
	if token.Token == token.Password ||
		strings.HasPrefix(token.Password, token.Token) ||
		strings.HasSuffix(token.Password, token.Token) {
		mr.BadRequest(fmt.Errorf(`New password contains old password`))
		mr.Super.Audit.
			WithField(`Code`, mr.Code).
			Warningln(mr.Error)
		return false
	}

	return true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
