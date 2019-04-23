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

	request := msg.New(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionList
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// GroupSearch function
func (x *Rest) GroupSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionSearch

	cReq := proto.NewGroupFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	if cReq.Filter.Group.Name == `` {
		x.replyBadRequest(&w, &request, nil) // XXX
		return
	}

	if cReq.Filter.Group.BucketID != `` && cReq.Filter.Group.BucketID != params.ByName(`bucketID`) {
		x.replyBadRequest(&w, &request, nil) // XXX
		return
	}

	if cReq.Filter.Group.RepositoryID != `` && cReq.Filter.Group.RepositoryID != params.ByName(`repositoryID`) {
		x.replyBadRequest(&w, &request, nil) // XXX
		return
	}
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
	x.send(&w, &result)
}

// GroupShow function
func (x *Rest) GroupShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionShow
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Group.ID = params.ByName(`groupID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// GroupTree function
func (x *Rest) GroupTree(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionTree
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Group.ID = params.ByName(`groupID`)
	request.Tree.ID = params.ByName(`groupID`)
	request.Tree.Type = msg.EntityGroup

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// GroupCreate function
func (x *Rest) GroupCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionCreate

	cReq := proto.NewGroupRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	nameLen := utf8.RuneCountInString(cReq.Group.Name)
	if nameLen < 2 || nameLen > 256 {
		x.replyBadRequest(&w, &request, fmt.Errorf(`Illegal group name length (2 <= x <= 256)`))
		return
	}
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Group = cReq.Group.Clone()

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// GroupDestroy function
func (x *Rest) GroupDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionDestroy
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Group.ID = params.ByName(`groupID`)
	request.Group.RepositoryID = params.ByName(`repositoryID`)
	request.Group.BucketID = params.ByName(`bucketID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// GroupMemberList function
func (x *Rest) GroupMemberList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionMemberList
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Group.ID = params.ByName(`groupID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// GroupMemberAssign function
func (x *Rest) GroupMemberAssign(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionMemberAssign

	cReq := proto.NewGroupRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}
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
			x.replyBadRequest(&w, &request, nil)
			return
		}
	case msg.EntityCluster:
		request.TargetEntity = msg.EntityCluster
		cReq.Group.MemberGroups = nil
		cReq.Group.MemberNodes = nil
		if cReq.Group.MemberClusters == nil || len(*cReq.Group.MemberClusters) != 1 {
			x.replyBadRequest(&w, &request, nil)
			return
		}
	case msg.EntityNode:
		request.TargetEntity = msg.EntityNode
		cReq.Group.MemberGroups = nil
		cReq.Group.MemberClusters = nil
		if cReq.Group.MemberNodes == nil || len(*cReq.Group.MemberNodes) != 1 {
			x.replyBadRequest(&w, &request, nil)
			return
		}
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// GroupMemberUnassign function
func (x *Rest) GroupMemberUnassign(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
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
		x.replyBadRequest(&w, &request, nil)
		return
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// GroupPropertyCreate function
func (x *Rest) GroupPropertyCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionPropertyCreate

	cReq := proto.NewGroupRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	switch {
	case params.ByName(`groupID`) != cReq.Group.ID:
		x.replyBadRequest(&w, &request, fmt.Errorf(
			"Mismatched group ids: %s, %s",
			params.ByName(`groupID`),
			cReq.Group.ID))
		return
	case len(*cReq.Group.Properties) != 1:
		x.replyBadRequest(&w, &request, fmt.Errorf(
			"Expected property count 1, actual count: %d",
			len(*cReq.Group.Properties)))
		return
	case (*cReq.Group.Properties)[0].Type == `service` && (*cReq.Group.Properties)[0].Service.Name == ``:
		x.replyBadRequest(&w, &request, fmt.Errorf(
			`Empty service name is invalid`))
		return
	}
	request.TargetEntity = msg.EntityGroup
	request.Group = cReq.Group.Clone()
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Property.Type = (*cReq.Group.Properties)[0].Type

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// GroupPropertyDestroy function
func (x *Rest) GroupPropertyDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionGroup
	request.Action = msg.ActionPropertyDestroy

	request.TargetEntity = msg.EntityGroup
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Property.Type = params.ByName(`propertyType`)
	request.Group.ID = params.ByName(`groupID`)
	request.Group.BucketID = params.ByName(`bucketID`)
	request.Group.RepositoryID = params.ByName(`repositoryID`)
	request.Group.Properties = &[]proto.Property{proto.Property{
		Type:             params.ByName(`propertyType`),
		RepositoryID:     params.ByName(`repositoryID`),
		BucketID:         params.ByName(`bucketID`),
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

// GroupPropertyUpdate function
func (x *Rest) GroupPropertyUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	// XXX BUG TODO
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
