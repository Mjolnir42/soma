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

	router.GET(`/attribute/:attribute`, x.Verify(x.AttributeShow))
	router.GET(`/attribute/`, x.Verify(x.AttributeList))
	router.GET(`/bucket/:bucket`, x.Verify(x.BucketShow))
	router.GET(`/bucket/`, x.Verify(x.BucketList))
	router.GET(`/capability/:capabilityID`, x.Verify(x.CapabilityShow))
	router.GET(`/capability/`, x.Verify(x.CapabilityList))
	router.GET(`/category/:category`, x.Verify(x.CategoryShow))
	router.GET(`/category/`, x.Verify(x.CategoryList))
	router.GET(`/checkconfig/:repositoryID/:checkID`, x.Verify(x.CheckConfigShow))
	router.GET(`/checkconfig/:repositoryID/`, x.Verify(x.CheckConfigList))
	router.GET(`/cluster/:clusterID/members/`, x.Verify(x.ClusterMemberList))
	router.GET(`/cluster/:clusterID`, x.Verify(x.ClusterShow))
	router.GET(`/cluster/`, x.Verify(x.ClusterList))
	router.GET(`/sections/:section/actions/:action`, x.Verify(x.ActionShow))
	router.GET(`/sections/:section/actions/`, x.Verify(x.ActionList))
	router.GET(`/sync/node/`, x.Verify(x.NodeMgmtSync))
	router.HEAD(`/authenticate/validate`, x.Verify(x.SupervisorValidate))
	router.POST(`/filter/actions/`, x.Verify(x.ActionSearch))
	router.POST(`/search/bucket/`, x.Verify(x.BucketSearch))
	router.POST(`/search/capability/`, x.Verify(x.CapabilitySearch))
	router.POST(`/search/checkconfig/:repositoryID/`, x.Verify(x.CheckConfigSearch))
	router.POST(`/search/cluster/`, x.Verify(x.ClusterSearch))

	if !x.conf.ReadOnly {
		if !x.conf.Observer {
			router.DELETE(`/accounts/tokens/:account`, x.Verify(x.SupervisorTokenInvalidateAccount))
			router.DELETE(`/attribute/:attribute`, x.Verify(x.AttributeRemove))
			router.DELETE(`/bucket/:bucket/property/:type/:source`, x.Verify(x.BucketPropertyDestroy))
			router.DELETE(`/capability/:capabilityID`, x.Verify(x.CapabilityRemove))
			router.DELETE(`/category/:category`, x.Verify(x.CategoryRemove))
			router.DELETE(`/checkconfig/:repositoryID/:checkID`, x.Verify(x.CheckConfigDestroy))
			router.DELETE(`/cluster/:clusterID/property/:propertyType/:sourceID`, x.Verify(x.ClusterPropertyDestroy))
			router.DELETE(`/node/:nodeID`, x.Verify(x.NodeMgmtRemove))
			router.DELETE(`/sections/:section/actions/:action`, x.Verify(x.ActionRemove))
			router.DELETE(`/tokens/global`, x.Verify(x.SupervisorTokenInvalidateGlobal))
			router.DELETE(`/tokens/self/active`, x.Verify(x.SupervisorTokenInvalidate))
			router.DELETE(`/tokens/self/all`, x.Verify(x.SupervisorTokenInvalidateSelf))
			router.PATCH(`/accounts/password/:kexID`, x.CheckShutdown(x.SupervisorPasswordChange))
			router.POST(`/attribute/`, x.Verify(x.AttributeAdd))
			router.POST(`/bucket/:bucket/property/:type/`, x.Verify(x.BucketPropertyCreate))
			router.POST(`/bucket/`, x.Verify(x.BucketCreate))
			router.POST(`/capability/`, x.Verify(x.CapabilityAdd))
			router.POST(`/category/`, x.Verify(x.CategoryAdd))
			router.POST(`/checkconfig/:repositoryID/`, x.Verify(x.CheckConfigCreate))
			router.POST(`/cluster/:clusterID/members/`, x.Verify(x.ClusterMemberAssign))
			router.POST(`/cluster/:clusterID/property/:propertyType/`, x.Verify(x.ClusterPropertyCreate))
			router.POST(`/cluster/`, x.Verify(x.ClusterCreate))
			router.POST(`/kex/`, x.CheckShutdown(x.SupervisorKex))
			router.POST(`/node/`, x.Verify(x.NodeMgmtAdd))
			router.POST(`/sections/:section/actions/`, x.Verify(x.ActionAdd))
			router.PUT(`/accounts/activate/root/:kexID`, x.CheckShutdown(x.SupervisorActivateRoot))
			router.PUT(`/accounts/activate/user/:kexID`, x.CheckShutdown(x.SupervisorActivateUser))
			router.PUT(`/accounts/password/:kexID`, x.CheckShutdown(x.SupervisorPasswordReset))
			router.PUT(`/node/:nodeID`, x.Verify(x.NodeMgmtUpdate))
			router.PUT(`/tokens/request/:kexID`, x.CheckShutdown(x.SupervisorTokenRequest))
		}
	}
	return router
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
