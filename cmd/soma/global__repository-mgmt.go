/*-
 * Copyright (c) 2015-2018, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerRepository(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// repository
			{
				Name:        `repository`,
				Usage:       `SUBCOMMANDS for repository management`,
				Description: help.Text(`repository::`),
				Subcommands: []cli.Command{
					{
						Name:         `create`,
						Usage:        `Create a new repository for a team`,
						Description:  help.Text(`repository-mgmt::create`),
						Action:       runtime(repositoryMgmtCreate),
						BashComplete: cmpl.Team,
					},
					{
						Name:        `list`,
						Usage:       `List existing repositories`,
						Description: help.Text(`repository-config::list`),
						Action:      runtime(repositoryConfigList),
					},
					{
						Name:         `show`,
						Usage:        `Show information about a specific repository`,
						Description:  help.Text(`repository::show`),
						Action:       runtime(repositoryShow),
						BashComplete: cmpl.Team,
					},
					{
						Name:         `search`,
						Usage:        `Search for repositories matching specific conditions`,
						Description:  help.Text(`repository-config::search`),
						Action:       runtime(repositoryConfigSearch),
						BashComplete: cmpl.RepositoryConfigSearch,
					},
				},
			},
		}...,
	)
	return &app
}

// repositoryMgmtCreate function
// soma repository create ${repository} team ${team}
func repositoryMgmtCreate(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`team`}
	mandatoryOptions := []string{`team`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var teamID string
	if err := adm.LookupTeamID(opts[`team`][0], &teamID); err != nil {
		return err
	}

	req := proto.NewRepositoryRequest()
	req.Repository.Name = c.Args().First()
	req.Repository.TeamID = teamID

	if err := adm.ValidateRuneCountRange(req.Repository.Name,
		4, 128); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/repository/`, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
