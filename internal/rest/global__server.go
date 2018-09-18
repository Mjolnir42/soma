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

// ServerList function
func (x *Rest) ServerList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionList

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ServerShow function
func (x *Rest) ServerShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionShow
	request.Server.ID = params.ByName(`serverID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ServerSearch function
func (x *Rest) ServerSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewServerFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if cReq.Filter.Server.Name == `` && cReq.Filter.Server.AssetID == 0 {
		dispatchBadRequest(&w, fmt.Errorf(`Bad search request with empty `+
			`Server.Name and Server.AssetID`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionSearch
	request.Search.Server.Name = cReq.Filter.Server.Name
	request.Search.Server.AssetID = cReq.Filter.Server.AssetID

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ServerSync function
func (x *Rest) ServerSync(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionSync

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ServerAdd function
func (x *Rest) ServerAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewServerRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionAdd
	request.Server.AssetID = cReq.Server.AssetID
	request.Server.Datacenter = cReq.Server.Datacenter
	request.Server.Location = cReq.Server.Location
	request.Server.Name = cReq.Server.Name
	request.Server.IsOnline = cReq.Server.IsOnline
	request.Server.IsDeleted = false

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ServerRemove function
func (x *Rest) ServerRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewServerRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionServer
	switch cReq.Flags.Purge {
	case true:
		request.Action = msg.ActionPurge
	case false:
		request.Action = msg.ActionRemove
	}
	request.Server.ID = params.ByName(`serverID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ServerUpdate function
func (x *Rest) ServerUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewServerRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if cReq.Server.ID != params.ByName(`serverID`) {
		dispatchBadRequest(&w, fmt.Errorf(`Mismatching server UUIDs`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionUpdate
	request.Server.ID = cReq.Server.ID
	request.Update.Server.AssetID = cReq.Server.AssetID
	request.Update.Server.Datacenter = cReq.Server.Datacenter
	request.Update.Server.Location = cReq.Server.Location
	request.Update.Server.Name = cReq.Server.Name
	request.Update.Server.IsOnline = cReq.Server.IsOnline
	request.Update.Server.IsDeleted = cReq.Server.IsDeleted

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ServerAddNull function
func (x *Rest) ServerAddNull(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewServerRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if cReq.Server.ID != `00000000-0000-0000-0000-000000000000` ||
		params.ByName(`serverID`) != `null` {
		dispatchBadRequest(&w, fmt.Errorf(`not null server`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionServer
	request.Action = msg.ActionInsertNullID
	request.Server.Datacenter = cReq.Server.Datacenter

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
