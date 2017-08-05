/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/mjolnir42/soma/internal/msg"
)

func (g *GuidePost) validateRequest(q *msg.Request) (error, bool) {
	switch q.Section {
	case `check`:
		if err, nf := g.validateCheckObjectInBucket(q); err != nil {
			return err, nf
		}
	case `node`:
		if err, nf := g.validateNodeConfig(q); err != nil {
			return err, nf
		}
		fallthrough
	case `cluster`, `group`:
		if err, nf := g.validateCorrectBucket(q); err != nil {
			return err, nf
		}
	case `bucket`:
		if err, nf := g.validateBucketInRepository(
			q.Bucket.RepositoryId,
			q.Bucket.Id,
		); err != nil {
			return err, nf
		}
	case `repository`:
		// since repository ids are the routing information,
		// it is unnecessary to check that the object is where the
		// routing would point to
	default:
		return fmt.Errorf("Invalid request type %s", q.Section), false
	}

	switch q.Action {
	case
		`add_node_to_cluster`,
		`add_node_to_group`,
		`add_cluster_to_group`,
		`add_group_to_group`:
		return g.validateObjectMatch(q)
	case
		`add_check_to_bucket`,
		`add_check_to_cluster`,
		`add_check_to_group`,
		`add_check_to_node`,
		`add_check_to_repository`:
		return g.validateCheckThresholds(q)
	case
		`create_bucket`:
		return g.validateBucketName(q)
	case
		`add_custom_property_to_bucket`,
		`add_custom_property_to_cluster`,
		`add_custom_property_to_group`,
		`add_custom_property_to_node`,
		`add_custom_property_to_repository`,
		`add_oncall_property_to_bucket`,
		`add_oncall_property_to_cluster`,
		`add_oncall_property_to_group`,
		`add_oncall_property_to_node`,
		`add_oncall_property_to_repository`,
		`add_service_property_to_bucket`,
		`add_service_property_to_cluster`,
		`add_service_property_to_group`,
		`add_service_property_to_node`,
		`add_service_property_to_repository`,
		`add_system_property_to_bucket`,
		`add_system_property_to_cluster`,
		`add_system_property_to_group`,
		`add_system_property_to_node`,
		`add_system_property_to_repository`,
		`assign_node`,
		`create_cluster`,
		`create_group`,
		`delete_custom_property_from_bucket`,
		`delete_custom_property_from_cluster`,
		`delete_custom_property_from_group`,
		`delete_custom_property_from_node`,
		`delete_custom_property_from_repository`,
		`delete_oncall_property_from_bucket`,
		`delete_oncall_property_from_cluster`,
		`delete_oncall_property_from_group`,
		`delete_oncall_property_from_node`,
		`delete_oncall_property_from_repository`,
		`delete_service_property_from_bucket`,
		`delete_service_property_from_cluster`,
		`delete_service_property_from_group`,
		`delete_service_property_from_node`,
		`delete_service_property_from_repository`,
		`delete_system_property_from_bucket`,
		`delete_system_property_from_cluster`,
		`delete_system_property_from_group`,
		`delete_system_property_from_node`,
		`delete_system_property_from_repository`,
		`remove_check`:
		// actions are accepted, but require no further validation
		return nil, false
	default:
		return fmt.Errorf("Unimplemented GuidePost/%s", q.Action), false
	}
}

