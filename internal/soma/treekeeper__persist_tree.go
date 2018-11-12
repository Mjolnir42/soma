/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"fmt"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/tree"
)

func (tk *TreeKeeper) txTree(a *tree.Action,
	stm map[string]*sql.Stmt, user string) error {
	switch a.Action {
	case tree.ActionCreate:
		return tk.txTreeCreate(a, stm, user)
	case tree.ActionRename:
		return tk.txTreeRename(a, stm, user)
	case tree.ActionUpdate:
		return tk.txTreeUpdate(a, stm)
	case tree.ActionDelete:
		return tk.txTreeDelete(a, stm)
	case tree.ActionMemberNew, tree.ActionNodeAssignment:
		return tk.txTreeMemberNew(a, stm)
	case tree.ActionMemberRemoved:
		return tk.txTreeMemberRemoved(a, stm)
	default:
		return fmt.Errorf("Illegal tree action: %s", a.Action)
	}
}

func (tk *TreeKeeper) txTreeCreate(a *tree.Action,
	stm map[string]*sql.Stmt, user string) error {
	var err error
	switch a.Type {
	case msg.EntityBucket:
		_, err = stm[`CreateBucket`].Exec(
			a.Bucket.ID,
			a.Bucket.Name,
			a.Bucket.IsFrozen,
			a.Bucket.IsDeleted,
			a.Bucket.RepositoryID,
			a.Bucket.Environment,
			a.Bucket.TeamID,
			user,
		)
	case msg.EntityGroup:
		_, err = stm[`GroupCreate`].Exec(
			a.Group.ID,
			a.Group.BucketID,
			a.Group.Name,
			a.Group.ObjectState,
			a.Group.TeamID,
			user,
		)
	case msg.EntityCluster:
		_, err = stm[`ClusterCreate`].Exec(
			a.Cluster.ID,
			a.Cluster.Name,
			a.Cluster.BucketID,
			a.Cluster.ObjectState,
			a.Cluster.TeamID,
			user,
		)
	}
	return err
}

func (tk *TreeKeeper) txTreeRename(a *tree.Action,
	stm map[string]*sql.Stmt, user string) error {
	var err error
	switch a.Type {
	case msg.EntityRepository:
		_, err = stm[`repository::rename`].Exec(
			a.Repository.ID,
			a.Repository.Name,
			user,
		)
	case msg.EntityBucket:
		_, err = stm[`bucket::rename`].Exec(
			a.Bucket.ID,
			a.Bucket.Name,
			user,
		)
	}
	return err
}

func (tk *TreeKeeper) txTreeUpdate(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var (
		err          error
		statement    *sql.Stmt
		id, newState string
	)
	switch a.Type {
	case msg.EntityGroup:
		statement = stm[`GroupUpdate`]
		id = a.Group.ID
		newState = a.Group.ObjectState
	case msg.EntityCluster:
		statement = stm[`ClusterUpdate`]
		id = a.Cluster.ID
		newState = a.Cluster.ObjectState
	case msg.EntityNode:
		statement = stm[`UpdateNodeState`]
		id = a.Node.ID
		newState = a.Node.State
	}
	_, err = statement.Exec(
		id,
		newState,
	)
	return err
}

func (tk *TreeKeeper) txTreeDelete(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var err error
	switch a.Type {
	case msg.EntityRepository:
		_, err = stm[`repository::destroy`].Exec(
			a.Repository.ID,
		)
	case msg.EntityBucket:
		_, err = stm[`bucket::destroy`].Exec(
			a.Bucket.ID,
		)
	case msg.EntityGroup:
		_, err = stm[`GroupDelete`].Exec(
			a.Group.ID,
		)
	case msg.EntityCluster:
		_, err = stm[`ClusterDelete`].Exec(
			a.Cluster.ID,
		)
	case msg.EntityNode:
		if _, err = stm[`NodeUnassignFromBucket`].Exec(
			a.Node.ID,
			a.Node.Config.BucketID,
			a.Node.TeamID,
		); err != nil {
			return err
		}
		// node unassign requires state update
		err = tk.txTreeUpdate(a, stm)
	}
	return err
}

func (tk *TreeKeeper) txTreeMemberNew(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var (
		err               error
		id, child, bucket string
		statement         *sql.Stmt
	)
	switch a.Type {
	case msg.EntityBucket:
		_, err = stm[`BucketAssignNode`].Exec(
			a.ChildNode.ID,
			a.Bucket.ID,
			a.Bucket.TeamID,
		)
		goto exit
	case msg.EntityGroup:
		id = a.Group.ID
		bucket = a.Group.BucketID
		switch a.ChildType {
		case msg.EntityGroup:
			statement = stm[`GroupMemberNewGroup`]
			child = a.ChildGroup.ID
		case msg.EntityCluster:
			statement = stm[`GroupMemberNewCluster`]
			child = a.ChildCluster.ID
		case msg.EntityNode:
			statement = stm[`GroupMemberNewNode`]
			child = a.ChildNode.ID
		}
	case msg.EntityCluster:
		id = a.Cluster.ID
		bucket = a.Cluster.BucketID
		child = a.ChildNode.ID
		statement = stm[`ClusterMemberNew`]
	}
	_, err = statement.Exec(
		id,
		child,
		bucket,
	)
exit:
	return err
}

func (tk *TreeKeeper) txTreeMemberRemoved(a *tree.Action,
	stm map[string]*sql.Stmt) error {
	var (
		err       error
		id, child string
		statement *sql.Stmt
	)
	switch a.Type {
	case msg.EntityGroup:
		id = a.Group.ID
		switch a.ChildType {
		case msg.EntityGroup:
			statement = stm[`GroupMemberRemoveGroup`]
			child = a.ChildGroup.ID
		case msg.EntityCluster:
			statement = stm[`GroupMemberRemoveCluster`]
			child = a.ChildCluster.ID
		case msg.EntityNode:
			statement = stm[`GroupMemberRemoveNode`]
			child = a.ChildNode.ID
		}
	case msg.EntityCluster:
		id = a.Cluster.ID
		child = a.ChildNode.ID
		statement = stm[`ClusterMemberRemove`]
	}
	_, err = statement.Exec(
		id,
		child,
	)
	return err
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
