/*-
 * Copyright (c) 2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"net/http"

	"github.com/mjolnir42/soma/internal/msg"
)

// replyNoContent returns a 204 HTTP statuscode reply with no content
func (x *Rest) replyNoContent(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusNoContent)
	(*w).Write(nil)
}

// replyBadRequest returns a 400 application error inside the returned
// JSON body
func (x *Rest) replyBadRequest(w *http.ResponseWriter, q *msg.Request, err error) {
	result := msg.FromRequest(q)
	result.BadRequest(err, q.Section)
	x.send(w, &result)
}

// replyForbidden returns a 403 application error inside the returned
// JSON body
func (x *Rest) replyForbidden(w *http.ResponseWriter, q *msg.Request) {
	result := msg.FromRequest(q)
	result.Forbidden(nil, q.Section)
	x.send(w, &result)
}

// replyServerError returns a 500 application error inside the
// returned JSON body
func (x *Rest) replyServerError(w *http.ResponseWriter, q *msg.Request, err error) {
	result := msg.FromRequest(q)
	result.ServerError(err, q.Section)
	x.send(w, &result)
}

// replyNotImplemented returns a 501 application error inside the
// returned JSON body
func (x *Rest) replyNotImplemented(w *http.ResponseWriter, q *msg.Request, err error) {
	result := msg.FromRequest(q)
	result.NotImplemented(err, q.Section)
	x.send(w, &result)
}

// hardConflict returns a 409 HTTP error with a provided error text.
// This function is intended to be used only for special policy
// violations
func (x *Rest) hardConflict(w *http.ResponseWriter, err error) {
	if err != nil {
		http.Error(*w, err.Error(), http.StatusConflict)
		return
	}
	http.Error(*w, http.StatusText(http.StatusConflict), http.StatusConflict)
}

// hardServerError returns a 500 HTTP error with no application data
// body. This function is intended to be used only if normal response
// generation itself fails
func (x *Rest) hardServerError(w *http.ResponseWriter) {
	http.Error(*w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
