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

// MonitoringMgmtAll function
func (x *Rest) MonitoringMgmtAll(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionMonitoringMgmt
	request.Action = msg.ActionAll
	request.Flag.Unscoped = true

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// MonitoringMgmtSearchAll function
func (x *Rest) MonitoringMgmtSearchAll(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewMonitoringFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if cReq.Filter.Monitoring.Name == `` {
		dispatchBadRequest(&w, fmt.Errorf(
			`Empty search request: name missing`))
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionMonitoringMgmt
	request.Action = msg.ActionSearchAll
	request.Flag.Unscoped = true
	request.Search.Monitoring.Name = cReq.Filter.Monitoring.Name

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// MonitoringMgmtAdd function
func (x *Rest) MonitoringMgmtAdd(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	cReq := proto.NewMonitoringRequest()
	if err := decodeJSONBody(r, &cReq); err != nil {
		dispatchBadRequest(&w, err)
		return
	}
	if strings.Contains(cReq.Monitoring.Name, `.`) {
		dispatchBadRequest(&w, fmt.Errorf(
			`Invalid monitoring system`+
				` name containing . character`))
		return
	}

	request := msg.New(r, params)
	request.Section = msg.SectionMonitoringMgmt
	request.Action = msg.ActionAdd
	request.Monitoring.Name = cReq.Monitoring.Name
	request.Monitoring.Mode = cReq.Monitoring.Mode
	request.Monitoring.Contact = cReq.Monitoring.Contact
	request.Monitoring.TeamID = cReq.Monitoring.TeamID
	request.Monitoring.Callback = cReq.Monitoring.Callback

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// MonitoringMgmtRemove function
func (x *Rest) MonitoringMgmtRemove(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionMonitoringMgmt
	request.Action = msg.ActionRemove
	request.Monitoring.ID = params.ByName(`monitoringID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
