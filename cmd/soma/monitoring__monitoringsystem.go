/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

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
