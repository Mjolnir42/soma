/*-
 * Copyright (c) 2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/soma/internal/stmt"

const AdminAdd = `
INSERT INTO auth.admin (
            id,
            uid,
            user_uid,
            is_active)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar,
       'no'::boolean;`

func init() {
	m[AdminAdd] = `AdminAdd`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
