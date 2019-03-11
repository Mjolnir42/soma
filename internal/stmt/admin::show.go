/*-
 * Copyright (c) 2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/soma/internal/stmt"

const AdminShow = `
SELECT auth.admin.id,
       auth.admin.uid,
       inventory.user.id,
       inventory.user.uid
FROM   auth.admin
JOIN   inventory.user
  ON   auth.admin.user_uid = inventory.user.uid
WHERE  inventory.user.id = $1::uuid;`

func init() {
	m[AdminShow] = `AdminShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
