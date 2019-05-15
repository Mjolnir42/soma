/*-
 * Copyright (c) 2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 Internet SE
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

// PropertyMgmtServiceAdd function
func (x *Rest) PropertyMgmtServiceAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionPropertyMgmt
	request.Action = msg.ActionAdd

	cReq := proto.NewPropertyRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	switch {
	case params.ByName(`propertyType`) != msg.PropertyService:
		x.replyBadRequest(&w, &request, fmt.Errorf("Invalid property type: %s", params.ByName(`propertyType`)))
		return
	case cReq.Property.Type != msg.PropertyService:
		x.replyBadRequest(&w, &request, fmt.Errorf("Invalid property type: %s", params.ByName(`propertyType`)))
		return
	case cReq.Property.Service.TeamID != params.ByName(`teamID`):
		x.replyBadRequest(&w, &request, fmt.Errorf("Mismatching team IDs: %s vs %s",
			cReq.Property.Service.TeamID, params.ByName(`teamID`)))
		return
	case cReq.Property.Service.Name == ``:
		x.replyBadRequest(&w, &request, fmt.Errorf(`Invalid empty service property name`))
		return
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	request.Section = msg.SectionPropertyService
	request.Property = cReq.Property.Clone()

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// PropertyMgmtServiceUpdate function
func (x *Rest) PropertyMgmtServiceUpdate(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionPropertyMgmt
	request.Action = msg.ActionUpdate

	cReq := proto.NewPropertyRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	switch {
	case params.ByName(`propertyType`) != msg.PropertyService:
		x.replyBadRequest(&w, &request, fmt.Errorf("Invalid property type: %s", params.ByName(`propertyType`)))
		return
	case cReq.Property.Type != msg.PropertyService:
		x.replyBadRequest(&w, &request, fmt.Errorf("Invalid property type: %s", params.ByName(`propertyType`)))
		return
	case cReq.Property.Service.TeamID != params.ByName(`teamID`):
		x.replyBadRequest(&w, &request, fmt.Errorf("Mismatching team IDs: %s vs %s",
			cReq.Property.Service.TeamID, params.ByName(`teamID`)))
		return
	case cReq.Property.Service.Name == ``:
		x.replyBadRequest(&w, &request, fmt.Errorf(`Invalid empty service property name`))
		return
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	request.Section = msg.SectionPropertyService
	request.Property = cReq.Property.Clone()

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// PropertyMgmtServiceRemove function
func (x *Rest) PropertyMgmtServiceRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionPropertyMgmt
	request.Action = msg.ActionRemove

	switch {
	case params.ByName(`propertyType`) != msg.PropertyService:
		x.replyBadRequest(&w, &request, fmt.Errorf("Invalid property type: %s", params.ByName(`propertyType`)))
		return
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	request.Section = msg.SectionPropertyService
	request.Property.Type = msg.PropertyService
	request.Property.Service = &proto.PropertyService{
		ID:     params.ByName(`propertyID`),
		TeamID: params.ByName(`teamID`),
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)

}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
