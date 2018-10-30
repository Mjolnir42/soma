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
	"net/url"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerTeams(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// teams
			{
				Name:        `team-mgmt`,
				Usage:       `SUBCOMMANDS for team management`,
				Description: help.Text(`team-mgmt::`),
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Register a new team`,
						Description:  help.Text(`team-mgmt::add`),
						Action:       runtime(teamMgmtAdd),
						BashComplete: cmpl.TeamCreate,
					},
					{
						Name:         `update`,
						Usage:        `Update information for a team`,
						Description:  help.Text(`team-mgmt::update`),
						Action:       runtime(teamMgmtUpdate),
						BashComplete: cmpl.TeamUpdate,
					},
					{
						Name:        `remove`,
						Usage:       `Remove an existing team`,
						Description: help.Text(`team-mgmt::remove`),
						Action:      runtime(teamMgmtRemove),
					},
					{
						Name:        `show`,
						Usage:       `Show information about a team`,
						Description: help.Text(`team-mgmt::show`),
						Action:      runtime(teamMgmtShow),
					},
					{
						Name:        `list`,
						Usage:       `List all teams`,
						Description: help.Text(`team-mgmt::list`),
						Action:      runtime(teamMgmtList),
					},
					{
						Name:        `sync`,
						Usage:       `Export a list of all teams suitable for sync`,
						Description: help.Text(`team-mgmt::sync`),
						Action:      runtime(teamMgmtSync),
					},
				},
			},
		}...,
	)
	return &app
}

// teamMgmtAdd function
// soma team-mgmt add ${team} ldap ${ldapID} [system ${bool}]
func teamMgmtAdd(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`ldap`, `system`}
	mandatoryOptions := []string{`ldap`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewTeamRequest()
	req.Team.Name = c.Args().First()
	req.Team.LdapID = opts[`ldap`][0]
	if len(opts[`system`]) > 0 {
		if err := adm.ValidateBool(opts[`system`][0],
			&req.Team.IsSystem); err != nil {
			return fmt.Errorf("Argument to system parameter must"+
				" be boolean: %s", err.Error())
		}
	}

	return adm.Perform(`postbody`, `/team/`, `command`, req, c)
}

// teamMgmtUpdate function
// soma team-mgmt update ${team} name ${name} ldap ${ldapID} [system ${bool}]
func teamMgmtUpdate(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`name`, `ldap`, `system`}
	mandatoryOptions := []string{`name`, `ldap`}

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
	if err := adm.LookupTeamID(c.Args().First(), &teamID); err != nil {
		return err
	}

	req := proto.NewTeamRequest()
	req.Team.ID = teamID
	req.Team.Name = opts[`name`][0]
	req.Team.LdapID = opts[`ldap`][0]
	if len(opts[`system`]) > 0 {
		if err := adm.ValidateBool(opts["system"][0],
			&req.Team.IsSystem); err != nil {
			return fmt.Errorf("Argument to system parameter must"+
				" be boolean: %s", err.Error())
		}
	}
	path := fmt.Sprintf("/team/%s",
		url.QueryEscape(teamID),
	)
	return adm.Perform(`putbody`, path, `command`, req, c)
}

// teamMgmtRemove function
// soma team-mgmt remove ${team}
func teamMgmtRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	var teamID string
	if err := adm.LookupTeamID(c.Args().First(), &teamID); err != nil {
		return err
	}

	path := fmt.Sprintf("/team/%s",
		url.QueryEscape(teamID),
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// teamMgmtList function
// soma team-mgmt list
func teamMgmtList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/team/`, `list`, nil, c)
}

// teamMgmtSync function
// soma team-mgmt sync
func teamMgmtSync(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/sync/team/`, `list`, nil, c)
}

// teamMgmtShow function
// soma team-mgmt show ${team}
func teamMgmtShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	var teamID string
	if err := adm.LookupTeamID(c.Args().First(), &teamID); err != nil {
		return err
	}

	path := fmt.Sprintf("/team/%s",
		url.QueryEscape(teamID),
	)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
