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

	request := newRequest(r, params)
	request.Section = msg.SectionAction
	request.Action = msg.ActionList
	request.ActionObj = proto.Action{
		SectionID: params.ByName(`section`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`supervisor`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// ActionShow accepts requests to show details about a specific
// action
func (x *Rest) ActionShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionAction
	request.Action = msg.ActionShow
	request.ActionObj = proto.Action{
		ID:        params.ByName(`action`),
		SectionID: params.ByName(`section`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`supervisor`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// ActionSearch accepts requests to look up actions by name and
// sectionId
func (x *Rest) ActionSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewActionRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Action.SectionID == `` || cReq.Action.Name == `` {
		dispatchBadRequest(&w,
			fmt.Errorf(`Invalid action search specification`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionAction
	request.Action = msg.ActionSearch
	request.ActionObj = proto.Action{
		Name:      cReq.Action.Name,
		SectionID: cReq.Action.SectionID,
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`supervisor`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// ActionAdd accepts requests to add a new action to a section
func (x *Rest) ActionAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewActionRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Action.SectionID != params.ByName(`section`) {
		dispatchBadRequest(&w, fmt.Errorf("SectionId mismatch: %s, %s",
			cReq.Action.SectionID, params.ByName(`section`)))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionAction
	request.Action = msg.ActionAdd
	request.ActionObj = proto.Action{
		Name:      cReq.Action.Name,
		SectionID: cReq.Action.SectionID,
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`supervisor`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// ActionRemove accepts requests to remove an action form a section
func (x *Rest) ActionRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionAction
	request.Action = msg.ActionRemove
	request.ActionObj = proto.Action{
		ID:        params.ByName(`action`),
		SectionID: params.ByName(`section`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`supervisor`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
