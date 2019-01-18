/*-
 * Copyright (c) 2016-2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/lib/proto"
)

// clusterConfigPropertyCreateSystem function
func clusterConfigPropertyCreateSystem(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeSystem, proto.EntityCluster)
}

// clusterConfigPropertyCreateCustom function
func clusterConfigPropertyCreateCustom(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeCustom, proto.EntityCluster)
}

// clusterConfigPropertyCreateService function
func clusterConfigPropertyCreateService(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeService, proto.EntityCluster)
}

// clusterConfigPropertyCreateOncall function
func clusterConfigPropertyCreateOncall(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeOncall, proto.EntityCluster)
}

// clusterConfigPropertyDestroySystem function
func clusterConfigPropertyDestroySystem(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeSystem, proto.EntityCluster)
}

// clusterConfigPropertyDestroyCustom function
func clusterConfigPropertyDestroyCustom(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeCustom, proto.EntityCluster)
}

// clusterConfigPropertyDestroyService function
func clusterConfigPropertyDestroyService(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeService, proto.EntityCluster)
}

// clusterConfigPropertyDestroyOncall function
func clusterConfigPropertyDestroyOncall(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeOncall, proto.EntityCluster)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
