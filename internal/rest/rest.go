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
	"github.com/mjolnir42/soma/internal/config"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/soma"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/satori/go.uuid"
)

// ShutdownInProgress indicates a pending service shutdown
var ShutdownInProgress bool

// Metrics is the map of runtime metric registries
var Metrics = make(map[string]metrics.Registry)

// Rest holds the required state for the REST interface
type Rest struct {
	isAuthorized func(*msg.Request) bool
	handlerMap   *soma.HandlerMap
	conf         *config.Config
	restricted   bool
}

// New returns a new REST interface
func New(
	authorizationFunction func(*msg.Request) bool,
	appHandlerMap *soma.HandlerMap,
	conf *config.Config,
) *Rest {
	x := Rest{}
	x.isAuthorized = authorizationFunction
	x.restricted = false
	x.handlerMap = appHandlerMap
	x.conf = conf
	return &x
}

// Run is the event server for Rest
func (x *Rest) Run() {
	router := httprouter.New()

	router.GET(`/sync/node/`, x.Check(x.BasicAuth(x.NodeSync)))
	router.HEAD(`/authenticate/validate/`, x.Check(x.BasicAuth(x.SupervisorValidate)))

	if !x.conf.ReadOnly {
		router.POST(`/authenticate/`, x.Check(x.SupervisorKex))
		router.PUT(`/authenticate/token/:uuid`, x.Check(x.SupervisorTokenRequest))

		if !x.conf.Observer {
			router.DELETE(`/node/:nodeID`, x.Check(x.BasicAuth(x.NodeRemove)))
			router.PATCH(`/authenticate/user/password/:uuid`, x.Check(x.SupervisorPasswordChange))
			router.POST(`/node/`, x.Check(x.BasicAuth(x.NodeAdd)))
			router.PUT(`/authenticate/activate/:uuid`, x.Check(x.SupervisorActivateUser))
			router.PUT(`/authenticate/bootstrap/:uuid`, x.Check(x.SupervisorBootstrap))
			router.PUT(`/authenticate/user/password/:uuid`, x.Check(x.SupervisorPasswordReset))
			router.PUT(`/node/:nodeID`, x.Check(x.BasicAuth(x.NodeUpdate)))
		}
	}

	// TODO switch to new abortable interface
	if x.conf.Daemon.TLS {
		// XXX log.Fatal
		http.ListenAndServeTLS(
			x.conf.Daemon.URL.Host,
			x.conf.Daemon.Cert,
			x.conf.Daemon.Key,
			router,
		)
	} else {
		// XXX log.Fatal
		http.ListenAndServe(x.conf.Daemon.URL.Host, router)
	}
}

// requestID extracts the RequestID set by Basic Authentication, making
// the ID consistent between all logs
func requestID(params httprouter.Params) (id uuid.UUID) {
	id, _ = uuid.FromString(params.ByName(`RequestID`))
	return
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
