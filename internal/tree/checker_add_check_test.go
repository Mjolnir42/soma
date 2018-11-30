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

func TestCheckerAddCheck(t *testing.T) {
	deterministicInheritanceOrder = true

	sTree, actionC, errC := testSpawnCheckTree()

	chkConfigID := uuid.Must(uuid.NewV4())
	capID := uuid.Must(uuid.NewV4())

	chk := Check{
		ID:            uuid.Nil,
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
		[]string{`repository`, ActionCreate},
		[]string{`fault`, ActionCreate},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, ActionCreate},
		[]string{`group`, ActionCreate},
		[]string{`group`, ActionCreate},
		[]string{`cluster`, ActionCreate},
		[]string{`cluster`, ActionCreate},
		[]string{`bucket`, ActionNodeAssignment}, // NewNode
		[]string{`node`, ActionUpdate},
		[]string{`bucket`, ActionNodeAssignment}, // NewNode
		[]string{`node`, ActionUpdate},
		[]string{`bucket`, ActionNodeAssignment}, // NewNode
		[]string{`node`, ActionUpdate},
		[]string{`bucket`, ActionNodeAssignment}, // NewNode
		[]string{`node`, ActionUpdate},
		[]string{`group`, ActionMemberNew}, // MoveGroupToGroup
		[]string{`group`, ActionUpdate},
		[]string{`group`, ActionMemberNew}, // MoveClusterToGroup
		[]string{`cluster`, ActionUpdate},
		[]string{`cluster`, ActionMemberNew}, // MoveNodeToCluster
		[]string{`node`, ActionUpdate},
		[]string{`group`, ActionMemberNew}, // MoveNodeToGroup
		[]string{`node`, ActionUpdate},
		[]string{`cluster`, ActionMemberNew}, // MoveNodeToCluster
		[]string{`node`, ActionUpdate},
		[]string{`node`, ActionCheckNew}, // SetCheck
		[]string{`cluster`, ActionCheckNew},
		[]string{`node`, ActionCheckNew},
		[]string{`group`, ActionCheckNew},
		[]string{`group`, ActionCheckNew},
		[]string{`node`, ActionCheckNew},
		[]string{`cluster`, ActionCheckNew},
		[]string{`node`, ActionCheckNew},
		[]string{`bucket`, ActionCheckNew},
		[]string{`repository`, ActionCheckNew},
		[]string{`node`, ActionCheckInstanceCreate}, // ComputeInstances
		[]string{`cluster`, ActionCheckInstanceCreate},
		[]string{`node`, ActionCheckInstanceCreate},
		[]string{`group`, ActionCheckInstanceCreate},
		[]string{`group`, ActionCheckInstanceCreate},
		[]string{`node`, ActionCheckInstanceCreate},
		[]string{`cluster`, ActionCheckInstanceCreate},
		[]string{`node`, ActionCheckInstanceCreate},
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
	deterministicInheritanceOrder = false
}

