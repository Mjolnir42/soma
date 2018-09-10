/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"testing"

	"github.com/satori/go.uuid"
)

// Invalid Attach
func TestInvalidRepositoryAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal attach on repository did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		ID:      repoID,
		Name:    `test`,
		Team:    teamID,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})
}

func TestInvalidBucketAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal attach on bucket did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create bucket
	NewBucket(BucketSpec{
		ID:          buckID,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamID,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentID:   clrID,
	})
}

func TestInvalidGroupAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal attach on group did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create group
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentID:   clrID,
	})
}

func TestInvalidClusterAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal attach on cluster did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create cluster
	NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentID:   repoID,
	})

}

func TestInvalidNodeAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal attach on node did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		ID:      repoID,
		Name:    `test`,
		Team:    teamID,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentID:   rootID,
	})

	// create new node
	NewNode(NodeSpec{
		ID:       nodeID,
		AssetID:  1,
		Name:     `testnode`,
		Team:     teamID,
		ServerID: servID,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentID:   repoID,
	})

}

// Double Attach
func TestDoubleRepositoryAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Double attach on repository did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	repo := NewRepository(RepositorySpec{
		ID:      repoID,
		Name:    `test`,
		Team:    teamID,
		Deleted: false,
		Active:  true,
	})

	repo.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentID:   rootID,
	})
	repo.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentID:   rootID,
	})
}

func TestDoubleBucketAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Double attach on bucket did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		ID:      repoID,
		Name:    `test`,
		Team:    teamID,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentID:   rootID,
	})

	// create bucket
	buck := NewBucket(BucketSpec{
		ID:          buckID,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamID,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoID,
	})
	buck.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentID:   repoID,
	})
	buck.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentID:   repoID,
	})
}

func TestDoubleGroupAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Double attach on group did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		ID:      repoID,
		Name:    `test`,
		Team:    teamID,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentID:   rootID,
	})

	// create bucket
	NewBucket(BucketSpec{
		ID:          buckID,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamID,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentID:   repoID,
	})

	// create group
	grp := NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	})
	grp.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})
	grp.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})
}

func TestDoubleClusterAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Double attach on cluster did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		ID:      repoID,
		Name:    `test`,
		Team:    teamID,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentID:   rootID,
	})

	// create bucket
	NewBucket(BucketSpec{
		ID:          buckID,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamID,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentID:   repoID,
	})

	// create cluster
	clr := NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	})
	clr.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})
	clr.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})
}

func TestDoubleNodeAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Double attach on node did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	NewRepository(RepositorySpec{
		ID:      repoID,
		Name:    `test`,
		Team:    teamID,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentID:   rootID,
	})

	// create bucket
	NewBucket(BucketSpec{
		ID:          buckID,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamID,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `repository`,
		ParentID:   repoID,
	})

	// create new node
	node := NewNode(NodeSpec{
		ID:       nodeID,
		AssetID:  1,
		Name:     `testnode`,
		Team:     teamID,
		ServerID: servID,
		Online:   true,
		Deleted:  false,
	})
	node.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})
	node.Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})
}

// Invalid Destroy
func TestInvalidRepositoryDestroy(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal destroy on repository did not panic`)
		}
	}()

	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()

	// create repository
	NewRepository(RepositorySpec{
		ID:      repoID,
		Name:    `test`,
		Team:    teamID,
		Deleted: false,
		Active:  true,
	}).Destroy()
}

func TestInvalidBucketDestroy(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal destroy on bucket did not panic`)
		}
	}()

	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()

	// create bucket
	NewBucket(BucketSpec{
		ID:          buckID,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamID,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoID,
	}).Destroy()
}

func TestInvalidGroupDestroy(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal destroy on group did not panic`)
		}
	}()

	teamID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()

	// create group
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Destroy()
}

func TestInvalidClusterDestroy(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal destroy on cluster did not panic`)
		}
	}()

	teamID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()

	// create cluster
	NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	}).Destroy()
}

func TestInvalidNodeDestroy(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal destroy on node did not panic`)
		}
	}()

	teamID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()

	// create new node
	NewNode(NodeSpec{
		ID:       nodeID,
		AssetID:  1,
		Name:     `testnode`,
		Team:     teamID,
		ServerID: servID,
		Online:   true,
		Deleted:  false,
	}).Destroy()
}

// Invalid Detach
func TestInvalidRepositoryDetach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal detach on repository did not panic`)
		}
	}()

	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()

	// create repository
	NewRepository(RepositorySpec{
		ID:      repoID,
		Name:    `test`,
		Team:    teamID,
		Deleted: false,
		Active:  true,
	}).Detach()
}

func TestInvalidBucketDetach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal detach on bucket did not panic`)
		}
	}()

	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()

	// create bucket
	NewBucket(BucketSpec{
		ID:          buckID,
		Name:        `test_master`,
		Environment: `testing`,
		Team:        teamID,
		Deleted:     false,
		Frozen:      false,
		Repository:  repoID,
	}).Detach()
}

func TestInvalidGroupDetach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal detach on group did not panic`)
		}
	}()

	teamID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()

	// create group
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Detach()
}

func TestInvalidClusterDetach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal detach on cluster did not panic`)
		}
	}()

	teamID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()

	// create cluster
	NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	}).Detach()
}

func TestInvalidNodeDetach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal detach on node did not panic`)
		}
	}()

	teamID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()

	// create new node
	NewNode(NodeSpec{
		ID:       nodeID,
		AssetID:  1,
		Name:     `testnode`,
		Team:     teamID,
		ServerID: servID,
		Online:   true,
		Deleted:  false,
	}).Detach()
}

// Invalid ReAttach
func TestInvalidGroupReAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal reattach on group did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create group
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})
}

func TestInvalidClusterReAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal reattach on cluster did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create cluster
	clr := NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	})
	clr.ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})
}

func TestInvalidNodeReAttach(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(`Illegal reattach on node did not panic`)
		}
	}()

	actionC := make(chan *Action, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create new node
	node := NewNode(NodeSpec{
		ID:       nodeID,
		AssetID:  1,
		Name:     `testnode`,
		Team:     teamID,
		ServerID: servID,
		Online:   true,
		Deleted:  false,
	})
	node.ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
