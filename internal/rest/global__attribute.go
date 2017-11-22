/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2017, Jörg Pernfuß
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

// AttributeList function
func (x *Rest) AttributeList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionAttribute
	request.Action = msg.ActionList
	request.ActionObj = proto.Action{
		SectionID: params.ByName(`section`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`attribute_r`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// AttributeShow function
func (x *Rest) AttributeShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionAttribute
	request.Action = msg.ActionShow
	request.Attribute = proto.Attribute{
		Name: params.ByName(`attribute`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`attribute_r`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// AttributeAdd function
func (x *Rest) AttributeAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewAttributeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionAttribute
	request.Action = msg.ActionAdd
	request.Attribute = proto.Attribute{
		Name:        cReq.Attribute.Name,
		Cardinality: cReq.Attribute.Cardinality,
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`attribute_w`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// AttributeRemove function
func (x *Rest) AttributeRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionAttribute
	request.Action = msg.ActionRemove
	request.Attribute = proto.Attribute{
		Name: params.ByName(`attribute`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`attribute_w`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
