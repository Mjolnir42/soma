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
	SectionStatements = ``

	SectionList = `
SELECT soma.section.id,
       soma.section.name
FROM   soma.section
WHERE  category = $1::varchar;`

	SectionSearch = `
SELECT soma.section.id,
       soma.section.name,
       soma.section.category
FROM   soma.section
WHERE  (soma.section.name = $1::varchar OR $1::varchar IS NULL)
  AND  (soma.section.id = $2::uuid OR $2::uuid IS NULL);`

	SectionShow = `
SELECT soma.section.id,
       soma.section.name,
       soma.section.category,
       inventory.user.uid,
       soma.section.created_at
FROM   soma.section
JOIN   inventory.user
  ON   soma.section.created_by = inventory.user.id
WHERE  soma.section.id = $1::uuid;`

	SectionLoad = `
SELECT soma.section.id,
       soma.section.name,
       soma.section.category
FROM   soma.section;`

	SectionRemoveFromMap = `
DELETE FROM soma.permission_map
WHERE       section_id = $1::uuid
  AND       action_id IS NULL;`

	SectionRemove = `
DELETE FROM soma.section
WHERE       soma.section.id = $1::uuid;`

	SectionListActions = `
SELECT soma.action.id
FROM   soma.action
WHERE  soma.action.section_id = $1::uuid;`

	SectionAdd = `
INSERT INTO soma.section (
            id,
            name,
            category,
            created_by)
SELECT      $1::uuid,
            $2::varchar,
            $3::varchar,
            ( SELECT inventory.user.id
              FROM   inventory.user
              WHERE  inventory.user.uid = $4::varchar)
WHERE       NOT EXISTS (
     SELECT soma.section.id
     FROM   soma.section
     WHERE  soma.section.name = $2::varchar);`
)

func init() {
	m[SectionAdd] = `SectionAdd`
	m[SectionListActions] = `SectionListActions`
	m[SectionList] = `SectionList`
	m[SectionLoad] = `SectionLoad`
	m[SectionRemoveFromMap] = `SectionRemoveFromMap`
	m[SectionRemove] = `SectionRemove`
	m[SectionSearch] = `SectionSearch`
	m[SectionShow] = `SectionShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
