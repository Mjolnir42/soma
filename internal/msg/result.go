/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package msg

import (
	"fmt"

	"github.com/mjolnir42/soma/lib/proto"
)

type Result struct {
	Section string
	Action  string
	Code    uint16
	Error   error
	JobId   string

	Super *Supervisor

	ActionObj   []proto.Action
	Attribute   []proto.Attribute
	Bucket      []proto.Bucket
	Capability  []proto.Capability
	Category    []proto.Category
	Cluster     []proto.Cluster
	Datacenter  []proto.Datacenter
	Entity      []proto.Entity
	Environment []proto.Environment
	Grant       []proto.Grant
	Group       []proto.Group
	Instance    []proto.Instance
	Job         []proto.Job
	Level       []proto.Level
	Metric      []proto.Metric
	Mode        []proto.Mode
	Monitoring  []proto.Monitoring
	Node        []proto.Node
	Oncall      []proto.Oncall
	Permission  []proto.Permission
	Predicate   []proto.Predicate
	Property    []proto.Property
	Provider    []proto.Provider
	Repository  []proto.Repository
	SectionObj  []proto.Section
	Server      []proto.Server
	State       []proto.State
	Status      []proto.Status
	System      []proto.SystemOperation
	Team        []proto.Team
	Tree        proto.Tree
	Unit        []proto.Unit
	User        []proto.User
	Validity    []proto.Validity
	View        []proto.View
	Workflow    []proto.Workflow
}

func FromRequest(rq *Request) Result {
	return Result{
		Section: rq.Section,
		Action:  rq.Action,
	}
}

func CacheUpdateFromRequest(rq *Request) Request {
	return Request{
		Section:      `cache`,
		Action:       `update`,
		CacheRequest: rq,
	}
}

func (r *Result) IsOK() bool {
	switch r.Code {
	case 200:
		return true
	default:
		return false
	}
}

func (r *Result) RowCnt(i int64, err error) bool {
	if err != nil {
		r.ServerError(err)
		return false
	}
	switch i {
	case 0:
		r.OK()
		r.SetError(fmt.Errorf(`No rows affected`))
	case 1:
		r.OK()
		return true
	default:
		r.ServerError(fmt.Errorf("Too many rows affected: %d", i))
	}
	return false
}

func (r *Result) Clear(s string) {
	switch s {
	case `action`:
		r.ActionObj = []proto.Action{}
	case `attribute`:
		r.Attribute = []proto.Attribute{}
	case `bucket`:
		r.Bucket = []proto.Bucket{}
	case `capability`:
		r.Capability = []proto.Capability{}
	case `category`:
		r.Category = []proto.Category{}
	case `cluster`:
		r.Cluster = []proto.Cluster{}
	case `datacenter`:
		r.Datacenter = []proto.Datacenter{}
	case `entity`:
		r.Entity = []proto.Entity{}
	case `environment`:
		r.Environment = []proto.Environment{}
	case `grant`:
		r.Grant = []proto.Grant{}
	case `group`:
		r.Group = []proto.Group{}
	case `instance`:
		r.Instance = []proto.Instance{}
	case `job`:
		r.Job = []proto.Job{}
	case `level`:
		r.Level = []proto.Level{}
	case `metric`:
		r.Metric = []proto.Metric{}
	case `mode`:
		r.Mode = []proto.Mode{}
	case `monitoringsystem`:
		r.Monitoring = []proto.Monitoring{}
	case `node`:
		r.Node = []proto.Node{}
	case `oncall`:
		r.Oncall = []proto.Oncall{}
	case `permission`:
		r.Permission = []proto.Permission{}
	case `predicate`:
		r.Predicate = []proto.Predicate{}
	case `property`:
		r.Property = []proto.Property{}
	case `provider`:
		r.Provider = []proto.Provider{}
	case `repository`:
		r.Repository = []proto.Repository{}
	case `section`:
		r.SectionObj = []proto.Section{}
	case `server`:
		r.Server = []proto.Server{}
	case `state`:
		r.State = []proto.State{}
	case `status`:
		r.Status = []proto.Status{}
	case `system`:
		r.System = []proto.SystemOperation{}
	case `team`:
		r.Team = []proto.Team{}
	case `unit`:
		r.Unit = []proto.Unit{}
	case `user`:
		r.User = []proto.User{}
	case `validity`:
		r.Validity = []proto.Validity{}
	case `view`:
		r.View = []proto.View{}
	case `workflow`:
		r.Workflow = []proto.Workflow{}
	}
}

func (r *Result) SetError(err error) {
	if err != nil {
		r.Error = err
	}
}

func (r *Result) OK() {
	r.Code = 200
	r.Error = nil
}

func (r *Result) Accepted() {
	r.Code = 202
	r.Error = nil
}

func (r *Result) Partial() {
	r.Code = 206
	r.Error = nil
}

func (r *Result) BadRequest(err error, section ...string) {
	r.Code = 400
	r.SetError(err)
	if len(section) > 0 {
		r.Clear(section[0])
	}
}

func (r *Result) Unauthorized(err error) {
	r.Code = 401
	r.SetError(err)
}

func (r *Result) Forbidden(err error) {
	r.Code = 403
	r.SetError(err)
}

func (r *Result) NotFound(err error, section ...string) {
	r.Code = 404
	r.SetError(err)
	if len(section) > 0 {
		r.Clear(section[0])
	}
}

func (r *Result) Conflict(err error) {
	r.Code = 406
	r.SetError(err)
}

func (r *Result) ServerError(err error, section ...string) {
	r.Code = 500
	r.SetError(err)
	if len(section) > 0 {
		r.Clear(section[0])
	}
}

func (r *Result) NotImplemented(err error) {
	r.Code = 501
	r.SetError(err)
}

func (r *Result) Unavailable(err error) {
	r.Code = 503
	r.SetError(err)
}

func (r *Result) UnknownRequest(q *Request) {
	r.NotImplemented(fmt.Errorf("Unknown requested action:"+
		" %s/%s", q.Section, q.Action))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
