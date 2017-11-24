package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// PropertyList function
func PropertyList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	prType, _ := GetPropertyTypeFromUrl(r.URL)
	var section string
	switch prType {
	case `native`, `system`, `custom`, `service`, `template`:
		section = fmt.Sprintf("property_%s", prType)
	default:
		DispatchBadRequest(&w, fmt.Errorf(`Invalid property type`))
		return
	}

	if !fixIsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    section,
		Action:     `list`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "list",
		reply:  returnChannel,
	}
	switch prType {
	case "native":
		req.prType = prType
	case "system":
		req.prType = prType
	case "custom":
		req.prType = prType
		req.Custom.RepositoryID = params.ByName("repository")
	case "service":
		req.prType = prType
		req.Service.TeamID = params.ByName("team")
	case "template":
		req.prType = prType
	default:
		SendPropertyReply(&w, &somaResult{})
	}

	handler := handlerMap["propertyReadHandler"].(*somaPropertyReadHandler)
	handler.input <- req
	result := <-returnChannel

	// declare here since goto does not jump over declarations
	cReq := proto.NewPropertyFilter()
	if result.Failure() {
		goto skip
	}

	_ = DecodeJsonBody(r, &cReq)
	if (cReq.Filter.Property.Type == "custom") && (cReq.Filter.Property.Name != "") &&
		(cReq.Filter.Property.RepositoryID != "") {
		filtered := []somaPropertyResult{}
		for _, i := range result.Properties {
			if (i.Custom.Name == cReq.Filter.Property.Name) &&
				(i.Custom.RepositoryID == cReq.Filter.Property.RepositoryID) {
				filtered = append(filtered, i)
			}
		}
		result.Properties = filtered
	}
	if (cReq.Filter.Property.Type == "system") && (cReq.Filter.Property.Name != "") {
		filtered := []somaPropertyResult{}
		for _, i := range result.Properties {
			if i.System.Name == cReq.Filter.Property.Name {
				filtered = append(filtered, i)
			}
		}
		result.Properties = filtered
	}
	if (cReq.Filter.Property.Type == "service") && (cReq.Filter.Property.Name != "") {
		filtered := []somaPropertyResult{}
		for _, i := range result.Properties {
			if (i.Service.Name == cReq.Filter.Property.Name) &&
				(i.Service.TeamID == params.ByName("team")) {
				filtered = append(filtered, i)
			}
		}
		result.Properties = filtered
	}
	if (cReq.Filter.Property.Type == "template") && (cReq.Filter.Property.Name != "") {
		filtered := []somaPropertyResult{}
		for _, i := range result.Properties {
			if i.Service.Name == cReq.Filter.Property.Name {
				filtered = append(filtered, i)
			}
		}
		result.Properties = filtered
	}

skip:
	SendPropertyReply(&w, &result)
}

// PropertyShow function
func PropertyShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	prType, _ := GetPropertyTypeFromUrl(r.URL)
	var section string
	switch prType {
	case `native`, `system`, `custom`, `service`, `template`:
		section = fmt.Sprintf("property_%s", prType)
	default:
		DispatchBadRequest(&w, fmt.Errorf(`Invalid property type`))
		return
	}

	if !fixIsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    section,
		Action:     `show`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "show",
		reply:  returnChannel,
	}
	switch prType {
	case "native":
		req.prType = prType
		req.Native.Name = params.ByName("native")
	case "system":
		req.prType = prType
		req.System.Name = params.ByName("system")
	case "custom":
		req.prType = prType
		req.Custom.ID = params.ByName("custom")
		req.Custom.RepositoryID = params.ByName("repository")
	case "service":
		req.prType = prType
		req.Service.Name = params.ByName("service")
		req.Service.TeamID = params.ByName("team")
	case "template":
		req.prType = prType
		req.Service.Name = params.ByName("service")
	default:
		SendPropertyReply(&w, &somaResult{})
	}

	handler := handlerMap["propertyReadHandler"].(*somaPropertyReadHandler)
	handler.input <- req
	result := <-returnChannel
	SendPropertyReply(&w, &result)
}

