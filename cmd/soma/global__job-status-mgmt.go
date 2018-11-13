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

// jobStatusMgmtAdd function
// soma job status-mgmt add ${status}
func jobStatusMgmtAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewJobStatusRequest()
	req.JobStatus.Name = c.Args().First()
	if err := adm.ValidateNotUUID(req.JobStatus.Name); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/job/status-mgmt/`, `command`, req, c)
}

// jobStatusMgmtRemove function
// soma job status-mgmt remove ${status}
func jobStatusMgmtRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	jobStatusID, err := adm.LookupJobStatusID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/job/status-mgmt/%s", url.QueryEscape(jobStatusID))
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// jobStatusMgmtShow function
// soma job status-mgmt show ${status}
func jobStatusMgmtShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	jobStatusID, err := adm.LookupJobStatusID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/job/status-mgmt/%s", url.QueryEscape(jobStatusID))
	return adm.Perform(`get`, path, `command`, nil, c)
}

// jobStatusMgmtList function
// soma job status-mgmt list
func jobStatusMgmtList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/job/status-mgmt/`, `list`, nil, c)
}

// jobStatusMgmtSearch function
// soma job status-mgmt search [id ${uuid}] [name ${status}]
func jobStatusMgmtSearch(c *cli.Context) error {
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
	req := proto.NewJobStatusFilter()
	if _, ok := opts[`id`]; ok {
		if err := adm.ValidateUUID(opts[`id`][0]); err != nil {
			return err
		}
		req.Filter.JobStatus.ID = opts[`id`][0]
		valid = true
	}
	if _, ok := opts[`name`]; ok {
		if err := adm.ValidateNotUUID(opts[`name`][0]); err != nil {
			return err
		}
		req.Filter.JobStatus.Name = opts[`name`][0]
		valid = true
	}
	if !valid {
		return fmt.Errorf(`Syntax error, must specify 'id', 'name' or both.`)
	}

	return adm.Perform(`postbody`, `/search/jobStatus/`, `show`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
