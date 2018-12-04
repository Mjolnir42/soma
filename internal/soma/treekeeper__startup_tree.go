package soma

import (
	"database/sql"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/tree"
	"github.com/mjolnir42/soma/lib/proto"
)

func (tk *TreeKeeper) startupBuckets(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		rows                                      *sql.Rows
		bucketID, bucketName, environment, teamID string
		frozen, deleted                           bool
		err                                       error
	)

	tk.startLog.Printf("TK[%s]: loading buckets", tk.meta.repoName)
	rows, err = stMap[`LoadBucket`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading buckets: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

	super := tk.soma.getSupervisor()

bucketloop:
	for rows.Next() {
		err = rows.Scan(
			&bucketID,
			&bucketName,
			&frozen,
			&deleted,
			&environment,
			&teamID,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break bucketloop
			}
			tk.startLog.Printf("TK[%s] Error: %s", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}
		tree.NewBucket(tree.BucketSpec{
			ID:          bucketID,
			Name:        bucketName,
			Environment: environment,
			Team:        teamID,
			Deleted:     deleted,
			Frozen:      frozen,
			Repository:  tk.meta.repoID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "repository",
			ParentID:   tk.meta.repoID,
			ParentName: tk.meta.repoName,
		})
		tk.drain(`action`)
		tk.drain(`error`)

		// very explicitly ensure that the go routine is receiving
		// actual copies of the value of the strings updated in rows.Next()
		bucket := proto.Bucket{
			ID:           bucketID,
			Name:         bucketName,
			TeamID:       teamID,
			RepositoryID: tk.meta.repoID,
			Environment:  environment,
		}
		req := msg.Request{
			Section: msg.SectionBucket,
			Action:  msg.ActionCreate,
			Bucket:  bucket.Clone(),
		}
		go func(q *msg.Request) {
			super.Update <- msg.CacheUpdateFromRequest(q)
		}(&req)
	}
}

func (tk *TreeKeeper) startupGroups(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		rows                                 *sql.Rows
		groupID, groupName, bucketID, teamID string
		err                                  error
	)

	tk.startLog.Printf("TK[%s]: loading groups", tk.meta.repoName)
	rows, err = stMap[`LoadGroup`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading groups: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

	super := tk.soma.getSupervisor()

grouploop:
	for rows.Next() {
		err = rows.Scan(
			&groupID,
			&groupName,
			&bucketID,
			&teamID,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break grouploop
			}
			tk.startLog.Printf("TK[%s] Error: %s", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}
		tree.NewGroup(tree.GroupSpec{
			ID:   groupID,
			Name: groupName,
			Team: teamID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentID:   bucketID,
		})
		tk.drain(`action`)
		tk.drain(`error`)

		// very explicitly ensure that the go routine is receiving
		// actual copies of the value of the strings updated in rows.Next()
		group := proto.Group{
			ID:           groupID,
			Name:         groupName,
			TeamID:       teamID,
			RepositoryID: tk.meta.repoID,
			BucketID:     bucketID,
		}
		req := msg.Request{
			Section: msg.SectionGroup,
			Action:  msg.ActionCreate,
			Group:   group.Clone(),
		}
		go func(q *msg.Request) {
			super.Update <- msg.CacheUpdateFromRequest(q)
		}(&req)
	}
}