func TestCheckerDeleteCheck(t *testing.T) {
	deterministicInheritanceOrder = true

	sTree, actionC, errC := testSpawnCheckTree()

	chkConfigID := uuid.Must(uuid.NewV4())
	capID := uuid.Must(uuid.NewV4())
	chkID := uuid.Must(uuid.NewV4())

	chk := Check{
		ID:            chkID,
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
		ID:            uuid.Nil,
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
		[]string{`repository`, ActionCreate},
		[]string{`fault`, ActionCreate},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, ActionCreate},
		[]string{`group`, ActionCreate},
		[]string{`group`, ActionCreate},
		[]string{`cluster`, ActionCreate},
		[]string{`cluster`, ActionCreate},
		[]string{`bucket`, ActionNodeAssignment}, // NewNode
		[]string{`node`, ActionUpdate},
		[]string{`bucket`, ActionNodeAssignment}, // NewNode
		[]string{`node`, ActionUpdate},
		[]string{`bucket`, ActionNodeAssignment}, // NewNode
		[]string{`node`, ActionUpdate},
		[]string{`bucket`, ActionNodeAssignment}, // NewNode
		[]string{`node`, ActionUpdate},
		[]string{`group`, ActionMemberNew}, // MoveGroupToGroup
		[]string{`group`, ActionUpdate},
		[]string{`group`, ActionMemberNew}, // MoveClusterToGroup
		[]string{`cluster`, ActionUpdate},
		[]string{`cluster`, ActionMemberNew}, // MoveNodeToCluster
		[]string{`node`, ActionUpdate},
		[]string{`group`, ActionMemberNew}, // MoveNodeToGroup
		[]string{`node`, ActionUpdate},
		[]string{`cluster`, ActionMemberNew}, // MoveNodeToCluster
		[]string{`node`, ActionUpdate},
		[]string{`node`, ActionCheckNew}, // SetCheck
		[]string{`cluster`, ActionCheckNew},
		[]string{`node`, ActionCheckNew},
		[]string{`group`, ActionCheckNew},
		[]string{`group`, ActionCheckNew},
		[]string{`node`, ActionCheckNew},
		[]string{`cluster`, ActionCheckNew},
		[]string{`node`, ActionCheckNew},
		[]string{`bucket`, ActionCheckNew},
		[]string{`repository`, ActionCheckNew},
		[]string{`node`, ActionCheckInstanceCreate}, // ComputeInstances
		[]string{`cluster`, ActionCheckInstanceCreate},
		[]string{`node`, ActionCheckInstanceCreate},
		[]string{`group`, ActionCheckInstanceCreate},
		[]string{`group`, ActionCheckInstanceCreate},
		[]string{`node`, ActionCheckInstanceCreate},
		[]string{`cluster`, ActionCheckInstanceCreate},
		[]string{`node`, ActionCheckInstanceCreate},
		[]string{`node`, ActionCheckRemoved}, // DeleteCheck
		[]string{`cluster`, ActionCheckRemoved},
		[]string{`node`, ActionCheckRemoved},
		[]string{`group`, ActionCheckRemoved},
		[]string{`group`, ActionCheckRemoved},
		[]string{`node`, ActionCheckRemoved},
		[]string{`cluster`, ActionCheckRemoved},
		[]string{`node`, ActionCheckRemoved},
		[]string{`bucket`, ActionCheckRemoved},
		[]string{`repository`, ActionCheckRemoved},
		[]string{`node`, ActionCheckInstanceDelete}, // ComputeInstances
		[]string{`cluster`, ActionCheckInstanceDelete},
		[]string{`node`, ActionCheckInstanceDelete},
		[]string{`group`, ActionCheckInstanceDelete},
		[]string{`group`, ActionCheckInstanceDelete},
		[]string{`node`, ActionCheckInstanceDelete},
		[]string{`cluster`, ActionCheckInstanceDelete},
		[]string{`node`, ActionCheckInstanceDelete},
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
	deterministicInheritanceOrder = false
}

