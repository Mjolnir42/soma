package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerModes(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// modes
			{
				Name:  "modes",
				Usage: "SUBCOMMANDS for monitoring system modes",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Create a new monitoring system mode",
						Action: runtime(cmdModeCreate),
					},
					{
						Name:   "delete",
						Usage:  "Delete a monitoring system mode",
						Action: runtime(cmdModeDelete),
					},
					{
						Name:   "list",
						Usage:  "List monitoring system modes",
						Action: runtime(cmdModeList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a monitoring mode",
						Action: runtime(cmdModeShow),
					},
				},
			}, // end modes
		}...,
	)
	return &app
}

func cmdModeCreate(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.Request{}
	req.Mode = &proto.Mode{}
	req.Mode.Mode = c.Args().First()

	return adm.Perform(`postbody`, `/mode/`, `command`, req, c)
}

func cmdModeDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/mode/%s", c.Args().First())
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdModeList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/mode/`, `list`, nil, c)
}

func cmdModeShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/mode/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
