/*-
 * Copyright (c) 2016-2018, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2018, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	TeamStatements = ``

	TeamList = `
SELECT id,
       name
FROM   inventory.team;`

	TeamShow = `
SELECT inventory.team.id,
       inventory.team.name,
       inventory.team.ldap_id,
       inventory.team.is_system,
       inventory.dictionary.id,
       inventory.dictionary.name,
       inventory.user.uid,
       inventory.team.created_at
FROM   inventory.team
JOIN   inventory.user
  ON   inventory.team.created_by = inventory.user.id
JOIN   inventory.dictionary
  ON   inventory.team.dictionary_id = inventory.dictionary.id
WHERE  inventory.team.id = $1::uuid;`

	TeamSync = `
SELECT id,
       name,
       ldap_id,
       is_system
FROM   inventory.team
WHERE  NOT is_system;`

	TeamAdd = `
INSERT INTO inventory.team (
            id,
            name,
            ldap_id,
            is_system,
            dictionary_id,
            created_by)
SELECT      $1::uuid,
            $2::varchar,
            $3::numeric,
            $4::boolean,
            -- hardcoded dictionary system
            '00000000-0000-0000-0000-000000000000'::uuid,
            ( SELECT inventory.user.id FROM inventory.user
              LEFT JOIN auth.admin
              ON inventory.user.uid = auth.admin.user_uid
              WHERE (   inventory.user.uid = $5::varchar
                     OR auth.admin.uid     = $5::varchar ))
WHERE  NOT EXISTS (
   SELECT id
   FROM   inventory.team
   WHERE  id = $1::uuid
      OR  name = $2::varchar
      OR  ldap_id = $3::numeric);`

	TeamUpdate = `
UPDATE inventory.team
SET    name = $1::varchar,
       ldap_id = $2::numeric,
       is_system = $3::boolean
WHERE  inventory.team.id = $4::uuid;`

	TeamRemove = `
DELETE FROM inventory.team
WHERE       inventory.team.id = $1;`

	TeamMembers = `
SELECT inventory.user.id,
       inventory.user.uid
FROM   inventory.team
JOIN   inventory.user
  ON   inventory.team.id = inventory.user.team_id
WHERE  inventory.team.id = $1::uuid
  AND  NOT inventory.user.is_deleted;`

	TeamLoad = `
SELECT id,
       name,
       ldap_id,
       is_system
FROM   inventory.team;`
)

func init() {
	m[TeamAdd] = `TeamAdd`
	m[TeamList] = `TeamList`
	m[TeamLoad] = `TeamLoad`
	m[TeamMembers] = `TeamMembers`
	m[TeamRemove] = `TeamRemove`
	m[TeamShow] = `TeamShow`
	m[TeamSync] = `TeamSync`
	m[TeamUpdate] = `TeamUpdate`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
