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

	request := msg.New(r, params)
	request.Section = msg.SectionNodeMgmt
	request.Action = msg.ActionAdd

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	var serverID string
	if cReq.Node.ServerID != `` {
		serverID = cReq.Node.ServerID
	} else {
		serverID = `00000000-0000-0000-0000-000000000000`
	}
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
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// NodeMgmtSync function
func (x *Rest) NodeMgmtSync(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionNodeMgmt
	request.Action = msg.ActionSync

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// NodeMgmtUpdate function
func (x *Rest) NodeMgmtUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionNodeMgmt
	request.Action = msg.ActionUpdate

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	request.Node = proto.Node{
		ID: cReq.Node.ID,
	}
	request.Update.Node = proto.Node{
		AssetID:   cReq.Node.AssetID,
		Name:      cReq.Node.Name,
		TeamID:    cReq.Node.TeamID,
		ServerID:  cReq.Node.ServerID,
		IsOnline:  cReq.Node.IsOnline,
		IsDeleted: cReq.Node.IsDeleted,
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// NodeMgmtRemove function
func (x *Rest) NodeMgmtRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionNodeMgmt
	action := msg.ActionRemove

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	if cReq.Flags.Purge {
		action = msg.ActionPurge
		switch params.ByName(`nodeID`) {
		case ``:
			request.Flag.Unscoped = true
		default:
			request.Node.ID = params.ByName(`nodeID`)
		}
	} else {
		request.Node.ID = params.ByName(`nodeID`)
	}
	request.Action = action

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
