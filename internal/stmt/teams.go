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

	ListTeams = `
SELECT organizational_team_id,
       organizational_team_name
FROM   inventory.organizational_teams;`

	ShowTeams = `
SELECT organizational_team_id,
       organizational_team_name,
       organizational_team_ldap_id,
       organizational_team_system
FROM   inventory.organizational_teams
WHERE  organizational_team_id = $1;`

	SyncTeams = `
SELECT organizational_team_id,
       organizational_team_name,
       organizational_team_ldap_id,
       organizational_team_system
FROM   inventory.organizational_teams
WHERE  NOT organizational_team_system;`

	TeamAdd = `
INSERT INTO inventory.organizational_teams (
            organizational_team_id,
            organizational_team_name,
            organizational_team_ldap_id,
            organizational_team_system)
SELECT $1::uuid, $2::varchar, $3::numeric, $4
WHERE  NOT EXISTS (
   SELECT organizational_team_id
   FROM   inventory.organizational_teams
   WHERE  organizational_team_id = $1::uuid
      OR  organizational_team_name = $2::varchar
      OR  organizational_team_ldap_id = $3::numeric);`

	TeamUpdate = `
UPDATE inventory.organizational_teams
SET    organizational_team_name = $1::varchar,
       organizational_team_ldap_id = $2::numeric,
       organizational_team_system = $3::boolean
WHERE  organizational_team_id = $4::uuid;`

	TeamDel = `
DELETE FROM inventory.organizational_teams
WHERE       organizational_team_id = $1;`

	TeamMembers = `
SELECT iu.user_id,
       iu.user_uid
FROM   inventory.organizational_teams iot
JOIN   inventory.users iu
  ON   iot.organizational_team_id = iu.organizational_team_id
WHERE  iot.organizational_team_id = $1::uuid
  AND  NOT iu.user_is_deleted;`

	TeamLoad = `
SELECT organizational_team_id,
       organizational_team_name,
       organizational_team_ldap_id,
       organizational_team_system
FROM   inventory.organizational_teams;`
)

func init() {
	m[ListTeams] = `ListTeams`
	m[ShowTeams] = `ShowTeams`
	m[SyncTeams] = `SyncTeams`
	m[TeamAdd] = `TeamAdd`
	m[TeamDel] = `TeamDel`
	m[TeamLoad] = `TeamLoad`
	m[TeamMembers] = `TeamMembers`
	m[TeamUpdate] = `TeamUpdate`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
