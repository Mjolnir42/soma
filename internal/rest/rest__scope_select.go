/*-
 * Copyright (c) 2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
)

// ScopeSelectMonitoringList function
func (x *Rest) ScopeSelectMonitoringList(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {

	request := newRequest(r, params)
	request.Section = msg.SectionMonitoringMgmt
	request.Action = msg.ActionAll
	request.Flag.Unscoped = true

	if x.isAuthorized(&request) {
		x.MonitoringMgmtAll(w, r, params)
		return
	}

	x.MonitoringList(w, r, params)
}

// ScopeSelectMonitoringSearch function
func (x *Rest) ScopeSelectMonitoringSearch(w http.ResponseWriter,
	r *http.Request, params httprouter.Params) {

	request := newRequest(r, params)
	request.Section = msg.SectionMonitoringMgmt
	request.Action = msg.ActionSearchAll
	request.Flag.Unscoped = true

	if x.isAuthorized(&request) {
		x.MonitoringMgmtSearchAll(w, r, params)
		return
	}

	x.MonitoringSearch(w, r, params)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