// PropertyAdd function
func PropertyAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	prType, _ := GetPropertyTypeFromUrl(r.URL)
	var section string
	switch prType {
	case `native`, `system`, `custom`, `service`, `template`:
		section = fmt.Sprintf("property_%s", prType)
	default:
		DispatchBadRequest(&w, fmt.Errorf(`Invalid property type`))
		return
	}

	if !fixIsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    section,
		Action:     `add`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	cReq := proto.NewPropertyRequest()
	err := DecodeJsonBody(r, &cReq)
	if err != nil {
		DispatchBadRequest(&w, err)
		return
	}
	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "add",
		reply:  returnChannel,
	}
	switch prType {
	case "native":
		req.prType = prType
		req.Native = *cReq.Property.Native
	case "system":
		req.prType = prType
		req.System = *cReq.Property.System
	case "custom":
		if params.ByName("repository") != cReq.Property.Custom.RepositoryID {
			DispatchBadRequest(&w, errors.New("Body and URL repositories do not match"))
			return
		}
		req.prType = prType
		req.Custom = *cReq.Property.Custom
		req.Custom.RepositoryID = params.ByName("repository")
	case "service":
		if params.ByName("team") != cReq.Property.Service.TeamID {
			DispatchBadRequest(&w, errors.New("Body and URL teams do not match"))
			return
		}
		req.prType = prType
		req.Service = *cReq.Property.Service
		req.Service.TeamID = params.ByName("team")
	case "template":
		req.prType = prType
		req.Service = *cReq.Property.Service
	default:
		SendPropertyReply(&w, &somaResult{})
	}

	handler := handlerMap["propertyWriteHandler"].(*somaPropertyWriteHandler)
	handler.input <- req
	result := <-returnChannel
	SendPropertyReply(&w, &result)
}

// PropertyRemove function
func PropertyRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer PanicCatcher(w)

	prType, _ := GetPropertyTypeFromUrl(r.URL)
	var section string
	switch prType {
	case `native`, `system`, `custom`, `service`, `template`:
		section = fmt.Sprintf("property_%s", prType)
	default:
		DispatchBadRequest(&w, fmt.Errorf(`Invalid property type`))
		return
	}

	if !fixIsAuthorized(&msg.Authorization{
		AuthUser:   params.ByName(`AuthenticatedUser`),
		RemoteAddr: extractAddress(r.RemoteAddr),
		Section:    section,
		Action:     `remove`,
	}) {
		DispatchForbidden(&w, nil)
		return
	}

	returnChannel := make(chan somaResult)
	req := somaPropertyRequest{
		action: "delete",
		reply:  returnChannel,
	}
	switch prType {
	case "native":
		req.prType = prType
		req.Native.Name = params.ByName("native")
	case "system":
		req.prType = prType
		req.System.Name = params.ByName("system")
	case "custom":
		req.prType = prType
		req.Custom.ID = params.ByName("custom")
		req.Custom.RepositoryID = params.ByName("repository")
	case "service":
		req.prType = prType
		req.Service.Name = params.ByName("service")
		req.Service.TeamID = params.ByName("team")
	case "template":
		req.prType = prType
		req.Service.Name = params.ByName("service")
	default:
		SendPropertyReply(&w, &somaResult{})
	}

	handler := handlerMap["propertyWriteHandler"].(*somaPropertyWriteHandler)
	handler.input <- req
	result := <-returnChannel
	SendPropertyReply(&w, &result)
}

// SendPropertyReply function
func SendPropertyReply(w *http.ResponseWriter, r *somaResult) {
	result := proto.NewPropertyResult()
	if r.MarkErrors(&result) {
		goto dispatch
	}
	for _, i := range (*r).Properties {
		switch i.prType {
		case "system":
			*result.Properties = append(*result.Properties, proto.Property{Type: "system",
				System: &proto.PropertySystem{
					Name:  i.System.Name,
					Value: i.System.Value,
				}})
		case "native":
			*result.Properties = append(*result.Properties, proto.Property{Type: "native",
				Native: &proto.PropertyNative{
					Name:  i.Native.Name,
					Value: i.Native.Value,
				}})
		case "custom":
			*result.Properties = append(*result.Properties, proto.Property{Type: "custom",
				Custom: &proto.PropertyCustom{
					ID:           i.Custom.ID,
					Name:         i.Custom.Name,
					Value:        i.Custom.Value,
					RepositoryID: i.Custom.RepositoryID,
				}})
		case "service":
			prop := proto.Property{
				Type: "service",
				Service: &proto.PropertyService{
					Name:       i.Service.Name,
					TeamID:     i.Service.TeamID,
					Attributes: []proto.ServiceAttribute{},
				}}
			for _, a := range i.Service.Attributes {
				prop.Service.Attributes = append(prop.Service.Attributes, proto.ServiceAttribute{
					Name:  a.Name,
					Value: a.Value,
				})
			}
			*result.Properties = append(*result.Properties, prop)
		case "template":
			prop := proto.Property{
				Type: "template",
				Service: &proto.PropertyService{
					Name:       i.Service.Name,
					Attributes: []proto.ServiceAttribute{},
				}}
			for _, a := range i.Service.Attributes {
				prop.Service.Attributes = append(prop.Service.Attributes, proto.ServiceAttribute{
					Name:  a.Name,
					Value: a.Value,
				})
			}
			*result.Properties = append(*result.Properties, prop)
		}
		if i.ResultError != nil {
			*result.Errors = append(*result.Errors, i.ResultError.Error())
		}
	}

dispatch:
	json, err := json.Marshal(result)
	if err != nil {
		DispatchInternalError(w, err)
		return
	}
	DispatchJsonReply(w, &json)
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