func (g *GuidePost) validateObjectMatch(q *msg.Request) (error, bool) {
	var (
		nodeId, clusterId, groupId, childGroupId              string
		valNodeBId, valClusterBId, valGroupBId, valChGroupBId string
	)

	switch q.Action {
	case `add_node_to_cluster`:
		nodeId = (*q.Cluster.Members)[0].Id
		clusterId = q.Cluster.Id
	case `add_node_to_group`:
		nodeId = (*q.Group.MemberNodes)[0].Id
		groupId = q.Group.Id
	case `add_cluster_to_group`:
		clusterId = (*q.Group.MemberClusters)[0].Id
		groupId = q.Group.Id
	case `add_group_to_group`:
		childGroupId = (*q.Group.MemberGroups)[0].Id
		groupId = q.Group.Id
	default:
		return fmt.Errorf("Incorrect validation attempted for %s",
			q.Action), false
	}

	if nodeId != `` {
		if err := g.stmtBucketForNodeID.QueryRow(nodeId).Scan(
			&valNodeBId,
		); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("Unknown node %s", nodeId), true
			}
			return err, false
		}
	}
	if clusterId != `` {
		if err := g.stmtBucketForClusterID.QueryRow(clusterId).Scan(
			&valClusterBId,
		); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("Unknown cluster %s", clusterId), true
			}
			return err, false
		}
	}
	if groupId != `` {
		if err := g.stmtBucketForGroupID.QueryRow(groupId).Scan(
			&valGroupBId,
		); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("Unknown group %s", groupId), true
			}
			return err, false
		}
	}
	if childGroupId != `` {
		if err := g.stmtBucketForGroupID.QueryRow(childGroupId).Scan(
			&valChGroupBId,
		); err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("Unknown group %s", childGroupId), true
			}
			return err, false
		}
	}

	switch q.Action {
	case `add_node_to_cluster`:
		if valNodeBId != valClusterBId {
			return fmt.Errorf(
				"Node and Cluster are in different buckets (%s/%s)",
				valNodeBId, valClusterBId,
			), false
		}
	case `add_node_to_group`:
		if valNodeBId != valGroupBId {
			return fmt.Errorf(
				"Node and Group are in different buckets (%s/%s)",
				valNodeBId, valGroupBId,
			), false
		}
	case `add_cluster_to_group`:
		if valClusterBId != valGroupBId {
			return fmt.Errorf(
				"Cluster and Group are in different buckets (%s/%s)",
				valClusterBId, valGroupBId,
			), false
		}
	case `add_group_to_group`:
		if valChGroupBId != valGroupBId {
			return fmt.Errorf(
				"Groups are in different buckets (%s/%s)",
				valGroupBId, valChGroupBId,
			), false
		}
	}
	return nil, false
}

// Verify that an object is assigned to the specified bucket.
func (g *GuidePost) validateCorrectBucket(q *msg.Request) (error, bool) {
	switch q.Action {
	case `assign_node`:
		return g.validateNodeUnassigned(q)
	case `create_cluster`, `create_group`:
		return nil, false
	}
	var bid string
	var err error
	switch q.Section {
	case `node`:
		err = g.stmtBucketForNodeID.QueryRow(
			q.Node.Id,
		).Scan(
			&bid,
		)
	case `cluster`:
		err = g.stmtBucketForClusterID.QueryRow(
			q.Cluster.Id,
		).Scan(
			&bid,
		)
	case `group`:
		err = g.stmtBucketForGroupID.QueryRow(
			q.Group.Id,
		).Scan(
			&bid,
		)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			// unassigned
			return fmt.Errorf("%s is not assigned to any bucket",
				q.Section), true
		}
		return err, false
	}
	switch q.Section {
	case `node`:
		if bid != q.Node.Config.BucketId {
			return fmt.Errorf("Node assigned to different bucket %s",
				bid), false
		}
	case `cluster`:
		if bid != q.Cluster.BucketId {
			return fmt.Errorf("Cluster in different bucket %s",
				bid), false
		}
	case `group`:
		if bid != q.Group.BucketId {
			return fmt.Errorf("Group in different bucket %s",
				bid), false
		}
	}
	return nil, false
}

// Verify that a node is not yet assigned to a bucket. Returns nil
// on success.
func (g *GuidePost) validateNodeUnassigned(q *msg.Request) (error, bool) {
	var bid string
	if err := g.stmtBucketForNodeID.QueryRow(q.Node.Id).Scan(
		&bid,
	); err != nil {
		if err == sql.ErrNoRows {
			// unassigned is not an error here
			return nil, false
		}
		return err, false
	}
	return fmt.Errorf("Node already assigned to bucket %s", bid), false
}

// Verify the node has a Config section
func (g *GuidePost) validateNodeConfig(q *msg.Request) (error, bool) {
	if q.Node.Config == nil {
		return fmt.Errorf("NodeConfig subobject missing"), false
	}
	return g.validateBucketInRepository(
		q.Node.Config.RepositoryId,
		q.Node.Config.BucketId,
	)
}

