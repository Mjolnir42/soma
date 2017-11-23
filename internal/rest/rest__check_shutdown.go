/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	metrics "github.com/rcrowley/go-metrics"
)

// CheckShutdown denies the request if a shutdown is in progress
func (x *Rest) CheckShutdown(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {

		if !ShutdownInProgress {
			metrics.GetOrRegisterCounter(`.requests`, Metrics[`soma`]).Inc(1)
			start := time.Now()

			h(w, r, ps)

			metrics.GetOrRegisterTimer(`.requests.latency`,
				Metrics[`soma`]).UpdateSince(start)
			return
		}

		http.Error(w, `Shutdown in progress`,
			http.StatusServiceUnavailable)
	}
}

// Verify is a wrapper for CheckShutdown and BasicAuth checks
func (x *Rest) Verify(h httprouter.Handle) httprouter.Handle {
	return x.CheckShutdown(
		x.BasicAuth(
			func(w http.ResponseWriter, r *http.Request,
				ps httprouter.Params) {
				h(w, r, ps)
			},
		),
	)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
