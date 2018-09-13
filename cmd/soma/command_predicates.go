package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerPredicates(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// predicates
			{
				Name:  "predicates",
				Usage: "SUBCOMMANDS for threshold predicates",
				Subcommands: []cli.Command{
					{
						Name:   "create",
						Usage:  "Add a predicate",
						Action: runtime(cmdPredicateCreate),
					},
					{
						Name:   "delete",
						Usage:  "Delete a predicate",
						Action: runtime(cmdPredicateDelete),
					},
					{
						Name:   "list",
						Usage:  "List predicates",
						Action: runtime(cmdPredicateList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a predicate",
						Action: runtime(cmdPredicateShow),
					},
				},
			}, // end predicates
		}...,
	)
	return &app
}

func cmdPredicateCreate(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.Request{}
	req.Predicate = &proto.Predicate{}
	req.Predicate.Symbol = c.Args().First()

	return adm.Perform(`postbody`, `/predicate/`, `command`, req, c)
}

func cmdPredicateDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/predicate/%s", c.Args().First())
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdPredicateList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/predicate/`, `list`, nil, c)
}

func cmdPredicateShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/predicate/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
