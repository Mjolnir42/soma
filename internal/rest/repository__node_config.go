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

// NodeConfigAssign function
func (x *Rest) NodeConfigAssign(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionNode
	request.Action = msg.ActionShow

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	// XXX check params.ByName(`nodeID`) == cReq.Node.ID
	request.Node.ID = cReq.Node.ID
	request.Node.Config = &proto.NodeConfig{
		RepositoryID: cReq.Node.Config.RepositoryID,
		BucketID:     cReq.Node.Config.BucketID,
	}
	request.Repository.ID = cReq.Node.Config.RepositoryID
	request.Bucket.ID = cReq.Node.Config.BucketID

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	// if there is no result we do not need to authorize it
	if len(result.Node) == 0 {
		x.send(&w, &result)
		return
	}
	// set the TeamID based of the result before we authorize the request
	// this allows authorization based on trusted information
	request.Node.TeamID = result.Node[0].TeamID
	// check if the user is allowed to assign nodes from this team
	request.Action = msg.ActionAssign
	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	// check if the user is allowed to assign nodes to the target repo
	request.Section = msg.SectionNodeConfig
	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result = <-request.Reply
	x.send(&w, &result)
}

// NodeConfigUnassign function
func (x *Rest) NodeConfigUnassign(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Node.ID = params.ByName(`nodeID`)
	request.Node.Config = &proto.NodeConfig{
		RepositoryID: params.ByName(`repositoryID`),
		BucketID:     params.ByName(`bucketID`),
	}
	request.Section = msg.SectionNode
	request.Action = msg.ActionShow
	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	// if there is no result we do not need to authorize it
	if len(result.Node) == 0 {
		x.send(&w, &result)
		return
	}
	// set the TeamID based of the result before we authorize the request
	// this allows authorization based on trusted information
	request.Node.TeamID = result.Node[0].TeamID
	// check if the user is allowed to unassign nodes from this team
	request.Action = msg.ActionUnassign
	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	// check if the user is allowed to unassign nodes from the target
	// repo
	request.Section = msg.SectionNodeConfig
	request.Action = msg.ActionUnassign
	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result = <-request.Reply
	x.send(&w, &result)
}

// NodeConfigPropertyCreate function
func (x *Rest) NodeConfigPropertyCreate(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionNodeConfig
	request.Action = msg.ActionPropertyCreate

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	switch {
	case params.ByName(`nodeID`) != cReq.Node.ID:
		x.replyBadRequest(&w, &request, fmt.Errorf(
			"Mismatched node ids: %s, %s",
			params.ByName(`nodeID`),
			cReq.Node.ID))
		return
	case len(*cReq.Node.Properties) != 1:
		x.replyBadRequest(&w, &request, fmt.Errorf(
			"Expected property count 1, actual count: %d",
			len(*cReq.Node.Properties)))
		return
	case (*cReq.Node.Properties)[0].Type == "service":
		if (*cReq.Node.Properties)[0].Service.Name == `` {
			x.replyBadRequest(&w, &request, fmt.Errorf(
				"Empty service name is invalid"))
			return
		}
	}
	request.TargetEntity = msg.EntityNode
	request.Node = cReq.Node.Clone()
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Node.ID = params.ByName(`nodeID`)
	request.Property.Type = params.ByName(`propertyType`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// NodeConfigPropertyDestroy function
func (x *Rest) NodeConfigPropertyDestroy(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionNodeConfig
	request.Action = msg.ActionPropertyDestroy

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	switch {
	case params.ByName(`nodeID`) != cReq.Node.ID:
		x.replyBadRequest(&w, &request, fmt.Errorf(
			"Mismatched node ids: %s, %s",
			params.ByName(`nodeID`),
			cReq.Node.ID))
		return
	case cReq.Node.Config == nil:
		x.replyBadRequest(&w, &request, fmt.Errorf(
			`Node configuration data missing`))
		return
	}
	// outside switch: _after_ nil test
	if cReq.Node.Config.RepositoryID == `` ||
		cReq.Node.Config.BucketID == `` {
		x.replyBadRequest(&w, &request, fmt.Errorf(
			`Node configuration data incomplete`))
		return
	}
	request.TargetEntity = msg.EntityNode
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Property.Type = params.ByName(`propertyType`)
	request.Node.ID = params.ByName(`nodeID`)
	request.Node.Config = &proto.NodeConfig{
		RepositoryID: request.Repository.ID,
		BucketID:     request.Bucket.ID,
	}
	request.Node.Properties = &[]proto.Property{proto.Property{
		Type:             params.ByName(`propertyType`),
		RepositoryID:     request.Repository.ID,
		BucketID:         request.Bucket.ID,
		SourceInstanceID: params.ByName(`sourceID`),
	}}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// NodeConfigPropertyUpdate function
func (x *Rest) NodeConfigPropertyUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	// XXX BUG TODO
}

// NodeConfigTree function
func (x *Rest) NodeConfigTree(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionNodeConfig
	request.Action = msg.ActionTree
	request.Tree = proto.Tree{
		ID:   params.ByName(`nodeID`),
		Type: msg.EntityNode,
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
