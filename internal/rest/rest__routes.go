/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package rest // import "github.com/mjolnir42/soma/internal/rest"

import (
	"github.com/julienschmidt/httprouter"
)

// setupRouter returns a configured httprouter
func (x *Rest) setupRouter() *httprouter.Router {
	router := httprouter.New()

	router.GET(`/attribute/:attribute`, x.Check(x.BasicAuth(x.AttributeShow)))
	router.GET(`/attribute/`, x.Check(x.BasicAuth(x.AttributeList)))
	router.GET(`/bucket/:bucket`, x.Check(x.BasicAuth(x.BucketShow)))
	router.GET(`/bucket/`, x.Check(x.BasicAuth(x.BucketList)))
	router.GET(`/capability/:capabilityID`, x.Check(x.BasicAuth(x.CapabilityShow)))
	router.GET(`/capability/`, x.Check(x.BasicAuth(x.CapabilityList)))
	router.GET(`/sections/:section/actions/:action`, x.Check(x.BasicAuth(x.ActionShow)))
	router.GET(`/sections/:section/actions/`, x.Check(x.BasicAuth(x.ActionList)))
	router.GET(`/sync/node/`, x.Check(x.BasicAuth(x.NodeMgmtSync)))
	router.HEAD(`/authenticate/validate`, x.Check(x.BasicAuth(x.SupervisorValidate)))
	router.POST(`/filter/actions/`, x.Check(x.BasicAuth(x.ActionSearch)))
	router.POST(`/search/bucket/`, x.Check(x.BasicAuth(x.BucketSearch)))
	router.POST(`/search/capability/`, x.Check(x.BasicAuth(x.CapabilitySearch)))

	if !x.conf.ReadOnly {
		if !x.conf.Observer {
			router.DELETE(`/accounts/tokens/:account`, x.Check(x.BasicAuth(x.SupervisorTokenInvalidateAccount)))
			router.DELETE(`/attribute/:attribute`, x.Check(x.BasicAuth(x.AttributeRemove)))
			router.DELETE(`/bucket/:bucket/property/:type/:source`, x.Check(x.BasicAuth(x.BucketPropertyDestroy)))
			router.DELETE(`/capability/:capabilityID`, x.Check(x.BasicAuth(x.CapabilityRemove)))
			router.DELETE(`/node/:nodeID`, x.Check(x.BasicAuth(x.NodeMgmtRemove)))
			router.DELETE(`/sections/:section/actions/:action`, x.Check(x.BasicAuth(x.ActionRemove)))
			router.DELETE(`/tokens/global`, x.Check(x.BasicAuth(x.SupervisorTokenInvalidateGlobal)))
			router.DELETE(`/tokens/self/active`, x.Check(x.BasicAuth(x.SupervisorTokenInvalidate)))
			router.DELETE(`/tokens/self/all`, x.Check(x.BasicAuth(x.SupervisorTokenInvalidateSelf)))
			router.PATCH(`/accounts/password/:kexID`, x.Check(x.SupervisorPasswordChange))
			router.POST(`/attribute/`, x.Check(x.BasicAuth(x.AttributeAdd)))
			router.POST(`/bucket/:bucket/property/:type/`, x.Check(x.BasicAuth(x.BucketPropertyCreate)))
			router.POST(`/bucket/`, x.Check(x.BasicAuth(x.BucketCreate)))
			router.POST(`/capability/`, x.Check(x.BasicAuth(x.CapabilityAdd)))
			router.POST(`/kex/`, x.Check(x.SupervisorKex))
			router.POST(`/node/`, x.Check(x.BasicAuth(x.NodeMgmtAdd)))
			router.POST(`/sections/:section/actions/`, x.Check(x.BasicAuth(x.ActionAdd)))
			router.PUT(`/accounts/activate/root/:kexID`, x.Check(x.SupervisorActivateRoot))
			router.PUT(`/accounts/activate/user/:kexID`, x.Check(x.SupervisorActivateUser))
			router.PUT(`/accounts/password/:kexID`, x.Check(x.SupervisorPasswordReset))
			router.PUT(`/node/:nodeID`, x.Check(x.BasicAuth(x.NodeMgmtUpdate)))
			router.PUT(`/tokens/request/:kexID`, x.Check(x.SupervisorTokenRequest))
		}
	}
	return router
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
