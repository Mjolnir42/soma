/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * All rights reserved
 */

package stmt

const LoadPermissions = `
SELECT permission_id,
       permission_name
FROM   soma.permissions;`

const AddPermissionCategory = `
INSERT INTO soma.permission_types (
            permission_type,
            created_by
)
SELECT $1::varchar,
       $2::uuid
WHERE NOT EXISTS (
      SELECT permission_type
      FROM   soma.permission_types
      WHERE  permission_type = $1::varchar
);`

const DeletePermissionCategory = `
DELETE FROM soma.permission_types
WHERE permission_type = $1::varchar;`

const ListPermissionCategory = `
SELECT spt.permission_type
FROM   soma.permission_types spt;`

const ShowPermissionCategory = `
SELECT spt.permission_type,
       iu.user_uid,
       spt.created_at
FROM   soma.permission_types spt
JOIN   inventory.users iu
ON     spt.created_by = iu.user_id
WHERE  spt.permission_type = $1::varchar;`

const AddPermission = `
INSERT INTO soma.permissions (
            permission_id,
            permission_name,
            permission_type,
            created_by
)
SELECT $1::uuid,
       $2::varchar,
       $3::varchar,
       $4::uuid
WHERE NOT EXISTS (
      SELECT permission_name
      FROM   soma.permissions
      WHERE  permission_name = $2::varchar
);`

const DeletePermission = `
DELETE FROM soma.permissions
WHERE permission_id = $1::uuid;`

const ListPermission = `
SELECT permission_id,
       permission_name
FROM   soma.permissions;`

const ShowPermission = `
SELECT sp.permission_id,
       sp.permission_name,
       sp.permission_type,
       iu.user_uid,
	   sp.created_at
FROM   soma.permissions sp
JOIN   inventory.users iu
ON     sp.created_by = iu.user_id
WHERE  sp.permission_name = $1::varchar;`

const SearchPermissionByName = `
SELECT permission_id,
       permission_name
FROM   soma.permissions
WHERE  permission_name = $1::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix