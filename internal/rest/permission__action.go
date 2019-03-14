/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2017, Jörg Pernfuß
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

// ActionList accepts requests to list actions in a specific section
func (x *Rest) ActionList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionAction
	request.Action = msg.ActionList
	request.ActionObj = proto.Action{
		SectionID: params.ByName(`sectionID`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ActionShow accepts requests to show details about a specific
// action
func (x *Rest) ActionShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionAction
	request.Action = msg.ActionShow
	request.ActionObj = proto.Action{
		ID:        params.ByName(`actionID`),
		SectionID: params.ByName(`sectionID`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ActionSearch accepts requests to look up actions by name and
// sectionId
func (x *Rest) ActionSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionAction
	request.Action = msg.ActionSearch

	cReq := proto.NewActionRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	if cReq.Filter.Action.SectionID == `` || cReq.Filter.Action.Name == `` {
		x.replyBadRequest(&w, &request,
			fmt.Errorf(`Invalid action search specification`))
		return
	}
	request.Search.ActionObj.Name = cReq.Filter.Action.Name
	request.Search.ActionObj.SectionID = cReq.Filter.Action.SectionID

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ActionAdd accepts requests to add a new action to a section
func (x *Rest) ActionAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionAction
	request.Action = msg.ActionAdd

	cReq := proto.NewActionRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	if cReq.Action.SectionID != params.ByName(`sectionID`) {
		x.replyBadRequest(&w, &request, fmt.Errorf("SectionId mismatch: %s, %s",
			cReq.Action.SectionID, params.ByName(`sectionID`)))
		return
	}
	request.ActionObj = proto.Action{
		Name:      cReq.Action.Name,
		SectionID: cReq.Action.SectionID,
		Category:  params.ByName(`category`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ActionRemove accepts requests to remove an action form a section
func (x *Rest) ActionRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionAction
	request.Action = msg.ActionRemove
	request.ActionObj = proto.Action{
		ID:        params.ByName(`actionID`),
		SectionID: params.ByName(`sectionID`),
		Category:  params.ByName(`category`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
