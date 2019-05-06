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

// variousPropertyUpdate is the generic function for updating
// properties on tree objects
func variousPropertyUpdate(c *cli.Context, propertyType, entity string) error {
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
		repositoryID, bucketID     string
		property, sourceID, path   string
		groupID, clusterID, nodeID string
		objectID, repoID           string
		config                     *proto.NodeConfig
		req                        proto.Request
		err                        error
	)
	// id lookup
	switch entity {
	case proto.EntityNode:
		if objectID, err = adm.LookupNodeID(opts[`on`][0]); err != nil {
			return err
		}
		if config, err = adm.LookupNodeConfig(objectID); err != nil {
			return err
		}
		repoID = config.RepositoryID
		bucketID = config.BucketID
	case proto.EntityCluster:
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
	case proto.EntityGroup:
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
	case proto.EntityBucket:
		bucketID, err = adm.LookupBucketID(opts[`on`][0])
		if err != nil {
			return err
		}
		objectID = bucketID
		if repoID, err = adm.LookupRepoByBucket(bucketID); err != nil {
			return err
		}
	case proto.EntityRepository:
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
		req.Bucket.RepositoryID = repoID
		req.Bucket.Properties = &[]proto.Property{prop}
	case proto.EntityRepository:
		req = proto.NewRepositoryRequest()
		req.Repository.ID = repoID
		req.Repository.Properties = &[]proto.Property{prop}
	}

	switch entity {
	case proto.EntityRepository:
		if err = adm.FindRepoPropSrcID(propertyType, property,
			opts[`view`][0], repositoryID, &sourceID); err != nil {
			return err
		}
		path = fmt.Sprintf("/repository/%s/property/%s/%s",
			url.QueryEscape(repositoryID),
			url.QueryEscape(propertyType),
			url.QueryEscape(sourceID),
		)
	case proto.EntityBucket:
		if err = adm.FindBucketPropSrcID(propertyType, property,
			opts[`view`][0], repositoryID, bucketID, &sourceID); err != nil {
			return err
		}
		path = fmt.Sprintf("/repository/%s/bucket/%s/property/%s/%s",
			url.QueryEscape(repositoryID),
			url.QueryEscape(bucketID),
			url.QueryEscape(propertyType),
			url.QueryEscape(sourceID),
		)
	case proto.EntityGroup:
		if err = adm.FindGroupPropSrcID(propertyType, property,
			opts[`view`][0], repositoryID, bucketID, groupID, &sourceID); err != nil {
			return err
		}
		path = fmt.Sprintf("/repository/%s/bucket/%s/group/%s/property/%s/%s",
			url.QueryEscape(repositoryID),
			url.QueryEscape(bucketID),
			url.QueryEscape(groupID),
			url.QueryEscape(propertyType),
			url.QueryEscape(sourceID),
		)
	case proto.EntityCluster:
		if err = adm.FindClusterPropSrcID(propertyType, property,
			opts[`view`][0], repositoryID, bucketID, clusterID, &sourceID); err != nil {
			return err
		}
		path = fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s/property/%s/%s",
			url.QueryEscape(repositoryID),
			url.QueryEscape(bucketID),
			url.QueryEscape(clusterID),
			url.QueryEscape(propertyType),
			url.QueryEscape(sourceID),
		)
	case proto.EntityNode:
		if err = adm.FindNodePropSrcID(propertyType, property,
			opts[`view`][0], nodeID, &sourceID); err != nil {
			return err
		}
		path = fmt.Sprintf("/repository/%s/bucket/%s/node/%s/property/%s/%s",
			url.QueryEscape(repositoryID),
			url.QueryEscape(bucketID),
			url.QueryEscape(nodeID),
			url.QueryEscape(propertyType),
			url.QueryEscape(sourceID),
		)
	default:
		return fmt.Errorf("Unknown entity: %s", entity)
	}
	command := fmt.Sprintf("%s-config::property-update", entity)
	return adm.Perform(`putbody`, path, command, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
