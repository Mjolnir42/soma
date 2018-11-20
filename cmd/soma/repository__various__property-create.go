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
	"os"

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

// cmdPropertyAdd function
func cmdPropertyAdd(c *cli.Context, pType, oType string) error {
	switch oType {
	case `node`, `bucket`, `repository`, `group`, `cluster`:
		switch pType {
		case `system`, `custom`, `service`, `oncall`:
		default:
			return fmt.Errorf("Unknown property type: %s", pType)
		}
	default:
		return fmt.Errorf("Unknown object type: %s", oType)
	}

	// argument parsing
	multiple := []string{}
	required := []string{`to`, `view`}
	unique := []string{`to`, `in`, `view`, `inheritance`, `childrenonly`}

	switch pType {
	case `system`:
		if err := adm.ValidateSystemProperty(
			c.Args().First()); err != nil {
			return err
		}
		fallthrough
	case `custom`:
		required = append(required, `value`)
		unique = append(unique, `value`)
	}
	switch oType {
	case `group`, `cluster`:
		required = append(required, `in`)
	}
	opts := map[string][]string{}
	if err := adm.ParseVariadicArguments(
		opts,
		multiple,
		unique,
		required,
		c.Args().Tail(),
	); err != nil {
		return err
	}

	// deprecation warning
	switch oType {
	case `repository`, `bucket`, `node`:
		if _, ok := opts[`in`]; ok {
			fmt.Fprintf(
				os.Stderr,
				"Hint: Keyword `in` is DEPRECATED for %s objects,"+
					" since they are global objects. Ignoring.",
				oType,
			)
		}
	}

	var (
		objectID, repoID, bucketID string
		config                     *proto.NodeConfig
		req                        proto.Request
		err                        error
	)
	// id lookup
	switch oType {
	case `node`:
		if objectID, err = adm.LookupNodeID(opts[`to`][0]); err != nil {
			return err
		}
		if config, err = adm.LookupNodeConfig(objectID); err != nil {
			return err
		}
		repoID = config.RepositoryID
		bucketID = config.BucketID
	case `cluster`:
		bucketID, err = adm.LookupBucketID(opts["in"][0])
		if err != nil {
			return err
		}
		if objectID, err = adm.LookupClusterID(opts[`to`][0],
			bucketID); err != nil {
			return err
		}
		if repoID, err = adm.LookupRepoByBucket(bucketID); err != nil {
			return err
		}
	case `group`:
		bucketID, err = adm.LookupBucketID(opts["in"][0])
		if err != nil {
			return err
		}
		if objectID, err = adm.LookupGroupID(opts[`to`][0],
			bucketID); err != nil {
			return err
		}
		if repoID, err = adm.LookupRepoByBucket(bucketID); err != nil {
			return err
		}
	case `bucket`:
		bucketID, err = adm.LookupBucketID(opts["to"][0])
		if err != nil {
			return err
		}
		objectID = bucketID
		if repoID, err = adm.LookupRepoByBucket(bucketID); err != nil {
			return err
		}
	case `repository`:
		repoID, err = adm.LookupRepoID(opts[`to`][0])
		if err != nil {
			return err
		}
		objectID = repoID
	}

	// property assembly
	prop := proto.Property{
		Type: pType,
		View: opts[`view`][0],
	}
	// property assembly, optional arguments
	if _, ok := opts[`childrenonly`]; ok {
		if err = adm.ValidateBool(opts[`childrenonly`][0],
			&prop.ChildrenOnly); err != nil {
			return err
		}
	} else {
		prop.ChildrenOnly = false
	}
	if _, ok := opts[`inheritance`]; ok {
		if err = adm.ValidateBool(opts[`inheritance`][0],
			&prop.Inheritance); err != nil {
			return err
		}
	} else {
		prop.Inheritance = true
	}
	switch pType {
	case `system`:
		prop.System = &proto.PropertySystem{
			Name:  c.Args().First(),
			Value: opts[`value`][0],
		}
	case `service`:
		var teamID string
		switch oType {
		case `repository`:
			if err = adm.LookupTeamByRepo(repoID, &teamID); err != nil {
				return err
			}
		default:
			if teamID, err = adm.LookupTeamByBucket(
				bucketID); err != nil {
				return err
			}
		}
		// no reason to fill out the attributes, client-provided
		// attributes are discarded by the server
		prop.Service = &proto.PropertyService{
			Name:       c.Args().First(),
			TeamID:     teamID,
			Attributes: []proto.ServiceAttribute{},
		}
	case `oncall`:
		oncallID, err := adm.LookupOncallID(c.Args().First())
		if err != nil {
			return err
		}
		prop.Oncall = &proto.PropertyOncall{
			ID: oncallID,
		}
		prop.Oncall.Name, prop.Oncall.Number, err = adm.LookupOncallDetails(
			oncallID,
		)
		if err != nil {
			return err
		}
	case `custom`:
		customID, err := adm.LookupCustomPropertyID(
			c.Args().First(), repoID)
		if err != nil {
			return err
		}

		prop.Custom = &proto.PropertyCustom{
			ID:           customID,
			Name:         c.Args().First(),
			RepositoryID: repoID,
			Value:        opts[`value`][0],
		}
	}

	// request assembly
	switch oType {
	case `node`:
		req = proto.NewNodeRequest()
		req.Node.ID = objectID
		req.Node.Config = config
		req.Node.Properties = &[]proto.Property{prop}
	case `cluster`:
		req = proto.NewClusterRequest()
		req.Cluster.ID = objectID
		req.Cluster.RepositoryID = repoID
		req.Cluster.BucketID = bucketID
		req.Cluster.Properties = &[]proto.Property{prop}
	case `group`:
		req = proto.NewGroupRequest()
		req.Group.ID = objectID
		req.Group.RepositoryID = repoID
		req.Group.BucketID = bucketID
		req.Group.Properties = &[]proto.Property{prop}
	case `bucket`:
		req = proto.NewBucketRequest()
		req.Bucket.ID = objectID
		req.Bucket.Properties = &[]proto.Property{prop}
	case `repository`:
		req = proto.NewRepositoryRequest()
		req.Repository.ID = repoID
		req.Repository.Properties = &[]proto.Property{prop}
	}

	var path string
	switch oType {
	case `cluster`, `group`, `node`:
		path = fmt.Sprintf("/repository/%s/bucket/%s/%s/%s/property/",
			repoID, bucketID, oType, objectID)
	case `bucket`:
		path = fmt.Sprintf("/repository/%s/%s/%s/property/",
			repoID, oType, objectID)
	case `repository`:
		path = fmt.Sprintf("/%s/%s/property/",
			oType, objectID)
	}
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
