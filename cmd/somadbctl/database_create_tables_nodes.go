package main

func createTablesNodes(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 10)

	queryMap["createTableNodes"] = `
create table if not exists soma.nodes (
    node_id                     uuid            PRIMARY KEY,
    node_asset_id               numeric(16,0)   UNIQUE NOT NULL,
    node_name                   varchar(256)    NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.team ( id ) DEFERRABLE,
    server_id                   uuid            NOT NULL REFERENCES inventory.servers ( server_id ) DEFERRABLE,
    object_state                varchar(64)     NOT NULL DEFAULT 'unassigned' REFERENCES soma.object_states ( object_state ) DEFERRABLE,
    node_online                 boolean         NOT NULL DEFAULT 'yes',
    node_deleted                boolean         NOT NULL DEFAULT 'no',
    created_by                  uuid            NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000' REFERENCES inventory.user ( id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    UNIQUE ( node_id, organizational_team_id )
);`
	queries[idx] = "createTableNodes"
	idx++

	queryMap["createTableNodeBucketAssignment"] = `
create table if not exists soma.node_bucket_assignment (
    node_id                     uuid            NOT NULL,
    bucket_id                   uuid            NOT NULL,
    organizational_team_id      uuid            NOT NULL REFERENCES inventory.team ( id ) DEFERRABLE,
    UNIQUE ( node_id ),
    UNIQUE ( node_id, bucket_id ),
    FOREIGN KEY ( node_id, organizational_team_id ) REFERENCES soma.nodes ( node_id, organizational_team_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, organizational_team_id ) REFERENCES soma.buckets ( bucket_id, organizational_team_id ) DEFERRABLE
);`
	queries[idx] = "createTableNodeBucketAssignment"
	idx++

	queryMap["createUniqueIndexNodeOnline"] = `
create unique index _unique_node_online
    on soma.nodes ( node_name, node_online )
    where node_online
;`
	queries[idx] = "createUniqueIndexNodeOnline"
	idx++

	queryMap["createTableNodeOncallProperty"] = `
create table if not exists soma.node_oncall_property (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    node_id                     uuid            NOT NULL REFERENCES soma.nodes ( node_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_team ( id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repository (id) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    UNIQUE ( node_id ),
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableNodeOncallProperty"
	idx++

	queryMap["createTableNodeServiceProperty"] = `
create table if not exists soma.node_service_property (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    node_id                     uuid            NOT NULL REFERENCES soma.nodes ( node_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    service_id                  uuid            NOT NULL,
    team_id                     uuid            NOT NULL REFERENCES inventory.team ( id ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repository (id) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    CONSTRAINT __node_service_property_unique_assignment UNIQUE ( node_id, service_id, view ) DEFERRABLE,
    CONSTRAINT __node_service_property_service_exists FOREIGN KEY ( team_id, service_id ) REFERENCES soma.service_property ( team_id, id ) DEFERRABLE,
    CONSTRAINT __node_service_property_service_owner FOREIGN KEY ( node_id, team_id ) REFERENCES soma.nodes ( node_id, organizational_team_id ) DEFERRABLE,
    CONSTRAINT __node_service_property_same_repository FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableNodeServiceProperty"
	idx++

	queryMap["createTableNodeSystemProperties"] = `
create table if not exists soma.node_system_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    node_id                     uuid            NOT NULL REFERENCES nodes ( node_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES views ( view ) DEFERRABLE,
    system_property             varchar(64)     NOT NULL REFERENCES soma.system_properties ( system_property ) DEFERRABLE,
    source_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    object_type                 varchar(64)     NOT NULL DEFAULT 'node' REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    repository_id               uuid            NOT NULL REFERENCES soma.repository (id) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    inherited                   boolean         NOT NULL DEFAULT 'yes',
    value                       text            NOT NULL,
    FOREIGN KEY ( system_property, object_type, inherited ) REFERENCES soma.system_property_validity ( system_property, object_type, inherited ) DEFERRABLE,
    CHECK( inherited OR object_type = 'node' ),
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableNodeSystemProperties"
	idx++

	// restrict all system properties to once per cluster+view, except
	// tags which would be silly if limited to once
	queryMap["createUniqueIndexNodeSystemProperties"] = `
create unique index _unique_node_system_properties
    on soma.node_system_properties ( node_id, system_property, view )
    where system_property != 'tag'
;`
	queries[idx] = "createUniqueIndexNodeSystemProperties"
	idx++

	queryMap["createTableNodeCustomProperties"] = `
create table if not exists soma.node_custom_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    node_id                     uuid            NOT NULL REFERENCES soma.nodes ( node_id ) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    custom_property_id          uuid            NOT NULL,
    bucket_id                   uuid            NOT NULL,
    repository_id               uuid            NOT NULL,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    UNIQUE ( node_id, custom_property_id, view ),
    -- ensure node is in this bucket
    -- ensure bucket is in this repository
    -- ensure custom_property is defined for this repository
    FOREIGN KEY ( node_id, bucket_id ) REFERENCES soma.node_bucket_assignment ( node_id, bucket_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE,
    FOREIGN KEY ( repository_id, custom_property_id ) REFERENCES soma.custom_properties ( repository_id, custom_property_id ) DEFERRABLE,
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableNodeCustomProperties"

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
