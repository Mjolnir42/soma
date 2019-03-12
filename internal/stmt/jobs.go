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
	JobStatements = ``

	ListAllOutstandingJobs = `
SELECT id,
       type
FROM   soma.job
WHERE  status != 'processed';`

	ListScopedOutstandingJobs = `
SELECT sj.id,
       sj.type
FROM   inventory.user iu
JOIN   soma.job sj
  ON   iu.id = sj.user_id
WHERE  iu.uid = $1::varchar
UNION
SELECT sj.id,
       sj.type
FROM   inventory.user iu
JOIN   soma.job sj
  ON   iu.team_id = sj.team_id
WHERE  iu.uid = $1::varchar
  AND  sj.user_id NOT IN
  (    SELECT id FROM inventory.user
       WHERE uid = $1::varchar);`

	JobResultForID = `
SELECT id,
       status,
       result,
       type,
       serial,
       repository_id,
       user_id,
       team_id,
       queued_at,
       started_at,
       finished_at,
       error,
       job
FROM   soma.job
WHERE  id = $1::uuid;`

	JobResultsForList = `
SELECT id,
       status,
       result,
       type,
       serial,
       repository_id,
       user_id,
       team_id,
       queued_at,
       started_at,
       finished_at,
       error,
       job
FROM   soma.job
WHERE  id = any($1::uuid[]);`

	JobSave = `
INSERT INTO soma.job (
            id,
            status,
            result,
            type,
            repository_id,
            user_id,
            team_id,
            job)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar,
       $4::varchar,
       $5::uuid,
       iu.id,
       iu.team_id,
       $7::jsonb
FROM   inventory.user iu
WHERE  iu.uid = $6::varchar;`

	JobTypeMgmtList = `
SELECT id
FROM   soma.job_type;`

	JobTypeMgmtShow = `
SELECT sjt.id,
       sjt.name,
       iu.uid,
       sjt.created_at
FROM   soma.job_type sjt
JOIN   inventory.user iu
  ON   sjt.created_by = iu.id
WHERE  sjt.id = $1::uuid;`

	JobTypeMgmtAdd = `
INSERT INTO soma.job_type (
            id,
            name,
            created_by)
SELECT $1::uuid,
       $2::varchar,
       ( SELECT inventory.user.id FROM inventory.user
         LEFT JOIN auth.admin
         ON inventory.user.uid = auth.admin.user_uid
         WHERE (   inventory.user.uid = $3::varchar
                OR auth.admin.uid     = $3::varchar ))
WHERE  NOT EXISTS (
   SELECT  id
   FROM    soma.job_type
   WHERE   id = $1::uuid
      OR   name = $2::varchar);`

	JobTypeMgmtRemove = `
DELETE FROM soma.job_type
WHERE       id = $1::uuid;`

	JobTypeMgmtSearch = `
SELECT id,
       name
FROM   soma.job_type
WHERE  ( id = $1::uuid OR $1::uuid IS NULL )
  AND  ( name = $2::varchar OR $2::varchar IS NULL )
  AND  NOT ( $1::uuid IS NULL AND $2::varchar IS NULL );`

	JobResultMgmtList = `
SELECT id
FROM   soma.job_result;`

	JobResultMgmtShow = `
SELECT sjr.id,
       sjr.name,
       iu.uid,
       sjr.created_at
FROM   soma.job_result sjr
JOIN   inventory.user iu
  ON   sjr.created_by = iu.id
WHERE  sjr.id = $1::uuid;`

	JobResultMgmtAdd = `
INSERT INTO soma.job_result (
            id,
            name,
            created_by)
SELECT $1::uuid,
       $2::varchar,
       ( SELECT inventory.user.id FROM inventory.user
         LEFT JOIN auth.admin
         ON inventory.user.uid = auth.admin.user_uid
         WHERE (   inventory.user.uid = $3::varchar
                OR auth.admin.uid     = $3::varchar ))
WHERE  NOT EXISTS (
   SELECT  id
   FROM    soma.job_result
   WHERE   id = $1::uuid
      OR   name = $2::varchar);`

	JobResultMgmtRemove = `
DELETE FROM soma.job_result
WHERE       id = $1::uuid;`

	JobResultMgmtSearch = `
SELECT id,
       name
FROM   soma.job_result
WHERE  ( id = $1::uuid OR $1::uuid IS NULL )
  AND  ( name = $2::varchar OR $2::varchar IS NULL )
  AND  NOT ( $1::uuid IS NULL AND $2::varchar IS NULL );`

	JobStatusMgmtList = `
SELECT id
FROM   soma.job_status;`

	JobStatusMgmtShow = `
SELECT sjs.id,
       sjs.name,
       iu.uid,
       sjs.created_at
FROM   soma.job_status sjs
JOIN   inventory.user iu
  ON   sjs.created_by = iu.id
WHERE  sjs.id = $1::uuid;`

	JobStatusMgmtAdd = `
INSERT INTO soma.job_status (
            id,
            name,
            created_by)
SELECT $1::uuid,
       $2::varchar,
       ( SELECT inventory.user.id FROM inventory.user
         LEFT JOIN auth.admin
         ON inventory.user.uid = auth.admin.user_uid
         WHERE (   inventory.user.uid = $3::varchar
                OR auth.admin.uid     = $3::varchar ))
WHERE  NOT EXISTS (
   SELECT  id
   FROM    soma.job_status
   WHERE   id = $1::uuid
      OR   name = $2::varchar);`

	JobStatusMgmtRemove = `
DELETE FROM soma.job_status
WHERE       id = $1::uuid;`

	JobStatusMgmtSearch = `
SELECT id,
       name
FROM   soma.job_status
WHERE  ( id = $1::uuid OR $1::uuid IS NULL )
  AND  ( name = $2::varchar OR $2::varchar IS NULL )
  AND  NOT ( $1::uuid IS NULL AND $2::varchar IS NULL );`
)

func init() {
	m[JobResultForID] = `JobResultForID`
	m[JobResultsForList] = `JobResultsForList`
	m[JobSave] = `JobSave`
	m[ListAllOutstandingJobs] = `ListAllOutstandingJobs`
	m[ListScopedOutstandingJobs] = `ListScopedOutstandingJobs`
	m[JobTypeMgmtList] = `JobTypeMgmtList`
	m[JobTypeMgmtShow] = `JobTypeMgmtShow`
	m[JobTypeMgmtAdd] = `JobTypeMgmtAdd`
	m[JobTypeMgmtRemove] = `JobTypeMgmtRemove`
	m[JobTypeMgmtSearch] = `JobTypeMgmtSearch`
	m[JobResultMgmtList] = `JobResultMgmtList`
	m[JobResultMgmtShow] = `JobResultMgmtShow`
	m[JobResultMgmtAdd] = `JobResultMgmtAdd`
	m[JobResultMgmtRemove] = `JobResultMgmtRemove`
	m[JobResultMgmtSearch] = `JobResultMgmtSearch`
	m[JobStatusMgmtList] = `JobStatusMgmtList`
	m[JobStatusMgmtShow] = `JobStatusMgmtShow`
	m[JobStatusMgmtAdd] = `JobStatusMgmtAdd`
	m[JobStatusMgmtRemove] = `JobStatusMgmtRemove`
	m[JobStatusMgmtSearch] = `JobStatusMgmtSearch`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
