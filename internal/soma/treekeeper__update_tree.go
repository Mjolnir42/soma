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

func (tk *TreeKeeper) treeRepository(q *msg.Request) {
	if q.Section == msg.SectionRepository && q.Action == msg.ActionRename {
		tk.tree.Find(tree.FindRequest{
			ElementID:   q.Repository.ID,
			ElementType: `repository`,
		}, true).SetName(
			q.Update.Repository.Name,
		)
	}
}

func (tk *TreeKeeper) treeBucket(q *msg.Request) {
	//XXX BUG validate bucket request
	//XXX BUG generate Bucket.UUID in Guidepost
	switch {
	case q.Section == msg.SectionBucket && q.Action == msg.ActionCreate:
		tree.NewBucket(tree.BucketSpec{
			ID:          uuid.Must(uuid.NewV4()).String(),
			Name:        q.Bucket.Name,
			Environment: q.Bucket.Environment,
			Team:        tk.meta.teamID,
			Deleted:     q.Bucket.IsDeleted,
			Frozen:      q.Bucket.IsFrozen,
			Repository:  q.Bucket.RepositoryID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: msg.EntityRepository,
			ParentID:   tk.meta.repoID,
			ParentName: tk.meta.repoName,
		})
	case q.Section == msg.SectionBucket && q.Action == msg.ActionRename:
		tk.tree.Find(tree.FindRequest{
			ElementID:   q.Bucket.ID,
			ElementType: `bucket`,
		}, true).SetName(
			q.Update.Bucket.Name,
		)
	}
}

func (tk *TreeKeeper) treeGroup(q *msg.Request) {
	if q.Section == msg.SectionGroup {
		switch q.Action {
		case msg.ActionCreate:
			tree.NewGroup(tree.GroupSpec{
				ID:   uuid.Must(uuid.NewV4()).String(),
				Name: q.Group.Name,
				Team: tk.meta.teamID,
			}).Attach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: msg.EntityBucket,
				ParentID:   q.Group.BucketID,
			})
		case msg.ActionDestroy:
			tk.tree.Find(tree.FindRequest{
				ElementType: msg.EntityGroup,
				ElementID:   q.Group.ID,
			}, true).(tree.BucketAttacher).Destroy()
		}
	}

	if q.Action == msg.ActionMemberUnassign && q.TargetEntity == msg.EntityGroup {
		switch q.Section {
		case msg.SectionGroup:
			tk.tree.Find(tree.FindRequest{
				ElementType: msg.EntityGroup,
				ElementID:   (*q.Group.MemberGroups)[0].ID,
			}, true).(tree.BucketAttacher).Detach()
		}
	}

	if q.Action == msg.ActionMemberAssign && q.TargetEntity == msg.EntityGroup {
		switch q.Section {
		case msg.SectionGroup:
			tk.tree.Find(tree.FindRequest{
				ElementType: msg.EntityGroup,
				ElementID:   (*q.Group.MemberGroups)[0].ID,
			}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: msg.EntityGroup,
				ParentID:   q.Group.ID,
			})
		}
	}
}

func (tk *TreeKeeper) treeCluster(q *msg.Request) {
	if q.Section == msg.SectionCluster {
		switch q.Action {
		case msg.ActionCreate:
			tree.NewCluster(tree.ClusterSpec{
				ID:   uuid.Must(uuid.NewV4()).String(),
				Name: q.Cluster.Name,
				Team: tk.meta.teamID,
			}).Attach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: msg.EntityBucket,
				ParentID:   q.Cluster.BucketID,
			})
		case msg.ActionDestroy:
			tk.tree.Find(tree.FindRequest{
				ElementType: msg.EntityCluster,
				ElementID:   q.Cluster.ID,
			}, true).(tree.BucketAttacher).Destroy()
		}
	}

	if q.Action == msg.ActionMemberUnassign && q.TargetEntity == msg.EntityCluster {
		switch q.Section {
		case msg.SectionGroup:
			tk.tree.Find(tree.FindRequest{
				ElementType: msg.EntityCluster,
				ElementID:   (*q.Group.MemberClusters)[0].ID,
			}, true).(tree.BucketAttacher).Detach()
		}
	}

	if q.Action == msg.ActionMemberAssign && q.TargetEntity == msg.EntityCluster {
		switch q.Section {
		case msg.SectionGroup:
			tk.tree.Find(tree.FindRequest{
				ElementType: msg.EntityCluster,
				ElementID:   (*q.Group.MemberClusters)[0].ID,
			}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: msg.EntityGroup,
				ParentID:   q.Group.ID,
			})
		}
	}
}

func (tk *TreeKeeper) treeNode(q *msg.Request) {
	if q.Section == msg.SectionNodeConfig {
		switch q.Action {
		case msg.ActionAssign:
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
				ParentType: msg.EntityBucket,
				ParentID:   q.Node.Config.BucketID,
			})
		case msg.ActionUnassign:
			tk.tree.Find(tree.FindRequest{
				ElementType: msg.EntityNode,
				ElementID:   q.Node.ID,
			}, true).(tree.BucketAttacher).Destroy()
		}
	}

	if q.Action == msg.ActionMemberUnassign && q.TargetEntity == msg.EntityNode {
		switch q.Section {
		case msg.SectionCluster:
			tk.tree.Find(tree.FindRequest{
				ElementType: msg.EntityNode,
				ElementID:   (*q.Cluster.Members)[0].ID,
			}, true).(tree.BucketAttacher).Detach()
		case msg.SectionGroup:
			tk.tree.Find(tree.FindRequest{
				ElementType: msg.EntityNode,
				ElementID:   (*q.Group.MemberNodes)[0].ID,
			}, true).(tree.BucketAttacher).Detach()
		}
	}

	if q.Action == msg.ActionMemberAssign && q.TargetEntity == msg.EntityNode {
		switch q.Section {
		case msg.SectionGroup:
			tk.tree.Find(tree.FindRequest{
				ElementType: msg.EntityNode,
				ElementID:   (*q.Group.MemberNodes)[0].ID,
			}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: msg.EntityGroup,
				ParentID:   q.Group.ID,
			})
		case msg.SectionCluster:
			tk.tree.Find(tree.FindRequest{
				ElementType: msg.EntityNode,
				ElementID:   (*q.Cluster.Members)[0].ID,
			}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: msg.EntityCluster,
				ParentID:   q.Cluster.ID,
			})
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
