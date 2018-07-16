/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// TeamMgmtList function
func (x *Rest) TeamMgmtList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionTeamMgmt
	request.Action = msg.ActionList
	request.Flag.Unscoped = true

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// TeamMgmtSearch function
func (x *Rest) TeamMgmtSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewTeamFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Filter.Team.Name == `` {
		dispatchBadRequest(&w, fmt.Errorf(
			`TeamMgmtSearch request missing Team.Name`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionTeamMgmt
	request.Action = msg.ActionSearch
	request.Search.Team.Name = cReq.Filter.Team.Name
	request.Flag.Unscoped = true

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply

	// XXX BUG filter in SQL statement
	filtered := []proto.Team{}
	for _, i := range result.Team {
		if i.Name == cReq.Filter.Team.Name {
			filtered = append(filtered, i)
		}
	}
	result.Team = filtered
	sendMsgResult(&w, &result)
}

// TeamMgmtShow function
func (x *Rest) TeamMgmtShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionTeamMgmt
	request.Action = msg.ActionShow
	request.Team = proto.Team{
		ID: params.ByName(`teamID`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// TeamMgmtSync function
func (x *Rest) TeamMgmtSync(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionTeamMgmt
	request.Action = msg.ActionSync
	request.Flag.Unscoped = true

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// TeamMgmtAdd function
func (x *Rest) TeamMgmtAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewTeamRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionTeamMgmt
	request.Action = msg.ActionAdd
	request.Team = cReq.Team.Clone()

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// TeamMgmtUpdate function
func (x *Rest) TeamMgmtUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewTeamRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if params.ByName(`teamID`) != cReq.Team.ID {
		dispatchBadRequest(&w, fmt.Errorf(
			"Mismatched teamID: %s vs %s",
			params.ByName(`teamID`),
			cReq.Team.ID,
		))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionTeamMgmt
	request.Action = msg.ActionUpdate
	request.Team = cReq.Team.Clone()

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// TeamMgmtRemove function
func (x *Rest) TeamMgmtRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionCapability
	request.Action = msg.ActionRemove
	request.Team = proto.Team{
		ID: params.ByName(`teamID`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
