/*-
 * Copyright (c) 2015-2019, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2018-2019, 1&1 IONOS SE
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

func registerServer(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `server`,
				Usage:       `SUBCOMMANDS for physical server management`,
				Description: help.Text(`server::`),
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a new physical server`,
						Description:  help.Text(`server::add`),
						Action:       runtime(serverAdd),
						BashComplete: cmpl.ServerAdd,
					},
					{
						Name:         `remove`,
						Usage:        `Remove a physical server`,
						Description:  help.Text(`server::remove`),
						Action:       runtime(serverRemove),
						BashComplete: comptime(bashCompServer),
					},
					{
						Name:         `purge`,
						Usage:        `Purge a removed physical server`,
						Description:  help.Text(`server::purge`),
						Action:       runtime(serverPurge),
						BashComplete: comptime(bashCompServer),
					},
					{
						Name:         `update`,
						Usage:        `Full update of server attributes (replace, not merge)`,
						Description:  help.Text(`server::update`),
						Action:       runtime(serverUpdate),
						BashComplete: cmpl.ServerUpdate,
					},
					{
						Name:         `list`,
						Usage:        `List all servers, see full description for possible filters`,
						Description:  help.Text(`server::list`),
						Action:       runtime(serverList),
						BashComplete: cmpl.None,
					},
					{
						Name:         `show`,
						Usage:        `Show details about a specific server`,
						Description:  help.Text(`server::show`),
						Action:       runtime(serverShow),
						BashComplete: comptime(bashCompServer),
					},
					{
						Name:         `sync`,
						Usage:        "Export a list of all servers suitable for sync",
						Description:  help.Text(`server::sync`),
						Action:       runtime(serverSync),
						BashComplete: cmpl.None,
					},
					{
						Name:         `null`,
						Usage:        `Bootstrap the null server`,
						Description:  help.Text(`server::null`),
						Action:       runtime(serverNull),
						BashComplete: cmpl.Datacenter,
					},
				},
			},
		}...,
	)
	return &app
}

// serverAdd function
// soma server add ${name} \
//      assetid ${num} \
//      datacenter ${locode} \
//      location ${loc} \
//      [online ${bool}]
func serverAdd(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{
		`assetid`,
		`datacenter`,
		`location`,
		`online`,
	}
	mandatoryOptions := []string{
		`assetid`,
		`datacenter`,
		`location`,
	}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewServerRequest()
	req.Server.Name = c.Args().First()
	req.Server.Datacenter = opts[`datacenter`][0]
	req.Server.Location = opts[`location`][0]

	if err := adm.ValidateNotUUID(req.Server.Name); err != nil {
		return err
	}
	if err := adm.ValidateLBoundUint64(opts[`assetid`][0],
		&req.Server.AssetID, 1); err != nil {
		return err
	}

	// optional argument: online, defaults to true
	req.Server.IsOnline = true
	if _, ok := opts[`online`]; ok {
		if err := adm.ValidateBool(opts[`online`][0],
			&req.Server.IsOnline); err != nil {
			return err
		}
	}

	return adm.Perform(`postbody`, `/server/`, `command`, req, c)
}

// serverRemove function
// soma server remove ${name}
func serverRemove(c *cli.Context) error {
	// check deferred errors
	if err := popError(); err != nil {
		return err
	}

	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	sid, err := adm.LookupServerID(c.Args().First())
	if err != nil {
		return err
	}

	// request must include a body, SOMA server
	// checks req.Flags.Purge
	req := proto.NewServerRequest()
	req.Server.ID = sid
	path := fmt.Sprintf("/server/%s", sid)
	return adm.Perform(`deletebody`, path, `command`, req, c)
}

// serverPurge function
// soma server purge ${name}
func serverPurge(c *cli.Context) error {
	// check deferred errors
	if err := popError(); err != nil {
		return err
	}

	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	sid, err := adm.LookupServerID(c.Args().First())
	if err != nil {
		return err
	}

	// request must include a body, SOMA server
	// checks req.Flags.Purge
	req := proto.NewServerRequest()
	req.Flags.Purge = true
	req.Server.ID = sid
	path := fmt.Sprintf("/server/%s", sid)
	return adm.Perform(`deletebody`, path, `command`, req, c)
}

// serverUpdate function
// soma server update ${serverID} \
//      name ${name} \
//      assetid ${assetID} \
//      datacenter ${locode} \
//      location ${loc} \
//      [online ${bool}] \
//      [deleted ${bool}]
func serverUpdate(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{
		`name`,
		`assetid`,
		`datacenter`,
		`location`,
		`online`,
		`deleted`,
	}
	mandatoryOptions := []string{
		`name`,
		`assetid`,
		`datacenter`,
		`location`,
	}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if !adm.IsUUID(c.Args().First()) {
		return fmt.Errorf("Server to update not referenced by"+
			" UUID: %s", c.Args().First())
	}

	req := proto.NewServerRequest()
	req.Server.ID = c.Args().First()
	req.Server.Name = opts[`name`][0]
	req.Server.Datacenter = opts[`datacenter`][0]
	req.Server.Location = opts[`location`][0]
	if err := adm.ValidateLBoundUint64(opts[`assetid`][0],
		&req.Server.AssetID, 1); err != nil {
		return err
	}
	// IsOnline defaults to true
	req.Server.IsDeleted = true
	if _, ok := opts[`online`]; ok {
		if err := adm.ValidateBool(opts[`online`][0],
			&req.Server.IsOnline); err != nil {
			return err
		}
	}
	// IsDeleted defaults to false
	req.Server.IsDeleted = false
	if _, ok := opts[`deleted`]; ok {
		if err := adm.ValidateBool(opts[`deleted`][0],
			&req.Server.IsDeleted); err != nil {
			return err
		}
	}
	if err := adm.ValidateNotUUID(req.Server.Name); err != nil {
		return err
	}

	path := fmt.Sprintf("/server/%s", c.Args().First())
	return adm.Perform(`putbody`, path, `command`, req, c)
}

// serverList function
// soma server list
func serverList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/server/`, `list`, nil, c)
}

// serverList function
// soma server sync
func serverSync(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/sync/server/`, `list`, nil, c)
}

// serverList function
// soma server show ${name}
func serverShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	// check deferred errors
	if err := popError(); err != nil {
		return err
	}

	serverID, err := adm.LookupServerID(c.Args().First())
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/server/%s",
		url.QueryEscape(serverID))
	return adm.Perform(`get`, path, `show`, nil, c)
}

// serverNull function
// soma server null datacenter ${locode}
func serverNull(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`datacenter`}
	mandatoryOptions := []string{`datacenter`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		adm.AllArguments(c),
	); err != nil {
		return err
	}

	req := proto.NewServerRequest()
	req.Server.ID = `00000000-0000-0000-0000-000000000000`
	req.Server.Datacenter = opts[`datacenter`][0]

	return adm.Perform(`postbody`, `/server/null`, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
