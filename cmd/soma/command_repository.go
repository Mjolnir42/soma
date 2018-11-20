/*-
 * Copyright (c) 2015-2018, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
)

func cmdRepositoryDestroy(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoID(c.Args().First())
	if err != nil {
		return err
	}
	var teamID string
	if err := adm.LookupTeamByRepo(id, &teamID); err != nil {
		return err
	}
	path := fmt.Sprintf("/team/%s/repository/%s", teamID, id)

	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdRepositoryRename(c *cli.Context) error {
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		[]string{`to`},
		[]string{`to`},
		c.Args().Tail()); err != nil {
		return err
	}
	id, err := adm.LookupRepoID(c.Args().First())
	if err != nil {
		return err
	}
	var teamID string
	if err := adm.LookupTeamByRepo(id, &teamID); err != nil {
		return err
	}
	path := fmt.Sprintf("/team/%s/repository/%s", teamID, id)

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.Name = opts[`to`][0]

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdRepositoryInstance(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/instance/", id)
	return adm.Perform(`get`, path, `list`, nil, c)
}

func cmdRepositoryTree(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/tree", id)
	return adm.Perform(`get`, path, `tree`, nil, c)
}

func cmdRepositorySystemPropertyAdd(c *cli.Context) error {
	return cmdRepositoryPropertyAdd(c, `system`)
}

func cmdRepositoryServicePropertyAdd(c *cli.Context) error {
	return cmdRepositoryPropertyAdd(c, `service`)
}

func cmdRepositoryOncallPropertyAdd(c *cli.Context) error {
	return cmdRepositoryPropertyAdd(c, `oncall`)
}

func cmdRepositoryCustomPropertyAdd(c *cli.Context) error {
	return cmdRepositoryPropertyAdd(c, `custom`)
}

func cmdRepositoryPropertyAdd(c *cli.Context, pType string) error {
	return cmdPropertyAdd(c, pType, `repository`)
}

func cmdRepositorySystemPropertyDelete(c *cli.Context) error {
	return cmdRepositoryPropertyDelete(c, `system`)
}

func cmdRepositoryServicePropertyDelete(c *cli.Context) error {
	return cmdRepositoryPropertyDelete(c, `service`)
}

func cmdRepositoryOncallPropertyDelete(c *cli.Context) error {
	return cmdRepositoryPropertyDelete(c, `oncall`)
}

func cmdRepositoryCustomPropertyDelete(c *cli.Context) error {
	return cmdRepositoryPropertyDelete(c, `custom`)
}

func cmdRepositoryPropertyDelete(c *cli.Context, pType string) error {
	unique := []string{`from`, `view`}
	required := []string{`from`, `view`}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}
	repositoryID, err := adm.LookupRepoID(opts[`from`][0])
	if err != nil {
		return err
	}

	if pType == `system` {
		if err := adm.ValidateSystemProperty(
			c.Args().First()); err != nil {
			return err
		}
	}
	var sourceID string
	if err := adm.FindRepoPropSrcID(pType, c.Args().First(),
		opts[`view`][0], repositoryID, &sourceID); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/property/%s/%s",
		repositoryID, pType, sourceID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
