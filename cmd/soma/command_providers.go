package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerProviders(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// providers
			{
				Name:  "providers",
				Usage: "SUBCOMMANDS for metric providers",
				Subcommands: []cli.Command{
					{
						Name:   "add",
						Usage:  "Add a new metric provider",
						Action: runtime(cmdProviderCreate),
					},
					{
						Name:   "remove",
						Usage:  "Remove a metric provider",
						Action: runtime(cmdProviderDelete),
					},
					{
						Name:   "list",
						Usage:  "List metric providers",
						Action: runtime(cmdProviderList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a metric provider",
						Action: runtime(cmdProviderShow),
					},
				},
			}, // end providers
		}...,
	)
	return &app
}

func cmdProviderCreate(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.Request{}
	req.Provider = &proto.Provider{}
	req.Provider.Name = c.Args().First()

	return adm.Perform(`postbody`, `/provider/`, `command`, req, c)
}

func cmdProviderDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/provider/%s", c.Args().First())
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdProviderList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/provider/`, `list`, nil, c)
}

func cmdProviderShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/provider/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
