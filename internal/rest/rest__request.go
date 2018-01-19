/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

// Package rest implements the REST routes to access SOMA.
package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mjolnir42/soma/internal/msg"
)

func newRequest(r *http.Request, params httprouter.Params) msg.Request {
	returnChannel := make(chan msg.Result, 1)
	return msg.Request{
		ID:         requestID(params),
		RemoteAddr: extractAddress(r.RemoteAddr),
		AuthUser:   params.ByName(`AuthenticatedUser`),
		Reply:      returnChannel,
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