// Verify that the ObjectId->BucketId->RepositoryId chain is part of
// the same tree.
func (g *GuidePost) validateCheckObjectInBucket(q *msg.Request) (error, bool) {
	var err error
	var bid string
	switch q.CheckConfig.ObjectType {
	case `repository`:
		if q.CheckConfig.RepositoryId !=
			q.CheckConfig.ObjectId {
			return fmt.Errorf("Conflicting repository ids: %s, %s",
				q.CheckConfig.RepositoryId,
				q.CheckConfig.ObjectId,
			), false
		}
		return nil, false
	case `bucket`:
		bid = q.CheckConfig.ObjectId
	case `group`:
		err = g.stmtBucketForGroupID.QueryRow(
			q.CheckConfig.ObjectId,
		).Scan(&bid)
	case `cluster`:
		err = g.stmtBucketForClusterID.QueryRow(
			q.CheckConfig.ObjectId,
		).Scan(&bid)
	case `node`:
		err = g.stmtBucketForNodeID.QueryRow(
			q.CheckConfig.ObjectId,
		).Scan(&bid)
	default:
		return fmt.Errorf("Unknown object type: %s",
			q.CheckConfig.ObjectType,
		), false
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("No bucketId for object found"), true
		}
		return err, false
	}
	if bid != q.CheckConfig.BucketId {
		return fmt.Errorf("Object is in bucket %s, not %s",
			bid, q.CheckConfig.BucketId,
		), false
	}
	return g.validateBucketInRepository(
		q.CheckConfig.RepositoryId,
		q.CheckConfig.BucketId,
	)
}

// Verify that the bucket is part of the specified repository
func (g *GuidePost) validateBucketInRepository(
	repo, bucket string) (error, bool) {
	var repoID, repoName string
	if err := g.stmtRepoForBucketID.QueryRow(
		bucket,
	).Scan(
		&repoID,
		&repoName,
	); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("No repository found for bucket %s",
				bucket), true
		}
		return err, false
	}
	if repo != repoID {
		return fmt.Errorf("Bucket is in different repository: %s",
			repoID), false
	}
	return nil, false
}

// check the check configuration to contain fewer thresholds than
// the limit for the capability
func (g *GuidePost) validateCheckThresholds(q *msg.Request) (error, bool) {
	var (
		thrLimit int
		err      error
	)

	if err = g.stmtCapabilityThresholds.QueryRow(
		q.CheckConfig.CapabilityId,
	).Scan(
		&thrLimit,
	); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf(
				"Capability %s not found",
				q.CheckConfig.CapabilityId), true
		}
		return err, false
	}
	if len(q.CheckConfig.Thresholds) > thrLimit {
		return fmt.Errorf(
			"Specified %d thresholds exceed limit of %d for capability",
			len(q.CheckConfig.Thresholds),
			thrLimit), false
	}
	return nil, false
}

// check the naming schema for the bucket (global unique object)
func (g *GuidePost) validateBucketName(q *msg.Request) (error, bool) {
	_, repoName, _, _ := g.extractRouting(q)

	if !strings.HasPrefix(
		q.Bucket.Name,
		fmt.Sprintf("%s_", repoName),
	) {
		return fmt.Errorf("Illegal bucket name format, " +
			"requires reponame_ prefix"), false
	}
	return nil, false
}

// validate current treekeeper state
func (g *GuidePost) validateKeeper(repoName string) (error, bool) {
	// check we have a treekeeper for that repository
	keeper := fmt.Sprintf("repository_%s", repoName)
	if _, ok := g.soma.handlerMap.Get(keeper).(*TreeKeeper); !ok {
		return fmt.Errorf(
			"No handler for repository %s currently registered.",
			repoName), true
	}
	handler := g.soma.handlerMap.Get(keeper).(*TreeKeeper)

	// check the treekeeper has not been stopped
	if handler.isStopped() {
		return fmt.Errorf(
			"Repository %s is currently stopped.", repoName), false
	}

	// check the treekeeper has not encountered a broken tree
	if handler.isBroken() {
		return fmt.Errorf(
			"Repository %s is broken.", repoName), false
	}

	// check the treekeeper has finished loading
	if !handler.isReady() {
		return fmt.Errorf(
			"Repository %s not fully loaded yet.", repoName), false
	}
	return nil, false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
