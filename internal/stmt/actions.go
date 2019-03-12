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
	ActionStatements = ``

	ActionList = `
SELECT soma.action.id,
       soma.action.name,
       soma.action.section_id
FROM   soma.action
WHERE  section_id = $1::uuid;`

	ActionSearch = `
SELECT soma.action.id,
       soma.action.name,
       soma.action.section_id
FROM   soma.action
WHERE  soma.action.name = $1::varchar
  AND  soma.action.section_id = $2::uuid;`

	ActionShow = `
SELECT soma.action.id,
       soma.action.name,
       soma.action.section_id,
       soma.section.name,
       soma.action.category,
       inventory.user.uid,
       soma.action.created_at
FROM   soma.action
JOIN   inventory.user
  ON   soma.action.created_by = inventory.user.id
JOIN   soma.section
  ON   soma.action.section_id = soma.section.id
WHERE  soma.action.id = $1::uuid;`

	ActionLoad = `
SELECT soma.action.id,
       soma.action.name,
       soma.action.section_id,
       soma.section.name,
       soma.action.category
FROM   soma.action
JOIN   soma.section
  ON   soma.action.section_id = soma.section.id;`

	ActionRemoveFromMap = `
DELETE FROM soma.permission_map
WHERE       action_id = $1::uuid;`

	ActionRemove = `
DELETE FROM soma.action
WHERE       id = $1::uuid;`

	ActionAdd = `
INSERT INTO soma.action (
            id,
            name,
            section_id,
            category,
            created_by)
SELECT      $1::uuid,
            $2::varchar,
            $3::uuid,
            ( SELECT soma.section.category
              FROM   soma.section
              WHERE  soma.section.id = $3::uuid),
            ( SELECT inventory.user.id FROM inventory.user
              LEFT JOIN auth.admin
              ON inventory.user.uid = auth.admin.user_uid
              WHERE (   inventory.user.uid = $4::varchar
                     OR auth.admin.uid     = $4::varchar ))
WHERE       NOT EXISTS (
     SELECT soma.action.id
     FROM   soma.action
     WHERE  soma.action.name = $2::varchar
     AND    soma.action.section_id = $3::uuid);`
)

func init() {
	m[ActionAdd] = `ActionAdd`
	m[ActionList] = `ActionList`
	m[ActionLoad] = `ActionLoad`
	m[ActionRemoveFromMap] = `ActionRemoveFromMap`
	m[ActionRemove] = `ActionRemove`
	m[ActionSearch] = `ActionSearch`
	m[ActionShow] = `ActionShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
