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
	"os"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerProperty(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			{
				Name:        `property-mgmt`,
				Usage:       `SUBCOMMANDS for property management`,
				Description: help.Text(`property-mgmt::`),
				Subcommands: []cli.Command{
					{
						Name:        `native`,
						Usage:       `SUBCOMMANDS for native introspection property management`,
						Description: help.Text(`property-native::`),
						Subcommands: []cli.Command{
							{
								Name:        `add`,
								Usage:       `Add a new global native introspection property`,
								Description: help.Text(`property-native::add`),
								Action:      runtime(propertyMgmtNativeAdd),
							},
							{
								Name:        `remove`,
								Usage:       `Remove a native introspection property`,
								Description: help.Text(`property-native::remove`),
								Action:      runtime(propertyMgmtNativeRemove),
							},
							{
								Name:        `show`,
								Usage:       `Show details about a native introspection property`,
								Description: help.Text(`property-native::show`),
								Action:      runtime(propertyMgmtNativeShow),
							},
							{
								Name:        `list`,
								Usage:       `List all native introspection properties`,
								Description: help.Text(`property-native::list`),
								Action:      runtime(propertyMgmtNativeList),
							},
						},
					},
					{
						Name:        `system`,
						Usage:       `SUBCOMMANDS for global system property management`,
						Description: help.Text(`property-system::`),
						Subcommands: []cli.Command{
							{
								Name:        `add`,
								Usage:       `Add a new global system property`,
								Description: help.Text(`property-system::add`),
								Action:      runtime(propertyMgmtSystemAdd),
							},
							{
								Name:        `remove`,
								Usage:       `Remove a global system property`,
								Description: help.Text(`property-system::remove`),
								Action:      runtime(propertyMgmtSystemRemove),
							},
							{
								Name:        `show`,
								Usage:       `Show details about a global system property`,
								Description: help.Text(`property-system::show`),
								Action:      runtime(propertyMgmtSystemShow),
							},
							{
								Name:        `list`,
								Usage:       `List all global system properties`,
								Description: help.Text(`property-system::list`),
								Action:      runtime(propertyMgmtSystemList),
							},
						},
					},
					{
						Name:        `template`,
						Usage:       `SUBCOMMANDS for global service property templates`,
						Description: help.Text(`property-template::`),
						Subcommands: []cli.Command{
							{
								Name:         `add`,
								Usage:        `Add a new global service template`,
								Description:  help.Text(`property-template::add`),
								Action:       runtime(propertyMgmtTemplateAdd),
								BashComplete: comptime(bashCompServiceAdd),
							},
							{
								Name:        `remove`,
								Usage:       `Remove a global service property template`,
								Description: help.Text(`property-template::remove`),
								Action:      runtime(propertyMgmtTemplateRemove),
							},
							{
								Name:        `show`,
								Usage:       `Show details about a service property template`,
								Description: help.Text(`property-template::show`),
								Action:      runtime(propertyMgmtTemplateShow),
							},
							{
								Name:        `list`,
								Usage:       `List all service property templates`,
								Description: help.Text(`property-template::list`),
								Action:      runtime(propertyMgmtTemplateList),
							},
						},
					},
				},
			},
		}...,
	)
	return &app
}

// SYSTEM PROPERTIES

// propertyMgmtSystemAdd function
// soma property-mgmt system add ${property}
func propertyMgmtSystemAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewSystemPropertyRequest()
	req.Property.System.Name = c.Args().First()

	path := fmt.Sprintf("/property-mgmt/%s/", proto.PropertyTypeSystem)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// propertyMgmtSystemRemove function
// soma property-mgmt system remove ${property}
func propertyMgmtSystemRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/property-mgmt/%s/%s",
		proto.PropertyTypeSystem,
		esc,
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// propertyMgmtSystemShow function
// soma property-mgmt system show ${property}
func propertyMgmtSystemShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/property-mgmt/%s/%s",
		proto.PropertyTypeSystem,
		esc,
	)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// propertyMgmtSystemList function
