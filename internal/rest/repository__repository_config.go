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

// RepositoryConfigList function
func (x *Rest) RepositoryConfigList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepositoryConfig
	request.Action = msg.ActionList

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// RepositoryConfigSearch function
func (x *Rest) RepositoryConfigSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepositoryConfig
	request.Action = msg.ActionSearch

	cReq := proto.NewRepositoryFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	switch {
	case cReq.Filter.Repository.ID != ``:
	case cReq.Filter.Repository.Name != ``:
	case cReq.Filter.Repository.TeamID != ``:
	case cReq.Filter.Repository.FilterOnIsDeleted:
	case cReq.Filter.Repository.FilterOnIsActive:
	default:
		x.replyBadRequest(&w, &request, fmt.Errorf(`RepositorySearch request without condition`))
		return
	}
	request.Search.Repository.ID = cReq.Filter.Repository.ID
	request.Search.Repository.Name = cReq.Filter.Repository.Name
	request.Search.Repository.TeamID = cReq.Filter.Repository.TeamID
	request.Search.Repository.IsDeleted = cReq.Filter.Repository.IsDeleted
	request.Search.Repository.IsActive = cReq.Filter.Repository.IsActive
	request.Search.Repository.FilterOnIsDeleted = cReq.Filter.Repository.FilterOnIsDeleted
	request.Search.Repository.FilterOnIsActive = cReq.Filter.Repository.FilterOnIsActive

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// RepositoryConfigShow function
func (x *Rest) RepositoryConfigShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepositoryConfig
	request.Action = msg.ActionShow
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Repository.TeamID = params.ByName(`teamID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// RepositoryConfigTree function
func (x *Rest) RepositoryConfigTree(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepositoryConfig
	request.Action = msg.ActionTree
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Tree = proto.Tree{
		ID:   params.ByName(`repositoryID`),
		Type: msg.EntityRepository,
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// RepositoryConfigPropertyCreate function
func (x *Rest) RepositoryConfigPropertyCreate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepositoryConfig
	request.Action = msg.ActionPropertyCreate

	cReq := proto.NewRepositoryRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	switch {
	case params.ByName(`repositoryID`) != cReq.Repository.ID:
		x.replyBadRequest(&w, &request, fmt.Errorf("Mismatched repository ids: %s, %s",
			params.ByName(`repositoryID`), cReq.Repository.ID))
		return
	case len(*cReq.Repository.Properties) != 1:
		x.replyBadRequest(&w, &request, fmt.Errorf("Expected property count 1, actual count: %d",
			len(*cReq.Repository.Properties)))
		return
	}
	switch (*cReq.Repository.Properties)[0].Type {
	case `service`:
		if (*cReq.Repository.Properties)[0].Service.Name == `` {
			x.replyBadRequest(&w, &request, fmt.Errorf(`Invalid service name: empty string`))
			return
		}
	}
	request.Repository = cReq.Repository.Clone()
	request.TargetEntity = msg.EntityRepository
	request.Property.Type = (*cReq.Repository.Properties)[0].Type

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// RepositoryConfigPropertyDestroy function
func (x *Rest) RepositoryConfigPropertyDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepositoryConfig
	request.Action = msg.ActionPropertyDestroy
	request.TargetEntity = msg.EntityRepository
	request.Property.Type = params.ByName(`propertyType`)
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Repository.Properties = &[]proto.Property{
		proto.Property{
			Type:             params.ByName(`propertyType`),
			RepositoryID:     params.ByName(`repositoryID`),
			SourceInstanceID: params.ByName(`sourceID`),
		},
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// RepositoryConfigPropertyUpdate function
func (x *Rest) RepositoryConfigPropertyUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepositoryConfig
	request.Action = msg.ActionPropertyUpdate

	cReq := proto.NewRepositoryRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	switch {
	case params.ByName(`repositoryID`) != cReq.Repository.ID:
		x.replyBadRequest(&w, &request, fmt.Errorf("Mismatched repository ids: %s, %s",
			params.ByName(`repositoryID`), cReq.Repository.ID))
		return
	case len(*cReq.Repository.Properties) != 1:
		x.replyBadRequest(&w, &request, fmt.Errorf("Expected property count 1, actual count: %d",
			len(*cReq.Repository.Properties)))
		return
	}
	switch (*cReq.Repository.Properties)[0].Type {
	case `service`:
		if (*cReq.Repository.Properties)[0].Service.Name == `` {
			x.replyBadRequest(&w, &request, fmt.Errorf(`Invalid service name: empty string`))
			return
		}
	}
	request.Repository = cReq.Repository.Clone()
	request.TargetEntity = msg.EntityRepository
	request.Property.Type = (*cReq.Repository.Properties)[0].Type

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
