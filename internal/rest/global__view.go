/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2015-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// ViewList function
func (x *Rest) ViewList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionView
	request.Action = msg.ActionList

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`view_r`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// ViewShow function
func (x *Rest) ViewShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionView
	request.Action = msg.ActionShow
	request.View = proto.View{
		Name: params.ByName(`view`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`view_r`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// ViewAdd function
func (x *Rest) ViewAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewViewRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if strings.Contains(cReq.View.Name, `.`) {
		dispatchBadRequest(&w, fmt.Errorf(`Invalid view name containing . character`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionAttribute
	request.Action = msg.ActionAdd
	request.View = proto.View{
		Name: cReq.View.Name,
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`view_w`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// ViewRemove function
func (x *Rest) ViewRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionView
	request.Action = msg.ActionRemove
	request.View = proto.View{
		Name: params.ByName(`view`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`view_w`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// ViewRename function
func (x *Rest) ViewRename(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewViewRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if strings.Contains(cReq.View.Name, `.`) {
		dispatchBadRequest(&w, fmt.Errorf(
			`Invalid view name containing . character`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionAttribute
	request.Action = msg.ActionAdd
	request.View = proto.View{
		Name: params.ByName(`view`),
	}
	request.Update.View = proto.View{
		Name: cReq.View.Name,
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`view_w`)
	handler.Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
