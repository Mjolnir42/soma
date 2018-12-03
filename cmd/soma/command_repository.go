/*-
 * Copyright (c) 2015-2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/soma/cmd/soma"

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
)

func cmdRepositoryInstance(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}
	id, err := adm.LookupRepoID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/instance/", id)
	return adm.Perform(`get`, path, `list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
