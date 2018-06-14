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

func TestCheckerAddCheck(t *testing.T) {
	sTree, actionC, errC := testSpawnCheckTree()

	chkConfigID := uuid.Must(uuid.NewV4())
	capID := uuid.Must(uuid.NewV4())

	chk := Check{
		Id:            uuid.Nil,
		SourceID:      uuid.Nil,
		InheritedFrom: uuid.Nil,
		Inheritance:   true,
		ChildrenOnly:  false,
		Interval:      60,
		ConfigID:      chkConfigID,
		CapabilityID:  capID,
		View:          `any`,
		Thresholds: []CheckThreshold{
			{
				Predicate: `>=`,
				Level:     1,
				Value:     100,
			},
			{
				Predicate: `>=`,
				Level:     3,
				Value:     450,
			},
		},
		Constraints: []CheckConstraint{},
	}

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementName: `checkTest`,
	}, true).SetCheck(chk)

	sTree.ComputeCheckInstances()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`group`, `create`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `create`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`group`, `member_new`}, // MoveGroupToGroup
		[]string{`group`, `update`},
		[]string{`group`, `member_new`}, // MoveClusterToGroup
		[]string{`cluster`, `update`},
		[]string{`cluster`, `member_new`}, // MoveNodeToCluster
		[]string{`node`, `update`},
		[]string{`group`, `member_new`}, // MoveNodeToGroup
		[]string{`node`, `update`},
		[]string{`cluster`, `member_new`}, // MoveNodeToCluster
		[]string{`node`, `update`},
		[]string{`node`, `check_new`}, // SetCheck
		[]string{`cluster`, `check_new`},
		[]string{`node`, `check_new`},
		[]string{`group`, `check_new`},
		[]string{`group`, `check_new`},
		[]string{`node`, `check_new`},
		[]string{`cluster`, `check_new`},
		[]string{`node`, `check_new`},
		[]string{`bucket`, `check_new`},
		[]string{`repository`, `check_new`},
		[]string{`node`, `check_instance_create`}, // ComputeInstances
		[]string{`cluster`, `check_instance_create`},
		[]string{`node`, `check_instance_create`},
		[]string{`group`, `check_instance_create`},
		[]string{`group`, `check_instance_create`},
		[]string{`node`, `check_instance_create`},
		[]string{`cluster`, `check_instance_create`},
		[]string{`node`, `check_instance_create`},
	}
	for a := range actionC {
		if elem >= len(actions) {
			t.Error(
				`Received additional action`,
				a.Type, a.Action,
			)
			elem++
			continue
		}

		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action`, elem, `. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}
}

func TestCheckerDeleteCheck(t *testing.T) {
	sTree, actionC, errC := testSpawnCheckTree()

	chkConfigID := uuid.Must(uuid.NewV4())
	capID := uuid.Must(uuid.NewV4())
	chkID := uuid.Must(uuid.NewV4())

	chk := Check{
		Id:            chkID,
		SourceID:      uuid.Nil,
		InheritedFrom: uuid.Nil,
		Inheritance:   true,
		ChildrenOnly:  false,
		Interval:      60,
		ConfigID:      chkConfigID,
		CapabilityID:  capID,
		View:          `any`,
		Thresholds: []CheckThreshold{
			{
				Predicate: `>=`,
				Level:     1,
				Value:     100,
			},
			{
				Predicate: `>=`,
				Level:     3,
				Value:     450,
			},
		},
		Constraints: []CheckConstraint{},
	}

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementName: `checkTest`,
	}, true).SetCheck(chk)

	sTree.ComputeCheckInstances()

	delChk := Check{
		Id:            uuid.Nil,
		InheritedFrom: uuid.Nil,
		SourceID:      chkID,
		ConfigID:      chkConfigID,
	}

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementName: `checkTest`,
	}, true).DeleteCheck(delChk)

	sTree.ComputeCheckInstances()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, `create`},
		[]string{`fault`, `create`},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, `create`},
		[]string{`group`, `create`},
		[]string{`group`, `create`},
		[]string{`cluster`, `create`},
		[]string{`cluster`, `create`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`bucket`, `node_assignment`}, // NewNode
		[]string{`node`, `update`},
		[]string{`group`, `member_new`}, // MoveGroupToGroup
		[]string{`group`, `update`},
		[]string{`group`, `member_new`}, // MoveClusterToGroup
		[]string{`cluster`, `update`},
		[]string{`cluster`, `member_new`}, // MoveNodeToCluster
		[]string{`node`, `update`},
		[]string{`group`, `member_new`}, // MoveNodeToGroup
		[]string{`node`, `update`},
		[]string{`cluster`, `member_new`}, // MoveNodeToCluster
		[]string{`node`, `update`},
		[]string{`node`, `check_new`}, // SetCheck
		[]string{`cluster`, `check_new`},
		[]string{`node`, `check_new`},
		[]string{`group`, `check_new`},
		[]string{`group`, `check_new`},
		[]string{`node`, `check_new`},
		[]string{`cluster`, `check_new`},
		[]string{`node`, `check_new`},
		[]string{`bucket`, `check_new`},
		[]string{`repository`, `check_new`},
		[]string{`node`, `check_instance_create`}, // ComputeInstances
		[]string{`cluster`, `check_instance_create`},
		[]string{`node`, `check_instance_create`},
		[]string{`group`, `check_instance_create`},
		[]string{`group`, `check_instance_create`},
		[]string{`node`, `check_instance_create`},
		[]string{`cluster`, `check_instance_create`},
		[]string{`node`, `check_instance_create`},
		[]string{`node`, `check_removed`}, // DeleteCheck
		[]string{`cluster`, `check_removed`},
		[]string{`node`, `check_removed`},
		[]string{`group`, `check_removed`},
		[]string{`group`, `check_removed`},
		[]string{`node`, `check_removed`},
		[]string{`cluster`, `check_removed`},
		[]string{`node`, `check_removed`},
		[]string{`bucket`, `check_removed`},
		[]string{`repository`, `check_removed`},
		[]string{`node`, `check_instance_delete`}, // ComputeInstances
		[]string{`cluster`, `check_instance_delete`},
		[]string{`node`, `check_instance_delete`},
		[]string{`group`, `check_instance_delete`},
		[]string{`group`, `check_instance_delete`},
		[]string{`node`, `check_instance_delete`},
		[]string{`cluster`, `check_instance_delete`},
		[]string{`node`, `check_instance_delete`},
	}
	for a := range actionC {
		if elem >= len(actions) {
			t.Error(
				`Received additional action`,
				elem, a.Type, a.Action,
			)
			elem++
			continue
		}

		if a.Type != actions[elem][0] || a.Action != actions[elem][1] {
			t.Error(
				`Received incorrect action`, elem, `. Expected`,
				actions[elem][0], actions[elem][1],
				`and received`, a.Type, a.Action,
			)
		}
		elem++
	}
}

