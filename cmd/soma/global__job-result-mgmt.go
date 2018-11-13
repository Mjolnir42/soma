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

// jobResultMgmtAdd function
// soma job result-mgmt add ${result}
func jobResultMgmtAdd(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	req := proto.NewJobResultRequest()
	req.JobResult.Name = c.Args().First()
	if err := adm.ValidateNotUUID(req.JobResult.Name); err != nil {
		return err
	}

	return adm.Perform(`postbody`, `/job/result-mgmt/`, `command`, req, c)
}

// jobResultMgmtRemove function
// soma job result-mgmt remove ${result}
func jobResultMgmtRemove(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	jobResultID, err := adm.LookupJobResultID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/job/result-mgmt/%s", url.QueryEscape(jobResultID))
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// jobResultMgmtShow function
// soma job result-mgmt show ${result}
func jobResultMgmtShow(c *cli.Context) error {
	if err := adm.VerifySingleArgument(c); err != nil {
		return err
	}

	jobResultID, err := adm.LookupJobResultID(c.Args().First())
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/job/result-mgmt/%s", url.QueryEscape(jobResultID))
	return adm.Perform(`get`, path, `command`, nil, c)
}

// jobResultMgmtList function
// soma job result-mgmt list
func jobResultMgmtList(c *cli.Context) error {
	if err := adm.VerifyNoArgument(c); err != nil {
		return err
	}

	return adm.Perform(`get`, `/job/result-mgmt/`, `list`, nil, c)
}

// jobResultMgmtSearch function
// soma job result-mgmt search [id ${uuid}] [name ${result}]
func jobResultMgmtSearch(c *cli.Context) error {
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
	req := proto.NewJobResultFilter()
	if _, ok := opts[`id`]; ok {
		if err := adm.ValidateUUID(opts[`id`][0]); err != nil {
			return err
		}
		req.Filter.JobResult.ID = opts[`id`][0]
		valid = true
	}
	if _, ok := opts[`name`]; ok {
		if err := adm.ValidateNotUUID(opts[`name`][0]); err != nil {
			return err
		}
		req.Filter.JobResult.Name = opts[`name`][0]
		valid = true
	}
	if !valid {
		return fmt.Errorf(`Syntax error, must specify 'id', 'name' or both.`)
	}

	return adm.Perform(`postbody`, `/search/jobResult/`, `show`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
