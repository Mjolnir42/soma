/*-
 * Copyright (c) 2016-2018, Jörg Pernfuß
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

func registerCapability(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `capability`,
				Usage:       `SUBCOMMANDS for monitoring system capability declarations`,
				Description: help.Text(`capability::`),
				Subcommands: []cli.Command{
					{
						Name:         `declare`,
						Usage:        `Declare a new monitoring system capability`,
						Description:  help.Text(`capability::declare`),
						Action:       runtime(capabilityDeclare),
						BashComplete: cmpl.CapabilityDeclare,
					},
					{
						Name:        `revoke`,
						Usage:       `Revoke a monitoring system capability`,
						Description: help.Text(`capability::revoke`),
						Action:      runtime(capabilityRevoke),
					},
					{
						Name:        `list`,
						Usage:       `List monitoring system capabilities`,
						Description: help.Text(`capability::list`),
						Action:      runtime(capabilityList),
					},
					{
						Name:        `show`,
						Usage:       `Show details about a monitoring system capability`,
						Description: help.Text(`capability::show`),
						Action:      runtime(capabilityShow),
					},
				},
			},
		}...,
	)
	return &app
}

// capabilityDeclare function
// soma capability declare ${monitoring} view ${view} \
//      metric ${metric} thresholds ${num}
func capabilityDeclare(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{
		`metric`,
		`view`,
		`thresholds`,
	}
	mandatoryOptions := []string{
		`metric`,
		`view`,
		`thresholds`,
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

	var thresholds uint64
	var err error

	if err = adm.ValidateLBoundUint64(
		opts[`thresholds`][0],
		&thresholds, 1,
	); err != nil {
		return err
	}
	if err = adm.ValidateNotUUID(c.Args().First()); err != nil {
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
	if err = adm.ValidateNoSlash(opts[`metric`][0]); err != nil {
		return err
	}

	req := proto.NewCapabilityRequest()
	req.Capability.Metric = opts[`metric`][0]
	req.Capability.View = opts[`view`][0]
	req.Capability.Thresholds = thresholds
	if req.Capability.MonitoringID, err = adm.LookupMonitoringID(
		c.Args().First()); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/capability/`, `command`, req, c)
}

// capabilityRevoke function
// soma capability revoke ${capability}
func capabilityRevoke(c *cli.Context) (err error) {
	if err = adm.VerifySingleArgument(c); err != nil {
		return err
	}

	var id, path string
	if id, err = adm.LookupCapabilityID(
		c.Args().First()); err != nil {
		return err
	}
	path = fmt.Sprintf("/capability/%s", url.QueryEscape(id))

	return adm.Perform(`delete`, path, `command`, nil, c)
}

// capabilityList function
// soma capability list
func capabilityList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/capability/`, `list`, nil, c)
}

// capabilityShow function
// soma capability show ${capability}
func capabilityShow(c *cli.Context) (err error) {
	if err = adm.VerifySingleArgument(c); err != nil {
		return err
	}

	var id, path string
	if id, err = adm.LookupCapabilityID(
		c.Args().First()); err != nil {
		return err
	}
	path = fmt.Sprintf("/capability/%s", url.QueryEscape(id))

	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
