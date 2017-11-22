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

func (g *GuidePost) validateRequest(q *msg.Request) (bool, error) {
	switch q.Section {
	case `check`:
		if nf, err := g.validateCheckObjectInBucket(q); err != nil {
			return nf, err
		}
	case `node`:
		if nf, err := g.validateNodeConfig(q); err != nil {
			return nf, err
		}
		fallthrough
	case `cluster`, `group`:
		if nf, err := g.validateCorrectBucket(q); err != nil {
			return nf, err
		}
	case `bucket`:
		if nf, err := g.validateBucketInRepository(
			q.Bucket.RepositoryID,
			q.Bucket.ID,
		); err != nil {
			return nf, err
		}
	case `repository`:
		// since repository ids are the routing information,
		// it is unnecessary to check that the object is where the
		// routing would point to
	default:
		return false, fmt.Errorf("Invalid request type %s", q.Section)
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
		return false, nil
	default:
		return false, fmt.Errorf("Unimplemented GuidePost/%s", q.Action)
	}
}

func (g *GuidePost) validateObjectMatch(q *msg.Request) (bool, error) {
	var (
		nodeID, clusterID, groupID, childGroupID              string
		valNodeBId, valClusterBId, valGroupBId, valChGroupBId string
	)

	switch q.Action {
	case `add_node_to_cluster`:
		nodeID = (*q.Cluster.Members)[0].Id
		clusterID = q.Cluster.Id
	case `add_node_to_group`:
		nodeID = (*q.Group.MemberNodes)[0].Id
		groupID = q.Group.Id
	case `add_cluster_to_group`:
		clusterID = (*q.Group.MemberClusters)[0].Id
		groupID = q.Group.Id
	case `add_group_to_group`:
		childGroupID = (*q.Group.MemberGroups)[0].Id
		groupID = q.Group.Id
	default:
		return false, fmt.Errorf("Incorrect validation attempted for %s",
			q.Action)
	}

	if nodeID != `` {
		if err := g.stmtBucketForNodeID.QueryRow(nodeID).Scan(
			&valNodeBId,
		); err != nil {
			if err == sql.ErrNoRows {
				return true, fmt.Errorf("Unknown node %s", nodeID)
			}
			return false, err
		}
	}
	if clusterID != `` {
		if err := g.stmtBucketForClusterID.QueryRow(clusterID).Scan(
			&valClusterBId,
		); err != nil {
			if err == sql.ErrNoRows {
				return true, fmt.Errorf("Unknown cluster %s", clusterID)
			}
			return false, err
		}
	}
	if groupID != `` {
		if err := g.stmtBucketForGroupID.QueryRow(groupID).Scan(
			&valGroupBId,
		); err != nil {
			if err == sql.ErrNoRows {
				return true, fmt.Errorf("Unknown group %s", groupID)
			}
			return false, err
		}
	}
	if childGroupID != `` {
		if err := g.stmtBucketForGroupID.QueryRow(childGroupID).Scan(
			&valChGroupBId,
		); err != nil {
			if err == sql.ErrNoRows {
				return true, fmt.Errorf("Unknown group %s", childGroupID)
			}
			return false, err
		}
	}

	switch q.Action {
	case `add_node_to_cluster`:
		if valNodeBId != valClusterBId {
			return false, fmt.Errorf(
				"Node and Cluster are in different buckets (%s/%s)",
				valNodeBId, valClusterBId,
			)
		}
	case `add_node_to_group`:
		if valNodeBId != valGroupBId {
			return false, fmt.Errorf(
				"Node and Group are in different buckets (%s/%s)",
				valNodeBId, valGroupBId,
			)
		}
	case `add_cluster_to_group`:
		if valClusterBId != valGroupBId {
			return false, fmt.Errorf(
				"Cluster and Group are in different buckets (%s/%s)",
				valClusterBId, valGroupBId,
			)
		}
	case `add_group_to_group`:
		if valChGroupBId != valGroupBId {
			return false, fmt.Errorf(
				"Groups are in different buckets (%s/%s)",
				valGroupBId, valChGroupBId,
			)
		}
	}
	return false, nil
}

// Verify that an object is assigned to the specified bucket.
func (g *GuidePost) validateCorrectBucket(q *msg.Request) (bool, error) {
	switch q.Action {
	case `assign_node`:
		return g.validateNodeUnassigned(q)
	case `create_cluster`, `create_group`:
		return false, nil
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
			return true, fmt.Errorf("%s is not assigned to any bucket",
				q.Section)
		}
		return false, err
	}
	switch q.Section {
	case `node`:
		if bid != q.Node.Config.BucketId {
			return false, fmt.Errorf("Node assigned to different bucket %s",
				bid)
		}
	case `cluster`:
		if bid != q.Cluster.BucketId {
			return false, fmt.Errorf("Cluster in different bucket %s",
				bid)
		}
	case `group`:
		if bid != q.Group.BucketId {
			return false, fmt.Errorf("Group in different bucket %s",
				bid)
		}
	}
	return false, nil
}

