/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	SupervisorInventoryStatements = ``

	LoadUserTeamMapping = `
SELECT iu.id,
       iu.uid,
       it.id,
       it.name
FROM   inventory.user iu
JOIN   inventory.team it
ON     iu.team_id = it.id;`
)

func init() {
	m[LoadUserTeamMapping] = `LoadUserTeamMapping`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
