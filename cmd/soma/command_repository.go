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
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerRepository(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// repository
			{
				Name:  "repository",
				Usage: "SUBCOMMANDS for repository",
				Subcommands: []cli.Command{
					{
						Name:   "destroy",
						Usage:  "Destroy an existing repository",
						Action: runtime(cmdRepositoryDestroy),
					},
					{
						Name:         "rename",
						Usage:        "Rename an existing repository",
						Action:       runtime(cmdRepositoryRename),
						BashComplete: cmpl.To,
					},
					{
						Name:   "list",
						Usage:  "List all existing repositories",
						Action: runtime(cmdRepositoryList),
					},
					{
						Name:   "show",
						Usage:  "Show information about a specific repository",
						Action: runtime(cmdRepositoryShow),
					},
					{
						Name:   `tree`,
						Usage:  `Display the repository as tree`,
						Action: runtime(cmdRepositoryTree),
					},
					{
						Name:   `instances`,
						Usage:  `List check instances for a repository`,
						Action: runtime(cmdRepositoryInstance),
					},
					{
						Name:  "property",
						Usage: "SUBCOMMANDS for properties",
						Subcommands: []cli.Command{
							{
								Name:  "add",
								Usage: "SUBCOMMANDS for property add",
								Subcommands: []cli.Command{
									{
										Name:         "system",
										Usage:        "Add a system property to a repository",
										Action:       runtime(cmdRepositorySystemPropertyAdd),
										BashComplete: cmpl.PropertyAddValue,
									},
									{
										Name:         "service",
										Usage:        "Add a service property to a repository",
										Action:       runtime(cmdRepositoryServicePropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         "oncall",
										Usage:        "Add an oncall property to a repository",
										Action:       runtime(cmdRepositoryOncallPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
									{
										Name:         "custom",
										Usage:        "Add a custom property to a repository",
										Action:       runtime(cmdRepositoryCustomPropertyAdd),
										BashComplete: cmpl.PropertyAdd,
									},
								},
							},
							{
								Name:  `delete`,
								Usage: `SUBCOMMANDS for property delete`,
								Subcommands: []cli.Command{
									{
										Name:         `system`,
										Usage:        `Delete a system property from a repository`,
										Action:       runtime(cmdRepositorySystemPropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `service`,
										Usage:        `Delete a service property from a repository`,
										Action:       runtime(cmdRepositoryServicePropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `oncall`,
										Usage:        `Delete an oncall property from a repository`,
										Action:       runtime(cmdRepositoryOncallPropertyDelete),
										BashComplete: cmpl.FromView,
									},
									{
										Name:         `custom`,
										Usage:        `Delete a custom property from a repository`,
										Action:       runtime(cmdRepositoryCustomPropertyDelete),
										BashComplete: cmpl.FromView,
									},
								},
							},
						},
					},
				},
			}, // end repository
		}...,
	)
	return &app
}

func cmdRepositoryDestroy(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoID(c.Args().First())
	if err != nil {
		return err
	}
	teamID, err := adm.LookupTeamByRepo(id)
	if err != nil {
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
	teamID, err := adm.LookupTeamByRepo(id)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/team/%s/repository/%s", teamID, id)

	var req proto.Request
	req.Repository = &proto.Repository{}
	req.Repository.Name = opts[`to`][0]

	return adm.Perform(`patchbody`, path, `command`, req, c)
}

func cmdRepositoryList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/repository/`, `list`, nil, c)
}

func cmdRepositoryShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s", id)
	return adm.Perform(`get`, path, `show`, nil, c)
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
