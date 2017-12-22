/*-
Copyright (c) 2016-2017, Jörg Pernfuß <code.jpe@gmail.com>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package soma

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

// TreeRead handles all read requests for deep tree exports
type TreeRead struct {
	Input    chan msg.Request
	Shutdown chan struct{}
	conn     *sql.DB
	// object details
	stmtShowRepository *sql.Stmt
	stmtShowBucket     *sql.Stmt
	stmtShowGroup      *sql.Stmt
	stmtShowCluster    *sql.Stmt
	stmtShowNode       *sql.Stmt
	// object tree
	stmtListRepositoryMemberBuckets *sql.Stmt
	stmtListBucketMemberGroups      *sql.Stmt
	stmtListBucketMemberClusters    *sql.Stmt
	stmtListBucketMemberNodes       *sql.Stmt
	stmtListGroupMemberGroups       *sql.Stmt
	stmtListGroupMemberClusters     *sql.Stmt
	stmtListGroupMemberNodes        *sql.Stmt
	stmtListClusterMemberNodes      *sql.Stmt
	appLog                          *logrus.Logger
	reqLog                          *logrus.Logger
	errLog                          *logrus.Logger
}

// newTreeRead return a new TreeRead handler with input buffer of
// length
func newTreeRead(length int) (r *TreeRead) {
	r = &TreeRead{}
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return
}

// Register initializes resources provided by the Soma app
func (r *TreeRead) Register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// Intake exposes the Input channel as part of the handler interface
func (r *TreeRead) Intake() chan msg.Request {
	return r.Input
}

// Run is the event loop for TreeRead
func (r *TreeRead) Run() {
	var err error

	// single-object return statements
	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.TreeShowRepository:      r.stmtShowRepository,
		stmt.TreeShowBucket:          r.stmtShowBucket,
		stmt.TreeShowGroup:           r.stmtShowGroup,
		stmt.TreeShowCluster:         r.stmtShowCluster,
		stmt.TreeShowNode:            r.stmtShowNode,
		stmt.TreeBucketsInRepository: r.stmtListRepositoryMemberBuckets,
		stmt.TreeGroupsInBucket:      r.stmtListBucketMemberGroups,
		stmt.TreeClustersInBucket:    r.stmtListBucketMemberClusters,
		stmt.TreeNodesInBucket:       r.stmtListBucketMemberNodes,
		stmt.TreeGroupsInGroup:       r.stmtListGroupMemberGroups,
		stmt.TreeClustersInGroup:     r.stmtListGroupMemberClusters,
		stmt.TreeNodesInGroup:        r.stmtListGroupMemberNodes,
		stmt.TreeNodesInCluster:      r.stmtListClusterMemberNodes,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`tree_r`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-r.Shutdown:
			break runloop
		case req := <-r.Input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

// process is the request dispatcher
func (r *TreeRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	var err error
	tree := proto.Tree{
		ID:   q.Tree.ID,
		Type: q.Tree.Type,
	}

	switch tree.Type {
	case msg.EntityRepository:
		tree.Repository, err = r.repository(tree.ID, 0)
	case msg.EntityBucket:
		tree.Bucket, err = r.bucket(tree.ID, 0)
	case msg.EntityGroup:
		tree.Group, err = r.group(tree.ID, 0)
	case msg.EntityCluster:
		tree.Cluster, err = r.cluster(tree.ID, 0)
	case msg.EntityNode:
		tree.Node, err = r.node(tree.ID, 0)
	default:
		result.UnknownRequest(q)
		goto skip
	}
	if err == sql.ErrNoRows {
		result.NotFound(
			fmt.Errorf(`Tree starting point not found`),
			q.Section,
		)
		goto skip
	} else if err != nil {
		result.ServerError(err, q.Section)
		goto skip
	}
	result.Tree = tree
	result.OK()

skip:
	q.Reply <- result
}

// repository returns the tree below a specifc repository
func (r *TreeRead) repository(id string, depth int) (*proto.Repository, error) {
	var (
		name, teamID, createdBy string
		isActive, isDeleted     bool
		createdAt               time.Time
	)

	if depth >= 64 {
		return nil, fmt.Errorf(`Maximum recursion depth exceeded`)
	}

	if err := r.stmtShowRepository.QueryRow(
		id,
	).Scan(
		&name,
		&isActive,
		&teamID,
		&isDeleted,
		&createdBy,
		&createdAt,
	); err != nil {
		r.errLog.Printf("Error in tree_r.repository() for %s: %s",
			id, err.Error())
		return nil, err
	}

	repo := proto.Repository{
		ID:        id,
		Name:      name,
		TeamID:    teamID,
		IsDeleted: isDeleted,
		IsActive:  isActive,
		Details: &proto.Details{
			CreatedBy: createdBy,
			CreatedAt: createdAt.UTC().Format(time.RFC3339),
		},
	}

	depth++
	repo.Members = &[]proto.Bucket{}

	buckets, err := r.bucketsInRepository(id)
	if err != nil {
		return nil, err
	}
	for i := range buckets {
		b, err := r.bucket(buckets[i], depth)
		if err != nil {
			return nil, err
		}
		*repo.Members = append(*repo.Members, *b)
	}
	return &repo, nil
}

// bucket returns the tree below a specific bucket
func (r *TreeRead) bucket(id string, depth int) (*proto.Bucket, error) {
	var (
		name, repositoryID, environment string
		teamID, createdBy               string
		isFrozen, isDeleted             bool
		createdAt                       time.Time
		groups, clusters, nodes         []string
		err                             error
		g                               *proto.Group
		c                               *proto.Cluster
		n                               *proto.Node
	)

	if depth >= 64 {
		return nil, fmt.Errorf(`Maximum recursion depth exceeded`)
	}

	if err = r.stmtShowBucket.QueryRow(
		id,
	).Scan(
		&name,
		&isFrozen,
		&isDeleted,
		&repositoryID,
		&environment,
		&teamID,
		&createdBy,
		&createdAt,
	); err != nil {
		r.errLog.Printf("Error in tree_r.bucket() for %s: %s",
			id, err.Error())
		return nil, err
	}

	bucket := proto.Bucket{
		ID:           id,
		Name:         name,
		RepositoryID: repositoryID,
		TeamID:       teamID,
		Environment:  environment,
		IsDeleted:    isDeleted,
		IsFrozen:     isFrozen,
		Details: &proto.Details{
			CreatedBy: createdBy,
			CreatedAt: createdAt.UTC().Format(time.RFC3339),
		},
	}

	depth++
	bucket.MemberGroups = &[]proto.Group{}
	bucket.MemberClusters = &[]proto.Cluster{}
	bucket.MemberNodes = &[]proto.Node{}

	if groups, err = r.groupsInBucket(id); err != nil {
		return nil, err
	}
	for i := range groups {
		if g, err = r.group(groups[i], depth); err != nil {
			return nil, err
		}
		*bucket.MemberGroups = append(*bucket.MemberGroups, *g)
	}

	if clusters, err = r.clustersInBucket(id); err != nil {
		return nil, err
	}
	for i := range clusters {
		if c, err = r.cluster(clusters[i], depth); err != nil {
			return nil, err
		}
		*bucket.MemberClusters = append(*bucket.MemberClusters, *c)
	}

	if nodes, err = r.nodesInBucket(id); err != nil {
		return nil, err
	}
	for i := range nodes {
		if n, err = r.node(nodes[i], depth); err != nil {
			return nil, err
		}
		*bucket.MemberNodes = append(*bucket.MemberNodes, *n)
	}
	return &bucket, nil
}

// group returns the tree below a specific group
func (r *TreeRead) group(id string, depth int) (*proto.Group, error) {
	var (
		err                     error
		bucketID, name, state   string
		teamID, createdBy       string
		createdAt               time.Time
		groups, clusters, nodes []string
		g                       *proto.Group
		c                       *proto.Cluster
		n                       *proto.Node
	)

	if depth >= 64 {
		return nil, fmt.Errorf(`Maximum recursion depth exceeded`)
	}

	if err = r.stmtShowGroup.QueryRow(
		id,
	).Scan(
		&bucketID,
		&name,
		&state,
		&teamID,
		&createdBy,
		&createdAt,
	); err != nil {
		r.errLog.Printf("Error in tree_r.group() for %s: %s",
			id, err.Error())
		return nil, err
	}

	group := proto.Group{
		ID:          id,
		BucketID:    bucketID,
		Name:        name,
		ObjectState: state,
		TeamID:      teamID,
		Details: &proto.Details{
			CreatedBy: createdBy,
			CreatedAt: createdAt.UTC().Format(time.RFC3339),
		},
	}

	depth++
	group.MemberGroups = &[]proto.Group{}
	group.MemberClusters = &[]proto.Cluster{}
	group.MemberNodes = &[]proto.Node{}

	if groups, err = r.groupsInGroup(id); err != nil {
		return nil, err
	}
	for i := range groups {
		if g, err = r.group(groups[i], depth); err != nil {
			return nil, err
		}
		*group.MemberGroups = append(*group.MemberGroups, *g)
	}

	if clusters, err = r.clustersInGroup(id); err != nil {
		return nil, err
	}
	for i := range clusters {
		if c, err = r.cluster(clusters[i], depth); err != nil {
			return nil, err
		}
		*group.MemberClusters = append(*group.MemberClusters, *c)
	}

	if nodes, err = r.nodesInGroup(id); err != nil {
		return nil, err
	}
	for i := range nodes {
		if n, err = r.node(nodes[i], depth); err != nil {
			return nil, err
		}
		*group.MemberNodes = append(*group.MemberNodes, *n)
	}
	return &group, nil
}

// cluster returns the tree below a specific cluster
func (r *TreeRead) cluster(id string, depth int) (*proto.Cluster, error) {
	var (
		err                                      error
		name, bucketID, state, teamID, createdBy string
		createdAt                                time.Time
		nodes                                    []string
		n                                        *proto.Node
	)

	if depth >= 64 {
		return nil, fmt.Errorf(`Maximum recursion depth exceeded`)
	}

	if err = r.stmtShowCluster.QueryRow(
		id,
	).Scan(
		&name,
		&bucketID,
		&state,
		&teamID,
		&createdBy,
		&createdAt,
	); err != nil {
		r.errLog.Printf("Error in tree_r.cluster() for %s: %s",
			id, err.Error())
		return nil, err
	}

	cluster := proto.Cluster{
		ID:          id,
		Name:        name,
		BucketID:    bucketID,
		ObjectState: state,
		TeamID:      teamID,
		Details: &proto.Details{
			CreatedBy: createdBy,
			CreatedAt: createdAt.UTC().Format(time.RFC3339),
		},
	}

	depth++
	cluster.Members = &[]proto.Node{}

	if nodes, err = r.nodesInCluster(id); err != nil {
		return nil, err
	}
	for i := range nodes {
		if n, err = r.node(nodes[i], depth); err != nil {
			return nil, err
		}
		*cluster.Members = append(*cluster.Members, *n)
	}
	return &cluster, nil
}

// node returns the tree below a specific node
func (r *TreeRead) node(id string, depth int) (*proto.Node, error) {
	var (
		assetID                           int
		name, teamID, serverID, state     string
		createdBy, repositoryID, bucketID string
		isOnline, isDeleted               bool
		createdAt                         time.Time
	)

	if depth >= 64 {
		return nil, fmt.Errorf(`Maximum recursion depth exceeded`)
	}

	if err := r.stmtShowNode.QueryRow(
		id,
	).Scan(
		&assetID,
		&name,
		&teamID,
		&serverID,
		&state,
		&isOnline,
		&isDeleted,
		&createdBy,
		&createdAt,
		&repositoryID,
		&bucketID,
	); err != nil {
		r.errLog.Printf("Error in tree_r.node() for %s: %s",
			id, err.Error())
		return nil, err
	}

	node := proto.Node{
		ID:        id,
		AssetID:   uint64(assetID),
		Name:      name,
		TeamID:    teamID,
		ServerID:  serverID,
		State:     state,
		IsOnline:  isOnline,
		IsDeleted: isDeleted,
		Details: &proto.Details{
			CreatedBy: createdBy,
			CreatedAt: createdAt.UTC().Format(time.RFC3339),
		},
		Config: &proto.NodeConfig{
			RepositoryID: repositoryID,
			BucketID:     bucketID,
		},
	}
	return &node, nil
}

// ShutdownNow signals the handler to shut down
func (r *TreeRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
