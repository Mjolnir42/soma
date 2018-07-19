/*-
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Ping is the function for HEAD requests on the base API that
// reports facts about the running application
func (x *Rest) Ping(w http.ResponseWriter, _ *http.Request,
	_ httprouter.Params) {
	defer panicCatcher(w)

	w.Header().Set(`X-Powered-By`, `SOMA Configuration System`)
	w.Header().Set(`X-Version`, x.conf.Version)
	switch {
	case x.conf.Observer == true:
		w.Header().Set(`X-SOMA-Mode`, `Observer`)
	case x.conf.ReadOnly == true:
		w.Header().Set(`X-SOMA-Mode`, `ReadOnly`)
	case x.conf.ReadOnly == false:
		w.Header().Set(`X-SOMA-Mode`, `Master`)
	}
	w.WriteHeader(http.StatusNoContent)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
