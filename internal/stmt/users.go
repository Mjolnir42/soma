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
	UserStatements = ``

	UserList = `
SELECT id,
       uid
FROM   inventory.user
WHERE  NOT is_deleted;`

	UserSearch = `
SELECT id,
       uid
FROM   inventory.user
WHERE  uid = $1::varchar;`

	UserShow = `
SELECT inventory.user.id,
       inventory.user.uid,
       inventory.user.first_name,
       inventory.user.last_name,
       inventory.user.employee_number,
       inventory.user.mail_address,
       inventory.user.is_active,
       inventory.user.is_system,
       inventory.user.is_deleted,
       inventory.user.team_id,
       inventory.user.dictionary_id,
       inventory.dictionary.name,
       creator.uid,
       inventory.user.created_at
FROM   inventory.user
JOIN   inventory.dictionary
  ON   inventory.user.dictionary_id = inventory.dictionary.id
JOIN   inventory.user creator
  ON   inventory.user.created_by = creator.id
WHERE  inventory.user.id = $1::uuid;`

	UserSync = `
SELECT id,
       uid,
       first_name,
       last_name,
       employee_number,
       mail_address,
       is_deleted,
       team_id
FROM   inventory.user
WHERE  NOT is_system;`

	UserAdd = `
INSERT INTO inventory.user (
            id,
            uid,
            first_name,
            last_name,
            employee_number,
            mail_address,
            is_active,
            is_system,
            is_deleted,
            team_id,
            dictionary_id,
            created_by)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar,
       $4::varchar,
       $5::numeric,
       $6::text,
       $7::boolean,
       $8::boolean,
       $9::boolean,
       $10::uuid,
       -- hardcoded dictionary system
       '00000000-0000-0000-0000-000000000000'::uuid,
       ( SELECT inventory.user.id FROM inventory.user
         LEFT JOIN auth.admin
         ON inventory.user.uid = auth.admin.user_uid
         WHERE (   inventory.user.uid = $11::varchar
                OR auth.admin.uid     = $11::varchar ))
WHERE  NOT EXISTS (
  SELECT id
  FROM   inventory.user
  WHERE  id = $1::uuid
     OR  uid = $2::varchar
     OR  employee_number = $5::numeric);`

	UserUpdate = `
UPDATE inventory.user
SET    uid = $1::varchar,
       first_name = $2::varchar,
       last_name = $3::varchar,
       employee_number = $4::numeric,
       mail_address = $5::text,
       is_deleted = $6::boolean,
       team_id = $7::uuid
WHERE  id = $8::uuid
  AND  ($6::boolean OR(is_deleted = $6::boolean));`

	UserRemove = `
UPDATE inventory.user
SET    is_deleted = 'yes',
       is_active = 'no'
WHERE  id = $1::uuid;`

	UserPurge = `
DELETE FROM inventory.user
WHERE  id = $1::uuid
AND    is_deleted;`

	UserLoad = `
SELECT id,
       uid,
       first_name,
       last_name,
       employee_number,
       mail_address,
       is_active,
       is_system,
       is_deleted,
       team_id
FROM   inventory.user;`
)

func init() {
	m[UserAdd] = `UserAdd`
	m[UserList] = `UserList`
	m[UserLoad] = `UserLoad`
	m[UserPurge] = `UserPurge`
	m[UserRemove] = `UserRemove`
	m[UserSearch] = `UserSearch`
	m[UserShow] = `UserShow`
	m[UserSync] = `UserSync`
	m[UserUpdate] = `UserUpdate`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
