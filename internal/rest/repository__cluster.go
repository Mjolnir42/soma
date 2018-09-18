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

// ClusterList function
func (x *Rest) ClusterList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionList
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ClusterSearch function
func (x *Rest) ClusterSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewClusterFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Filter.Cluster.Name == `` {
		dispatchBadRequest(&w, fmt.Errorf(`ClusterSearch on empty name`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionSearch
	request.Search.Cluster.Name = cReq.Filter.Cluster.Name
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply

	// XXX BUG filter in SQL statement
	filtered := []proto.Cluster{}
	for _, i := range result.Cluster {
		if i.Name == cReq.Filter.Cluster.Name &&
			i.BucketID == params.ByName(`bucketID`) {
			filtered = append(filtered, i)
		}
	}
	result.Cluster = filtered
	x.send(&w, &result)
}

// ClusterShow function
func (x *Rest) ClusterShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionShow
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Cluster.ID = params.ByName(`clusterID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ClusterTree function
func (x *Rest) ClusterTree(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionTree
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Cluster.ID = params.ByName(`clusterID`)
	request.Tree.ID = params.ByName(`clusterID`)
	request.Tree.Type = msg.EntityCluster

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ClusterCreate function
func (x *Rest) ClusterCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewClusterRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	nameLen := utf8.RuneCountInString(cReq.Cluster.Name)
	if nameLen < 4 || nameLen > 256 {
		dispatchBadRequest(&w,
			fmt.Errorf(`Illegal cluster name length (4 <= x <= 256)`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionCreate
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Cluster = cReq.Cluster.Clone()

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ClusterDestroy function
func (x *Rest) ClusterDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionDestroy
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Cluster.ID = params.ByName(`clusterID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ClusterMemberList function
func (x *Rest) ClusterMemberList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionMemberList
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Cluster.ID = params.ByName(`clusterID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ClusterMemberAssign function
func (x *Rest) ClusterMemberAssign(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewClusterRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionMemberAssign
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Cluster = cReq.Cluster.Clone()
	request.Cluster.ID = params.ByName(`clusterID`)
	request.TargetEntity = msg.EntityNode

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ClusterMemberUnassign function
func (x *Rest) ClusterMemberUnassign(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionMemberUnassign
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Cluster.ID = params.ByName(`clusterID`)

	switch params.ByName(`memberType`) {
	case msg.EntityNode:
		request.TargetEntity = msg.EntityNode
		request.Cluster.Members = &[]proto.Node{
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
	x.send(&w, &result)
}

// ClusterPropertyCreate function
func (x *Rest) ClusterPropertyCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewClusterRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	switch {
	case params.ByName(`clusterID`) != cReq.Cluster.ID:
		dispatchBadRequest(&w, fmt.Errorf(
			"Mismatched cluster ids: %s, %s",
			params.ByName(`clusterID`),
			cReq.Cluster.ID,
		))
		return
	case len(*cReq.Cluster.Properties) != 1:
		dispatchBadRequest(&w, fmt.Errorf(
			"Expected property count 1, actual count: %d",
			len(*cReq.Cluster.Properties),
		))
		return
	case (*cReq.Cluster.Properties)[0].Type == `service` && (*cReq.Cluster.Properties)[0].Service.Name == ``:
		dispatchBadRequest(&w, fmt.Errorf(
			`Invalid empty service name`,
		))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionPropertyCreate
	request.TargetEntity = msg.EntityCluster
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Cluster = cReq.Cluster.Clone()
	request.Cluster.ID = params.ByName(`clusterID`)
	request.Property.Type = params.ByName(`propertyType`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ClusterPropertyDestroy function
func (x *Rest) ClusterPropertyDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewClusterRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	switch {
	case params.ByName(`clusterID`) != cReq.Cluster.ID:
		dispatchBadRequest(&w, fmt.Errorf(
			"Mismatched cluster ids: %s, %s",
			params.ByName(`clusterID`),
			cReq.Cluster.ID,
		))
		return
	case cReq.Cluster.BucketID == ``:
		dispatchBadRequest(&w, fmt.Errorf(
			`Missing bucketID in bucket property delete request`,
		))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionPropertyDestroy
	request.TargetEntity = msg.EntityCluster
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Property.Type = params.ByName(`propertyType`)
	request.Cluster = proto.Cluster{
		ID:           params.ByName(`clusterID`),
		RepositoryID: params.ByName(`repositoryID`),
		BucketID:     params.ByName(`bucketID`),
		Properties: &[]proto.Property{
			proto.Property{
				Type:             params.ByName(`propertyType`),
				BucketID:         cReq.Cluster.BucketID,
				SourceInstanceID: params.ByName(`sourceID`),
			},
		},
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// ClusterPropertyUpdate function
func (x *Rest) ClusterPropertyUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	// XXX BUG TODO
}

// ClusterRename function
func (x *Rest) ClusterRename(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewClusterRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	nameLen := utf8.RuneCountInString(cReq.Cluster.Name)
	if nameLen < 4 || nameLen > 256 {
		dispatchBadRequest(&w,
			fmt.Errorf(`Illegal cluster name length (4 <= x <= 256)`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionRename
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Bucket.ID = params.ByName(`bucketID`)
	request.Cluster.ID = params.ByName(`clusterID`)
	request.Cluster.RepositoryID = params.ByName(`repositoryID`)
	request.Cluster.BucketID = params.ByName(`bucketID`)
	request.Update.Cluster.Name = cReq.Cluster.Name

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
