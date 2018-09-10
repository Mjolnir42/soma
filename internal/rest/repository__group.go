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
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// GroupList function
func (x *Rest) GroupList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionList
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// GroupSearch function
func (x *Rest) GroupSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewGroupFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Filter.Group.Name == `` {
		dispatchBadRequest(&w, nil) // XXX
		return
	}

	if cReq.Filter.Group.BucketID != `` && cReq.Filter.Group.BucketID != params.ByName(`bucketID`) {
		dispatchBadRequest(&w, nil) // XXX
		return
	}

	if cReq.Filter.Group.RepositoryID != `` && cReq.Filter.Group.RepositoryID != params.ByName(`repositoryID`) {
		dispatchBadRequest(&w, nil) // XXX
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionSearch
	request.Search.Group.Name = cReq.Filter.Group.Name
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply

	// XXX BUG filter in SQL statement
	filtered := []proto.Group{}
	for _, i := range result.Group {
		if i.Name == cReq.Filter.Group.Name && cReq.Filter.Group.BucketID == params.ByName(`bucketID`) {
			filtered = append(filtered, i)
		}
	}
	result.Group = filtered
	send(&w, &result)
}

// GroupShow function
func (x *Rest) GroupShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionShow
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Group.ID = params.ByName(`groupID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// GroupTree function
func (x *Rest) GroupTree(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionTree
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Group.ID = params.ByName(`groupID`)
	request.Tree.ID = params.ByName(`groupID`)
	request.Tree.Type = msg.EntityGroup

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// GroupCreate function
func (x *Rest) GroupCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewGroupRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	nameLen := utf8.RuneCountInString(cReq.Group.Name)
	if nameLen < 4 || nameLen > 256 {
		dispatchBadRequest(&w, fmt.Errorf(`Illegal group name length (4 <= x <= 256)`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionCreate
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Group = cReq.Group.Clone()

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// GroupDestroy function
func (x *Rest) GroupDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionDestroy
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Group.ID = params.ByName(`groupID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// GroupMemberList function
func (x *Rest) GroupMemberList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionMemberList
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Group.ID = params.ByName(`groupID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// GroupMemberAssign function
func (x *Rest) GroupMemberAssign(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewGroupRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionMemberAssign
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Group = cReq.Group.Clone()
	request.Group.ID = params.ByName(`groupID`)

	switch params.ByName(`memberType`) {
	case msg.EntityGroup:
		request.TargetEntity = msg.EntityGroup
		cReq.Group.MemberClusters = nil
		cReq.Group.MemberNodes = nil
		if cReq.Group.MemberGroups == nil || len(*cReq.Group.MemberGroups) != 1 {
			dispatchBadRequest(&w, nil)
			return
		}
	case msg.EntityCluster:
		request.TargetEntity = msg.EntityCluster
		cReq.Group.MemberGroups = nil
		cReq.Group.MemberNodes = nil
		if cReq.Group.MemberClusters == nil || len(*cReq.Group.MemberClusters) != 1 {
			dispatchBadRequest(&w, nil)
			return
		}
	case msg.EntityNode:
		request.TargetEntity = msg.EntityNode
		cReq.Group.MemberGroups = nil
		cReq.Group.MemberClusters = nil
		if cReq.Group.MemberNodes == nil || len(*cReq.Group.MemberNodes) != 1 {
			dispatchBadRequest(&w, nil)
			return
		}
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// GroupMemberUnassign function
func (x *Rest) GroupMemberUnassign(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionMemberUnassign
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Group.ID = params.ByName(`groupID`)

	switch params.ByName(`memberType`) {
	case msg.EntityGroup:
		request.TargetEntity = msg.EntityGroup
		request.Group.MemberGroups = &[]proto.Group{
			proto.Group{ID: params.ByName(`memberID`)},
		}
	case msg.EntityCluster:
		request.TargetEntity = msg.EntityCluster
		request.Group.MemberClusters = &[]proto.Cluster{
			proto.Cluster{ID: params.ByName(`memberID`)},
		}
	case msg.EntityNode:
		request.TargetEntity = msg.EntityNode
		request.Group.MemberNodes = &[]proto.Node{
			proto.Node{ID: params.ByName(`memberID`)},
		}
	default:
		dispatchBadRequest(&w, nil)
		return
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// GroupPropertyCreate function
func (x *Rest) GroupPropertyCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewGroupRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	switch {
	case params.ByName(`groupID`) != cReq.Group.ID:
		dispatchBadRequest(&w, fmt.Errorf(
			"Mismatched group ids: %s, %s",
			params.ByName(`groupID`),
			cReq.Group.ID))
		return
	case len(*cReq.Group.Properties) != 1:
		dispatchBadRequest(&w, fmt.Errorf(
			"Expected property count 1, actual count: %d",
			len(*cReq.Group.Properties)))
		return
	case (*cReq.Group.Properties)[0].Type == `service` && (*cReq.Group.Properties)[0].Service.Name == ``:
		dispatchBadRequest(&w, fmt.Errorf(
			`Empty service name is invalid`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionPropertyCreate
	request.Group = cReq.Group.Clone()
	request.Property.Type = params.ByName(`propertyType`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// GroupPropertyDestroy function
func (x *Rest) GroupPropertyDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewGroupRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	switch {
	case params.ByName(`groupID`) != cReq.Group.ID:
		dispatchBadRequest(&w, fmt.Errorf(
			"Mismatched group ids: %s, %s",
			params.ByName(`groupID`),
			cReq.Group.ID))
		return
	case cReq.Group.BucketID == ``:
		dispatchBadRequest(&w, fmt.Errorf(
			`Missing bucketId in group delete request`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionPropertyDestroy
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Property.Type = params.ByName(`propertyType`)
	request.Group.ID = params.ByName(`groupID`)
	request.Group.Properties = &[]proto.Property{proto.Property{
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
	send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
