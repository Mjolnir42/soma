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

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

// propertyMgmtCustomAdd function
// soma property-mgmt custom add ${property} to ${repository}
func propertyMgmtCustomAdd(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`}
	mandatoryOptions := []string{`to`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	repositoryID, err := adm.LookupRepoID(opts[`to`][0])
	if err != nil {
		return err
	}

	req := proto.NewCustomPropertyRequest()
	req.Property.Custom.Name = c.Args().First()
	req.Property.Custom.RepositoryID = repositoryID

	path := fmt.Sprintf("/repository/%s/property-mgmt/%s/",
		url.QueryEscape(repositoryID),
		url.QueryEscape(proto.PropertyTypeCustom),
	)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// propertyMgmtCustomRemove function
// soma property-mgmt custom remove ${property} from ${repository}
func propertyMgmtCustomRemove(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`from`}
	mandatoryOptions := []string{`from`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	repositoryID, err := adm.LookupRepoID(opts[`from`][0])
	if err != nil {
		return err
	}

	propertyID, err := adm.LookupCustomPropertyID(
		c.Args().First(),
		repositoryID,
	)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/property-mgmt/%s/%s",
		url.QueryEscape(repositoryID),
		url.QueryEscape(proto.PropertyTypeCustom),
		url.QueryEscape(propertyID),
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// propertyMgmtCustomShow function
// soma property-mgmt custom show ${property} in ${repository}
func propertyMgmtCustomShow(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{`in`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	repositoryID, err := adm.LookupRepoID(opts[`in`][0])
	if err != nil {
		return err
	}

	propertyID, err := adm.LookupCustomPropertyID(
		c.Args().First(),
		repositoryID,
	)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/property-mgmt/%s/%s",
		url.QueryEscape(repositoryID),
		url.QueryEscape(proto.PropertyTypeCustom),
		url.QueryEscape(propertyID),
	)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// propertyMgmtCustomList function
// soma property-mgmt custom list in ${repository}
func propertyMgmtCustomList(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`in`}
	mandatoryOptions := []string{`in`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		adm.AllArguments(c),
	); err != nil {
		return err
	}

	repositoryID, err := adm.LookupRepoID(opts[`in`][0])
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/property-mgmt/%s/",
		url.QueryEscape(repositoryID),
		url.QueryEscape(proto.PropertyTypeCustom),
	)
	return adm.Perform(`get`, path, `list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
