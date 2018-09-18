/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
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

// InstanceShow returns information about a check instance
func (x *Rest) InstanceShow(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	// BUG needs scope switch repository|bucket|group|cluster|node

	request := msg.New(r, params)
	request.Section = msg.SectionInstance
	request.Action = msg.ActionShow
	request.Instance.ID = params.ByName(`instanceID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// InstanceVersions returns information about a check instance's
// version history
func (x *Rest) InstanceVersions(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	// BUG needs scope switch repository|bucket|group|cluster|node

	request := msg.New(r, params)
	request.Section = msg.SectionInstance
	request.Action = msg.ActionVersions
	request.Instance.ID = params.ByName(`instanceID`)

	if !x.isAuthorized(&request) {
		dispatchForbidden(&w, nil)
		return
	}

	x.handlerMap.MustLookup(&request).Intake() <- request
	result := <-request.Reply
	x.send(&w, &result)
}

// InstanceList returns the list of instances in the subtree
// below the queried object.
// Currently only supports repositories and buckets as target.
func (x *Rest) InstanceList(w http.ResponseWriter, r *http.Request,
	params httprouter.Params) {
	defer panicCatcher(w)

	request := msg.New(r, params)
	request.Section = msg.SectionInstance
	request.Action = msg.ActionList

	switch {
	case params.ByName(`repository`) != ``:
		request.Instance.ObjectType = `repository`
		request.Instance.ObjectID = params.ByName(`repository`)
		request.Instance.RepositoryID = params.ByName(`repository`)
	case params.ByName(`bucket`) != ``:
		request.Instance.ObjectType = `bucket`
		request.Instance.ObjectID = params.ByName(`bucket`)
		request.Instance.BucketID = params.ByName(`bucket`)
	case params.ByName(`group`) != ``:
		fallthrough
	case params.ByName(`cluster`) != ``:
		fallthrough
	case params.ByName(`node`) != ``:
		fallthrough
	default:
		dispatchNotImplemented(&w, nil)
		return
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
