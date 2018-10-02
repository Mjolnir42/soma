/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg // import "github.com/mjolnir42/soma/internal/msg"

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

type Request struct {
	ID            uuid.UUID
	Section       string
	Action        string
	TargetEntity  string
	RemoteAddr    string
	AuthUser      string
	RequestURI    string
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
	System      proto.System
	Team        proto.Team
	Tree        proto.Tree
	Unit        proto.Unit
	User        proto.User
	Validity    proto.Validity
	View        proto.View
	Workflow    proto.Workflow
}

// New returns a Request
func New(r *http.Request, params httprouter.Params) Request {
	returnChannel := make(chan Result, 1)
	return Request{
		ID:         requestID(params),
		RequestURI: requestURI(params),
		RemoteAddr: remoteAddr(r),
		AuthUser:   authUser(params),
		Reply:      returnChannel,
	}
}

type Filter struct {
	IsDetailed bool
	ActionObj  proto.Action
	Bucket     proto.Bucket
	Cluster    proto.Cluster
	Grant      proto.Grant
	Group      proto.Group
	Job        proto.JobFilter
	Level      proto.Level
	Monitoring proto.Monitoring
	Node       proto.Node
	Oncall     proto.Oncall
	Permission proto.Permission
	Property   proto.Property
	Repository proto.Repository
	SectionObj proto.Section
	Server     proto.Server
	Team       proto.Team
	User       proto.User
}

type UpdateData struct {
	Cluster     proto.Cluster
	Datacenter  proto.Datacenter
	Entity      proto.Entity
	Environment proto.Environment
	Oncall      proto.Oncall
	Property    proto.Property
	Repository  proto.Repository
	Server      proto.Server
	State       proto.State
	User        proto.User
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
