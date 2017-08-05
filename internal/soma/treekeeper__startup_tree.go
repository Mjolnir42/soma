package soma

import (
	"database/sql"

	"github.com/mjolnir42/soma/internal/tree"
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

	tk.startLog.Printf("TK[%s]: loading buckets\n", tk.meta.repoName)
	rows, err = stMap[`LoadBucket`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading buckets: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}
		tree.NewBucket(tree.BucketSpec{
			Id:          bucketID,
			Name:        bucketName,
			Environment: environment,
			Team:        teamID,
			Deleted:     deleted,
			Frozen:      frozen,
			Repository:  tk.meta.repoID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "repository",
			ParentId:   tk.meta.repoID,
			ParentName: tk.meta.repoName,
		})
		tk.drain(`action`)
		tk.drain(`error`)
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

	tk.startLog.Printf("TK[%s]: loading groups\n", tk.meta.repoName)
	rows, err = stMap[`LoadGroup`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading groups: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}
		tree.NewGroup(tree.GroupSpec{
			Id:   groupID,
			Name: groupName,
			Team: teamID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentId:   bucketID,
		})
		tk.drain(`action`)
		tk.drain(`error`)
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

	tk.startLog.Printf("TK[%s]: loading group-member-groups\n", tk.meta.repoName)
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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		tk.tree.Find(tree.FindRequest{
			ElementType: "group",
			ElementId:   childGroupID,
		}, true).(tree.BucketAttacher).ReAttach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   groupID,
		})
	}
	tk.drain(`action`)
	tk.drain(`error`)
}

func (tk *TreeKeeper) startupGroupedClusters(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err                                     error
		rows                                    *sql.Rows
		clusterID, clusterName, teamID, groupID string
	)

	tk.startLog.Printf("TK[%s]: loading grouped-clusters\n", tk.meta.repoName)
	rows, err = stMap[`LoadGroupMbrCluster`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading clusters: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

clusterloop:
	for rows.Next() {
		err = rows.Scan(
			&clusterID,
			&clusterName,
			&teamID,
			&groupID,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				break clusterloop
			}
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		tree.NewCluster(tree.ClusterSpec{
			Id:   clusterID,
			Name: clusterName,
			Team: teamID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "group",
			ParentId:   groupID,
		})
	}
	tk.drain(`action`)
	tk.drain(`error`)
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

	tk.startLog.Printf("TK[%s]: loading clusters\n", tk.meta.repoName)
	rows, err = stMap[`LoadCluster`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading clusters: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		tree.NewCluster(tree.ClusterSpec{
			Id:   clusterID,
			Name: clusterName,
			Team: teamID,
		}).Attach(tree.AttachRequest{
			Root:       tk.tree,
			ParentType: "bucket",
			ParentId:   bucketID,
		})
	}
	tk.drain(`action`)
	tk.drain(`error`)
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

	tk.startLog.Printf("TK[%s]: loading nodes\n", tk.meta.repoName)
	rows, err = stMap[`LoadNode`].Query(tk.meta.repoID)
	if err != nil {
		tk.startLog.Printf("TK[%s] Error loading nodes: %s", tk.meta.repoName, err.Error())
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

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
			tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
			tk.status.isBroken = true
			return
		}

		node := tree.NewNode(tree.NodeSpec{
			Id:       nodeID,
			AssetId:  uint64(assetID),
			Name:     nodeName,
			Team:     teamID,
			ServerId: serverID,
			Online:   nodeOnline,
			Deleted:  nodeDeleted,
		})
		if clusterID.Valid {
			node.Attach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: "cluster",
				ParentId:   clusterID.String,
			})
		} else if groupID.Valid {
			node.Attach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: "group",
				ParentId:   groupID.String,
			})
		} else {
			node.Attach(tree.AttachRequest{
				Root:       tk.tree,
				ParentType: "bucket",
				ParentId:   bucketID,
			})
		}
	}
	tk.drain(`action`)
	tk.drain(`error`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
