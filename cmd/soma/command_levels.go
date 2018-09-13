package main

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerLevels(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// levels
			{
				Name:  "levels",
				Usage: "SUBCOMMANDS for notification levels",
				Subcommands: []cli.Command{
					{
						Name:         "create",
						Usage:        "Create a new notification level",
						Action:       runtime(cmdLevelCreate),
						BashComplete: cmpl.LevelCreate,
					},
					{
						Name:   "delete",
						Usage:  "Delete a notification level",
						Action: runtime(cmdLevelDelete),
					},
					{
						Name:   "list",
						Usage:  "List notification levels",
						Action: runtime(cmdLevelList),
					},
					{
						Name:   "show",
						Usage:  "Show details about a notification level",
						Action: runtime(cmdLevelShow),
					},
					{
						Name:   `search`,
						Usage:  `Lookup a notification level by name or shortname`,
						Action: runtime(cmdLevelSearch),
					},
				},
			}, // end levels
		}...,
	)
	return &app
}

func cmdLevelCreate(c *cli.Context) error {
	uniqKeys := []string{"shortname", "numeric"}
	opts := map[string][]string{}

	if err := adm.ParseVariadicArguments(
		opts,
		[]string{},
		uniqKeys,
		uniqKeys,
		c.Args().Tail()); err != nil {
		return err
	}

	req := proto.Request{}
	req.Level = &proto.Level{}
	req.Level.Name = c.Args().First()
	req.Level.ShortName = opts["shortname"][0]

	var l uint64
	if err := adm.ValidateLBoundUint64(opts["numeric"][0],
		&l, 0); err != nil {
		return err
	}
	req.Level.Numeric = uint16(l)

	return adm.Perform(`postbody`, `/level/`, `command`, req, c)
}

func cmdLevelDelete(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/level/%s", c.Args().First())
	return adm.Perform(`delete`, path, `command`, nil, c)
}

func cmdLevelList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/level/`, `list`, nil, c)
}

func cmdLevelShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/level/%s", c.Args().First())
	return adm.Perform(`get`, path, `show`, nil, c)
}

func cmdLevelSearch(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewLevelFilter()
	req.Filter.Level.Name = c.Args().First()
	req.Filter.Level.ShortName = c.Args().First()

	return adm.Perform(`postbody`, `/search/level/`, `list`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
