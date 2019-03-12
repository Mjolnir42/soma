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
	ForestCustodianStatements = ``

	ForestRebuildDeleteChecks = `
UPDATE soma.checks sc
SET    deleted = 'yes'::boolean
WHERE  sc.repository_id = $1::uuid;`

	ForestRebuildDeleteInstances = `
UPDATE soma.check_instances sci
SET    deleted = 'yes'::boolean
FROM   soma.checks sc
WHERE  sci.check_id = sc.check_id
AND    sc.repository_id = $1::uuid;`

	ForestRepoNameByID = `
SELECT name,
       team_id
FROM   soma.repository
WHERE  id = $1::uuid;`

	ForestLoadRepository = `
SELECT id,
       name,
       is_deleted,
       is_active,
       team_id
FROM   soma.repository;`

	ForestAddRepository = `
INSERT INTO soma.repository (
            id,
            name,
            is_active,
            is_deleted,
            team_id,
            created_by)
SELECT      $1::uuid,
            $2::varchar,
            $3::boolean,
            $4::boolean,
            $5::uuid,
            inventory.user.id
FROM        inventory.user
LEFT JOIN   auth.admin
  ON        inventory.user.uid = auth.admin.user_uid
WHERE       (   inventory.user.uid = $6::varchar
             OR auth.admin.uid     = $6::varchar )
AND NOT EXISTS (
	SELECT  id
	FROM    soma.repository
	WHERE   id   = $1::uuid
	  OR    name = $2::varchar);`
)

func init() {
	m[ForestAddRepository] = `ForestAddRepository`
	m[ForestLoadRepository] = `ForestLoadRepository`
	m[ForestRebuildDeleteChecks] = `ForestRebuildDeleteChecks`
	m[ForestRebuildDeleteInstances] = `ForestRebuildDeleteInstances`
	m[ForestRepoNameByID] = `ForestRepoNameByID`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
