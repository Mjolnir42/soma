/*-
 * Copyright (c) 2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// ScopeSelectMonitoringList function
func (x *Rest) ScopeSelectMonitoringList(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionMonitoringMgmt
	request.Action = msg.ActionAll
	request.Flag.Unscoped = true

	if x.isAuthorized(&request) {
		x.MonitoringMgmtAll(w, r, params)
		return
	}

	x.MonitoringList(w, r, params)
}

// ScopeSelectMonitoringSearch function
func (x *Rest) ScopeSelectMonitoringSearch(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionMonitoringMgmt
	request.Action = msg.ActionSearchAll
	request.Flag.Unscoped = true

	if x.isAuthorized(&request) {
		x.MonitoringMgmtSearchAll(w, r, params)
		return
	}

	x.MonitoringSearch(w, r, params)
}

// ScopeSelectInstanceList function
func (x *Rest) ScopeSelectInstanceList(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionInstanceMgmt
	request.Action = msg.ActionAll
	request.Flag.Unscoped = true

	if x.isAuthorized(&request) {
		x.InstanceMgmtAll(w, r, params)
		return
	}

	x.InstanceList(w, r, params)
}

// ScopeSelectInstanceShow function
func (x *Rest) ScopeSelectInstanceShow(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionInstanceMgmt
	request.Action = msg.ActionShow
	request.Flag.Unscoped = true

	if x.isAuthorized(&request) {
		x.InstanceMgmtShow(w, r, params)
		return
	}

	x.InstanceShow(w, r, params)
}

// ScopeSelectJobList function
func (x *Rest) ScopeSelectJobList(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionJobMgmt
	request.Action = msg.ActionList
	request.Flag.Unscoped = true

	if x.isAuthorized(&request) {
		x.JobMgmtList(w, r, params)
		return
	}

	x.JobList(w, r, params)
}

// ScopeSelectJobWait function
func (x *Rest) ScopeSelectJobWait(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionJobMgmt
	request.Action = msg.ActionWait
	request.Job.ID = params.ByName(`jobID`)
	request.Flag.Unscoped = true

	if x.isAuthorized(&request) {
		x.JobMgmtWait(w, r, params)
		return
	}

	x.JobWait(w, r, params)
}

// ScopeSelectUserShow function
func (x *Rest) ScopeSelectUserShow(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionUserMgmt
	request.Action = msg.ActionShow
	request.User.ID = params.ByName(`userID`)
	request.Flag.Unscoped = true

	if x.isAuthorized(&request) {
		x.UserMgmtShow(w, r, params)
		return
	}

	x.UserShow(w, r, params)
}

// ScopeSelectUserSearch function
func (x *Rest) ScopeSelectUserSearch(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewUserFilter()
	if err := peekJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionUserMgmt
	request.Action = msg.ActionSearch
	request.Search.User.UserName = cReq.Filter.User.UserName
	request.Flag.Unscoped = true

	if x.isAuthorized(&request) {
		x.UserMgmtSearch(w, r, params)
		return
	}

	x.UserSearch(w, r, params)
}

// ScopeSelectTeamShow function
func (x *Rest) ScopeSelectTeamShow(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {

	request := msg.New(r, params)
	request.Section = msg.SectionTeamMgmt
	request.Action = msg.ActionShow
	request.Team.ID = params.ByName(`teamID`)
	request.Flag.Unscoped = true

	if x.isAuthorized(&request) {
		x.TeamMgmtShow(w, r, params)
		return
	}

	x.TeamShow(w, r, params)
}

// ScopeSelectTeamSearch function
func (x *Rest) ScopeSelectTeamSearch(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewTeamFilter()
	if err := peekJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionTeamMgmt
	request.Action = msg.ActionSearch
	request.Search.Team.Name = cReq.Filter.Team.Name
	request.Flag.Unscoped = true

	if x.isAuthorized(&request) {
		x.TeamMgmtSearch(w, r, params)
		return
	}

	x.TeamSearch(w, r, params)
}

// ScopeSelectRepositorySearch function
func (x *Rest) ScopeSelectRepositorySearch(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewRepositoryFilter()
	if err := peekJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionRepository
	request.Action = msg.ActionSearch

	if x.isAuthorized(&request) {
		x.RepositorySearch(w, r, params)
		return
	}

	x.RepositoryConfigSearch(w, r, params)
}

// ScopeSelectRepositoryShow function
func (x *Rest) ScopeSelectRepositoryShow(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepository
	request.Action = msg.ActionShow
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Repository.TeamID = params.ByName(`teamID`)

	if x.isAuthorized(&request) {
		x.RepositoryShow(w, r, params)
		return
	}

	x.RepositoryConfigShow(w, r, params)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
