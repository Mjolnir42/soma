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

// PermissionList function
func (x *Rest) PermissionList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionPermission
	request.Action = msg.ActionList
	request.Permission.Category = params.ByName(`category`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// PermissionShow function
func (x *Rest) PermissionShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionPermission
	request.Action = msg.ActionShow
	request.Permission.ID = params.ByName(`permissionID`)
	request.Permission.Category = params.ByName(`category`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// PermissionSearch function
func (x *Rest) PermissionSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewPermissionFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Filter.Permission.Name == `` || cReq.Filter.Permission.Category == `` {
		dispatchBadRequest(&w, fmt.Errorf(
			`PermissionSearch request missing permission name or category`))
		return
	}
	request := msg.New(r, params)
	request.Section = msg.SectionPermission
	request.Action = msg.ActionSearch
	request.Search.Permission.Name = cReq.Filter.Permission.Name
	request.Search.Permission.Category = cReq.Filter.Permission.Category

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// PermissionAdd function
func (x *Rest) PermissionAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewPermissionRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Permission.Category != params.ByName(`category`) {
		dispatchBadRequest(&w, fmt.Errorf(`Category mismatch`))
		return
	}
	if strings.Contains(params.ByName(`category`), `:grant`) {
		dispatchBadRequest(&w, fmt.Errorf(
			`Permissions in :grant categories are auto-managed.`))
		return
	}
	if params.ByName(`category`) == msg.CategorySystem ||
		params.ByName(`category`) == msg.CategoryOmnipotence {
		dispatchForbidden(&w, nil)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionPermission
	request.Action = msg.ActionAdd
	request.Permission.Name = cReq.Permission.Name
	request.Permission.Category = cReq.Permission.Category

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// PermissionRemove function
func (x *Rest) PermissionRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	if strings.Contains(params.ByName(`category`), `:grant`) {
		dispatchBadRequest(&w, fmt.Errorf(
			`Permissions in :grant categories are auto-managed.`))
		return
	}
	if params.ByName(`category`) == msg.CategorySystem ||
		params.ByName(`category`) == msg.CategoryOmnipotence {
		dispatchForbidden(&w, nil)
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionPermission
	request.Action = msg.ActionRemove
	request.Permission.ID = params.ByName(`permissionID`)
	request.Permission.Category = params.ByName(`category`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// PermissionEdit function
func (x *Rest) PermissionEdit(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewPermissionRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}

	if cReq.Permission.Category != params.ByName(`category`) {
		dispatchBadRequest(&w, fmt.Errorf(`Category mismatch`))
		return
	}
	if cReq.Permission.ID != params.ByName(`permissionID`) {
		dispatchBadRequest(&w, fmt.Errorf(`PermissionID mismatch`))
		return
	}
	if strings.Contains(params.ByName(`category`), `:grant`) {
		dispatchBadRequest(&w, fmt.Errorf(
			`Permissions in :grant categories can not be mapped`))
		return
	}
	if params.ByName(`category`) == msg.CategorySystem ||
		params.ByName(`category`) == msg.CategoryOmnipotence {
		dispatchForbidden(&w, nil)
		return
	}
	// invalid: map+unmap at the same time
	if cReq.Flags.Add && cReq.Flags.Remove {
		dispatchBadRequest(&w, fmt.Errorf(`Ambiguous instruction`))
		return
	}
	// invalid: batched mapping
	if cReq.Permission.Actions != nil && cReq.Permission.Sections != nil {
		dispatchBadRequest(&w, fmt.Errorf(`Invalid batch mapping`))
		return
	}
	if cReq.Permission.Actions != nil {
		if len(*cReq.Permission.Actions) != 1 ||
			params.ByName(`category`) != (*cReq.Permission.Actions)[0].Category {
			dispatchBadRequest(&w, fmt.Errorf(`Invalid action specification`))
			return
		}
	}
	if cReq.Permission.Sections != nil {
		if len(*cReq.Permission.Sections) != 1 ||
			params.ByName(`category`) != (*cReq.Permission.Sections)[0].Category {
			dispatchBadRequest(&w, fmt.Errorf(`Invalid section specification`))
			return
		}
	}

	request := msg.New(r, params)
	request.Section = msg.SectionPermission
	request.Permission.ID = cReq.Permission.ID
	request.Permission.Name = cReq.Permission.Name
	request.Permission.Category = cReq.Permission.Category
	// XXX Clone
	request.Permission.Sections = cReq.Permission.Sections
	request.Permission.Actions = cReq.Permission.Actions

	switch {
	case cReq.Flags.Add:
		request.Action = msg.ActionMap
	case cReq.Flags.Remove:
		request.Action = msg.ActionUnmap
	}

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
