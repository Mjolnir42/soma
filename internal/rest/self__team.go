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

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// TeamShow function
func (x *Rest) TeamShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionShow
	request.User.ID = params.ByName(`teamID`)

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// TeamSearch function
func (x *Rest) TeamSearch(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionTeam
	request.Action = msg.ActionSearch

	cReq := proto.NewTeamFilter()
	if err := decodeJSONBody(r, &cReq); err != nil {
		x.replyBadRequest(&w, &request, err)
		return
	}

	if cReq.Filter.Team.Name == `` {
		x.replyBadRequest(&w, &request, fmt.Errorf(
			`TeamSearch request missing Team.Name`))
		return
	}
	request.Search.Team.Name = cReq.Filter.Team.Name

	if !x.isAuthorized(&request) {
		x.replyForbidden(&w, &request, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply

	// XXX BUG filter in SQL statement
	filtered := []proto.Team{}
	for _, i := range result.Team {
		if i.Name == cReq.Filter.Team.Name {
			filtered = append(filtered, i)
		}
	}
	result.Team = filtered
	x.send(&w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
