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

// nodeConfigPropertyCreateSystem function
func nodeConfigPropertyCreateSystem(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeSystem, proto.EntityNode)
}

// nodeConfigPropertyCreateService function
func nodeConfigPropertyCreateService(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeService, proto.EntityNode)
}

// nodeConfigPropertyCreateOncall function
func nodeConfigPropertyCreateOncall(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeOncall, proto.EntityNode)
}

// nodeConfigPropertyCreateCustom function
func nodeConfigPropertyCreateCustom(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeCustom, proto.EntityNode)
}

// nodeConfigPropertyDestroySystem function
func nodeConfigPropertyDestroySystem(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeSystem, proto.EntityNode)
}

// nodeConfigPropertyDestroyService function
func nodeConfigPropertyDestroyService(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeService, proto.EntityNode)
}

// nodeConfigPropertyDestroyOncall function
func nodeConfigPropertyDestroyOncall(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeOncall, proto.EntityNode)
}

// nodeConfigPropertyDestroyCustom function
func nodeConfigPropertyDestroyCustom(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeCustom, proto.EntityNode)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
