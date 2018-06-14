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

func TestErrorChannelNode(t *testing.T) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	nodeID := uuid.Must(uuid.NewV4()).String()
	servID := uuid.Must(uuid.NewV4()).String()

	// create tree
	sTree := New(TreeSpec{
		Id:     rootID,
		Name:   `root_testing`,
		Action: actionC,
	})

	// create repository
	repo := NewRepository(RepositorySpec{
		Id:      repoID,
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
	sTree.SetError(errC)
	if repo.Fault.Error == nil {
		t.Errorf(`Repository.Fault.Error is nil`)
	} else {
		repo.Fault.Error <- &Error{Action: `testmessage_repo`}
	}

	// create bucket
	buck := NewBucket(BucketSpec{
		Id:          buckID,
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
	if buck.Fault.Error == nil {
		t.Errorf(`Bucket.Fault.Error is nil`)
	} else {
		buck.Fault.Error <- &Error{Action: `testmessage_bucket`}
	}

	// create new node
	node := NewNode(NodeSpec{
		Id:       nodeID,
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

	if node.Fault.Error == nil {
		t.Errorf(`Node.Fault.Error is nil`)
	} else {
		node.Fault.Error <- &Error{Action: `testmessage_node`}
	}

	close(actionC)
	close(errC)

	if len(errC) != 3 {
		t.Error(len(errC), `elements in error channel`)
	}

	if len(actionC) != 6 {
		t.Error(len(actionC), `elements in action channel`)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