// soma property-mgmt system list
func propertyMgmtSystemList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/property-mgmt/%s/", proto.PropertyTypeSystem)
	return adm.Perform(`get`, path, `list`, nil, c)
}

// NATIVE PROPERTIES

// propertyMgmtNativeAdd function
// soma property-mgmt native add ${property}
func propertyMgmtNativeAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	req := proto.NewNativePropertyRequest()
	req.Property.Native.Name = c.Args().First()

	path := fmt.Sprintf("/property-mgmt/%s/",
		url.QueryEscape(proto.PropertyTypeNative),
	)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// propertyMgmtNativeRemove function
// soma property-mgmt native remove ${property}
func propertyMgmtNativeRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/property-mgmt/%s/%s",
		url.QueryEscape(proto.PropertyTypeNative),
		url.QueryEscape(c.Args().First()),
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// propertyMgmtNativeShow function
// soma property-mgmt native show ${property}
func propertyMgmtNativeShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	path := fmt.Sprintf("/property-mgmt/%s/%s",
		url.QueryEscape(proto.PropertyTypeNative),
		url.QueryEscape(c.Args().First()),
	)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// propertyMgmtNativeList function
// soma property-mgmt native list
func propertyMgmtNativeList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/property-mgmt/%s/",
		url.QueryEscape(proto.PropertyTypeNative),
	)
	return adm.Perform(`get`, path, `list`, nil, c)
}

// TEMPLATE PROPERTIES

// propertyMgmtTemplateAdd function
func propertyMgmtTemplateAdd(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{}
	mandatoryOptions := []string{}

	// sort attributes based on their cardinality so we can use them
	// for command line parsing
	for _, attr := range attributeFetch() {
		switch attr.Cardinality {
		case `once`:
			uniqueOptions = append(uniqueOptions, attr.Name)
		case `multi`:
			multipleAllowed = append(multipleAllowed, attr.Name)
		default:
			return fmt.Errorf("Unknown attribute cardinality: %s",
				attr.Cardinality)
		}
	}

	// check deferred errors
	if err := popError(); err != nil {
		return err
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

	// construct request body
	req := proto.NewTemplatePropertyRequest()
	req.Property.Service.Name = c.Args().First()
	if err := adm.ValidateRuneCount(
		req.Property.Service.Name, 128); err != nil {
		return err
	}

	// fill attributes into request body
	for oName := range opts {
		for _, oVal := range opts[oName] {
			if err := adm.ValidateRuneCount(oName, 128); err != nil {
				return err
			}
			if err := adm.ValidateRuneCount(oVal, 128); err != nil {
				return err
			}
			req.Property.Service.Attributes = append(
				req.Property.Service.Attributes,
				proto.ServiceAttribute{
					Name:  oName,
					Value: oVal,
				},
			)
		}
	}

	path := fmt.Sprintf("/property-mgmt/%s/",
		url.QueryEscape(proto.PropertyTypeTemplate),
	)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// propertyMgmtTemplateRemove function
func propertyMgmtTemplateRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	propertyID, err := adm.LookupTemplatePropertyID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/property-mgmt/%s/%s",
		url.QueryEscape(proto.PropertyTypeTemplate),
		url.QueryEscape(propertyID),
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// propertyMgmtTemplateShow function
// soma property-mgmt template show ${property}
func propertyMgmtTemplateShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	propertyID, err := adm.LookupTemplatePropertyID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/property-mgmt/%s/%s",
		url.QueryEscape(proto.PropertyTypeTemplate),
		url.QueryEscape(propertyID),
	)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// propertyMgmtTemplateList function
// soma property-mgmt template list
func propertyMgmtTemplateList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	path := fmt.Sprintf("/property-mgmt/%s/",
		url.QueryEscape(proto.PropertyTypeTemplate),
	)
	return adm.Perform(`get`, path, `list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
