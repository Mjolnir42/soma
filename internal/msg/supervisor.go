/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg

import (
	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/lib/auth"
	"github.com/mjolnir42/soma/lib/proto"
)

type Supervisor struct {
	Task               string
	Verdict            uint16
	RestrictedEndpoint bool
	// KeyExchange Data
	Kex auth.Kex
	// Fields for encrypted requests
	Encrypted struct {
		KexID string
		Data  []byte
	}
	// Fields for basic authentication requests
	BasicAuth struct {
		User  string
		Token string
	}
	// The active token to be invalidated
	AuthToken string
	// Request to be authorized
	Authorize *Request
	// AuditLog Entry for this supervisor task
	Audit *logrus.Entry
	// XXX Everything below is deprecated
	// Fields for map update notifications
	Object string
	User   proto.User
	Team   proto.Team
	// Fields for Grant revocation
	GrantId string
}

func (s *Supervisor) Clear() {
	// Task and Verdict are kept intact
	s.RestrictedEndpoint = true
	s.Kex = auth.Kex{}
	s.Encrypted = struct {
		KexID string
		Data  []byte
	}{}
	s.BasicAuth = struct {
		User  string
		Token string
	}{}
	s.Authorize = nil
	s.Object = ``
	s.User = proto.User{}
	s.Team = proto.Team{}
	s.GrantId = ``
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
