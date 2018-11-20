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
	RepositoryStatements = ``

	RepoOncProps = `
SELECT op.instance_id,
       op.source_instance_id,
       op.view,
       op.oncall_duty_id,
       iot.name
FROM   soma.repository_oncall_properties op
JOIN   inventory.oncall_team iot
  ON   op.oncall_duty_id = iot.id
WHERE  op.repository_id = $1::uuid;`

	RepoSvcProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.service_id
FROM   soma.repository_service_property sp
WHERE  sp.repository_id = $1::uuid;`

	RepoSysProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.system_property,
       sp.value
FROM   soma.repository_system_properties sp
WHERE  sp.repository_id = $1::uuid;`

	RepoCstProps = `
SELECT cp.instance_id,
       cp.source_instance_id,
       cp.view,
       cp.custom_property_id,
       cp.value,
       scp.custom_property
FROM   soma.repository_custom_properties cp
JOIN   soma.custom_properties scp
  ON   cp.custom_property_id = scp.custom_property_id
WHERE  cp.repository_id = $1::uuid;`

	RepoSystemPropertyForDelete = `
SELECT view,
       system_property,
       value
FROM   soma.repository_system_properties
WHERE  source_instance_id = $1::uuid
  AND  source_instance_id = instance_id;`

	RepoCustomPropertyForDelete = `
SELECT srcp.view,
       srcp.custom_property_id,
       srcp.value,
       scp.custom_property
FROM   soma.repository_custom_properties srcp
JOIN   soma.custom_properties scp
  ON   srcp.repository_id = scp.repository_id
 AND   srcp.custom_property_id = scp.custom_property_id
WHERE  source_instance_id = $1::uuid
  AND  source_instance_id = instance_id;`

	RepoOncallPropertyForDelete = `
SELECT srop.view,
       srop.oncall_duty_id,
       iot.name,
       iot.phone_number
FROM   soma.repository_oncall_properties srop
JOIN   inventory.oncall_team iot
  ON   srop.oncall_duty_id = iot.id
WHERE  source_instance_id = $1::uuid
  AND  source_instance_id = instance_id;`

	RepoServicePropertyForDelete = `
SELECT srsp.view,
       srsp.service_id
FROM   soma.repository_service_property srsp
JOIN   soma.service_property ssp
  ON   srsp.team_id = ssp.team_id
 AND   srsp.service_id = ssp.id
WHERE  srsp.source_instance_id = $1::uuid
  AND  srsp.source_instance_id = srsp.instance_id;`

	RepoNameByID = `
SELECT name
FROM   soma.repository
WHERE  id = $1::uuid;`

	RepoByBucketID = `
SELECT sb.repository_id,
       soma.repository.name
FROM   soma.buckets sb
JOIN   soma.repository
  ON   sb.repository_id = soma.repository.id
WHERE  sb.bucket_id = $1::uuid;`
)

func init() {
	m[RepoByBucketID] = `RepoByBucketID`
	m[RepoCstProps] = `RepoCstProps`
	m[RepoCustomPropertyForDelete] = `RepoCustomPropertyForDelete`
	m[RepoNameByID] = `RepoNameByID`
	m[RepoOncProps] = `RepoOncProps`
	m[RepoOncallPropertyForDelete] = `RepoOncallPropertyForDelete`
	m[RepoServicePropertyForDelete] = `RepoServicePropertyForDelete`
	m[RepoSvcProps] = `RepoSvcProps`
	m[RepoSysProps] = `RepoSysProps`
	m[RepoSystemPropertyForDelete] = `RepoSystemPropertyForDelete`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
