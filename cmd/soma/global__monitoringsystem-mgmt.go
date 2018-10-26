/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
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

func registerMonitoringMgmt(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `monitoringsystem-mgmt`,
				Usage:       `SUBCOMMANDS for monitoring system management`,
				Description: help.Text(`monitoringsystem-mgmt::`),
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a new monitoring system`,
						Description:  help.Text(`monitoringsystem-mgmt::add`),
						Action:       runtime(monitoringMgmtAdd),
						BashComplete: cmpl.MonitoringMgmtAdd,
					},
					{
						Name:        `remove`,
						Usage:       `Remove a monitoring system`,
						Description: help.Text(`monitoringsystem-mgmt::remove`),
						Action:      runtime(monitoringMgmtRemove),
					},
					{
						Name:        `list`,
						Usage:       `List monitoring systems`,
						Description: help.Text(`monitoringsystem::list`),
						Action:      runtime(monitoringList),
					},
					{
						Name:        `show`,
						Usage:       `Show details about a monitoring system`,
						Description: help.Text(`monitoringsystem::show`),
						Action:      runtime(monitoringShow),
					},
					{
						Name:        `search`,
						Usage:       `Lookup a monitoring system ID by name`,
						Description: help.Text(`monitoringsystem::search`),
						Action:      runtime(monitoringSearch),
					},
				},
			},
		}...,
	)
	return &app
}

// monitoringMgmtAdd function
// soma monitoringsystem-mgmt add ${name} mode ${mode} contact ${user} \
//      team ${team} [ callback ${callback} ]
func monitoringMgmtAdd(c *cli.Context) error {
	var err error
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`mode`, `contact`, `team`, `callback`}
	mandatoryOptions := []string{`mode`, `contact`, `team`}

	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if err = adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	if err = adm.ValidateNoColon(c.Args().First()); err != nil {
		return err
	}

	if err = adm.ValidateNoDot(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewMonitoringRequest()
	req.Monitoring.Name = c.Args().First()
	req.Monitoring.Mode = opts[`mode`][0]

	if req.Monitoring.Contact, err = adm.LookupUserID(
		opts[`contact`][0],
	); err != nil {
		return err
	}

	if req.Monitoring.TeamID, err = adm.LookupTeamID(
		opts[`team`][0],
	); err != nil {
		return err
	}

	// optional arguments
	if _, ok := opts[`callback`]; ok {
		req.Monitoring.Callback = opts[`callback`][0]
	}

	return adm.Perform(`postbody`, `/monitoringsystem/`, `command`, req, c)
}

// monitoringMgmtRemove function
// soma monitoringsystem-mgmt remove ${name}
func monitoringMgmtRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	monitoringID, err := adm.LookupMonitoringID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf(
		"/monitoringsystem/%s",
		url.QueryEscape(monitoringID),
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