func testSpawnCheckTree() (*Tree, chan *Action, chan *Error) {
	actionC := make(chan *Action, 128)
	errC := make(chan *Error, 128)

	rootID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()
	repoID := uuid.Must(uuid.NewV4()).String()
	buckID := uuid.Must(uuid.NewV4()).String()
	grp1Id := uuid.Must(uuid.NewV4()).String()
	grp2Id := uuid.Must(uuid.NewV4()).String()
	clr1Id := uuid.Must(uuid.NewV4()).String()
	clr2Id := uuid.Must(uuid.NewV4()).String()
	nod1Id := uuid.Must(uuid.NewV4()).String()
	srv1Id := uuid.Must(uuid.NewV4()).String()
	nod2Id := uuid.Must(uuid.NewV4()).String()
	srv2Id := uuid.Must(uuid.NewV4()).String()
	nod3Id := uuid.Must(uuid.NewV4()).String()
	srv3Id := uuid.Must(uuid.NewV4()).String()
	nod4Id := uuid.Must(uuid.NewV4()).String()
	srv4Id := uuid.Must(uuid.NewV4()).String()

	sTree := New(TreeSpec{
		Id:     rootID,
		Name:   `root_checkTest`,
		Action: actionC,
	})

	NewRepository(RepositorySpec{
		Id:      repoID,
		Name:    `checkTest`,
		Team:    teamID,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentID:   rootID,
	})
	sTree.SetError(errC)

	NewBucket(BucketSpec{
		Id:          buckID,
		Name:        `checkTest_master`,
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

	NewGroup(GroupSpec{
		Id:   grp1Id,
		Name: `testGroup1`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	NewGroup(GroupSpec{
		Id:   grp2Id,
		Name: `testGroup2`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	NewCluster(ClusterSpec{
		Id:   clr1Id,
		Name: `testcluster1`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	NewCluster(ClusterSpec{
		Id:   clr2Id,
		Name: `testcluster2`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	NewNode(NodeSpec{
		Id:       nod1Id,
		AssetID:  1,
		Name:     `testnode1`,
		Team:     teamID,
		ServerID: srv1Id,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	NewNode(NodeSpec{
		Id:       nod2Id,
		AssetID:  2,
		Name:     `testnode2`,
		Team:     teamID,
		ServerID: srv2Id,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	NewNode(NodeSpec{
		Id:       nod3Id,
		AssetID:  3,
		Name:     `testnode3`,
		Team:     teamID,
		ServerID: srv3Id,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	NewNode(NodeSpec{
		Id:       nod4Id,
		AssetID:  4,
		Name:     `testnode4`,
		Team:     teamID,
		ServerID: srv4Id,
		Online:   true,
		Deleted:  false,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	sTree.Find(FindRequest{
		ElementType: `group`,
		ElementID:   grp2Id,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grp1Id,
	})

	sTree.Find(FindRequest{
		ElementType: `cluster`,
		ElementID:   clr1Id,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grp2Id,
	})

	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementID:   nod1Id,
	}, true).(ClusterAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentID:   clr1Id,
	})

	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementID:   nod2Id,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `group`,
		ParentID:   grp2Id,
	})

	sTree.Find(FindRequest{
		ElementType: `node`,
		ElementID:   nod3Id,
	}, true).(GroupAttacher).ReAttach(AttachRequest{
		Root:       sTree,
		ParentType: `cluster`,
		ParentID:   clr2Id,
	})

	return sTree, actionC, errC
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
