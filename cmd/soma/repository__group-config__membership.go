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
// soma group member assign group ${child-group} to ${group} in ${bucket}
func groupConfigMemberAssignGroup(c *cli.Context) error {
	return groupConfigMemberAssign(c, proto.EntityGroup)
}

// groupConfigMemberAssignCluster function
// soma group member assign cluster ${cluster} to ${group} in ${bucket}
func groupConfigMemberAssignCluster(c *cli.Context) error {
	return groupConfigMemberAssign(c, proto.EntityCluster)
}

// groupConfigMemberAssignNode function
// soma group member assign node ${node} to ${group} [in ${bucket}]
func groupConfigMemberAssignNode(c *cli.Context) error {
	return groupConfigMemberAssign(c, proto.EntityNode)
}

// groupConfigMemberAssign function
func groupConfigMemberAssign(c *cli.Context, childEntity string) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`to`, `in`}
	mandatoryOptions := []string{`to`}

	switch childEntity {
	case proto.EntityGroup, proto.EntityCluster:
		mandatoryOptions = append(mandatoryOptions, `in`)
	case proto.EntityNode:
	default:
		return fmt.Errorf(
			"Unknown child entity type in group membership assignment: %s",
			childEntity)
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
	var bucketID, repositoryID, groupID, childID, controlID string

	switch childEntity {
	case proto.EntityGroup, proto.EntityCluster:
		// fetch bucketID via in argument
		if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
			return err
		}
		if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
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
		if childID, err = adm.LookupNodeID(c.Args().First()); err != nil {
			return err
		}
		nodeConfig := &proto.NodeConfig{}
		if nodeConfig, err = adm.LookupNodeConfig(childID); err != nil {
			return err
		}
		bucketID = nodeConfig.BucketID
		repositoryID = nodeConfig.RepositoryID

		// for proto.EntityNode optional argument must be correct if provided
		if _, ok := opts[`in`]; ok {
			if controlID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
				return err
			} else if controlID != bucketID {
				return fmt.Errorf(
					"Invalid bucket %s(%s), node %s is assigned to %s",
					opts[`in`][0],
					controlID,
					c.Args().First(),
					bucketID,
				)
			}
		}
	}
	if groupID, err = adm.LookupGroupID(opts[`to`][0], bucketID); err != nil {
		return err
	}

	req := proto.NewGroupRequest()
	req.Group.ID = groupID
	req.Group.BucketID = bucketID
	req.Group.RepositoryID = repositoryID
	switch childEntity {
	case proto.EntityGroup:
		req.Group.MemberGroups = &[]proto.Group{}
		*req.Group.MemberGroups = append(*req.Group.MemberGroups, proto.Group{
			ID: childID,
		})
	case proto.EntityCluster:
		req.Group.MemberClusters = &[]proto.Cluster{}
		*req.Group.MemberClusters = append(*req.Group.MemberClusters, proto.Cluster{
			ID: childID,
		})
	case proto.EntityNode:
		req.Group.MemberNodes = &[]proto.Node{}
		*req.Group.MemberNodes = append(*req.Group.MemberNodes, proto.Node{
			ID: childID,
		})
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/group/%s/member/%s/",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(groupID),
		url.QueryEscape(childEntity),
	)
	return adm.Perform(`postbody`, path, `group-config::member-assign`, req, c)
}

// groupConfigMemberUnassignGroup function
// soma group member unassign group ${child-group} from ${group} in ${bucket}
func groupConfigMemberUnassignGroup(c *cli.Context) error {
	return groupConfigMemberUnassign(c, proto.EntityGroup)
}

// groupConfigMemberUnassignCluster function
// soma group member unassign cluster ${cluster} from ${group} in ${bucket}
func groupConfigMemberUnassignCluster(c *cli.Context) error {
	return groupConfigMemberUnassign(c, proto.EntityCluster)
}

// groupConfigMemberUnassignNode function
// soma group member unassign node ${node} from ${group} [in ${bucket}]
func groupConfigMemberUnassignNode(c *cli.Context) error {
	return groupConfigMemberUnassign(c, proto.EntityNode)
}

// groupConfigMemberUnassign function
func groupConfigMemberUnassign(c *cli.Context, childEntity string) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`from`, `in`}
	mandatoryOptions := []string{`from`}

	switch childEntity {
	case proto.EntityGroup, proto.EntityCluster:
		mandatoryOptions = append(mandatoryOptions, `in`)
	case proto.EntityNode:
	default:
		return fmt.Errorf(
			"Unknown child entity type in group membership unassignment: %s",
			childEntity)
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
	var bucketID, repositoryID, groupID, childID, controlID string

	switch childEntity {
	case proto.EntityGroup, proto.EntityCluster:
		// fetch bucketID via mandatory `in` argument
		if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
			return err
		}
		if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
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
		if childID, err = adm.LookupNodeID(c.Args().First()); err != nil {
			return err
		}
		nodeConfig := &proto.NodeConfig{}
		if nodeConfig, err = adm.LookupNodeConfig(childID); err != nil {
			return err
		}
		bucketID = nodeConfig.BucketID
		repositoryID = nodeConfig.RepositoryID

		// for proto.EntityNode optional argument must be correct if provided
		if _, ok := opts[`in`]; ok {
			if controlID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
				return err
			} else if controlID != bucketID {
				return fmt.Errorf(
					"Invalid bucket %s(%s), node %s is assigned to %s",
					opts[`in`][0],
					controlID,
					c.Args().First(),
					bucketID,
				)
			}
		}
	}
	if groupID, err = adm.LookupGroupID(opts[`to`][0], bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/group/%s/member/%s/%s",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(groupID),
		url.QueryEscape(childEntity),
		url.QueryEscape(childID),
	)
	return adm.Perform(`delete`, path, `group-config::member-unassign`, nil, c)
}

// groupConfigMemberList function
// soma group member list of ${group} in ${bucket}
func groupConfigMemberList(c *cli.Context) error {
	opts := map[string][]string{}
	multipleAllowed := []string{}
	uniqueOptions := []string{`of`, `in`}
	mandatoryOptions := []string{`of`, `in`}

	if err := adm.ParseVariadicArguments(
		opts,
		multipleAllowed,
		uniqueOptions,
		mandatoryOptions,
		adm.AllArguments(c),
	); err != nil {
		return err
	}

	var err error
	var bucketID, repositoryID, groupID string
	if bucketID, err = adm.LookupBucketID(opts[`in`][0]); err != nil {
		return err
	}
	if repositoryID, err = adm.LookupRepoByBucket(bucketID); err != nil {
		return err
	}
	if groupID, err = adm.LookupGroupID(opts[`of`][0], bucketID); err != nil {
		return err
	}

	path := fmt.Sprintf("/repository/%s/bucket/%s/group/%s/member/",
		url.QueryEscape(repositoryID),
		url.QueryEscape(bucketID),
		url.QueryEscape(groupID),
	)
	return adm.Perform(`get`, path, `group-config::member-list`, nil, c)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
