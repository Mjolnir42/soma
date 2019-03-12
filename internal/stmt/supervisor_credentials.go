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
	SupervisorCredentialStatements = ``

	LoadAllUserCredentials = `
SELECT aua.user_id,
       aua.crypt,
       aua.reset_pending,
       aua.valid_from,
       aua.valid_until,
       iu.uid
FROM   inventory.user iu
JOIN   auth.user_authentication aua
ON     iu.id = aua.user_id
WHERE  iu.id != '00000000-0000-0000-0000-000000000000'::uuid
AND    NOW() < aua.valid_until
AND    NOT iu.is_deleted
AND    iu.is_active;`

	LoadAllAdminCredentials = `
SELECT auth.admin_authentication.admin_id,
       auth.admin_authentication.crypt,
       auth.admin_authentication.reset_pending,
       auth.admin_authentication.valid_from,
       auth.admin_authentication.valid_until,
       auth.admin.uid
FROM   auth.admin
JOIN   auth.admin_authentication
ON     auth.admin.id = auth.admin_authentication.admin_id
JOIN   inventory.user
ON     auth.admin.user_uid = inventory.user.uid
WHERE  inventory.user.id != '00000000-0000-0000-0000-000000000000'::uuid
AND    NOT inventory.user.is_deleted
AND    inventory.user.is_active
AND    auth.admin.is_active
AND    NOW() > auth.admin_authentication.valid_from
AND    NOW() < auth.admin_authentication.valid_until;`

	FindUserID = `
SELECT id
FROM   inventory.user
WHERE  uid = $1::varchar
AND    NOT is_deleted;`

	FindAdminID = `
SELECT id
FROM   auth.admin
WHERE  uid = $1::varchar;`

	FindUserName = `
SELECT uid
FROM   inventory.user
WHERE  id = $1::uuid
AND    NOT is_deleted;`

	CheckUserActive = `
SELECT is_active
FROM   inventory.user
WHERE  id = $1::uuid
AND    NOT is_deleted;`

	CheckAdminActive = `
SELECT is_active
FROM   auth.admin
WHERE  id = $1::uuid;`

	SetUserCredential = `
INSERT INTO auth.user_authentication (
            user_id,
            crypt,
            reset_pending,
            valid_from,
            valid_until
) VALUES (
            $1::uuid,
            $2::text,
            'no'::boolean,
            $3::timestamptz,
			$4::timestamptz
);`

	SetAdminCredential = `
INSERT INTO auth.admin_authentication (
            admin_id,
            crypt,
            reset_pending,
            valid_from,
            valid_until
) VALUES (
            $1::uuid,
            $2::text,
            'no'::boolean,
            $3::timestamptz,
			$4::timestamptz
);`

	ActivateUser = `
UPDATE inventory.user
SET    is_active = 'yes'::boolean
WHERE  id = $1::uuid;`

	ActivateAdminUser = `
UPDATE auth.admin
SET    is_active = 'yes'::boolean
WHERE  id = $1::uuid;`

	InvalidateUserCredential = `
UPDATE auth.user_authentication aua
SET    valid_until = $1::timestamptz
FROM   inventory.user iu
WHERE  aua.user_id = iu.id
  AND  aua.user_id = $2::uuid
  AND  NOW() < aua.valid_until
  AND  iu.is_active = 'yes'::boolean
  AND  NOT iu.is_deleted
  AND  iu.id != '00000000-0000-0000-0000-000000000000'::uuid;`
)

func init() {
	m[ActivateUser] = `ActivateUser`
	m[ActivateAdminUser] = `ActivateAdminUser`
	m[CheckUserActive] = `CheckUserActive`
	m[CheckAdminActive] = `CheckAdminActive`
	m[FindUserID] = `FindUserID`
	m[FindAdminID] = `FindAdminID`
	m[FindUserName] = `FindUserName`
	m[InvalidateUserCredential] = `InvalidateUserCredential`
	m[LoadAllUserCredentials] = `LoadAllUserCredentials`
	m[LoadAllAdminCredentials] = `LoadAllAdminCredentials`
	m[SetUserCredential] = `SetUserCredential`
	m[SetAdminCredential] = `SetAdminCredential`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
