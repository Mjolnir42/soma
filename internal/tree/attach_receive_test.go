/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"io/ioutil"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
)

func TestAttachRepository(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

	// create repository
	NewRepository(RepositorySpec{
		ID:      repoID,
		Name:    `test`,
		Team:    uuid.Must(uuid.NewV4()).String(),
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentID:   rootID,
	})
	sTree.SetError()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 3 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestAttachBucket(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

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
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 4 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestAttachGroup(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 5 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestAttachGroupToGroup(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	grpID2 := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// create group
	NewGroup(GroupSpec{
		ID:   grpID2,
		Name: `testgroup2`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grpID,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 7 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestAttachCluster(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

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
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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
	NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 5 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestAttachNode(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

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
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 6 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestAttachNodeToGroup(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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
		ParentType: `group`,
		ParentID:   grpID,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 7 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

func TestMoveNodeToGroup(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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

	// create node
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
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// create group
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// move node to group
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementID:   nodeID,
	}, true).(*Node).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grpID,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 9 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestMoveClusterToGroup(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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
	NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// create group
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// move cluster to group
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementID:   clrID,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grpID,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 8 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestMoveGroupToGroup(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	grp2Id := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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
	NewGroup(GroupSpec{
		ID:   grp2Id,
		Name: `testgroup2`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// create group
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// move group to group
	sTree.Find(FindRequest{
		ElementType: `group`,
		ElementID:   grp2Id,
	}, true).(*Group).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grpID,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 8 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestMoveNodeToCluster(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// create cluster
	NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// move node to cluster
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementID:   nodeID,
	}, true).(*Node).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentID:   clrID,
	})

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 9 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestDetachGroupToBucket(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	grp2Id := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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
	NewGroup(GroupSpec{
		ID:   grp2Id,
		Name: `testgroup2`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// create group
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// move group to group
	sTree.Find(FindRequest{
		ElementType: `group`,
		ElementID:   grp2Id,
	}, true).(*Group).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grpID,
	})

	// detach group
	sTree.Find(FindRequest{
		ElementType: `group`,
		ElementID:   grp2Id,
	}, true).(*Group).Detach()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 10 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestDetachClusterToBucket(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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
	NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// create group
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// move cluster to group
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementID:   clrID,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grpID,
	})

	// detach cluster
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementID:   clrID,
	}, true).(Attacher).Detach()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 10 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestDetachNodeToBucket(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
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
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// create cluster
	NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// move cluster to group
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementID:   clrID,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grpID,
	})

	// move node to cluster
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementID:   nodeID,
	}, true).(ClusterAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentID:   clrID,
	})

	// detach node
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementID:   nodeID,
	}, true).(Attacher).Detach()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 15 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestDestroyRepository(t *testing.T) {
	deterministicInheritanceOrder = true
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()
	sTree.SwitchLogger(newDiscardLogger())

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
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
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
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// create cluster
	NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// move cluster to group
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementID:   clrID,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grpID,
	})

	// move node to cluster
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementID:   nodeID,
	}, true).(ClusterAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentID:   clrID,
	})

	// destroy bucket
	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementID:   repoID,
	}, true).(Attacher).Destroy()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 20 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}

	if sTree.Child != nil {
		t.Error(`Destroy failed`)
	}
	deterministicInheritanceOrder = false
}

func TestDestroyBucket(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()
	sTree.SwitchLogger(newDiscardLogger())

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
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
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
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// create cluster
	NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// move cluster to group
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementID:   clrID,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grpID,
	})

	// move node to cluster
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementID:   nodeID,
	}, true).(ClusterAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentID:   clrID,
	})

	// destroy bucket
	sTree.Find(FindRequest{
		ElementType: `bucket`,
		ElementID:   buckID,
	}, true).(Attacher).Destroy()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 18 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}
}

func TestRollbackDetachNodeToBucket(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	grpID := uuid.Must(uuid.NewV4()).String()
	clrID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})
	sTree.RegisterErrChan(errC)

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
	sTree.SetError()

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
	NewGroup(GroupSpec{
		ID:   grpID,
		Name: `testgroup`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
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
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// create cluster
	NewCluster(ClusterSpec{
		ID:   clrID,
		Name: `testcluster`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	// move cluster to group
	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementID:   clrID,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grpID,
	})

	// move node to cluster
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementID:   nodeID,
	}, true).(ClusterAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentID:   clrID,
	})

	sTree.Begin()

	if sTree.Snap.ID.String() != repoID {
		t.Error(`Clone failure`)
	}
	if sTree.Snap.Children[buckID].(*Bucket).
		Children[grpID].(*Group).
		Children[clrID].(*Cluster).
		Children[nodeID].(*Node).Name != `testnode` {
		t.Error(`Deep clone failure`)
	}

	// detach node
	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementID:   nodeID,
	}, true).(Attacher).Detach()

	sTree.Rollback()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	if len(actionC) != 16 {
		t.Error(len(actionC), `elements in action channel`)
		for a := range actionC {
			t.Error(`Action:`, a.Type, a.Action)
		}
	}

	if sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementID:   clrID,
	}, true).(*Cluster).Children[nodeID] != sTree.Find(FindRequest{
		ElementType: `node`,
		ElementID:   nodeID,
	}, true) {
		t.Error(`Bad things`)
	}
}

func newDiscardLogger() *logrus.Logger {
	discardLog := logrus.New()
	discardLog.Out = ioutil.Discard
	return discardLog
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
