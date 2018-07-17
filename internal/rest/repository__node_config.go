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

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	// XXX check params.ByName(`nodeID`) == cReq.Node.ID

	request := newRequest(r, params)
	request.Node.ID = cReq.Node.ID
	request.Node.Config.RepositoryID = cReq.Node.Config.RepositoryID
	request.Node.Config.BucketID = cReq.Node.Config.BucketID
	request.Repository.ID = cReq.Node.Config.RepositoryID
	request.Bucket.ID = cReq.Node.Config.BucketID

	// check if the user is allowed to assign nodes from this team
	request.Section = msg.SectionNode
	request.Action = msg.ActionAssign
	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	// check if the user is allowed to assign nodes to the target repo
	request.Section = msg.SectionNodeConfig
	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// NodeConfigUnassign function
func (x *Rest) NodeConfigUnassign(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionNodeConfig
	request.Action = msg.ActionUnassign
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Node.ID = params.ByName(`nodeID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// NodeConfigPropertyCreate function
func (x *Rest) NodeConfigPropertyCreate(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	switch {
	case params.ByName(`nodeID`) != cReq.Node.ID:
		dispatchBadRequest(&w, fmt.Errorf(
			"Mismatched node ids: %s, %s",
			params.ByName(`nodeID`),
			cReq.Node.ID))
		return
	case len(*cReq.Node.Properties) != 1:
		dispatchBadRequest(&w, fmt.Errorf(
			"Expected property count 1, actual count: %d",
			len(*cReq.Node.Properties)))
		return
	case params.ByName(`propertyType`) != (*cReq.Node.Properties)[0].Type:
		dispatchBadRequest(&w, fmt.Errorf(
			"Mismatched property types: %s, %s",
			params.ByName(`propertyType`),
			(*cReq.Node.Properties)[0].Type))
		return
	case (params.ByName(`propertyType`) == "service"):
		if (*cReq.Node.Properties)[0].Service.Name == `` {
			dispatchBadRequest(&w, fmt.Errorf(
				"Empty service name is invalid"))
			return
		}
	}

	request := newRequest(r, params)
	request.Section = msg.SectionNodeConfig
	request.Action = msg.ActionPropertyCreate
	request.Node = cReq.Node.Clone()
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Node.ID = params.ByName(`nodeID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// NodeConfigPropertyDestroy function
func (x *Rest) NodeConfigPropertyDestroy(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewNodeRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	switch {
	case params.ByName(`nodeID`) != cReq.Node.ID:
		dispatchBadRequest(&w, fmt.Errorf(
			"Mismatched node ids: %s, %s",
			params.ByName(`nodeID`),
			cReq.Node.ID))
		return
	case cReq.Node.Config == nil:
		dispatchBadRequest(&w, fmt.Errorf(
			`Node configuration data missing`))
		return
	}
	// outside switch: _after_ nil test
	if cReq.Node.Config.RepositoryID == `` ||
		cReq.Node.Config.BucketID == `` {
		dispatchBadRequest(&w, fmt.Errorf(
			`Node configuration data incomplete`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionNodeConfig
	request.Action = msg.ActionPropertyDestroy
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
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
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
