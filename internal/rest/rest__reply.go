/*-
 * Copyright (c) 2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/auth"
	"github.com/mjolnir42/soma/lib/proto"
)

// replyForbidden returns a 403 error
func (x *Rest) replyForbidden(w *http.ResponseWriter, q *msg.Request, err error) {
	result := msg.FromRequest(q)
	result.Forbidden(err, q.Section)
	x.send(w, &result)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
