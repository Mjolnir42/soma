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
	AttributeStatements = ``

	AttributeList = `
SELECT attribute,
       cardinality
FROM   soma.attribute;`

	AttributeShow = `
SELECT attribute,
       cardinality
FROM   soma.attribute
WHERE  attribute = $1::varchar;`

	AttributeAdd = `
INSERT INTO soma.attribute (
            attribute,
            cardinality)
SELECT $1::varchar,
       $2::varchar
WHERE  NOT EXISTS (
    SELECT attribute
    FROM   soma.attribute
    WHERE  attribute = $1::varchar);`

	AttributeRemove = `
DELETE FROM soma.attribute
WHERE       attribute = $1::varchar;`
)

func init() {
	m[AttributeAdd] = `AttributeAdd`
	m[AttributeList] = `AttributeList`
	m[AttributeRemove] = `AttributeRemove`
	m[AttributeShow] = `AttributeShow`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
