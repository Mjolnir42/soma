/*-
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

// RightList function
func (x *Rest) RightList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRight
	request.Action = msg.ActionList
	request.Grant.Category = params.ByName(`category`)
	request.Grant.PermissionID = params.ByName(`permissionID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// RightShow function
func (x *Rest) RightShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRight
	request.Action = msg.ActionShow
	request.Grant.Category = params.ByName(`category`)
	request.Grant.ID = params.ByName(`grantID`)
	request.Grant.PermissionID = params.ByName(`permissionID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// RightSearch function
func (x *Rest) RightSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	crq := proto.NewGrantFilter()
	if err := decodeJSONBody(r, &crq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionRight
	request.Action = msg.ActionSearch
	request.Search.Grant.RecipientType = crq.Filter.Grant.RecipientType
	request.Search.Grant.RecipientID = crq.Filter.Grant.RecipientID
	request.Search.Grant.PermissionID = crq.Filter.Grant.PermissionID
	request.Search.Grant.Category = crq.Filter.Grant.Category
	request.Search.Grant.ObjectType = crq.Filter.Grant.ObjectType
	request.Search.Grant.ObjectID = crq.Filter.Grant.ObjectID

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// RightGrant function
func (x *Rest) RightGrant(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewGrantRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Grant.Category != params.ByName(`category`) ||
		cReq.Grant.PermissionID != params.ByName(`permissionID`) {
		dispatchBadRequest(&w, fmt.Errorf(
			`Category/PermissionId mismatch`))
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionRight
	request.Action = msg.ActionGrant
	request.Grant = cReq.Grant.Clone()

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// RightRevoke function
func (x *Rest) RightRevoke(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRight
	request.Action = msg.ActionRevoke
	request.Grant.ID = params.ByName(`grantID`)
	request.Grant.Category = params.ByName(`category`)
	request.Grant.PermissionID = params.ByName(`permissionID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
