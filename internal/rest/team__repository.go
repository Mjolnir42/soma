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

// RepositoryDestroy function
func (x *Rest) RepositoryDestroy(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepository
	request.Action = msg.ActionDestroy
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

// RepositoryAudit function
func (x *Rest) RepositoryAudit(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepository
	request.Action = msg.ActionAudit
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

// RepositoryRename function
func (x *Rest) RepositoryRename(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepository
	request.Action = msg.ActionRename

	cReq := proto.NewRepositoryRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	nameLen := utf8.RuneCountInString(cReq.Repository.Name)
	if nameLen < 4 || nameLen > 128 {
		x.replyBadRequest(&w, &request, fmt.Errorf(`Illegal new repository name length (4 < x <= 128)`))
		return
	}
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Repository.TeamID = params.ByName(`teamID`)
	request.Update.Repository.Name = cReq.Repository.Name

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// RepositoryRepossess function
func (x *Rest) RepositoryRepossess(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepository
	request.Action = msg.ActionRepossess

	cReq := proto.NewRepositoryRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}
	request.Repository.ID = params.ByName(`repositoryID`)
	request.Repository.TeamID = params.ByName(`teamID`)
	request.Update.Repository.TeamID = cReq.Repository.TeamID

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// RepositorySearch function
func (x *Rest) RepositorySearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepository
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

// RepositoryShow function
func (x *Rest) RepositoryShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionRepository
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