func TestCheckerDestroyRepoWithChecks(t *testing.T) {
	deterministicInheritanceOrder = true

	sTree, actionC, errC := testSpawnCheckTree()

	chkConfigID := uuid.Must(uuid.NewV4())
	capID := uuid.Must(uuid.NewV4())
	chkID := uuid.Must(uuid.NewV4())

	chk := Check{
		ID:            chkID,
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

	sTree.Find(FindRequest{
		ElementType: `repository`,
		ElementName: `checkTest`,
	}, true).Destroy()

	close(actionC)
	close(errC)

	if len(errC) > 0 {
		t.Error(`Error channel not empty`)
	}

	elem := 0
	actions := [][]string{
		[]string{`repository`, ActionCreate},
		[]string{`fault`, ActionCreate},
		[]string{`errorchannel`, `attached`},
		[]string{`bucket`, ActionCreate},
		[]string{`group`, ActionCreate},
		[]string{`group`, ActionCreate},
		[]string{`cluster`, ActionCreate},
		[]string{`cluster`, ActionCreate},
		[]string{`bucket`, ActionNodeAssignment}, // NewNode
		[]string{`node`, ActionUpdate},
		[]string{`bucket`, ActionNodeAssignment}, // NewNode
		[]string{`node`, ActionUpdate},
		[]string{`bucket`, ActionNodeAssignment}, // NewNode
		[]string{`node`, ActionUpdate},
		[]string{`bucket`, ActionNodeAssignment}, // NewNode
		[]string{`node`, ActionUpdate},
		[]string{`group`, ActionMemberNew}, // MoveGroupToGroup
		[]string{`group`, ActionUpdate},
		[]string{`group`, ActionMemberNew}, // MoveClusterToGroup
		[]string{`cluster`, ActionUpdate},
		[]string{`cluster`, ActionMemberNew}, // MoveNodeToCluster
		[]string{`node`, ActionUpdate},
		[]string{`group`, ActionMemberNew}, // MoveNodeToGroup
		[]string{`node`, ActionUpdate},
		[]string{`cluster`, ActionMemberNew}, // MoveNodeToCluster
		[]string{`node`, ActionUpdate},
		[]string{`node`, ActionCheckNew}, // SetCheck
		[]string{`cluster`, ActionCheckNew},
		[]string{`node`, ActionCheckNew},
		[]string{`group`, ActionCheckNew},
		[]string{`group`, ActionCheckNew},
		[]string{`node`, ActionCheckNew},
		[]string{`cluster`, ActionCheckNew},
		[]string{`node`, ActionCheckNew},
		[]string{`bucket`, ActionCheckNew},
		[]string{`repository`, ActionCheckNew},
		[]string{`node`, ActionCheckInstanceCreate}, // ComputeInstances
		[]string{`cluster`, ActionCheckInstanceCreate},
		[]string{`node`, ActionCheckInstanceCreate},
		[]string{`group`, ActionCheckInstanceCreate},
		[]string{`group`, ActionCheckInstanceCreate},
		[]string{`node`, ActionCheckInstanceCreate},
		[]string{`cluster`, ActionCheckInstanceCreate},
		[]string{`node`, ActionCheckInstanceCreate},
		[]string{`repository`, ActionDelete},
		[]string{`node`, ActionCheckRemoved},
		[]string{`cluster`, ActionCheckRemoved},
		[]string{`node`, ActionCheckRemoved},
		[]string{`group`, ActionCheckRemoved},
		[]string{`group`, ActionCheckRemoved},
		[]string{`node`, ActionCheckRemoved},
		[]string{`cluster`, ActionCheckRemoved},
		[]string{`node`, ActionCheckRemoved},
		[]string{`bucket`, ActionCheckRemoved},
		[]string{`repository`, ActionCheckRemoved},
		[]string{`cluster`, ActionDelete},
		[]string{`node`, ActionDelete},
		[]string{`group`, ActionDelete},
		[]string{`cluster`, ActionCheckInstanceDelete},
		[]string{`group`, ActionCheckInstanceDelete},
		[]string{`node`, ActionCheckInstanceDelete},
		[]string{`node`, ActionDelete},
		[]string{`group`, ActionDelete},
		[]string{`group`, ActionCheckInstanceDelete},
		[]string{`node`, ActionCheckInstanceDelete},
		[]string{`cluster`, ActionMemberRemoved},
		[]string{`cluster`, ActionDelete},
		[]string{`node`, ActionDelete},
		[]string{`cluster`, ActionCheckInstanceDelete},
		[]string{`node`, ActionCheckInstanceDelete},
		[]string{`group`, ActionMemberRemoved},
		[]string{`node`, ActionDelete},
		[]string{`node`, ActionCheckInstanceDelete},
		[]string{`cluster`, ActionMemberRemoved},
		[]string{`group`, ActionMemberRemoved},
		[]string{`group`, ActionMemberRemoved},
		[]string{`bucket`, ActionDelete},
		[]string{`fault`, `remove_actionchannel`},
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
	deterministicInheritanceOrder = false
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

	discardLog := logrus.New()
	discardLog.Out = ioutil.Discard

	sTree := New(Spec{
		ID:     rootID,
		Name:   `root_checkTest`,
		Action: actionC,
		Log:    discardLog,
	})
	sTree.RegisterErrChan(errC)

	NewRepository(RepositorySpec{
		ID:      repoID,
		Name:    `checkTest`,
		Team:    teamID,
		Deleted: false,
		Active:  true,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `root`,
		ParentID:   rootID,
	})
	sTree.SetError()

	NewBucket(BucketSpec{
		ID:          buckID,
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
		ID:   grp1Id,
		Name: `testGroup1`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	NewGroup(GroupSpec{
		ID:   grp2Id,
		Name: `testGroup2`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	NewCluster(ClusterSpec{
		ID:   clr1Id,
		Name: `testcluster1`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	NewCluster(ClusterSpec{
		ID:   clr2Id,
		Name: `testcluster2`,
		Team: teamID,
	}).Attach(AttachRequest{
		Root:       sTree,
		ParentType: `bucket`,
		ParentID:   buckID,
	})

	NewNode(NodeSpec{
		ID:       nod1Id,
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
		ID:       nod2Id,
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
		ID:       nod3Id,
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
		ID:       nod4Id,
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
