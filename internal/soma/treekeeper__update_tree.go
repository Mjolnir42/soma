/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/tree"
	"github.com/satori/go.uuid"
)

func (tk *TreeKeeper) treeBucket(q *msg.Request) {
	//XXX BUG convert to section/action model
	//XXX BUG validate bucket request
	//XXX BUG generate Bucket.UUID in Guidepost
	switch q.Action {
	case `create_bucket`:
		tree.NewBucket(tree.BucketSpec{
			ID:          uuid.NewV4().String(),
			Name:        q.Bucket.Name,
			Environment: q.Bucket.Environment,
			Team:        tk.meta.teamID,
			Deleted:     q.Bucket.IsDeleted,
			Frozen:      q.Bucket.IsFrozen,
			Repository:  q.Bucket.RepositoryID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `repository`,
			ParentID:   tk.meta.repoID,
			ParentName: tk.meta.repoName,
		})
	}
}

func (tk *TreeKeeper) treeGroup(q *msg.Request) {
	switch q.Action {
	case `create_group`:
		tree.NewGroup(tree.GroupSpec{
			ID:   uuid.NewV4().String(),
			Name: q.Group.Name,
			Team: tk.meta.teamID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `bucket`,
			ParentID:   q.Group.BucketID,
		})
	case `delete_group`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementID:   q.Group.ID,
		}, true).(tree.BucketAttacher).Destroy()
	case `reset_group_to_bucket`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementID:   q.Group.ID,
		}, true).(tree.BucketAttacher).Detach()
	case `add_group_to_group`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `group`,
			ElementID:   (*q.Group.MemberGroups)[0].ID,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `group`,
			ParentID:   q.Group.ID,
		})
	}
}

func (tk *TreeKeeper) treeCluster(q *msg.Request) {
	switch q.Action {
	case `create_cluster`:
		tree.NewCluster(tree.ClusterSpec{
			ID:   uuid.NewV4().String(),
			Name: q.Cluster.Name,
			Team: tk.meta.teamID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `bucket`,
			ParentID:   q.Cluster.BucketID,
		})
	case `delete_cluster`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementID:   q.Cluster.ID,
		}, true).(tree.BucketAttacher).Destroy()
	case `reset_cluster_to_bucket`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementID:   q.Cluster.ID,
		}, true).(tree.BucketAttacher).Detach()
	case `add_cluster_to_group`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `cluster`,
			ElementID:   (*q.Group.MemberClusters)[0].ID,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `group`,
			ParentID:   q.Group.ID,
		})
	}
}

func (tk *TreeKeeper) treeNode(q *msg.Request) {
	switch q.Action {
	case `assign_node`:
		tree.NewNode(tree.NodeSpec{
			ID:       q.Node.ID,
			AssetID:  q.Node.AssetID,
			Name:     q.Node.Name,
			Team:     q.Node.TeamID,
			ServerID: q.Node.ServerID,
			Online:   q.Node.IsOnline,
			Deleted:  q.Node.IsDeleted,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `bucket`,
			ParentID:   q.Node.Config.BucketID,
		})
	case `delete_node`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementID:   q.Node.ID,
		}, true).(tree.BucketAttacher).Destroy()
	case `reset_node_to_bucket`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementID:   q.Node.ID,
		}, true).(tree.BucketAttacher).Detach()
	case `add_node_to_group`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementID:   (*q.Group.MemberNodes)[0].ID,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `group`,
			ParentID:   q.Group.ID,
		})
	case `add_node_to_cluster`:
		tk.tree.Find(tree.FindRequest{
			ElementType: `node`,
			ElementID:   (*q.Cluster.Members)[0].ID,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: `cluster`,
			ParentID:   q.Cluster.ID,
		})
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
