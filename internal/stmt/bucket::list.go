/*-
 * Copyright (c) 2019, Jörg Pernfuß
 * Copyright (c) 2019, 1&1 IONOS SE
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt // import "github.com/mjolnir42/soma/internal/stmt"

const AuthorizedBucketList = `
-- $1 section.name       ::varchar
-- $2 action.name        ::varchar
-- $3 user.uid           ::varchar
-- $4 repository.id      ::uuid
-------------------------------
-- CASE1: root user has omnipotence permission
SELECT      soma.buckets.bucket_id,
            soma.buckets.bucket_name
FROM        inventory.user
JOIN        soma.authorizations_global
  ON        inventory.user.id = soma.authorizations_global.user_id
JOIN        soma.permission
  ON        soma.authorizations_global.permission_id = soma.permission.id
            -- unscoped, use carthesian product on all repositories
CROSS JOIN  soma.buckets
JOIN        soma.repository
  ON        soma.buckets.repository_id = soma.repository.id
WHERE       inventory.user.uid = $3::varchar
  AND       soma.authorizations_global.category = 'omnipotence'
  AND       soma.permission.name = 'omnipotence'
  AND       (   $1::varchar = 'bucket'
             OR $1::varchar = 'bucket-config')
  AND       $2::varchar = 'list'
  AND       soma.buckets.repository_id = $4::uuid
  AND NOT   soma.repository.is_deleted
  AND NOT   soma.buckets.bucket_deleted
UNION
-- CASE2: admin user has suitable system permission for requested section::action
SELECT      soma.buckets.bucket_id,
            soma.buckets.bucket_name
FROM        auth.admin
JOIN        soma.authorizations_global
  ON        auth.admin.id = soma.authorizations_global.admin_id
JOIN        soma.permission
  ON        soma.authorizations_global.permission_id = soma.permission.id
JOIN        soma.section
            -- system permissions have the category they grant as permission name
  ON        soma.permission.name = soma.section.category
JOIN        soma.action
  ON        soma.section.id = soma.action.section_id
            -- unscoped, use carthesian product on all repositories
CROSS JOIN  soma.buckets
JOIN        soma.repository
  ON        soma.buckets.repository_id = soma.repository.id
WHERE       auth.admin.uid = $3::varchar
  AND       auth.admin.is_active
  AND       soma.authorizations_global.category = 'system'
  AND       soma.section.name = $1::varchar
  AND       soma.action.name  = $2::varchar
  AND       (   $1::varchar = 'bucket'
             OR $1::varchar = 'bucket-config')
  AND       $2::varchar = 'list'
  AND       soma.buckets.repository_id = $4::uuid
  AND NOT   soma.repository.is_deleted
  AND NOT   soma.buckets.bucket_deleted
UNION
-- CASE 3:  regular user has repository scoped bucket::list
SELECT      soma.buckets.bucket_id,
            soma.buckets.bucket_name
FROM        inventory.user
JOIN        soma.authorizations_repository
            -- authorization could be on the user or inherited from the team
  ON        (   inventory.user.id      = soma.authorizations_repository.user_id
             OR inventory.user.team_id = soma.authorizations_repository.team_id)
JOIN        soma.permission_map
  ON        soma.authorizations_repository.permission_id = soma.permission_map.permission_id
JOIN        soma.section
  ON        soma.permission_map.section_id = soma.section.id
JOIN        soma.action
  ON        soma.section.id = soma.action.section_id
            -- grant must be scoped on target repository OR its buckets
JOIN        soma.buckets
  ON        (   soma.authorizations_repository.bucket_id     = soma.buckets.bucket_id
             OR soma.authorizations_repository.repository_id = soma.buckets.repository_id )
JOIN        soma.repository
  ON        soma.buckets.repository_id = soma.repository.id
WHERE       inventory.user.uid = $3::varchar
  AND       inventory.user.is_active
  AND NOT   inventory.user.is_deleted
  AND       soma.section.name = $1::varchar
  AND       soma.action.name  = $2::varchar
  AND       (   $1::varchar = 'bucket'
             OR $1::varchar = 'bucket-config' )
  AND       $2::varchar = 'list'
            -- section grant for all actions has soma.permission_map.action_id as NULL
  AND       (   soma.permission_map.action_id = soma.action.id
             OR soma.permission_map.action_id IS NULL                 )
  AND NOT   soma.repository.is_deleted
  AND NOT   soma.buckets.bucket_deleted
  AND       soma.buckets.repository_id = $4::uuid
UNION
-- CASE 4:  regular user has team scoped bucket::list, which allows to list all
--          buckets from that team
SELECT      soma.repository.id,
            soma.repository.name
FROM        inventory.user
JOIN        soma.authorizations_team
            -- authorization could be on the user or inherited from the team
  ON        (   inventory.user.id       = soma.authorizations_team.user_id
             OR inventory.user.team_id  = soma.authorizations_team.team_id)
JOIN        soma.permission_map
  ON        soma.authorizations_team.permission_id = soma.permission_map.permission_id
JOIN        soma.section
  ON        soma.permission_map.section_id = soma.section.id
JOIN        soma.action
  ON        soma.section.id = soma.action.section_id
            -- grant must be scoped on target bucket owner team
JOIN        soma.buckets
  ON        soma.authorizations_team.authorized_team_id = soma.buckets.organizational_team_id
JOIN        soma.repository
  ON        soma.buckets.repository_id = soma.repository.id
WHERE       inventory.user.uid = $3::varchar
  AND       inventory.user.is_active
  AND NOT   inventory.user.is_deleted
  AND       soma.section.name = $1::varchar
  AND       soma.action.name  = $2::varchar
  AND       (   $1::varchar = 'bucket'
             OR $1::varchar = 'bucket-config' )
  AND       $2::varchar = 'list'
            -- section grant for all actions has soma.permission_map.action_id as NULL
  AND       (   soma.permission_map.action_id = soma.action.id
             OR soma.permission_map.action_id IS NULL                 )
  AND NOT   soma.repository.is_deleted
  AND NOT   soma.buckets.bucket_deleted
  AND       soma.buckets.repository_id = $4::uuid;`

func init() {
	m[AuthorizedRepositoryList] = `AuthorizedRepositoryList`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
