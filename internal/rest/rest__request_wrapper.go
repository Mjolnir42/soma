/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/satori/go.uuid"
)

// Unauthenticated is a wrapper for unauthenticated or implicitly
// authenticated requests
func (x *Rest) Unauthenticated(h httprouter.Handle) httprouter.Handle {
	return x.checkShutdown(
		x.enrich(
			x.intakeLog(
				func(w http.ResponseWriter, r *http.Request,
					ps httprouter.Params) {
					h(w, r, ps)
				},
			),
		),
	)
}

// Authenticated is the standard request wrapper
func (x *Rest) Authenticated(h httprouter.Handle) httprouter.Handle {
	return x.Unauthenticated(
		x.basicAuth(
			func(w http.ResponseWriter, r *http.Request,
				ps httprouter.Params) {
				h(w, r, ps)
			},
		),
	)
}

// checkShutdown denies the request if a shutdown is in progress
func (x *Rest) checkShutdown(h httprouter.Handle) httprouter.Handle {
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

// enrich is a wrapper that adds metadata information to the request
func (x *Rest) enrich(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {

		// generate and record the requestID
		requestID := uuid.Must(uuid.NewV4())
		ps = append(ps, httprouter.Param{
			Key:   `RequestID`,
			Value: requestID.String(),
		})

		// record the request URI
		ps = append(ps, httprouter.Param{
			Key:   `RequestURI`,
			Value: r.RequestURI,
		})

		h(w, r, ps)
	}
}

// intakeLog writes the pre-authentication record into the request log
func (x *Rest) intakeLog(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		ps httprouter.Params) {

		x.reqLog.
			WithField(
				`RequestID`, ps.ByName(`RequestID`),
			).
			WithField(
				`RequestURI`, ps.ByName(`RequestURI`),
			).
			WithField(
				`Phase`, `pre-authentication`,
			).
			Debug(`received:rest`)

		h(w, r, ps)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
