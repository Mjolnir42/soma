/*-
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

func registerAttributes(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// attributes
			{
				Name:        `attribute`,
				Usage:       `SUBCOMMANDS for service attributes`,
				Description: help.Text(`attribute::`),
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a new service attribute`,
						Description:  help.Text(`attribute::add`),
						Action:       runtime(cmdAttributeAdd),
						BashComplete: cmpl.AttributeAdd,
					},
					{
						Name:        `remove`,
						Usage:       `Remove a service attribute`,
						Description: help.Text(`attribute::remove`),
						Action:      runtime(cmdAttributeRemove),
					},
					{
						Name:        `list`,
						Usage:       `List service attributes`,
						Description: help.Text(`attribute::list`),
						Action:      runtime(cmdAttributeList),
					},
					{
						Name:        `show`,
						Usage:       `Show details about a service attribute`,
						Description: help.Text(`attribute::show`),
						Action:      runtime(cmdAttributeShow),
					},
				},
			}, // end attributes
		}...,
	)
	return &app
}

// cmdAttributeAdd function
// soma attribute add ${attribute} cardinality ${cardinality}
func cmdAttributeAdd(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`cardinality`}
	mandatoryOptions := []string{`cardinality`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	switch opts[`cardinality`][0] {
	case `once`, `multi`:
	default:
		return fmt.Errorf("Illegal value for cardinality: %s."+
			" Accepted: once, multi", opts[`cardinality`][0])
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewAttributeRequest()
	req.Attribute.Name = c.Args().First()
	req.Attribute.Cardinality = opts[`cardinality`][0]

	// check attribute length
	if err := adm.ValidateRuneCount(
		req.Attribute.Name,
		128,
	); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/attribute/`, `command`, req, c)
}

// cmdAttributeRemove function
// soma attribute remove ${attribute}
func cmdAttributeRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/attribute/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdAttributeList function
// soma attribute list
func cmdAttributeList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/attribute/`, `list`, nil, c)
}

// cmdAttributeShow function
// soma attribute show ${attribute}
func cmdAttributeShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/attribute/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
