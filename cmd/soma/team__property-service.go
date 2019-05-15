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

// propertyMgmtServiceAdd function
// soma property service add ${property} team ${team} [${attribute} ${attrValue}, ...]
func propertyMgmtServiceAdd(c *cli.Context) error {
	var (
		teamID string
		err    error
	)
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`team`}
	mandatoryOptions := []string{`team`}

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
	if err = popError(); err != nil {
		return err
	}

	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if err = adm.LookupTeamID(opts[`team`][0], &teamID); err != nil {
		return err
	}

	// construct request body
	req := proto.NewServicePropertyRequest()
	req.Property.Service.Name = c.Args().First()
	req.Property.Service.TeamID = teamID

	if err = adm.ValidateRuneCount(
		req.Property.Service.Name, 128); err != nil {
		return err
	}

	// fill attributes into request body
attrConversionLoop:
	for oName := range opts {
		if oName == `team` {
			continue attrConversionLoop
		}
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

	path := fmt.Sprintf("/team/%s/property-mgmt/%s/",
		url.QueryEscape(teamID),
		url.QueryEscape(proto.PropertyTypeService),
	)
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// propertyMgmtServiceAdd function
// soma property service add ${property} team ${team} [${attribute} ${attrValue}, ...]
func propertyMgmtServiceUpdate(c *cli.Context) error {
	var (
		teamID string
		err    error
	)
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`team`}
	mandatoryOptions := []string{`team`}

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
	if err = popError(); err != nil {
		return err
	}

	if err = adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	if err = adm.LookupTeamID(opts[`team`][0], &teamID); err != nil {
		return err
	}

	// construct request body
	req := proto.NewServicePropertyRequest()
	req.Property.Service.Name = c.Args().First()
	req.Property.Service.TeamID = teamID

	if err = adm.ValidateRuneCount(
		req.Property.Service.Name, 128); err != nil {
		return err
	}
	serviceID, err := adm.LookupServicePropertyID(req.Property.Service.Name, teamID)
	if err != nil {
		return err
	}
	req.Property.Service.ID = serviceID
	// fill attributes into request body
attrConversionLoop:
	for oName := range opts {
		if oName == `team` {
			continue attrConversionLoop
		}
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
	path := fmt.Sprintf("/team/%s/property-mgmt/%s/%s",
		url.QueryEscape(teamID),
		url.QueryEscape(proto.PropertyTypeService),
		url.QueryEscape(serviceID),
	)
	return adm.Perform(`putbody`, path, `command`, req, c)
}

// propertyMgmtServiceRemove function
// soma property service remove ${property} team ${team}
func propertyMgmtServiceRemove(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`team`}
	mandatoryOptions := []string{`team`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var teamID string
	if err := adm.LookupTeamID(opts[`team`][0], &teamID); err != nil {
		return err
	}
	propertyID, err := adm.LookupServicePropertyID(
		c.Args().First(), teamID)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/team/%s/property-mgmt/%s/%s",
		url.QueryEscape(teamID),
		url.QueryEscape(proto.PropertyTypeService),
		url.QueryEscape(propertyID),
	)
	return adm.Perform(`delete`, path, `command`, nil, c)
}

// propertyMgmtServiceShow function
// soma property service show ${property} team ${team}
func propertyMgmtServiceShow(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`team`}
	mandatoryOptions := []string{`team`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	var teamID string
	if err := adm.LookupTeamID(opts[`team`][0], &teamID); err != nil {
		return err
	}
	propertyID, err := adm.LookupServicePropertyID(
		c.Args().First(), teamID)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/team/%s/property-mgmt/%s/%s",
		url.QueryEscape(teamID),
		url.QueryEscape(proto.PropertyTypeService),
		url.QueryEscape(propertyID),
	)
	return adm.Perform(`get`, path, `show`, nil, c)
}

// propertyMgmtServiceList function
// soma property service list team ${team}
func propertyMgmtServiceList(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`team`}
	mandatoryOptions := []string{`team`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		adm.AllArguments(c),
	); err != nil {
		return err
	}

	var teamID string
	if err := adm.LookupTeamID(opts[`team`][0], &teamID); err != nil {
		return err
	}

	path := fmt.Sprintf("/team/%s/property-mgmt/%s/",
		url.QueryEscape(teamID),
		url.QueryEscape(proto.PropertyTypeService),
	)
	return adm.Perform(`get`, path, `list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