// Verify that a node is not yet assigned to a bucket. Returns nil
// on success.
func (g *GuidePost) validateNodeUnassigned(q *msg.Request) (bool, error) {
	var bid string
	if err := g.stmtBucketForNodeID.QueryRow(q.Node.Id).Scan(
		&bid,
	); err != nil {
		if err == sql.ErrNoRows {
			// unassigned is not an error here
			return false, nil
		}
		return false, err
	}
	return false, fmt.Errorf("Node already assigned to bucket %s", bid)
}

// Verify the node has a Config section
func (g *GuidePost) validateNodeConfig(q *msg.Request) (bool, error) {
	if q.Node.Config == nil {
		return false, fmt.Errorf("NodeConfig subobject missing")
	}
	return g.validateBucketInRepository(
		q.Node.Config.RepositoryId,
		q.Node.Config.BucketId,
	)
}

// Verify that the ObjectId->BucketId->RepositoryId chain is part of
// the same tree.
func (g *GuidePost) validateCheckObjectInBucket(q *msg.Request) (bool, error) {
	var err error
	var bid string
	switch q.CheckConfig.ObjectType {
	case `repository`:
		if q.CheckConfig.RepositoryId !=
			q.CheckConfig.ObjectId {
			return false, fmt.Errorf("Conflicting repository ids: %s, %s",
				q.CheckConfig.RepositoryId,
				q.CheckConfig.ObjectId,
			)
		}
		return false, nil
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
		return false, fmt.Errorf("Unknown object type: %s",
			q.CheckConfig.ObjectType,
		)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return true, fmt.Errorf("No bucketID for object found")
		}
		return false, err
	}
	if bid != q.CheckConfig.BucketId {
		return false, fmt.Errorf("Object is in bucket %s, not %s",
			bid, q.CheckConfig.BucketId,
		)
	}
	return g.validateBucketInRepository(
		q.CheckConfig.RepositoryId,
		q.CheckConfig.BucketId,
	)
}

// Verify that the bucket is part of the specified repository
func (g *GuidePost) validateBucketInRepository(
	repo, bucket string) (bool, error) {
	var repoID, repoName string
	if err := g.stmtRepoForBucketID.QueryRow(
		bucket,
	).Scan(
		&repoID,
		&repoName,
	); err != nil {
		if err == sql.ErrNoRows {
			return true, fmt.Errorf("No repository found for bucket %s",
				bucket)
		}
		return false, err
	}
	if repo != repoID {
		return false, fmt.Errorf("Bucket is in different repository: %s",
			repoID)
	}
	return false, nil
}

// check the check configuration to contain fewer thresholds than
// the limit for the capability
func (g *GuidePost) validateCheckThresholds(q *msg.Request) (bool, error) {
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
			return true, fmt.Errorf(
				"Capability %s not found",
				q.CheckConfig.CapabilityId)
		}
		return false, err
	}
	if len(q.CheckConfig.Thresholds) > thrLimit {
		return false, fmt.Errorf(
			"Specified %d thresholds exceed limit of %d for capability",
			len(q.CheckConfig.Thresholds),
			thrLimit)
	}
	return false, nil
}

// check the naming schema for the bucket (global unique object)
func (g *GuidePost) validateBucketName(q *msg.Request) (bool, error) {
	_, repoName, _, _ := g.extractRouting(q)

	if !strings.HasPrefix(
		q.Bucket.Name,
		fmt.Sprintf("%s_", repoName),
	) {
		return false, fmt.Errorf("Illegal bucket name format, " +
			"requires reponame_ prefix")
	}
	return false, nil
}

// validate current treekeeper state
func (g *GuidePost) validateKeeper(repoName string) (bool, error) {
	// check we have a treekeeper for that repository
	keeper := fmt.Sprintf("repository_%s", repoName)
	if _, ok := g.soma.handlerMap.Get(keeper).(*TreeKeeper); !ok {
		return true, fmt.Errorf(
			"No handler for repository %s currently registered.",
			repoName)
	}
	handler := g.soma.handlerMap.Get(keeper).(*TreeKeeper)

	// check the treekeeper has not been stopped
	if handler.isStopped() {
		return false, fmt.Errorf(
			"Repository %s is currently stopped.", repoName)
	}

	// check the treekeeper has not encountered a broken tree
	if handler.isBroken() {
		return false, fmt.Errorf(
			"Repository %s is broken.", repoName)
	}

	// check the treekeeper has finished loading
	if !handler.isReady() {
		return false, fmt.Errorf(
			"Repository %s not fully loaded yet.", repoName)
	}
	return false, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
