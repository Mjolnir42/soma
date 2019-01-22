/*-
 * Copyright (c) 2016-2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
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

// groupConfigMemberAssignGroup function
func groupConfigMemberAssignGroup(c *cli.Context) error {
	return groupConfigMemberAssign(c, proto.EntityGroup)
}

// groupConfigMemberAssign function
func groupConfigMemberAssign(c *cli.Context, childEntity string) error {
	switch childEntity {
	case proto.EntityGroup, proto.EntityCluster, proto.EntityNode:
	default:
		return fmt.Errorf("Unknown child entity type in group membership assignment: %s",
			childEntity)
	}

	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`, `in`}
	mandatoryOptions := []string{`to`}

	switch childEntity {
	case proto.EntityGroup, proto.EntityCluster:
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

	var err error
	var bucketID, repositoryID, groupID, childID string

	switch childEntity {
	case proto.EntityGroup, proto.EntityCluster:
		// fetch bucketID via in argument
		if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
			return err
		}
		if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
			return err
		}
		if groupID, err = adm.LookupGroupID(opts[`to`][0], bucketID); err != nil {
			return err
		}
		switch childEntity {
		case proto.EntityGroup:
			if childID, err = adm.LookupGroupID(c.Args().First(), bucketID); err != nil {
				return err
			}
		case proto.EntityCluster:
			if childID, err = adm.LookupClusterID(c.Args().First(), bucketID); err != nil {
				return err
			}
		}
	case proto.EntityNode:
		// fetch bucketID via nodeConfig
	}

	req := proto.NewGroupRequest()
	req.Group.ID = groupID
	req.Group.BucketID = bucketID
	req.Group.RepositoryID = repositoryID
	switch childEntity {
	case proto.EntityGroup:
		*req.Group.MemberGroups = append(*req.Group.MemberGroups, proto.Group{
			ID: childID,
		})
	case proto.EntityCluster:
		*req.Group.MemberClusters = append(*req.Group.MemberClusters, proto.Cluster{
			ID: childID,
		})
	case proto.EntityNode:
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/group/%s/member/",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(groupID),
	)
	return adm.Perform(`postbody`, path, `group-config::member-assign`, req, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
