/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/lib/proto"
)

// deprovisionForUpdate checks if a new version for this ID is waiting
// and returns the task string that should be used
func (w *DeploymentWrite) deprovisionForUpdate(q *msg.Request) (string, error) {
	var (
		err       error
		hasUpdate bool
		task      string
	)

	// returns true if there is a updated version blocked, ie.
	// after this deprovisioning a new version will be rolled out
	// -- statement always returns true or false, never null
	if err = w.stmtDeprovisionForUpdate.QueryRow(
		q.Deployment.ID,
	).Scan(
		&hasUpdate,
	); err != nil {
		return ``, err
	}

	switch hasUpdate {
	case false:
		task = proto.TaskDelete
	default:
		task = proto.TaskDeprovision
	}

	return task, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
