/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
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

// NodeMgmtAdd function
func (x *Rest) NodeMgmtAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	var serverID string
	if cReq.Node.ServerID != `` {
		serverID = cReq.Node.ServerID
	} else {
		serverID = `00000000-0000-0000-0000-000000000000`
	}

	request := newRequest(r, params)
	request.Section = msg.SectionNodeMgmt
	request.Action = msg.ActionAdd
	request.Node = proto.Node{
		AssetID:   cReq.Node.AssetID,
		Name:      cReq.Node.Name,
		TeamID:    cReq.Node.TeamID,
		ServerID:  serverID,
		State:     `unassigned`,
		IsOnline:  cReq.Node.IsOnline,
		IsDeleted: false,
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// NodeMgmtSync function
func (x *Rest) NodeMgmtSync(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionNodeMgmt
	request.Action = msg.ActionSync

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// NodeMgmtUpdate function
func (x *Rest) NodeMgmtUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionNodeMgmt
	request.Action = msg.ActionUpdate
	request.Node = proto.Node{
		ID:        cReq.Node.ID,
		AssetID:   cReq.Node.AssetID,
		Name:      cReq.Node.Name,
		TeamID:    cReq.Node.TeamID,
		ServerID:  cReq.Node.ServerID,
		IsOnline:  cReq.Node.IsOnline,
		IsDeleted: cReq.Node.IsDeleted,
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// NodeMgmtRemove function
func (x *Rest) NodeMgmtRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	action := msg.ActionRemove
	if cReq.Flags.Purge {
		action = msg.ActionPurge
	}

	request := newRequest(r, params)
	request.Section = msg.SectionNodeMgmt
	request.Action = action
	request.Node = proto.Node{
		ID: params.ByName(`nodeID`),
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
