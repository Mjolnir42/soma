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

	"github.com/mjolnir42/soma/internal/msg"
)

// Extract the request routing information
func (g *GuidePost) extractRouting(q *msg.Request) (string, string, bool, error) {
	var repoID, repoName, bucketID string
	var err error

	repoID, bucketID = g.extractID(q)

	// lookup repository by bucket
	if bucketID != `` {
		if err = g.stmtRepoForBucketID.QueryRow(
			bucketID,
		).Scan(
			&repoID,
			&repoName,
		); err != nil {
			if err == sql.ErrNoRows {
				return ``, ``, true, fmt.Errorf(
					"No repository found for bucketID %s",
					bucketID,
				)
			}
			return ``, ``, false, err
		}
	}

	// lookup repository name
	if repoName == `` && repoID != `` {
		if err = g.stmtRepoNameByID.QueryRow(
			repoID,
		).Scan(
			&repoName,
		); err != nil {
			if err == sql.ErrNoRows {
				return ``, ``, true, fmt.Errorf(
					"No repository found with id %s",
					repoID,
				)
			}
			return ``, ``, false, err
		}
	}

	if repoName == `` {
		return ``, ``, true, fmt.Errorf(
			`GuidePost: unable find repository for request`,
		)
	}
	return repoID, repoName, false, nil
}

// Extract embedded IDs that can be used for routing
func (g *GuidePost) extractID(q *msg.Request) (string, string) {
	switch q.Action {
	case
		`add_system_property_to_repository`,
		`add_custom_property_to_repository`,
		`add_oncall_property_to_repository`,
		`add_service_property_to_repository`,
		`delete_system_property_from_repository`,
		`delete_custom_property_from_repository`,
		`delete_oncall_property_from_repository`,
		`delete_service_property_from_repository`:
		return q.Repository.Id, ``
	case
		`create_bucket`:
		return q.Bucket.RepositoryID, ``
	case
		`add_system_property_to_bucket`,
		`add_custom_property_to_bucket`,
		`add_oncall_property_to_bucket`,
		`add_service_property_to_bucket`,
		`delete_system_property_from_bucket`,
		`delete_custom_property_from_bucket`,
		`delete_oncall_property_from_bucket`,
		`delete_service_property_from_bucket`:
		return ``, q.Bucket.ID
	case
		`add_group_to_group`,
		`add_cluster_to_group`,
		`add_node_to_group`,
		`create_group`,
		`add_system_property_to_group`,
		`add_custom_property_to_group`,
		`add_oncall_property_to_group`,
		`add_service_property_to_group`,
		`delete_system_property_from_group`,
		`delete_custom_property_from_group`,
		`delete_oncall_property_from_group`,
		`delete_service_property_from_group`:
		return ``, q.Group.BucketId
	case
		`add_node_to_cluster`,
		`create_cluster`,
		`add_system_property_to_cluster`,
		`add_custom_property_to_cluster`,
		`add_oncall_property_to_cluster`,
		`add_service_property_to_cluster`,
		`delete_system_property_from_cluster`,
		`delete_custom_property_from_cluster`,
		`delete_oncall_property_from_cluster`,
		`delete_service_property_from_cluster`:
		return ``, q.Cluster.BucketId
	case
		`add_check_to_repository`,
		`add_check_to_bucket`,
		`add_check_to_group`,
		`add_check_to_cluster`,
		`add_check_to_node`,
		`remove_check`:
		return q.CheckConfig.RepositoryId, ``
	case
		`assign_node`,
		`add_system_property_to_node`,
		`add_custom_property_to_node`,
		`add_oncall_property_to_node`,
		`add_service_property_to_node`,
		`delete_system_property_from_node`,
		`delete_custom_property_from_node`,
		`delete_oncall_property_from_node`,
		`delete_service_property_from_node`:
		return q.Node.Config.RepositoryId, q.Node.Config.BucketId
	}
	return ``, ``
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
