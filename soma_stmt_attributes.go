package main

const stmtAttributeList = `
SELECT service_property_attribute
FROM   soma.service_property_attributes;`

const stmtAttributeShow = `
SELECT service_property_attribute
FROM   soma.service_property_attributes
WHERE  service_property_attribute = $1::varchar;`

const stmtAttributeAdd = `
INSERT INTO soma.service_property_attributes (
            service_property_attribute)
SELECT $1::varchar WHERE NOT EXISTS (
       SELECT service_property_attribute
       FROM   soma.service_property_attributes
       WHERE  service_property_attribute = $1::varchar);`

const stmtAttributeDelete = `
DELETE FROM soma.service_property_attributes
WHERE       service_property_attribute = $1::varchar;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
