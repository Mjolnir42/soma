/*-
 * Copyright (c) 2015-2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
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

// repositoryConfigPropertyCreateSystem function
// soma repository property create system ${system} on ${repository} view ${view} \
//      value ${value} [inheritance ${inherit}] [childrenonly ${child}]
func repositoryConfigPropertyCreateSystem(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeSystem, proto.EntityRepository)
}

// repositoryConfigPropertyUpdateSystem function
// soma repository property create system ${system} on ${repository} view ${view} \
//      value ${value} [inheritance ${inherit}] [childrenonly ${child}]
func repositoryConfigPropertyUpdateSystem(c *cli.Context) error {
	return variousPropertyUpdate(c, proto.PropertyTypeSystem, proto.EntityRepository)
}

// repositoryConfigPropertyCreateCustom function
// soma repository property create custom ${custom} on ${repository} view ${view} \
//      value ${value} [inheritance ${inherit}] [childrenonly ${child}]
func repositoryConfigPropertyCreateCustom(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeCustom, proto.EntityRepository)
}

// repositoryConfigPropertyUpdateCustom function
// soma repository property create custom ${custom} on ${repository} view ${view} \
//      value ${value} [inheritance ${inherit}] [childrenonly ${child}]
func repositoryConfigPropertyUpdateCustom(c *cli.Context) error {
	return variousPropertyUpdate(c, proto.PropertyTypeCustom, proto.EntityRepository)
}

// repositoryConfigPropertyCreateService function
// soma repository property create service ${service} on ${repository} view ${view}
//      [inheritance ${inherit}] [childrenonly ${child}]
func repositoryConfigPropertyCreateService(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeService, proto.EntityRepository)
}

// repositoryConfigPropertyCreateOncall function
// soma repository property create oncall ${oncall} on ${repository} view ${view}
//      [inheritance ${inherit}] [childrenonly ${child}]
func repositoryConfigPropertyCreateOncall(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeOncall, proto.EntityRepository)
}

// repositoryConfigPropertyDestroySystem function
// soma repository property destroy system ${system} on ${repository} view ${view}
func repositoryConfigPropertyDestroySystem(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeSystem, proto.EntityRepository)
}

// repositoryConfigPropertyDestroyCustom function
// soma repository property destroy custom ${custom} on ${repository} view ${view}
func repositoryConfigPropertyDestroyCustom(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeCustom, proto.EntityRepository)
}

// repositoryConfigPropertyDestroyService function
// soma repository property destroy service ${service} on ${repository} view ${view}
func repositoryConfigPropertyDestroyService(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeService, proto.EntityRepository)
}

// repositoryConfigPropertyDestroyOncall function
// soma repository property destroy oncall ${oncall} on ${repository} view ${view}
func repositoryConfigPropertyDestroyOncall(c *cli.Context) error {
	return variousPropertyDestroy(c, proto.PropertyTypeOncall, proto.EntityRepository)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