func (tk *TreeKeeper) startupGroupMemberGroups(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		rows                  *sql.Rows
		groupID, childGroupID string
		err                   error
	)

	tk.startLog.Printf("TK[%s]: loading group-member-groups", tk.meta.repoName)
	rows, err = stMap[`LoadGroupMbrGroup`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading groups: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

memberloop:
	for rows.Next() {
		err = rows.Scan(
			&groupID,
			&childGroupID,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break memberloop
			}
			tk.startLog.Printf("TK[%s] Error: %s", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		tk.tree.Find(tree.FindRequest{
			ElementType: "group",
			ElementID:   childGroupID,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentID:   groupID,
		})
		tk.drain(`action`)
		tk.drain(`error`)
	}
}

func (tk *TreeKeeper) startupGroupedClusters(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err                                     error
		rows                                    *sql.Rows
		clusterID, clusterName, teamID, groupID string
		bucketID                                string
	)

	tk.startLog.Printf("TK[%s]: loading grouped-clusters", tk.meta.repoName)
	rows, err = stMap[`LoadGroupMbrCluster`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading clusters: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

	super := tk.soma.getSupervisor()

clusterloop:
	for rows.Next() {
		err = rows.Scan(
			&clusterID,
			&clusterName,
			&teamID,
			&groupID,
			&bucketID,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break clusterloop
			}
			tk.startLog.Printf("TK[%s] Error: %s", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		tree.NewCluster(tree.ClusterSpec{
			ID:   clusterID,
			Name: clusterName,
			Team: teamID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentID:   groupID,
		})
		tk.drain(`action`)
		tk.drain(`error`)

		// very explicitly ensure that the go routine is receiving
		// actual copies of the value of the strings updated in rows.Next()
		cluster := proto.Cluster{
			ID:           clusterID,
			Name:         clusterName,
			TeamID:       teamID,
			RepositoryID: tk.meta.repoID,
			BucketID:     bucketID,
		}
		req := msg.Request{
			Section: msg.SectionCluster,
			Action:  msg.ActionCreate,
			Cluster: cluster.Clone(),
		}
		go func(q *msg.Request) {
			super.Update <- msg.CacheUpdateFromRequest(q)
		}(&req)
	}
}

func (tk *TreeKeeper) startupClusters(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err                                      error
		rows                                     *sql.Rows
		clusterID, clusterName, bucketID, teamID string
	)

	tk.startLog.Printf("TK[%s]: loading clusters", tk.meta.repoName)
	rows, err = stMap[`LoadCluster`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading clusters: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

	super := tk.soma.getSupervisor()

clusterloop:
	for rows.Next() {
		err = rows.Scan(
			&clusterID,
			&clusterName,
			&bucketID,
			&teamID,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break clusterloop
			}
			tk.startLog.Printf("TK[%s] Error: %s", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		tree.NewCluster(tree.ClusterSpec{
			ID:   clusterID,
			Name: clusterName,
			Team: teamID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentID:   bucketID,
		})
		tk.drain(`action`)
		tk.drain(`error`)

		cluster := proto.Cluster{
			ID:           clusterID,
			Name:         clusterName,
			TeamID:       teamID,
			RepositoryID: tk.meta.repoID,
			BucketID:     bucketID,
		}
		req := msg.Request{
			Section: msg.SectionCluster,
			Action:  msg.ActionCreate,
			Cluster: cluster.Clone(),
		}
		go func(q *msg.Request) {
			super.Update <- msg.CacheUpdateFromRequest(q)
		}(&req)
	}
}

func (tk *TreeKeeper) startupNodes(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err                                          error
		rows                                         *sql.Rows
		nodeID, nodeName, teamID, serverID, bucketID string
		assetID                                      int
		nodeOnline, nodeDeleted                      bool
		clusterID, groupID                           sql.NullString
	)

	tk.startLog.Printf("TK[%s]: loading nodes", tk.meta.repoName)
	rows, err = stMap[`LoadNode`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading nodes: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

	super := tk.soma.getSupervisor()

nodeloop:
	for rows.Next() {
		err = rows.Scan(
			&nodeID,
			&assetID,
			&nodeName,
			&teamID,
			&serverID,
			&nodeOnline,
			&nodeDeleted,
			&bucketID,
			&clusterID,
			&groupID,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break nodeloop
			}
			tk.startLog.Printf("TK[%s] Error: %s", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		node := tree.NewNode(tree.NodeSpec{
			ID:       nodeID,
			AssetID:  uint64(assetID),
			Name:     nodeName,
			Team:     teamID,
			ServerID: serverID,
			Online:   nodeOnline,
			Deleted:  nodeDeleted,
		})
		if clusterID.Valid {
			node.Attach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: "cluster",
				ParentID:   clusterID.String,
			})
		} else if groupID.Valid {
			node.Attach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: "group",
				ParentID:   groupID.String,
			})
		} else {
			node.Attach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: "bucket",
				ParentID:   bucketID,
			})
		}
		tk.drain(`action`)
		tk.drain(`error`)

		go func() {
			super.Update <- msg.CacheUpdateFromRequest(&msg.Request{
				Section: msg.SectionNodeConfig,
				Action:  msg.ActionAssign,
				Node: proto.Node{
					ID:        nodeID,
					AssetID:   uint64(assetID),
					Name:      nodeName,
					TeamID:    teamID,
					ServerID:  serverID,
					IsOnline:  nodeOnline,
					IsDeleted: nodeDeleted,
					Config: &proto.NodeConfig{
						RepositoryID: tk.meta.repoID,
						BucketID:     bucketID,
					},
				},
			})
		}()
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
