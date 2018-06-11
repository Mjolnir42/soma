/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg

import (
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

type Request struct {
	ID            uuid.UUID
	Section       string
	Action        string
	RemoteAddr    string
	AuthUser      string
	Reply         chan Result
	JobID         uuid.UUID
	Search        Filter
	Update        UpdateData
	Flag          Flags
	DeploymentIDs []string

	Super *Supervisor
	Cache *Request

	ActionObj   proto.Action
	Attribute   proto.Attribute
	Bucket      proto.Bucket
	Capability  proto.Capability
	Category    proto.Category
	CheckConfig proto.CheckConfig
	Cluster     proto.Cluster
	Datacenter  proto.Datacenter
	Deployment  proto.Deployment
	Entity      proto.Entity
	Environment proto.Environment
	Grant       proto.Grant
	Group       proto.Group
	Instance    proto.Instance
	Job         proto.Job
	Level       proto.Level
	Metric      proto.Metric
	Mode        proto.Mode
	Monitoring  proto.Monitoring
	Node        proto.Node
	Oncall      proto.Oncall
	Permission  proto.Permission
	Predicate   proto.Predicate
	Property    proto.Property
	Provider    proto.Provider
	Repository  proto.Repository
	SectionObj  proto.Section
	Server      proto.Server
	State       proto.State
	Status      proto.Status
	System      proto.SystemOperation
	Team        proto.Team
	Tree        proto.Tree
	Unit        proto.Unit
	User        proto.User
	Validity    proto.Validity
	View        proto.View
	Workflow    proto.Workflow
}

type Filter struct {
	IsDetailed bool
	Job        proto.JobFilter
	Server     proto.Server
}

type UpdateData struct {
	Datacenter  proto.Datacenter
	Entity      proto.Entity
	Environment proto.Environment
	State       proto.State
	View        proto.View
}

type Flags struct {
	JobDetail    bool
	Unscoped     bool
	Rebuild      bool
	RebuildLevel string
}

func CacheUpdateFromRequest(rq *Request) Request {
	return Request{
		Section: SectionSupervisor,
		Action:  ActionCacheUpdate,
		Cache:   rq,
	}
}

// Log logs the request to the provided logwriter
func (r *Request) Log(l *logrus.Logger) {
	l.Infof(LogStrSRq,
		r.Section,
		r.Action,
		r.AuthUser,
		r.RemoteAddr,
	)
}

// New returns a new request
func New(r *http.Request, params httprouter.Params) Request {
	returnChannel := make(chan Result, 1)
	return Request{
		ID:         requestID(params),
		RemoteAddr: remoteAddr(r),
		AuthUser:   authUser(params),
		Reply:      returnChannel,
	}
}

// requestID extracts the RequestID set by Basic Authentication, making
// the ID consistent between all logs
func requestID(params httprouter.Params) (id uuid.UUID) {
	id, _ = uuid.FromString(params.ByName(`RequestID`))
	return
}

// authUser returns the extracted authenticated user
func authUser(params httprouter.Params) string {
	return params.ByName(`AuthenticatedUser`)
}

// remoteAddr extracts the IP address part of the IP:port string
// set as net/http.Request.RemoteAddr. It handles IPv4 cases like
// 192.0.2.1:48467 and IPv6 cases like [2001:db8::1%lo0]:48467
func remoteAddr(r *http.Request) string {
	var addr string

	switch {
	case strings.Contains(r.RemoteAddr, `]`):
		// IPv6 address [2001:db8::1%lo0]:48467
		addr = strings.Split(r.RemoteAddr, `]`)[0]
		addr = strings.Split(addr, `%`)[0]
		addr = strings.TrimLeft(addr, `[`)
	default:
		// IPv4 address 192.0.2.1:48467
		addr = strings.Split(r.RemoteAddr, `:`)[0]
	}
	return addr
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
