/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2018, 1&1 IONOS SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
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

// OncallList function
func (x *Rest) OncallList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionOncall
	request.Action = msg.ActionList

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// OncallShow function
func (x *Rest) OncallShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionOncall
	request.Action = msg.ActionShow
	request.Oncall.ID = params.ByName(`oncallID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// OncallSearch function
func (x *Rest) OncallSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewOncallFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionOncall
	request.Action = msg.ActionSearch
	request.Search.Oncall.Name = cReq.Filter.Oncall.Name

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// OncallAdd function
func (x *Rest) OncallAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewOncallRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionOncall
	request.Action = msg.ActionAdd
	request.Oncall = cReq.Oncall.Clone()
	request.Oncall.Sanitize()

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// OncallUpdate function
func (x *Rest) OncallUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewOncallRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionOncall
	request.Action = msg.ActionUpdate
	request.Update.Oncall = cReq.Oncall.Clone()
	request.Update.Oncall.Sanitize()
	request.Oncall.ID = params.ByName(`oncallID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// OncallRemove function
func (x *Rest) OncallRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionOncall
	request.Action = msg.ActionRemove
	request.Oncall.ID = params.ByName(`oncallID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// OncallMemberAssign function
func (x *Rest) OncallMemberAssign(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewOncallRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Oncall.Members == nil || *cReq.Oncall.Members == nil || len(*cReq.Oncall.Members) != 1 {
		dispatchBadRequest(&w, nil)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionOncall
	request.Action = msg.ActionMemberAssign
	request.Oncall = cReq.Oncall.Clone()
	request.Oncall.ID = params.ByName(`oncallID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// OncallMemberUnassign function
func (x *Rest) OncallMemberUnassign(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionOncall
	request.Action = msg.ActionMemberUnassign
	request.Oncall.ID = params.ByName(`oncallID`)
	request.Oncall.Members = &[]proto.OncallMember{
		proto.OncallMember{UserID: params.ByName(`userID`)},
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// OncallMemberList function
func (x *Rest) OncallMemberList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionOncall
	request.Action = msg.ActionMemberList
	request.Oncall.ID = params.ByName(`oncallID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
