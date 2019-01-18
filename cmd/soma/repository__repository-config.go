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

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
