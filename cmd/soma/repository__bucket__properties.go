/*-
 * Copyright (c) 2015-2019, Jörg Pernfuß
 * Copyright (c) 2018-2019, 1&1 IONOS SE
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

// bucketPropertyCreateSystem function
func bucketPropertyCreateSystem(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeSystem, proto.EntityBucket)
}

// bucketPropertyUpdateSystem function
func bucketPropertyUpdateSystem(c *cli.Context) error {
	return variousPropertyUpdate(c, proto.PropertyTypeSystem, proto.EntityBucket)
}

// bucketPropertyCreateCustom function
func bucketPropertyCreateCustom(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeCustom, proto.EntityBucket)
}

// bucketPropertyUpdateCustom function
func bucketPropertyUpdateCustom(c *cli.Context) error {
	return variousPropertyUpdate(c, proto.PropertyTypeCustom, proto.EntityBucket)
}

// bucketPropertyCreateService function
func bucketPropertyCreateService(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeService, proto.EntityBucket)
}

// bucketPropertyCreateOncall function
func bucketPropertyCreateOncall(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeOncall, proto.EntityBucket)
}

// bucketPropertyDestroySystem function
func bucketPropertyDestroySystem(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeSystem, proto.EntityBucket)
}

// bucketPropertyDestroyCustom function
func bucketPropertyDestroyCustom(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeCustom, proto.EntityBucket)
}

// bucketPropertyDestroyService function
func bucketPropertyDestroyService(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeService, proto.EntityBucket)
}

// bucketPropertyDestroyOncall function
func bucketPropertyDestroyOncall(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeOncall, proto.EntityBucket)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
