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

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
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

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply

	filtered := []proto.Cluster{}
	for _, i := range result.Cluster {
		if i.Name == cReq.Filter.Cluster.Name &&
			i.BucketID == cReq.Filter.Cluster.BucketID {
			filtered = append(filtered, i)
		}
	}
	result.Cluster = filtered
	sendMsgResult(&w, &result)
}

// ClusterShow function
func (x *Rest) ClusterShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionShow
	request.Cluster = proto.Cluster{
		ID: params.ByName(`clusterID`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// ClusterTree function
func (x *Rest) ClusterTree(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionTree
	request.Tree = proto.Tree{
		ID:   params.ByName(`clusterID`),
		Type: msg.EntityCluster,
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
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
	request.Cluster = cReq.Cluster.Clone()

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
}

// ClusterMemberList function
func (x *Rest) ClusterMemberList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionMemberList
	request.Cluster = proto.Cluster{
		ID: params.ByName(`clusterID`),
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
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
	request.Cluster = cReq.Cluster.Clone()

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
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
	case params.ByName(`propertyType`) != (*cReq.Cluster.Properties)[0].Type:
		dispatchBadRequest(&w, fmt.Errorf(
			"Mismatched property types: %s, %s",
			params.ByName(`propertyType`),
			(*cReq.Cluster.Properties)[0].Type,
		))
		return
	case (params.ByName(`propertyType`) == `service`) && (*cReq.Cluster.Properties)[0].Service.Name == ``:
		dispatchBadRequest(&w, fmt.Errorf(
			`Invalid empty service name`,
		))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionCluster
	request.Action = msg.ActionPropertyCreate
	request.Cluster = cReq.Cluster.Clone()

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	sendMsgResult(&w, &result)
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
	request.Cluster = proto.Cluster{
		ID: params.ByName(`clusterID`),
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
	sendMsgResult(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
