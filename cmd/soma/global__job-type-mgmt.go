/*-
 * Copyright (c) 2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
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
	"github.com/mjolnir42/soma/lib/proto"
)

// jobTypeMgmtAdd function
// soma job type-mgmt add ${type}
func jobTypeMgmtAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewJobTypeRequest()
	req.JobType.Name = c.Args().First()
	if err := adm.ValidateNotUUID(req.JobType.Name); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/job/type-mgmt/`, `command`, req, c)
}

// jobTypeMgmtRemove function
// soma job type-mgmt remove ${type}
func jobTypeMgmtRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	jobTypeID, err := adm.LookupJobTypeID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/job/type-mgmt/%s", url.QueryEscape(jobTypeID))
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// jobTypeMgmtShow function
// soma job type-mgmt show ${type}
func jobTypeMgmtShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	jobTypeID, err := adm.LookupJobTypeID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/job/type-mgmt/%s", url.QueryEscape(jobTypeID))
	return adm.Perform(`get`, path, `command`, nil, c)
}

// jobTypeMgmtList function
// soma job type-mgmt list
func jobTypeMgmtList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/job/type-mgmt/`, `list`, nil, c)
}

// jobTypeMgmtSearch function
// soma job type-mgmt search [id ${uuid}] [name ${type}]
func jobTypeMgmtSearch(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`id`, `name`}
	mandatoryOptions := []string{}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		adm.AllArguments(c),
	); err != nil {
		return err
	}

	valid := false
	req := proto.NewJobTypeFilter()
	if _, ok := opts[`id`]; ok {
		if err := adm.ValidateUUID(opts[`id`][0]); err != nil {
			return err
		}
		req.Filter.JobType.ID = opts[`id`][0]
		valid = true
	}
	if _, ok := opts[`name`]; ok {
		if err := adm.ValidateNotUUID(opts[`name`][0]); err != nil {
			return err
		}
		req.Filter.JobType.Name = opts[`name`][0]
		valid = true
	}
	if !valid {
		return fmt.Errorf(`Syntax error, must specify 'id', 'name' or both.`)
	}

	return adm.Perform(`postbody`, `/search/jobType/`, `show`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
