/*-
 * Copyright (c) 2015-2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
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
		clusterID                                 string
		err                                       error
	)

	// id lookup
	switch entity {
	case proto.EntityRepository:
		if repositoryID, err = adm.LookupRepoID(opts[`on`][0]); err != nil {
			return err
		}
		if err = adm.LookupTeamByRepo(repositoryID, &teamID); err != nil {
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
	case proto.EntityGroup:
	case proto.EntityCluster:
		if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
			return err
		}
		if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
			return err
		}
		if clusterID, err = adm.LookupClusterID(opts[`on`][0], bucketID); err != nil {
			return err
		}
	case proto.EntityNode:
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
		command = `bucket-config::property-destroy`
	case proto.EntityGroup:
	case proto.EntityCluster:
		if err = adm.FindClusterPropSrcID(propertyType, property,
			opts[`view`][0], clusterID, &sourceID); err != nil {
			return err
		}
		path = fmt.Sprintf("/repository/%s/bucket/%s/cluster/%s/property/%s/%s",
			repositoryID, bucketID, clusterID, propertyType, sourceID,
		)
		command = `cluster-config::property-destroy`
	case proto.EntityNode:
	default:
		return fmt.Errorf("Unknown entity: %s", entity)
	}

	return adm.Perform(`delete`, path, command, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix