package main

/*
 * Statements for GROUP actions
 */

const tkStmtGroupCreate = `
INSERT INTO soma.groups (
            group_id,
            bucket_id,
            group_name,
            object_state,
            organizational_team_id,
            created_by)
SELECT $1::uuid,
       $2::uuid,
       $3::varchar,
       $4::varchar,
       $5::uuid,
       user_id
FROM   inventory.users iu
WHERE  iu.user_uid = $6::varchar;`

const tkStmtGroupUpdate = `
UPDATE soma.groups
SET    object_state = $2::varchar
WHERE  group_id = $1::uuid;`

const tkStmtGroupDelete = `
DELETE FROM soma.groups
WHERE       group_id = $1::uuid;`

const tkStmtGroupMemberNewNode = `
INSERT INTO soma.group_membership_nodes (
            group_id,
            child_node_id,
            bucket_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

const tkStmtGroupMemberNewCluster = `
INSERT INTO soma.group_membership_clusters (
            group_id,
            child_cluster_id,
            bucket_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

const tkStmtGroupMemberNewGroup = `
INSERT INTO soma.group_membership_groups (
            group_id,
            child_group_id,
            bucket_id)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid;`

const tkStmtGroupMemberRemoveNode = `
DELETE FROM soma.group_membership_nodes
WHERE       group_id = $1::uuid
AND         child_node_id = $2::uuid;`

const tkStmtGroupMemberRemoveCluster = `
DELETE FROM soma.group_membership_clusters
WHERE       group_id = $1::uuid
AND         child_cluster_id = $2::uuid;`

const tkStmtGroupMemberRemoveGroup = `
DELETE FROM soma.group_membership_groups
WHERE       group_id = $1::uuid
AND         child_group_id = $2::uuid;`

const tkStmtGroupPropertyOncallCreate = `
INSERT INTO soma.group_oncall_properties (
            instance_id,
            source_instance_id,
            group_id,
            view,
            oncall_duty_id,
            repository_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::uuid,
       $7::boolean,
       $8::boolean;`

const tkStmtGroupPropertyOncallDelete = `
DELETE FROM soma.group_oncall_properties
WHERE       instance_id = $1::uuid;`

const tkStmtGroupPropertyServiceCreate = `
INSERT INTO soma.group_service_properties (
            instance_id,
            source_instance_id,
            group_id,
            view,
            service_property,
            organizational_team_id,
            repository_id,
            inheritance_enabled,
            children_only)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::uuid,
       $7::uuid,
       $8::boolean,
       $9::boolean;`

const tkStmtGroupPropertyServiceDelete = `
DELETE FROM soma.group_service_properties
WHERE       instance_id = $1::uuid;`

const tkStmtGroupPropertySystemCreate = `
INSERT INTO soma.group_system_properties (
            instance_id,
            source_instance_id,
            group_id,
            view,
            system_property,
            source_type,
            repository_id,
            inheritance_enabled,
            children_only,
            value,
            inherited)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::varchar,
       $6::varchar,
       $7::uuid,
       $8::boolean,
       $9::boolean,
       $10::text,
       $11::boolean;`

const tkStmtGroupPropertySystemDelete = `
DELETE FROM soma.group_system_properties
WHERE       instance_id = $1::uuid;`

const tkStmtGroupPropertyCustomCreate = `
INSERT INTO soma.group_custom_properties (
            instance_id,
            source_instance_id,
            group_id,
            view,
            custom_property_id,
            bucket_id,
            repository_id,
            inheritance_enabled,
            children_only,
            value)
SELECT $1::uuid,
       $2::uuid,
       $3::uuid,
       $4::varchar,
       $5::uuid,
       $6::uuid,
       $7::uuid,
       $8::boolean,
       $9::boolean,
       $10::text;`

const tkStmtGroupPropertyCustomDelete = `
DELETE FROM soma.group_custom_properties
WHERE       instance_id = $1::uuid;`

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
