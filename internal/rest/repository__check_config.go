/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"fmt"
	"net/http"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"

	"github.com/julienschmidt/httprouter"
)

// CheckConfigList function
func (x *Rest) CheckConfigList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionCheckConfig
	request.Action = msg.ActionList
	request.CheckConfig = proto.CheckConfig{
		RepositoryID: params.ByName(`repositoryID`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// CheckConfigSearch function
func (x *Rest) CheckConfigSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionCheckConfig
	request.Action = msg.ActionSearch

	cReq := proto.NewCheckConfigRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	if cReq.Filter.CheckConfig.Name == `` {
		x.replyBadRequest(&w, &request, fmt.Errorf(`CheckConfigSearch on empty name`))
		return
	}
	request.CheckConfig = proto.CheckConfig{
		RepositoryID: params.ByName(`repositoryID`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply

	filtered := []proto.CheckConfig{}
	for _, i := range result.CheckConfig {
		if i.Name == cReq.Filter.CheckConfig.Name {
			filtered = append(filtered, i)
		}
	}
	result.CheckConfig = filtered
	x.send(&w, &result)
}

// CheckConfigShow function
func (x *Rest) CheckConfigShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionCheckConfig
	request.Action = msg.ActionShow
	request.CheckConfig = proto.CheckConfig{
		ID:           params.ByName(`checkID`),
		RepositoryID: params.ByName(`repositoryID`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// CheckConfigCreate function
func (x *Rest) CheckConfigCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionMonitoring
	request.Action = msg.ActionUse

	cReq := proto.NewCheckConfigRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}
	request.CheckConfig = cReq.CheckConfig.Clone()

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	request.Section = msg.SectionCheckConfig
	request.Action = msg.ActionCreate

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// CheckConfigDestroy function
func (x *Rest) CheckConfigDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionMonitoring
	request.Action = msg.ActionUse
	request.CheckConfig = proto.CheckConfig{
		ID:           params.ByName(`checkID`),
		RepositoryID: params.ByName(`repositoryID`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	request.Section = msg.SectionCheckConfig
	request.Action = msg.ActionDestroy

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
