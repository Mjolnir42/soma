/*-
 * Copyright (c) 2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/soma/internal/stmt"

const AdminRemove = `
DELETE FROM auth.admin
WHERE  id = $1::uuid;`

func init() {
	m[AdminRemove] = `AdminRemove`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
