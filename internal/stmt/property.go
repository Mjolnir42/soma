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
	PropertyStatements = ``

	ServiceLookup = `
SELECT soma.service_property.id,
       soma.service_property.name
FROM   soma.repository
JOIN   soma.service_property
ON     soma.repository.team_id = soma.service_property.team_id
WHERE  soma.repository.id = $1::uuid
AND    (soma.service_property.name = $3::varchar OR soma.service_property.id = $2::uuid)
AND    soma.repository.team_id = $4::uuid;`

	ServiceAttributes = `
SELECT soma.service_property_value.attribute,
       soma.service_property_value.value
FROM   soma.repository
JOIN   soma.service_property
ON     soma.repository.team_id = soma.service_property.team_id
JOIN   soma.service_property_value
ON     soma.service_property.team_id = soma.service_property_value.team_id
AND    soma.service_property.id = soma.service_property_value.service_id
WHERE  soma.repository.id = $1::uuid
AND    soma.service_property.id = $2::uuid
AND    soma.repository.team_id = $3::uuid;`

	PropertySystemList = `
SELECT system_property
FROM   soma.system_properties;`

	PropertyServiceList = `
SELECT id,
       name,
       team_id
FROM   soma.service_property
WHERE  team_id = $1::uuid;`

	PropertyNativeList = `
SELECT native_property
FROM   soma.native_properties;`

	PropertyTemplateList = `
SELECT id,
       name
FROM   soma.template_property;`

	PropertyCustomList = `
SELECT custom_property_id,
       repository_id,
       custom_property
FROM   soma.custom_properties
WHERE  repository_id = $1::uuid;`

	PropertySystemShow = `
SELECT system_property
FROM   soma.system_properties
WHERE  system_property = $1::varchar;`

	PropertyNativeShow = `
SELECT native_property
FROM   soma.native_properties
WHERE  native_property = $1::varchar;`

	PropertyCustomShow = `
SELECT custom_property_id,
       repository_id,
       custom_property
FROM   soma.custom_properties
WHERE  custom_property_id = $1::uuid
AND    repository_id = $2::uuid;`

	PropertyServiceShow = `
SELECT ssp.id,
       ssp.name,
       ssp.team_id,
       sspv.attribute,
       sspv.value
FROM   soma.service_property ssp
JOIN   soma.service_property_value sspv
ON     ssp.id = sspv.service_id
WHERE  ssp.id = $1::uuid;`

	PropertyTemplateShow = `
SELECT stp.id,
       stp.name,
       stpv.attribute,
       stpv.value
FROM   soma.template_property stp
JOIN   soma.template_property_value stpv
ON     stp.id = stpv.template_id
WHERE  stp.id = $1::uuid;`

	PropertySystemAdd = `
INSERT INTO soma.system_properties (
            system_property)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT system_property
   FROM   soma.system_properties
   WHERE  system_property = $1::varchar);`

	PropertyNativeAdd = `
INSERT INTO soma.native_properties (
            native_property)
SELECT $1::varchar
WHERE  NOT EXISTS (
   SELECT native_property
   FROM   soma.native_properties
   WHERE  native_property = $1::varchar);`

	PropertyCustomAdd = `
INSERT INTO soma.custom_properties (
            custom_property_id,
            repository_id,
            custom_property)
SELECT $1::uuid, $2::uuid, $3::varchar
WHERE  NOT EXISTS (
   SELECT custom_property
   FROM   soma.custom_properties
   WHERE  custom_property = $3::varchar
     AND  repository_id = $2::uuid);`

	PropertyServiceAdd = `
INSERT INTO soma.service_property (
            id,
            name,
            team_id
            )
SELECT $1::uuid, $2::varchar, $3::uuid
WHERE  NOT EXISTS (
   SELECT name
   FROM   soma.service_property
   WHERE  team_id = $3::uuid
   AND    name = $2::varchar);`

	PropertyServiceAttributeAdd = `
INSERT INTO soma.service_property_value (
            team_id,
            service_id,
            attribute,
            value)
SELECT $1::uuid, $2::uuid, $3::varchar, $4::varchar;`

	PropertyTemplateAdd = `
INSERT INTO soma.template_property (
            id,
            name)
SELECT $1::uuid, $2::varchar
WHERE  NOT EXISTS (
   SELECT name
   FROM   soma.template_property
   WHERE  name = $2::varchar);`

	PropertyTemplateAttributeAdd = `
INSERT INTO soma.template_property_value (
            template_id,
            attribute,
            value)
SELECT $1::uuid, $2::varchar, $3::varchar;`

	PropertySystemDel = `
DELETE FROM soma.system_properties
WHERE  system_property = $1::varchar;`

	PropertyNativeDel = `
DELETE FROM soma.native_properties
WHERE  native_property = $1::varchar;`

	PropertyCustomDel = `
DELETE FROM soma.custom_properties
WHERE  repository_id = $1::uuid
AND    custom_property_id = $2::uuid;`

	PropertyServiceDel = `
DELETE FROM soma.service_property
WHERE  id = $1::uuid;`

	PropertyServiceAttributeDel = `
DELETE FROM soma.service_property_value
WHERE  service_id = $1::uuid;`

	PropertyTemplateDel = `
DELETE FROM soma.template_property
WHERE  id = $1::uuid;`

	PropertyTemplateAttributeDel = `
DELETE FROM soma.template_property_value
WHERE  template_id = $1::uuid;`
)

func init() {
	m[PropertyCustomAdd] = `PropertyCustomAdd`
	m[PropertyCustomDel] = `PropertyCustomDel`
	m[PropertyCustomList] = `PropertyCustomList`
	m[PropertyCustomShow] = `PropertyCustomShow`
	m[PropertyNativeAdd] = `PropertyNativeAdd`
	m[PropertyNativeDel] = `PropertyNativeDel`
	m[PropertyNativeList] = `PropertyNativeList`
	m[PropertyNativeShow] = `PropertyNativeShow`
	m[PropertyServiceAdd] = `PropertyServiceAdd`
	m[PropertyServiceAttributeAdd] = `PropertyServiceAttributeAdd`
	m[PropertyServiceAttributeDel] = `PropertyServiceAttributeDel`
	m[PropertyServiceDel] = `PropertyServiceDel`
	m[PropertyServiceList] = `PropertyServiceList`
	m[PropertyServiceShow] = `PropertyServiceShow`
	m[PropertySystemAdd] = `PropertySystemAdd`
	m[PropertySystemDel] = `PropertySystemDel`
	m[PropertySystemList] = `PropertySystemList`
	m[PropertySystemShow] = `PropertySystemShow`
	m[PropertyTemplateAdd] = `PropertyTemplateAdd`
	m[PropertyTemplateAttributeAdd] = `PropertyTemplateAttributeAdd`
	m[PropertyTemplateAttributeDel] = `PropertyTemplateAttributeDel`
	m[PropertyTemplateDel] = `PropertyTemplateDel`
	m[PropertyTemplateList] = `PropertyTemplateList`
	m[PropertyTemplateShow] = `PropertyTemplateShow`
	m[ServiceAttributes] = `ServiceAttributes`
	m[ServiceLookup] = `ServiceLookup`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
