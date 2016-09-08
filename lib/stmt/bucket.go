/*-
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const BucketOncProps = `
SELECT op.instance_id,
       op.source_instance_id,
       op.view,
       op.oncall_duty_id,
       iodt.oncall_duty_name
FROM   soma.bucket_oncall_properties op
JOIN   inventory.oncall_duty_teams iodt
  ON   op.oncall_duty_id = iodt.oncall_duty_id
WHERE  op.bucket_id = $1::uuid;`

const BucketSvcProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.service_property
FROM   soma.bucket_service_properties sp
WHERE  sp.bucket_id = $1::uuid;`

const BucketSysProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.system_property,
       sp.value
FROM   soma.bucket_system_properties sp
WHERE  sp.bucket_id = $1::uuid;`

const BucketCstProps = `
SELECT cp.instance_id,
       cp.source_instance_id,
       cp.view,
       cp.custom_property_id,
       cp.value,
       scp.custom_property
FROM   soma.bucket_custom_properties cp
JOIN   soma.custom_properties scp
  ON   cp.custom_property_id = scp.custom_property_id
WHERE  cp.bucket_id = $1::uuid;`

const BucketSystemPropertyForDelete = `
SELECT view,
       system_property,
       value
FROM   soma.bucket_system_properties
WHERE  source_instance_id = $1::uuid
  AND  source_instance_id = instance_id;`

const BucketCustomPropertyForDelete = `
SELECT sbcp.view,
       sbcp.custom_property_id,
       sbcp.value,
       scp.custom_property
FROM   soma.bucket_custom_properties sbcp
JOIN   soma.custom_properties scp
  ON   sbcp.repository_id = scp.repository_id
 AND   sbcp.custom_property_id = scp.custom_property_id
WHERE  sbcp.source_instance_id = $1::uuid
  AND  sbcp.source_instance_id = sbcp.instance_id;`

const BucketOncallPropertyForDelete = `
SELECT sbop.view,
       sbop.oncall_duty_id,
       iodt.oncall_duty_name,
       iodt.oncall_duty_phone_number
FROM   soma.bucket_oncall_properties sbop
JOIN   inventory.oncall_duty_teams iodt
  ON   sbop.oncall_duty_id = iodt.oncall_duty_id
WHERE  sbop.source_instance_id = $1::uuid
  AND  sbop.source_instance_id = sbop.instance_id;`

const BucketServicePropertyForDelete = `
SELECT sbsp.view,
       sbsp.service_property
FROM   soma.bucket_service_properties sbsp
JOIN   soma.team_service_properties stsp
  ON   sbsp.organizational_team_id = stsp.organizational_team_id
 AND   sbsp.service_property = stsp.service_property
WHERE  sbsp.source_instance_id = $1::uuid
  AND  sbsp.source_instance_id = sbsp.instance_id;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix