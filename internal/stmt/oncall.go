/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2018, 1&1 IONOS SE
 * Copyright (c) 2016-2018, Jörg Pernfuß <code.jpe@gmail.com>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	OncallStatements = ``

	OncallList = `
SELECT inventory.oncall_team.id,
       inventory.oncall_team.name
FROM   inventory.oncall_team;`

	OncallShow = `
SELECT inventory.oncall_team.id,
       inventory.oncall_team.name,
       inventory.oncall_team.phone_number,
       inventory.dictionary.id,
       inventory.dictionary.name,
       inventory.user.uid,
       inventory.oncall_team.created_at
FROM   inventory.oncall_team
JOIN   inventory.user
  ON   inventory.oncall_team.created_by = inventory.user.id
JOIN   inventory.dictionary
  ON   inventory.oncall_team.dictionary_id = inventory.dictionary.id
WHERE  inventory.oncall_team.id = $1::uuid;`

	OncallSearch = `
SELECT inventory.oncall_team.id,
       inventory.oncall_team.name
FROM   inventory.oncall_team
WHERE  inventory.oncall_team.name = $1::varchar;`

	OncallAdd = `
INSERT INTO inventory.oncall_team (
            id,
            name,
            phone_number,
            dictionary_id,
            created_by)
SELECT $1::uuid,
       $2::varchar,
       $3::numeric,
       '00000000-0000-0000-0000-000000000000'::uuid,
       ( SELECT inventory.user.id FROM inventory.user
         LEFT JOIN auth.admin
         ON inventory.user.uid = auth.admin.user_uid
         WHERE (   inventory.user.uid = $4::varchar
                OR auth.admin.uid     = $4::varchar ))
WHERE  NOT EXISTS (
   SELECT inventory.oncall_team.id
   FROM   inventory.oncall_team
   WHERE  inventory.oncall_team.id = $1::uuid
      OR  inventory.oncall_team.name = $2::varchar
      OR  inventory.oncall_team.phone_number = $3::numeric);`

	OncallUpdate = `
UPDATE inventory.oncall_team
SET    name = CASE WHEN $1::varchar IS NOT NULL
              THEN      $1::varchar
              ELSE      inventory.oncall_team.name
              END,
       phone_number = CASE WHEN $2::numeric IS NOT NULL
                      THEN      $2::numeric
                      ELSE      inventory.oncall_team.phone_number
                      END
WHERE  inventory.oncall_team.id = $3::uuid;`

	OncallRemove = `
DELETE FROM inventory.oncall_team
WHERE  inventory.oncall_team.id = $1::uuid;`

	OncallMemberAssign = `
INSERT INTO inventory.oncall_membership (
            oncall_id,
            user_id,
            created_by)
SELECT $1::uuid,
       $2::uuid,
       ( SELECT inventory.user.id FROM inventory.user
         LEFT JOIN auth.admin
         ON inventory.user.uid = auth.admin.user_uid
         WHERE (   inventory.user.uid = $3::varchar
                OR auth.admin.uid     = $3::varchar ))
WHERE NOT EXISTS (
   SELECT inventory.oncall_membership.oncall_id
   FROM   inventory.oncall_membership
   WHERE  inventory.oncall_membership.oncall_id = $1::uuid
     AND  inventory.oncall_membership.user_id = $2::uuid);`

	OncallMemberUnassign = `
DELETE FROM inventory.oncall_membership
WHERE  inventory.oncall_membership.oncall_id = $1::uuid
  AND  inventory.oncall_membership.user_id = $2::uuid;`

	OncallMemberList = `
SELECT inventory.oncall_membership.user_id,
       inventory.user.uid
FROM   inventory.oncall_membership
JOIN   inventory.user
  ON   inventory.oncall_membership.user_id = inventory.user.id
WHERE  inventory.oncall_membership.oncall_id = $1::uuid;`
)

func init() {
	m[OncallAdd] = `OncallAdd`
	m[OncallList] = `OncallList`
	m[OncallMemberAssign] = `OncallMemberAssign`
	m[OncallMemberList] = `OncallMemberList`
	m[OncallMemberUnassign] = `OncallMemberUnassign`
	m[OncallRemove] = `OncallRemove`
	m[OncallSearch] = `OncallSearch`
	m[OncallShow] = `OncallShow`
	m[OncallUpdate] = `OncallUpdate`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
