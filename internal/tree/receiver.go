/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

type Receiver interface {
	Receive(r ReceiveRequest)
}

type Unlinker interface {
	Unlink(u UnlinkRequest)
}

// implemented by: root
type RepositoryReceiver interface {
	Receiver
	RepositoryUnlinker

	receiveRepository(r ReceiveRequest)
}

type RepositoryUnlinker interface {
	Unlinker

	unlinkRepository(u UnlinkRequest)
}

// implemented by: repositories
type BucketReceiver interface {
	Receiver
	BucketUnlinker

	receiveBucket(r ReceiveRequest)
	resyncProperty(srcID, pType, childID string)
}

type BucketUnlinker interface {
	Unlinker

	unlinkBucket(u UnlinkRequest)
}

type FaultReceiver interface {
	Receiver
	FaultUnlinker

	receiveFault(r ReceiveRequest)
}

type FaultUnlinker interface {
	Unlinker

	unlinkFault(u UnlinkRequest)
}

// implemented by: buckets, groups
type GroupReceiver interface {
	Receiver
	GroupUnlinker

	receiveGroup(r ReceiveRequest)
	resyncProperty(srcID, pType, childID string)
}

type GroupUnlinker interface {
	Unlinker

	unlinkGroup(u UnlinkRequest)
}

// implemented by: buckets, groups
type ClusterReceiver interface {
	Receiver
	ClusterUnlinker

	receiveCluster(r ReceiveRequest)
	resyncProperty(srcID, pType, childID string)
}

type ClusterUnlinker interface {
	Unlinker

	unlinkCluster(u UnlinkRequest)
}

// implemented by: buckets, groups, clusters
type NodeReceiver interface {
	Receiver
	NodeUnlinker

	receiveNode(r ReceiveRequest)
	resyncProperty(srcID, pType, childID string)
}

type NodeUnlinker interface {
	Unlinker

	unlinkNode(u UnlinkRequest)
}

//
type ReceiveRequest struct {
	ParentType string
	ParentID   string
	ParentName string
	ChildType  string
	Repository *Repository
	Bucket     *Bucket
	Group      *Group
	Cluster    *Cluster
	Node       *Node
	Fault      *Fault
}

type UnlinkRequest struct {
	ParentType string
	ParentID   string
	ParentName string
	ChildType  string
	ChildName  string
	ChildID    string
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
