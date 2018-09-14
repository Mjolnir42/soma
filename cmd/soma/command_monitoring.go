package main

import (
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerMonitoring(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// monitoringsystem
			{
				Name:  `monitoringsystem`,
				Usage: "SUBCOMMANDS for monitoring systems",
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a new monitoring system`,
						Action:       runtime(cmdMonitoringCreate),
						BashComplete: cmpl.MonitoringCreate,
					},
					{
						Name:   `remove`,
						Usage:  "Remove a monitoring system",
						Action: runtime(cmdMonitoringDelete),
					},
					{
						Name:   "list",
						Usage:  "List monitoring systems",
						Action: runtime(cmdMonitoringList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a monitoring system",
						Action: runtime(cmdMonitoringShow),
					},
					{
						Name:   "search",
						Usage:  "Lookup a monitoring system ID by name",
						Action: runtime(cmdMonitoringSearch),
					},
				},
			}, // end monitoringsystem
		}...,
	)
	return &app
}

func cmdMonitoringCreate(c *cli.Context) error {
	unique := []string{"mode", "contact", "team", "callback"}
	required := []string{"mode", "contact", "team"}
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

	req := proto.NewMonitoringRequest()
	req.Monitoring.Name = c.Args().First()
	req.Monitoring.Mode = opts["mode"][0]
	var err error
	if req.Monitoring.Contact, err = adm.LookupUserID(
		opts[`contact`][0]); err != nil {
		return err
	}
	req.Monitoring.TeamID, err = adm.LookupTeamID(opts[`team`][0])
	if err != nil {
		return err
	}
	if strings.Contains(req.Monitoring.Name, `.`) {
		return fmt.Errorf(
			`Monitoring system names must not contain` +
				` the character '.'`)
	}

	// optional arguments
	if _, ok := opts["callback"]; ok {
		req.Monitoring.Callback = opts["callback"][0]
	}

	return adm.Perform(`postbody`, `/monitoringsystem/`, `command`, req, c)
}

func cmdMonitoringDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	monitoringID, err := adm.LookupMonitoringID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/monitoringsystem/%s", monitoringID)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdMonitoringList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/monitoringsystem/`, `list`, nil, c)
}

func cmdMonitoringShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	monitoringID, err := adm.LookupMonitoringID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/monitoringsystem/%s", monitoringID)
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdMonitoringSearch(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	req := proto.NewMonitoringFilter()
	req.Filter.Monitoring.Name = c.Args().First()

	return adm.Perform(`postbody`, `/search/monitoringsystem/`, `list`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix