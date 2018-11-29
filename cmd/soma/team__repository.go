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

// repositoryDestroy function
// soma repository destroy ${repository} [from ${team}]
func repositoryDestroy(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`team`}
	mandatoryOptions := []string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	repositoryID, err := adm.LookupRepoID(c.Args().First())
	if err != nil {
		return err
	}

	var teamID string
	if _, ok := opts[`team`]; ok {
		if err := adm.LookupTeamID(opts[`team`][0], &teamID); err != nil {
			return err
		}
	} else {
		if err := adm.LookupTeamByRepo(repositoryID, &teamID); err != nil {
			return err
		}
	}

	path := fmt.Sprintf("/team/%s/repository/%s",
		url.QueryEscape(teamID),
		url.QueryEscape(repositoryID),
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// repositoryRename function
// soma repository rename ${repository} to ${newName} [from ${team}]
func repositoryRename(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`, `from`}
	mandatoryOptions := []string{`to`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewRepositoryRequest()
	req.Repository.Name = opts[`to`][0]

	repositoryID, err := adm.LookupRepoID(c.Args().First())
	if err != nil {
		return err
	}

	var teamID string
	if _, ok := opts[`team`]; ok {
		if err := adm.LookupTeamID(opts[`team`][0], &teamID); err != nil {
			return err
		}
	} else {
		if err := adm.LookupTeamByRepo(repositoryID, &teamID); err != nil {
			return err
		}
	}

	path := fmt.Sprintf(
		"/team/%s/repository/%s/name",
		url.QueryEscape(teamID),
		url.QueryEscape(repositoryID),
	)
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

// repositoryShow function
// soma repository show ${repository} [from ${team}]
func repositoryShow(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`team`}
	mandatoryOptions := []string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	repositoryID, err := adm.LookupRepoID(c.Args().First())
	if err != nil {
		return err
	}

	var teamID string
	if _, ok := opts[`team`]; ok {
		if err := adm.LookupTeamID(opts[`team`][0], &teamID); err != nil {
			return err
		}
	} else {
		if err := adm.LookupTeamByRepo(repositoryID, &teamID); err != nil {
			return err
		}
	}

	path := fmt.Sprintf("/team/%s/repository/%s",
		url.QueryEscape(teamID),
		url.QueryEscape(repositoryID),
	)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// TODO repositoryAudit

// repositoryRepossess function
// soma repository repossess ${repository} to ${newTeam} [from ${team}]
func repositoryRepossess(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`, `from`}
	mandatoryOptions := []string{`to`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var newOwnerTeamID string
	if err := adm.LookupTeamID(opts[`to`][0], &newOwnerTeamID); err != nil {
		return err
	}

	req := proto.NewRepositoryRequest()
	req.Repository.TeamID = newOwnerTeamID

	repositoryID, err := adm.LookupRepoID(c.Args().First())
	if err != nil {
		return err
	}

	var currentOwnerTeamID string
	if _, ok := opts[`team`]; ok {
		if err := adm.LookupTeamID(opts[`team`][0], &currentOwnerTeamID); err != nil {
			return err
		}
	} else {
		if err := adm.LookupTeamByRepo(repositoryID, &currentOwnerTeamID); err != nil {
			return err
		}
	}

	path := fmt.Sprintf(
		"/team/%s/repository/%s/owner",
		url.QueryEscape(currentOwnerTeamID),
		url.QueryEscape(repositoryID),
	)
	return adm.Perform(`patchbody`, path, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
