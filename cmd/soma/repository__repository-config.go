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
	"fmt"
	"net/url"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

// repositoryConfigList function
// soma repository list
func repositoryConfigList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/repository/`, `list`, nil, c)
}

// repositoryConfigSearch function
// soma repository search [id ${uuid}] [name ${repository}] [team ${team}] [deleted ${isDeleted}] [active ${isActive}]
func repositoryConfigSearch(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`id`, `name`, `team`, `deleted`, `active`}
	mandatoryOptions := []string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		adm.AllArguments(c),
	); err != nil {
		return err
	}

	validCondition := false
	req := proto.NewRepositoryFilter()

	if _, ok := opts[`id`]; ok {
		req.Filter.Repository.ID = opts[`id`][0]
		if err := adm.ValidateUUID(req.Filter.Repository.ID); err != nil {
			return err
		}
		validCondition = true
	}

	if _, ok := opts[`name`]; ok {
		req.Filter.Repository.Name = opts[`name`][0]
		if err := adm.ValidateNotUUID(req.Filter.Repository.Name); err != nil {
			return err
		}
		validCondition = true
	}

	if _, ok := opts[`team`]; ok {
		var teamID string
		if err := adm.LookupTeamID(opts[`team`][0], &teamID); err != nil {
			return err
		}
		req.Filter.Repository.TeamID = teamID
		validCondition = true
	}

	if _, ok := opts[`deleted`]; ok {
		if err := adm.ValidateBool(
			opts[`deleted`][0],
			&req.Filter.Repository.IsDeleted,
		); err != nil {
			return err
		}
		req.Filter.Repository.FilterOnIsDeleted = true
		validCondition = true
	}

	if _, ok := opts[`active`]; ok {
		if err := adm.ValidateBool(
			opts[`active`][0],
			&req.Filter.Repository.IsActive,
		); err != nil {
			return err
		}
		req.Filter.Repository.FilterOnIsActive = true
		validCondition = true
	}
	if !validCondition {
		return fmt.Errorf(`Syntax error: at least one search condition must be specified`)
	}

	return adm.Perform(`postbody`, `/search/repository/`, `list`, req, c)
}

// repositoryConfigTree function
// soma repository dumptree ${repository}
func repositoryConfigTree(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	repositoryID, err := adm.LookupRepoID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/tree", url.QueryEscape(
		repositoryID,
	))
	return adm.Perform(`get`, path, `tree`, nil, c)
}

// repositoryConfigPropertyCreateSystem function
// soma repository property create system ${system} on ${repository} view ${view} \
//      value ${value} [inheritance ${inherit}] [childrenonly ${child}]
func repositoryConfigPropertyCreateSystem(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeSystem, proto.EntityRepository)
}

// repositoryConfigPropertyCreateCustom function
// soma repository property create custom ${custom} on ${repository} view ${view} \
//      value ${value} [inheritance ${inherit}] [childrenonly ${child}]
func repositoryConfigPropertyCreateCustom(c *cli.Context) error {
	return variousPropertyCreate(c, proto.PropertyTypeCustom, proto.EntityRepository)
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
	return repositoryConfigPropertyDestroy(c, proto.PropertyTypeSystem)
}

// repositoryConfigPropertyDestroyCustom function
// soma repository property destroy custom ${custom} on ${repository} view ${view}
func repositoryConfigPropertyDestroyCustom(c *cli.Context) error {
	return repositoryConfigPropertyDestroy(c, proto.PropertyTypeCustom)
}

// repositoryConfigPropertyDestroyService function
// soma repository property destroy service ${service} on ${repository} view ${view}
func repositoryConfigPropertyDestroyService(c *cli.Context) error {
	return repositoryConfigPropertyDestroy(c, proto.PropertyTypeService)
}

// repositoryConfigPropertyDestroyOncall function
// soma repository property destroy oncall ${oncall} on ${repository} view ${view}
func repositoryConfigPropertyDestroyOncall(c *cli.Context) error {
	return repositoryConfigPropertyDestroy(c, proto.PropertyTypeOncall)
}

// repositoryConfigPropertyDestroy is the generic function for property destroy on
// repository objects
func repositoryConfigPropertyDestroy(c *cli.Context, propertyType string) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`on`, `view`}
	mandatoryOptions := []string{`on`, `view`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	repositoryID, err := adm.LookupRepoID(opts[`on`][0])
	if err != nil {
		return err
	}

	var teamID, serviceID string
	if err = adm.LookupTeamByRepo(repositoryID, &teamID); err != nil {
		return err
	}

	switch propertyType {
	case proto.PropertyTypeSystem:
		if err := adm.ValidateSystemProperty(
			c.Args().First()); err != nil {
			return err
		}
	case proto.PropertyTypeService:
		serviceID, err = adm.LookupServicePropertyID(c.Args().First(), teamID)
		if err != nil {
			return err
		}
	}

	var property string
	switch propertyType {
	case proto.PropertyTypeService:
		property = serviceID
	default:
		property = c.Args().First()
	}

	var sourceID string
	if err := adm.FindRepoPropSrcID(propertyType, property,
		opts[`view`][0], repositoryID, &sourceID); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/property/%s/%s",
		repositoryID, propertyType, sourceID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
