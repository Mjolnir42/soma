/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest

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
	if cReq.Node.ServerId != `` {
		serverID = cReq.Node.ServerId
	} else {
		serverID = `00000000-0000-0000-0000-000000000000`
	}

	returnChannel := make(chan msg.Result)
	request := msg.Request{
		ID:         requestID(params),
		Section:    msg.SectionNodeMgmt,
		Action:     msg.ActionAdd,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Node: proto.Node{
			AssetId:   cReq.Node.AssetId,
			Name:      cReq.Node.Name,
			TeamId:    cReq.Node.TeamId,
			ServerId:  serverID,
			State:     `unassigned`,
			IsOnline:  cReq.Node.IsOnline,
			IsDeleted: false,
		},
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`node_w`)
	handler.Intake() <- request
	result := <-returnChannel
	sendMsgResult(&w, &result)
}

// NodeMgmtSync function
func (x *Rest) NodeMgmtSync(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	returnChannel := make(chan msg.Result)
	request := msg.Request{
		ID:         requestID(params),
		Section:    msg.SectionNodeMgmt,
		Action:     msg.ActionSync,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`node_r`)
	handler.Intake() <- request
	result := <-returnChannel
	sendMsgResult(&w, &result)
}

// NodeMgmtUpdate function
func (x *Rest) NodeMgmtUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewNodeRequest()
	err := decodeJSONBody(r, &cReq)
	if err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	returnChannel := make(chan msg.Result)
	request := msg.Request{
		ID:         requestID(params),
		Section:    msg.SectionNodeMgmt,
		Action:     msg.ActionUpdate,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Node: proto.Node{
			Id:        cReq.Node.Id,
			AssetId:   cReq.Node.AssetId,
			Name:      cReq.Node.Name,
			TeamId:    cReq.Node.TeamId,
			ServerId:  cReq.Node.ServerId,
			IsOnline:  cReq.Node.IsOnline,
			IsDeleted: cReq.Node.IsDeleted,
		},
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`node_w`)
	handler.Intake() <- request
	result := <-returnChannel
	sendMsgResult(&w, &result)
}

// NodeMgmtRemove function
func (x *Rest) NodeMgmtRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewNodeRequest()
	err := decodeJSONBody(r, &cReq)
	if err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	action := msg.ActionRemove
	if cReq.Flags.Purge {
		action = msg.ActionPurge
	}

	returnChannel := make(chan msg.Result)
	request := msg.Request{
		ID:         requestID(params),
		Section:    msg.SectionNodeMgmt,
		Action:     action,
		Reply:      returnChannel,
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Node: proto.Node{
			Id: params.ByName(`nodeID`),
		},
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	handler := x.handlerMap.Get(`node_w`)
	handler.Intake() <- request
	result := <-returnChannel
	sendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
