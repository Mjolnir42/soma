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
	SupervisorPermissionStatements = ``

	PermissionLoad = `
SELECT id,
       name,
       category
FROM   soma.permission;`

	PermissionAdd = `
INSERT INTO soma.permission (
            id,
            name,
            category,
            created_by
)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar,
       ( SELECT id
         FROM   inventory.user
         WHERE  uid = $4::varchar)
WHERE NOT EXISTS (
      SELECT id
      FROM   soma.permission
      WHERE  name = $2::varchar
        AND  category = $3::varchar);`

	PermissionLinkGrant = `
INSERT INTO soma.permission_grant_map (
            category,
            permission_id,
            granted_category,
            granted_permission_id)
SELECT $1::varchar,
       $2::uuid,
       $3::varchar,
       $4::uuid
WHERE  NOT EXISTS (
   -- a permission can not have two grant records
   SELECT permission_id
   FROM   soma.permission_grant_map
   WHERE  permission_id = $2::uuid);`

	PermissionLookupGrantID = `
SELECT permission_id
FROM   soma.permission_grant_map
WHERE  granted_permission_id = $1::uuid;`

	PermissionRevokeGlobal = `
DELETE FROM soma.authorizations_global
WHERE       permission_id = $1::uuid;`

	PermissionRevokeRepository = `
DELETE FROM soma.authorizations_repository
WHERE       permission_id = $1::uuid;`

	PermissionRevokeTeam = `
DELETE FROM soma.authorizations_team
WHERE       permission_id = $1::uuid;`

	PermissionRevokeMonitoring = `
DELETE FROM soma.authorizations_monitoring
WHERE       permission_id = $1::uuid;`

	PermissionUnmapAll = `
DELETE FROM soma.permission_map
WHERE       permission_id = $1::uuid;`

	PermissionRemove = `
DELETE FROM soma.permission
WHERE       id = $1::uuid;`

	PermissionRemoveLink = `
DELETE FROM soma.permission_grant_map
WHERE       granted_permission_id = $1::uuid;`

	PermissionRemoveByName = `
DELETE FROM soma.permission
WHERE       name = $1::varchar
AND         category = $2::varchar;`

	PermissionList = `
SELECT id,
       name
FROM   soma.permission
WHERE  category = $1::varchar;`

	PermissionShow = `
SELECT soma.permission.id,
       soma.permission.name,
       soma.permission.category,
       inventory.user.uid,
       soma.permission.created_at
FROM   soma.permission
JOIN   inventory.user
ON     soma.permission.created_by = inventory.user.id
WHERE  soma.permission.id = $1::uuid
  AND  soma.permission.category = $2::varchar;`

	PermissionSearchByName = `
SELECT id,
       name
FROM   soma.permission
WHERE  name = $1::varchar
  AND  category= $2::varchar;`

	PermissionMapLoad = `
SELECT    soma.permission_map.id,
          soma.permission_map.category,
          soma.permission_map.permission_id,
          soma.permission.name,
          soma.permission_map.section_id,
          soma.section.name,
          soma.permission_map.action_id,
          soma.action.name
FROM      soma.permission_map
JOIN      soma.permission
  ON      soma.permission_map.permission_id = soma.permission.id
JOIN      soma.section
  ON      soma.permission_map.section_id = soma.section.id
LEFT JOIN soma.action
  ON      soma.permission_map.action_id = soma.action.id;`

	PermissionMappedActions = `
SELECT soma.action.id,
       soma.action.name,
       soma.section.id,
       soma.section.name,
       soma.action.category
FROM   soma.permission
JOIN   soma.permission_map
  ON   soma.permission.id = soma.permission_map.permission_id
JOIN   soma.section
  ON   soma.permission_map.section_id = soma.section.id
JOIN   soma.action
  ON   soma.permission_map.action_id = soma.action.id
WHERE  soma.permission.id = $1::uuid
  AND  soma.permission.category = $2::varchar
  AND  soma.permission_map.action_id IS NOT NULL
  AND  soma.action.section_id = soma.permission_map.section_id;`

	PermissionMappedSections = `
SELECT soma.section.id,
       soma.section.name,
       soma.section.category
FROM   soma.permission
JOIN   soma.permission_map spm
  ON   soma.permission.id = soma.permission_map.permission_id
JOIN   soma.section
  ON   soma.permission_map.section_id = soma.section.id
WHERE  soma.permission.id = $1::uuid
  AND  soma.permission.category = $2::varchar
  AND  soma.permission_map.action_id IS NULL;`

	PermissionMapEntry = `
INSERT INTO soma.permission_map (
            id,
            category,
            permission_id,
            section_id,
            action_id,
            created_by)
SELECT $1::uuid,
       $2::varchar,
       $3::uuid,
       $4::uuid,
       $5::uuid,
       ( SELECT inventory.user.id
         FROM   inventory.user
         WHERE  inventory.user.uid = $6::varchar);`

	PermissionUnmapEntry = `
DELETE FROM soma.permission_map
WHERE       permission_id = $1::uuid
  AND       category = $2::varchar
  AND       section_id = $3::uuid
  AND       (action_id = $4::uuid OR ($4::uuid IS NULL AND action_id IS NULL));`
)

func init() {
	m[PermissionAdd] = `PermissionAdd`
	m[PermissionLinkGrant] = `PermissionLinkGrant`
	m[PermissionList] = `PermissionList`
	m[PermissionLoad] = `PermissionLoad`
	m[PermissionLookupGrantID] = `PermissionLookupGrantID`
	m[PermissionMapEntry] = `PermissionMapEntry`
	m[PermissionMapLoad] = `PermissionMapLoad`
	m[PermissionMappedActions] = `PermissionMappedActions`
	m[PermissionMappedSections] = `PermissionMappedSections`
	m[PermissionRemoveByName] = `PermissionRemoveByName`
	m[PermissionRemoveLink] = `PermissionRemoveLink`
	m[PermissionRemove] = `PermissionRemove`
	m[PermissionRevokeGlobal] = `PermissionRevokeGlobal`
	m[PermissionRevokeMonitoring] = `PermissionRevokeMonitoring`
	m[PermissionRevokeRepository] = `PermissionRevokeRepository`
	m[PermissionRevokeTeam] = `PermissionRevokeTeam`
	m[PermissionSearchByName] = `PermissionSearchByName`
	m[PermissionShow] = `PermissionShow`
	m[PermissionUnmapAll] = `PermissionUnmapAll`
	m[PermissionUnmapEntry] = `PermissionUnmapEntry`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
