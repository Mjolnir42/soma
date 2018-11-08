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
	case msg.SectionCheckConfig:
		if nf, err := g.validateCheckObjectInBucket(q); err != nil {
			return nf, err
		}
	case msg.SectionNodeConfig:
		if nf, err := g.validateNodeConfig(q); err != nil {
			return nf, err
		}
		fallthrough
	case msg.SectionCluster, msg.SectionGroup:
		if nf, err := g.validateCorrectBucket(q); err != nil {
			return nf, err
		}
	case msg.SectionBucket:
		if nf, err := g.validateBucketInRepository(
			q.Bucket.RepositoryID,
			q.Bucket.ID,
		); err != nil {
			return nf, err
		}
	case msg.SectionRepositoryConfig, msg.SectionRepository:
		// since repository ids are the routing information,
		// it is unnecessary to check that the object is where the
		// routing would point to
	default:
		return false, fmt.Errorf("Invalid request type %s", q.Section)
	}

	switch q.Action {
	case msg.ActionMemberAssign:
		return g.validateObjectMatch(q)
	}

	switch q.Section {
	case msg.SectionCheckConfig:
		switch q.Action {
		case msg.ActionCreate:
			return g.validateCheckThresholds(q)
		}
	case msg.SectionBucket:
		switch q.Action {
		case msg.ActionCreate:
			return g.validateBucketName(q)
		case msg.ActionRename:
			return g.validateBucketName(q)
		}
	}

	// listed actions are accepted, but require no further validation
	switch q.Action {
	case msg.ActionPropertyCreate, msg.ActionPropertyDestroy:
		switch q.Section {
		case
			msg.SectionRepositoryConfig,
			msg.SectionBucket,
			msg.SectionGroup,
			msg.SectionCluster,
			msg.SectionNodeConfig:
			return false, nil
		}
	case msg.ActionAssign:
		switch q.Section {
		case msg.SectionNodeConfig:
			return false, nil
		}
	case msg.ActionCreate:
		switch q.Section {
		case msg.SectionGroup, msg.SectionCluster:
			return false, nil
		}
	case msg.ActionDestroy:
		switch q.Section {
		case msg.SectionCheckConfig:
			return false, nil
		}
	case msg.ActionRename:
		switch q.Section {
		case msg.SectionRepository:
			return false, nil
		}
	}
	return false, fmt.Errorf("Unimplemented guidepost/%s::%s", q.Section, q.Action)
}

