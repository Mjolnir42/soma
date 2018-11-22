/*-
 * Copyright (c) 2016-2018, Jörg Pernfuß
 * Copyright (c) 2016-2018, 1&1 Internet SE
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

// PropertyMgmtCustomAdd function
func (x *Rest) PropertyMgmtCustomAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionPropertyCustom
	request.Action = msg.ActionAdd

	cReq := proto.NewPropertyRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	switch {
	case params.ByName(`propertyType`) != msg.PropertyCustom:
		x.replyBadRequest(&w, &request, fmt.Errorf("Invalid property type: %s", params.ByName(`propertyType`)))
		return
	case cReq.Property.Type != msg.PropertyCustom:
		x.replyBadRequest(&w, &request, fmt.Errorf("Invalid property type: %s", params.ByName(`propertyType`)))
		return
	case cReq.Property.Custom.RepositoryID != params.ByName(`repositoryID`):
		x.replyBadRequest(&w, &request, fmt.Errorf("Mismatching repository IDs: %s vs %s",
			cReq.Property.Custom.RepositoryID, params.ByName(`repositoryID`)))
		return
	case cReq.Property.Custom.Name == ``:
		x.replyBadRequest(&w, &request, fmt.Errorf(`Invalid empty custom property name`))
		return
	}
	request.Property = cReq.Property.Clone()

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// PropertyMgmtCustomRemove function
func (x *Rest) PropertyMgmtCustomRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionPropertyCustom
	request.Action = msg.ActionRemove

	switch {
	case params.ByName(`propertyType`) != msg.PropertyCustom:
		x.replyBadRequest(&w, &request, fmt.Errorf("Invalid property type: %s", params.ByName(`propertyType`)))
		return
	}

	request.Property.Type = msg.PropertyCustom
	request.Property.RepositoryID = params.ByName(`repositoryID`)
	request.Property.Custom = &proto.PropertyCustom{
		ID:           params.ByName(`propertyID`),
		RepositoryID: params.ByName(`repositoryID`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
