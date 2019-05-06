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

// groupConfigPropertyCreateSystem function
func groupConfigPropertyCreateSystem(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeSystem, proto.EntityGroup)
}

// groupConfigPropertyUpdateSystem function
func groupConfigPropertyUpdateSystem(c *cli.Context) error {
	return variousPropertyUpdate(c, proto.PropertyTypeSystem, proto.EntityGroup)
}

// groupConfigPropertyCreateCustom function
func groupConfigPropertyCreateCustom(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeCustom, proto.EntityGroup)
}

// groupConfigPropertyUpdateCustom function
func groupConfigPropertyUpdateCustom(c *cli.Context) error {
	return variousPropertyUpdate(c, proto.PropertyTypeCustom, proto.EntityGroup)
}

// groupConfigPropertyCreateService function
func groupConfigPropertyCreateService(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeService, proto.EntityGroup)
}

// groupConfigPropertyCreateOncall function
func groupConfigPropertyCreateOncall(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeOncall, proto.EntityGroup)
}

// groupConfigPropertyDestroySystem function
func groupConfigPropertyDestroySystem(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeSystem, proto.EntityGroup)
}

// groupConfigPropertyDestroyCustom function
func groupConfigPropertyDestroyCustom(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeCustom, proto.EntityGroup)
}

// groupConfigPropertyDestroyService function
func groupConfigPropertyDestroyService(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeService, proto.EntityGroup)
}

// groupConfigPropertyDestroyOncall function
func groupConfigPropertyDestroyOncall(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeOncall, proto.EntityGroup)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
