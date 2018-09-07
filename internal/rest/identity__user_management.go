/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// UserMgmtList function
func (x *Rest) UserMgmtList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionUserMgmt
	request.Action = msg.ActionList
	request.Flag.Unscoped = true

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// UserMgmtShow function
func (x *Rest) UserMgmtShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionUserMgmt
	request.Action = msg.ActionShow
	request.User.ID = params.ByName(`userID`)
	request.Flag.Unscoped = true

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// UserMgmtSync function
func (x *Rest) UserMgmtSync(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := newRequest(r, params)
	request.Section = msg.SectionUserMgmt
	request.Action = msg.ActionSync
	request.Flag.Unscoped = true

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// UserMgmtSearch function
func (x *Rest) UserMgmtSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewUserFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionUserMgmt
	request.Action = msg.ActionSearch
	request.Search.User.UserName = cReq.Filter.User.UserName
	request.Flag.Unscoped = true

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// UserMgmtAdd function
func (x *Rest) UserMgmtAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewUserRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if strings.Contains(cReq.User.UserName, `:`) {
		dispatchBadRequest(&w, fmt.Errorf(
			`Invalid username containing : character`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionUserMgmt
	request.Action = msg.ActionAdd
	request.User.UserName = cReq.User.UserName
	request.User.FirstName = cReq.User.FirstName
	request.User.LastName = cReq.User.LastName
	request.User.EmployeeNumber = cReq.User.EmployeeNumber
	request.User.MailAddress = cReq.User.MailAddress
	request.User.IsActive = false
	request.User.IsSystem = cReq.User.IsSystem
	request.User.IsDeleted = false
	request.User.TeamID = cReq.User.TeamID

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// UserMgmtUpdate function
func (x *Rest) UserMgmtUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewUserRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if strings.Contains(cReq.User.UserName, `:`) {
		dispatchBadRequest(&w, fmt.Errorf(
			`Invalid username containing : character`))
		return
	}
	if params.ByName(`userID`) != cReq.User.ID {
		dispatchBadRequest(&w, fmt.Errorf(
			`Mismatching user UUIDs in body and URL`))
		return
	}

	request := newRequest(r, params)
	request.Section = msg.SectionUserMgmt
	request.Action = msg.ActionUpdate
	request.User.ID = cReq.User.ID
	request.Update.User.UserName = cReq.User.UserName
	request.Update.User.FirstName = cReq.User.FirstName
	request.Update.User.LastName = cReq.User.LastName
	request.Update.User.EmployeeNumber = cReq.User.EmployeeNumber
	request.Update.User.MailAddress = cReq.User.MailAddress
	request.Update.User.IsDeleted = cReq.User.IsDeleted
	request.Update.User.TeamID = cReq.User.TeamID

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// UserMgmtRemove function
func (x *Rest) UserMgmtRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewUserRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	request := newRequest(r, params)
	request.Section = msg.SectionUserMgmt
	switch cReq.Flags.Purge {
	case true:
		request.Action = msg.ActionPurge
	case false:
		request.Action = msg.ActionRemove
	}
	request.User.ID = params.ByName(`userID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
