/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const GrantGlobalOrSystemToUser = `
INSERT INTO soma.authorizations_global (
    grant_id,
    user_id,
    permission_id,
    permission_type,
    created_by
)
VALUES (
    $1::uuid,
    $2::uuid,
    $3::uuid,
    $4::varchar,
    $5::uuid
);`

const RevokeGlobalOrSystemFromUser = `
DELETE FROM soma.authorizations_global
WHERE grant_id = $1::uuid;`

const GrantLimitedRepoToUser = `
INSERT INTO soma.authorizations_repository (
	grant_id,
	user_id,
	repository_id,
	permission_id,
	permission_type,
	created_by
)
VALUES (
	$1::uuid,
	$2::uuid,
	$3::uuid,
	$4::uuid,
	$5::varchar,
	$6::uuid
);`

const RevokeLimitedRepoFromUser = `
DELETE FROM soma.authorizations_repository
WHERE grant_id = $1::uuid;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
