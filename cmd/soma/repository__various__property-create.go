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

	"github.com/codegangsta/cli"
	"github.com/mjolnir42/soma/internal/adm"
	"github.com/mjolnir42/soma/lib/proto"
)

// variousPropertyCreate is the generic function for creating properties
// on tree objects
func variousPropertyCreate(c *cli.Context, propertyType, entity string) error {
	switch entity {
	case proto.EntityRepository, proto.EntityBucket, proto.EntityGroup, proto.EntityCluster, proto.EntityNode:
	default:
		return fmt.Errorf("Unknown entity: %s", entity)
	}
	switch propertyType {
	case proto.PropertyTypeSystem, proto.PropertyTypeCustom, proto.PropertyTypeService, proto.PropertyTypeOncall:
	case proto.PropertyTypeNative:
		return fmt.Errorf(`Native properties are for introspection and can not be created on tree objects`)
	case proto.PropertyTypeTemplate:
		return fmt.Errorf(`Template properties can not be created on tree objects`)
	default:
		return fmt.Errorf("Unknown property type: %s", propertyType)
	}

	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`on`, `view`, `inheritance`, `childrenonly`}
	mandatoryOptions := []string{`on`, `view`}

	switch propertyType {
	case proto.PropertyTypeSystem:
		uniqueOptions = append(uniqueOptions, `value`)
		mandatoryOptions = append(mandatoryOptions, `value`)
	case proto.PropertyTypeCustom:
		uniqueOptions = append(uniqueOptions, `value`)
		mandatoryOptions = append(mandatoryOptions, `value`)
	}

	switch entity {
	case proto.EntityGroup, proto.EntityCluster:
		uniqueOptions = append(uniqueOptions, `in`)
		mandatoryOptions = append(mandatoryOptions, `in`)
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

	switch propertyType {
	case proto.PropertyTypeSystem:
		if err := adm.ValidateSystemProperty(c.Args().First()); err != nil {
			return err
		}
	}

	var (
		objectID, repoID, bucketID string
		config                     *proto.NodeConfig
		req                        proto.Request
		err                        error
	)
	// id lookup
	switch entity {
	case `node`:
		if objectID, err = adm.LookupNodeID(opts[`on`][0]); err != nil {
			return err
		}
		if config, err = adm.LookupNodeConfig(objectID); err != nil {
			return err
		}
		repoID = config.RepositoryID
		bucketID = config.BucketID
	case `cluster`:
		bucketID, err = adm.LookupBucketID(opts[`in`][0])
		if err != nil {
			return err
		}
		if objectID, err = adm.LookupClusterID(opts[`on`][0],
			bucketID); err != nil {
			return err
		}
		if repoID, err = adm.LookupRepoByBucket(bucketID); err != nil {
			return err
		}
	case `group`:
		bucketID, err = adm.LookupBucketID(opts[`in`][0])
		if err != nil {
			return err
		}
		if objectID, err = adm.LookupGroupID(opts[`on`][0],
			bucketID); err != nil {
			return err
		}
		if repoID, err = adm.LookupRepoByBucket(bucketID); err != nil {
			return err
		}
	case `bucket`:
		bucketID, err = adm.LookupBucketID(opts[`on`][0])
		if err != nil {
			return err
		}
		objectID = bucketID
		if repoID, err = adm.LookupRepoByBucket(bucketID); err != nil {
			return err
		}
	case `repository`:
		repoID, err = adm.LookupRepoID(opts[`on`][0])
		if err != nil {
			return err
		}
		objectID = repoID
	}

	// property assembly
	prop := proto.Property{
		Type: propertyType,
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
	switch propertyType {
	case proto.PropertyTypeSystem:
		prop.System = &proto.PropertySystem{
			Name:  c.Args().First(),
			Value: opts[`value`][0],
		}
	case proto.PropertyTypeService:
		var serviceID, teamID string
		switch entity {
		case proto.EntityRepository:
			if err = adm.LookupTeamByRepo(repoID, &teamID); err != nil {
				return err
			}
		default:
			if teamID, err = adm.LookupTeamByBucket(bucketID); err != nil {
				return err
			}
		}
		serviceID, err = adm.LookupServicePropertyID(
			c.Args().First(),
			teamID)
		if err != nil {
			return err
		}

		// no reason to fill out the attributes, client-provided
		// attributes are discarded by the server
		prop.Service = &proto.PropertyService{
			ID:         serviceID,
			Name:       c.Args().First(),
			TeamID:     teamID,
			Attributes: []proto.ServiceAttribute{},
		}
	case proto.PropertyTypeOncall:
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
	case proto.PropertyTypeCustom:
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
	switch entity {
	case proto.EntityNode:
		req = proto.NewNodeRequest()
		req.Node.ID = objectID
		req.Node.Config = config
		req.Node.Properties = &[]proto.Property{prop}
	case proto.EntityCluster:
		req = proto.NewClusterRequest()
		req.Cluster.ID = objectID
		req.Cluster.RepositoryID = repoID
		req.Cluster.BucketID = bucketID
		req.Cluster.Properties = &[]proto.Property{prop}
	case proto.EntityGroup:
		req = proto.NewGroupRequest()
		req.Group.ID = objectID
		req.Group.RepositoryID = repoID
		req.Group.BucketID = bucketID
		req.Group.Properties = &[]proto.Property{prop}
	case proto.EntityBucket:
		req = proto.NewBucketRequest()
		req.Bucket.ID = objectID
		req.Bucket.Properties = &[]proto.Property{prop}
	case proto.EntityRepository:
		req = proto.NewRepositoryRequest()
		req.Repository.ID = repoID
		req.Repository.Properties = &[]proto.Property{prop}
	}

	var path string
	switch entity {
	case `cluster`, `group`, `node`:
		path = fmt.Sprintf("/repository/%s/bucket/%s/%s/%s/property/",
			repoID, bucketID, entity, objectID)
	case `bucket`:
		path = fmt.Sprintf("/repository/%s/%s/%s/property/",
			repoID, entity, objectID)
	case proto.EntityRepository:
		path = fmt.Sprintf("/%s/%s/property/",
			proto.EntityRepository,
			objectID,
		)
	}
	return adm.Perform(`postbody`, path, `command`, req, c)
}

// variousPropertyDestroy is the generic function for destroying
// properties on tree objects
func variousPropertyDestroy(c *cli.Context, propertyType, entity string) error {
	switch entity {
	case proto.EntityRepository, proto.EntityBucket, proto.EntityGroup, proto.EntityCluster, proto.EntityNode:
	default:
		return fmt.Errorf("Unknown entity: %s", entity)
	}
	switch propertyType {
	case proto.PropertyTypeSystem, proto.PropertyTypeCustom, proto.PropertyTypeService, proto.PropertyTypeOncall:
	case proto.PropertyTypeNative:
		return fmt.Errorf(`Native properties are for introspection and can not be created on tree objects`)
	case proto.PropertyTypeTemplate:
		return fmt.Errorf(`Template properties can not be created on tree objects`)
	default:
		return fmt.Errorf("Unknown property type: %s", propertyType)
	}

	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`on`, `view`}
	mandatoryOptions := []string{`on`, `view`}

	switch entity {
	case proto.EntityGroup, proto.EntityCluster:
		uniqueOptions = append(uniqueOptions, `in`)
		mandatoryOptions = append(mandatoryOptions, `in`)
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

	var (
		repositoryID, bucketID, teamID, serviceID string
		property, sourceID, path, command         string
		err                                       error
	)

	// id lookup
	switch entity {
	case proto.EntityRepository:
		if repositoryID, err = adm.LookupRepoID(opts[`on`][0]); err != nil {
			return err
		}
	case proto.EntityBucket:
		if bucketID, err = adm.LookupBucketID(opts[`on`][0]); err != nil {
			return err
		}
		if teamID, err = adm.LookupTeamByBucket(bucketID); err != nil {
			return err
		}
		if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
			return err
		}
	case proto.EntityGroup, proto.EntityCluster, proto.EntityNode:
	default:
		return fmt.Errorf("Unknown entity: %s", entity)
	}

	switch propertyType {
	case proto.PropertyTypeService:
		if serviceID, err = adm.LookupServicePropertyID(c.Args().First(), teamID); err != nil {
			return err
		}
		property = serviceID
	case proto.PropertyTypeSystem:
		if err := adm.ValidateSystemProperty(
			c.Args().First()); err != nil {
			return err
		}
		property = c.Args().First()
	default:
		property = c.Args().First()
	}

	switch entity {
	case proto.EntityRepository:
		if err = adm.FindRepoPropSrcID(propertyType, property,
			opts[`view`][0], repositoryID, &sourceID); err != nil {
			return err
		}
		path = fmt.Sprintf("/repository/%s/property/%s/%s",
			repositoryID, propertyType, sourceID,
		)
		command = `repository-config::property-destroy`
	case proto.EntityBucket:
		if err = adm.FindBucketPropSrcID(propertyType, property,
			opts[`view`][0], bucketID, &sourceID); err != nil {
			return err
		}
		path = fmt.Sprintf("/repository/%s/bucket/%s/property/%s/%s",
			repositoryID, bucketID, propertyType, sourceID,
		)
	case proto.EntityGroup, proto.EntityCluster, proto.EntityNode:
	default:
		return fmt.Errorf("Unknown entity: %s", entity)
	}

	return adm.Perform(`delete`, path, command, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
