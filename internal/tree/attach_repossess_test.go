/*-
 * Copyright (c) 2018, Jörg Pernfuß
 * Copyright (c) 2018, 1&1 IONOS SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree // import "github.com/mjolnir42/soma/internal/tree"

import (
	"testing"

	"github.com/satori/go.uuid"
)

func TestRepossessRepository(t *testing.T) {
	deterministicInheritanceOrder = true

	testTree, actionChan, errorChan := testSpawnCheckTree()

	newTeamID := uuid.Must(uuid.NewV4()).String()

	testTree.SetTeamID(
		newTeamID,
	)

	close(actionChan)
	close(errorChan)

	if len(errorChan) > 0 {
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

		// SetTeamID
		[]string{`node`, ActionRepossess},
		[]string{`cluster`, ActionRepossess},
		[]string{`node`, ActionRepossess},
		[]string{`group`, ActionRepossess},
		[]string{`group`, ActionRepossess},
		[]string{`node`, ActionRepossess},
		[]string{`cluster`, ActionRepossess},
		[]string{`node`, ActionRepossess},
		[]string{`bucket`, ActionRepossess},
		[]string{`repository`, ActionRepossess},
	}
	for a := range actionChan {
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
	for len(actions) > elem {
		t.Error(`missing action:`, actions[elem][0], actions[elem][1])
		elem++
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
