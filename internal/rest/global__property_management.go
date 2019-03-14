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

// PropertyMgmtList function
func (x *Rest) PropertyMgmtList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionPropertyMgmt
	request.Action = msg.ActionList

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	request.Property.Type = params.ByName(`propertyType`)

	switch request.Property.Type {
	case msg.PropertyNative:
		request.Section = msg.SectionPropertyNative
	case msg.PropertyTemplate:
		request.Section = msg.SectionPropertyTemplate
	case msg.PropertySystem:
		request.Section = msg.SectionPropertySystem
	case msg.PropertyCustom:
		request.Section = msg.SectionPropertyCustom
		request.Repository.ID = params.ByName(`repositoryID`)
		request.Property.RepositoryID = request.Repository.ID
		request.Property.Custom = &proto.PropertyCustom{}
		request.Property.Custom.RepositoryID = request.Repository.ID
	case msg.PropertyService:
		request.Section = msg.SectionPropertyService
		request.Team.ID = params.ByName(`teamID`)
		request.Property.Service = &proto.PropertyService{}
		request.Property.Service.TeamID = request.Team.ID
	default:
		x.replyBadRequest(&w, &request, fmt.Errorf("Invalid property type: %s", request.Property.Type))
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

// PropertyMgmtShow function
func (x *Rest) PropertyMgmtShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionPropertyMgmt
	request.Action = msg.ActionShow

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	request.Property.Type = params.ByName(`propertyType`)

	switch request.Property.Type {
	case msg.PropertyNative:
		request.Section = msg.SectionPropertyNative
		request.Property.Native = &proto.PropertyNative{}
		request.Property.Native.Name = params.ByName(`propertyID`)
	case msg.PropertyTemplate:
		request.Section = msg.SectionPropertyTemplate
		request.Property.Service = &proto.PropertyService{}
		request.Property.Service.ID = params.ByName(`propertyID`)
	case msg.PropertySystem:
		request.Section = msg.SectionPropertySystem
		request.Property.System = &proto.PropertySystem{}
		request.Property.System.Name = params.ByName(`propertyID`)
	case msg.PropertyCustom:
		request.Section = msg.SectionPropertyCustom
		request.Repository.ID = params.ByName(`repositoryID`)
		request.Property.RepositoryID = request.Repository.ID
		request.Property.Custom = &proto.PropertyCustom{}
		request.Property.Custom.RepositoryID = request.Repository.ID
		request.Property.Custom.ID = params.ByName(`propertyID`)
	case msg.PropertyService:
		request.Section = msg.SectionPropertyService
		request.Team.ID = params.ByName(`teamID`)
		request.Property.Service = &proto.PropertyService{}
		request.Property.Service.TeamID = request.Team.ID
		request.Property.Service.ID = params.ByName(`propertyID`)
	default:
		x.replyBadRequest(&w, &request, fmt.Errorf("Invalid property type: %s", request.Property.Type))
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

// PropertyMgmtSearch function
func (x *Rest) PropertyMgmtSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionPropertyMgmt
	request.Action = msg.ActionSearch

	cReq := proto.NewPropertyFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	if cReq.Filter.Property.Type != params.ByName(`propertyType`) {
		x.replyBadRequest(&w, &request, fmt.Errorf(
			"Mismatched propertyType %s vs %s",
			cReq.Filter.Property.Type, params.ByName(`propertyType`),
		))
		return
	}

	request.Property.Type = params.ByName(`propertyType`)

	switch cReq.Filter.Property.Type {
	case msg.PropertyNative:
		request.Section = msg.SectionPropertyNative
		request.Search.Property.Native = &proto.PropertyNative{}
		request.Search.Property.Native.Name = cReq.Filter.Property.Name
	case msg.PropertyTemplate:
		request.Section = msg.SectionPropertyTemplate
		request.Search.Property.Service = &proto.PropertyService{}
		request.Search.Property.Service.Name = cReq.Filter.Property.Name
	case msg.PropertySystem:
		request.Section = msg.SectionPropertySystem
		request.Search.Property.System = &proto.PropertySystem{}
		request.Search.Property.System.Name = cReq.Filter.Property.Name
	case msg.PropertyCustom:
		request.Section = msg.SectionPropertyCustom
		request.Repository.ID = cReq.Filter.Property.RepositoryID
		request.Search.Property.RepositoryID = cReq.Filter.Property.RepositoryID
		request.Search.Property.Custom = &proto.PropertyCustom{}
		request.Search.Property.Custom.Name = cReq.Filter.Property.Name
		request.Search.Property.Custom.RepositoryID = cReq.Filter.Property.RepositoryID
		if cReq.Filter.Property.RepositoryID != request.Repository.ID {
			x.replyBadRequest(&w, &request, fmt.Errorf(
				"PropertyMgmtSearch with mismatched repositoryID %s vs %s",
				cReq.Filter.Property.RepositoryID, request.Repository.ID,
			))
			return
		}
		if request.Repository.ID == `` {
			x.replyBadRequest(&w, &request, fmt.Errorf(
				"PropertyMgmtSearch with empty repositoryID",
			))
			return
		}
	case msg.PropertyService:
		request.Section = msg.SectionPropertyService
		request.Team.ID = params.ByName(`teamID`)
		request.Search.Property.Service = &proto.PropertyService{}
		request.Search.Property.Service.TeamID = request.Team.ID
		request.Search.Property.Service.Name = cReq.Filter.Property.Name

	default:
		x.replyBadRequest(&w, &request, fmt.Errorf(
			"PropertyMgmtSearch request has unknown property type: %s",
			cReq.Filter.Property.Type,
		))
		return
	}

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply

	// XXX BUG filter in SQL statement
	filtered := []proto.Property{}

	switch result.Section {
	case msg.SectionPropertyNative:
		for _, i := range result.Property {
			if i.Native.Name == cReq.Filter.Property.Name {
				filtered = append(filtered, i)
			}
		}
	case msg.SectionPropertyTemplate:
		for _, i := range result.Property {
			if i.Service.Name == cReq.Filter.Property.Name {
				filtered = append(filtered, i)
			}
		}
	case msg.SectionPropertySystem:
		for _, i := range result.Property {
			if i.System.Name == cReq.Filter.Property.Name {
				filtered = append(filtered, i)
			}
		}
	case msg.SectionPropertyCustom:
		for _, i := range result.Property {
			if (i.Custom.Name == cReq.Filter.Property.Name) &&
				(i.Custom.RepositoryID == cReq.Filter.Property.RepositoryID) {
				filtered = append(filtered, i)
			}
		}
	case msg.SectionPropertyService:
		for _, i := range result.Property {
			if (i.Service.Name == cReq.Filter.Property.Name) &&
				(i.Service.TeamID == params.ByName(`teamID`)) {
				filtered = append(filtered, i)
			}
		}
	}
	result.Property = filtered
	x.send(&w, &result)
}

// PropertyMgmtAdd function
func (x *Rest) PropertyMgmtAdd(w http.ResponseWriter, r *http.Request,
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

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}
	request.Property.Type = params.ByName(`propertyType`)

	if request.Property.Type != cReq.Property.Type {
		x.replyBadRequest(&w, &request, fmt.Errorf("Mismatching property types in URI and body: %s vs %s",
			request.Property.Type,
			cReq.Property.Type,
		))
		return
	}

	switch request.Property.Type {
	case msg.PropertyNative:
		request.Section = msg.SectionPropertyNative
		request.Property = cReq.Property.Clone()
		if request.Property.Native.Name == `` {
			x.replyBadRequest(&w, &request, fmt.Errorf(`Invalid empty property name`))
			return
		}
	case msg.PropertyTemplate:
		request.Section = msg.SectionPropertyTemplate
		request.Property = cReq.Property.Clone()
		if request.Property.Service.Name == `` {
			x.replyBadRequest(&w, &request, fmt.Errorf(`Invalid empty property name`))
			return
		}
	case msg.PropertySystem:
		request.Section = msg.SectionPropertySystem
		request.Property = cReq.Property.Clone()
		if request.Property.System.Name == `` {
			x.replyBadRequest(&w, &request, fmt.Errorf(`Invalid empty property name`))
			return
		}
	case msg.PropertyCustom:
		fallthrough
	case msg.PropertyService:
		x.replyServerError(&w, &request, fmt.Errorf(
			"Request routing error. Type %s request appeared in global"+
				" handler for types: native, template, system",
			request.Property.Type))
		return
	default:
		x.replyBadRequest(&w, &request, fmt.Errorf("Invalid property type: %s", request.Property.Type))
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

// PropertyMgmtRemove function
func (x *Rest) PropertyMgmtRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionPropertyMgmt
	request.Action = msg.ActionRemove

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request)
		return
	}
	request.Property.Type = params.ByName(`propertyType`)

	switch request.Property.Type {
	case msg.PropertyNative:
		request.Section = msg.SectionPropertyNative
		request.Property.Native = &proto.PropertyNative{}
		request.Property.Native.Name = params.ByName(`propertyID`)
	case msg.PropertyTemplate:
		request.Section = msg.SectionPropertyTemplate
		request.Property.Service = &proto.PropertyService{}
		request.Property.Service.ID = params.ByName(`propertyID`)
	case msg.PropertySystem:
		request.Section = msg.SectionPropertySystem
		request.Property.System = &proto.PropertySystem{}
		request.Property.System.Name = params.ByName(`propertyID`)
	case msg.PropertyCustom:
		fallthrough
	case msg.PropertyService:
		x.replyServerError(&w, &request, fmt.Errorf(
			"Request routing error. Type %s request appeared in global"+
				" handler for types: native, template, system",
			request.Property.Type))
		return
	default:
		x.replyBadRequest(&w, &request, fmt.Errorf("Invalid property type: %s", request.Property.Type))
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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