func (g *GuidePost) validateObjectMatch(q *msg.Request) (bool, error) {
	var (
		nodeID, clusterID, groupID, childGroupID              string
		valNodeBId, valClusterBId, valGroupBId, valChGroupBId string
	)

	switch q.Action {
	case msg.ActionMemberAssign:
		switch q.Section {
		case msg.SectionCluster:
			switch q.TargetEntity {
			case msg.EntityNode:
				nodeID = (*q.Cluster.Members)[0].ID
				clusterID = q.Cluster.ID
			default:
				return false, fmt.Errorf("Incorrect validation attempted for %s::%s(%s)",
					q.Section, q.Action, q.TargetEntity)
			}
		case msg.SectionGroup:
			switch q.TargetEntity {
			case msg.EntityNode:
				nodeID = (*q.Group.MemberNodes)[0].ID
				groupID = q.Group.ID
			case msg.EntityCluster:
				clusterID = (*q.Group.MemberClusters)[0].ID
				groupID = q.Group.ID
			case msg.EntityGroup:
				childGroupID = (*q.Group.MemberGroups)[0].ID
				groupID = q.Group.ID
			default:
				return false, fmt.Errorf("Incorrect validation attempted for %s::%s(%s)",
					q.Section, q.Action, q.TargetEntity)
			}
		default:
			return false, fmt.Errorf("Incorrect validation attempted for %s::%s(%s)",
				q.Section, q.Action, q.TargetEntity)
		}
	default:
		return false, fmt.Errorf("Incorrect validation attempted for %s::%s(%s)",
			q.Section, q.Action, q.TargetEntity)
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

	if q.Section == msg.SectionCluster && q.Action == msg.ActionMemberAssign && q.TargetEntity == msg.EntityNode {
		if valNodeBId != valClusterBId {
			return false, fmt.Errorf(
				"Node and Cluster are in different buckets (%s/%s)",
				valNodeBId, valClusterBId,
			)
		}
	}
	if q.Section == msg.SectionGroup && q.Action == msg.ActionMemberAssign {
		switch q.TargetEntity {
		case msg.EntityNode:
			if valNodeBId != valGroupBId {
				return false, fmt.Errorf(
					"Node and Group are in different buckets (%s/%s)",
					valNodeBId, valGroupBId,
				)
			}
		case msg.EntityCluster:
			if valClusterBId != valGroupBId {
				return false, fmt.Errorf(
					"Cluster and Group are in different buckets (%s/%s)",
					valClusterBId, valGroupBId,
				)
			}
		case msg.EntityGroup:
			if valChGroupBId != valGroupBId {
				return false, fmt.Errorf(
					"Groups are in different buckets (%s/%s)",
					valGroupBId, valChGroupBId,
				)
			}
		}
	}
	return false, nil
}

// Verify that an object is assigned to the specified bucket.
func (g *GuidePost) validateCorrectBucket(q *msg.Request) (bool, error) {
	switch q.Section {
	case msg.SectionCluster:
		switch q.Action {
		case msg.ActionCreate:
			return false, nil
		}
	case msg.SectionGroup:
		switch q.Action {
		case msg.ActionCreate:
			return false, nil
		}
	case msg.SectionNodeConfig:
		switch q.Action {
		case msg.ActionAssign:
			return g.validateNodeUnassigned(q)
		}
	}

	var bid string
	var err error
	switch q.Section {
	case msg.SectionNodeConfig:
		err = g.stmtBucketForNodeID.QueryRow(
			q.Node.ID,
		).Scan(
			&bid,
		)
	case msg.SectionCluster:
		err = g.stmtBucketForClusterID.QueryRow(
			q.Cluster.ID,
		).Scan(
			&bid,
		)
	case msg.SectionGroup:
		err = g.stmtBucketForGroupID.QueryRow(
			q.Group.ID,
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
	case msg.SectionNodeConfig:
		if bid != q.Node.Config.BucketID {
			return false, fmt.Errorf("Node assigned to different bucket %s",
				bid)
		}
	case msg.SectionCluster:
		if bid != q.Cluster.BucketID {
			return false, fmt.Errorf("Cluster in different bucket %s",
				bid)
		}
	case msg.SectionGroup:
		if bid != q.Group.BucketID {
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
	if err := g.stmtBucketForNodeID.QueryRow(q.Node.ID).Scan(
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
		q.Node.Config.RepositoryID,
		q.Node.Config.BucketID,
	)
}

// Verify that the ObjectId->BucketId->RepositoryId chain is part of
// the same tree.
func (g *GuidePost) validateCheckObjectInBucket(q *msg.Request) (bool, error) {
	var err error
	var bid string
	switch q.CheckConfig.ObjectType {
	case msg.EntityRepository:
		if q.CheckConfig.RepositoryID !=
			q.CheckConfig.ObjectID {
			return false, fmt.Errorf("Conflicting repository ids: %s, %s",
				q.CheckConfig.RepositoryID,
				q.CheckConfig.ObjectID,
			)
		}
		return false, nil
	case msg.EntityBucket:
		bid = q.CheckConfig.ObjectID
	case msg.EntityGroup:
		err = g.stmtBucketForGroupID.QueryRow(
			q.CheckConfig.ObjectID,
		).Scan(&bid)
	case msg.EntityCluster:
		err = g.stmtBucketForClusterID.QueryRow(
			q.CheckConfig.ObjectID,
		).Scan(&bid)
	case msg.EntityNode:
		err = g.stmtBucketForNodeID.QueryRow(
			q.CheckConfig.ObjectID,
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
	if bid != q.CheckConfig.BucketID {
		return false, fmt.Errorf("Object is in bucket %s, not %s",
			bid, q.CheckConfig.BucketID,
		)
	}
	return g.validateBucketInRepository(
		q.CheckConfig.RepositoryID,
		q.CheckConfig.BucketID,
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
		q.CheckConfig.CapabilityID,
	).Scan(
		&thrLimit,
	); err != nil {
		if err == sql.ErrNoRows {
			return true, fmt.Errorf(
				"Capability %s not found",
				q.CheckConfig.CapabilityID)
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

	var name string
	switch q.Action {
	case msg.ActionCreate:
		name = q.Bucket.Name
	case msg.ActionRename:
		name = q.Update.Bucket.Name
	default:
		name = `...INVALID....`
	}

	if !strings.HasPrefix(
		name,
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
