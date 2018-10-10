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
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/internal/cmpl"
	"github.com/mjolnir42/soma/internal/help"
	"github.com/mjolnir42/soma/lib/proto"
)

func registerMetrics(app cli.App) *cli.App {
	app.Commands = append(app.Commands,
		[]cli.Command{
			// metrics
			{
				Name:        `metric`,
				Usage:       `SUBCOMMANDS for metrics`,
				Description: help.Text(`metric::`),
				Subcommands: []cli.Command{
					{
						Name:         `add`,
						Usage:        `Add a new metric`,
						Description:  help.Text(`metric::add`),
						Action:       runtime(cmdMetricAdd),
						BashComplete: cmpl.MetricAdd,
					},
					{
						Name:        `remove`,
						Usage:       `Remove a metric`,
						Description: help.Text(`metric::remove`),
						Action:      runtime(cmdMetricRemove),
					},
					{
						Name:        `list`,
						Usage:       `List metrics`,
						Description: help.Text(`metric::list`),
						Action:      runtime(cmdMetricList),
					},
					{
						Name:        `show`,
						Usage:       `Show details about a metric`,
						Description: help.Text(`metric::show`),
						Action:      runtime(cmdMetricShow),
					},
				},
			},
		}...,
	)
	return &app
}

// cmdMetricAdd function
// soma metric add ${metric} unit ${unit} description ${text}
// [package ${provider}::${package}, ...]
func cmdMetricAdd(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{`package`}
	uniqueOptions := []string{`unit`, `description`}
	mandatoryOptions := []string{`unit`, `description`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	req := proto.NewMetricRequest()
	req.Metric.Path = c.Args().First()
	req.Metric.Unit = opts[`unit`][0]
	req.Metric.Description = opts[`description`][0]

	if err := adm.ValidateUnit(req.Metric.Unit); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(req.Metric.Path); err != nil {
		return err
	}

	pkgs := []proto.MetricPackage{}
	if _, ok := opts[`package`]; ok {
		for _, p := range opts[`package`] {
			split := strings.SplitN(p, `::`, 2)
			if len(split) != 2 { // coult not split
				return fmt.Errorf(
					"Package spec error, contains no :: divisor: %s", p)
			}

			if err := adm.ValidateProvider(split[0]); err != nil {
				return err
			}

			pkgs = append(pkgs, proto.MetricPackage{
				Provider: split[0],
				Name:     split[1],
			})
		}
		req.Metric.Packages = &pkgs
	}

	return adm.Perform(`postbody`, `/metric/`, `command`, req, c)
}

// cmdMetricRemove function
// soma metric remove ${metric}
func cmdMetricRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/metric/%s", esc)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// cmdMetricList function
// soma metric list
func cmdMetricList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/metric/`, `list`, nil, c)
}

// cmdMetricShow function
// soma metric show ${metric}
func cmdMetricShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	if err := adm.ValidateNoSlash(c.Args().First()); err != nil {
		return err
	}

	esc := url.QueryEscape(c.Args().First())
	path := fmt.Sprintf("/metric/%s", esc)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
