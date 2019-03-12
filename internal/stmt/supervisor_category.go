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
	SupervisorCategoryStatements = ``

	CategoryAdd = `
INSERT INTO soma.category (
            name,
            created_by
)
SELECT $1::varchar,
       ( SELECT inventory.user.id FROM inventory.user
         LEFT JOIN auth.admin
         ON inventory.user.uid = auth.admin.user_uid
         WHERE (   inventory.user.uid = $2::varchar
                OR auth.admin.uid     = $2::varchar ))
WHERE NOT EXISTS (
      SELECT soma.category.name
      FROM   soma.category
      WHERE  soma.category.name = $1::varchar);`

	CategoryRemove = `
DELETE FROM soma.category
WHERE soma.category.name = $1::varchar;`

	CategoryList = `
SELECT soma.category.name
FROM   soma.category;`

	CategoryShow = `
SELECT soma.category.name,
       inventory.user.uid,
       soma.category.created_at
FROM   soma.category
JOIN   inventory.user
ON     soma.category.created_by = inventory.user.id
WHERE  soma.category.name = $1::varchar;`

	CategoryListSections = `
SELECT soma.section.id
FROM   soma.section
WHERE  soma.section.category = $1::varchar;`

	CategoryListPermissions = `
SELECT soma.permission.id
FROM   soma.permission
WHERE  soma.permission.category = $1::varchar;`
)

func init() {
	m[CategoryAdd] = `CategoryAdd`
	m[CategoryListPermissions] = `CategoryListPermissions`
	m[CategoryListSections] = `CategoryListSections`
	m[CategoryList] = `CategoryList`
	m[CategoryRemove] = `CategoryRemove`
	m[CategoryShow] = `CategoryShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
