package main

func createTableRepositories(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 5)

	queryMap[`create__soma.repository`] = `
create table if not exists soma.repository (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    name                        varchar(128)    NOT NULL,
    is_deleted                  boolean         NOT NULL DEFAULT 'no',
    is_active                   boolean         NOT NULL DEFAULT 'yes',
    team_id                     uuid            NOT NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT _repository_primary_key PRIMARY KEY (id),
    CONSTRAINT _repository_creator_exists FOREIGN KEY (created_by) REFERENCES inventory.user (id) DEFERRABLE,
    CONSTRAINT _repository_team_exists FOREIGN KEY (team_id) REFERENCES inventory.team (id) DEFERRABLE,
    CONSTRAINT _repository_timezone_utc CHECK( EXTRACT( TIMEZONE FROM created_at ) = '0' )
);`
	queries[idx] = `create__soma.repository`
	idx++

	queryMap[`create__soma.repository__index__unique_name`] = `
CREATE UNIQUE INDEX _repository_unique_name
    ON soma.repository ( name, is_deleted )
    WHERE NOT is_deleted;`
	queries[idx] = `create__soma.repository__index__unique_name`
	idx++

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

func createTablesRepositoryProperties(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 5)

	queryMap["createTableRepositoryOncallProperty"] = `
create table if not exists soma.repository_oncall_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    repository_id               uuid            NOT NULL REFERENCES soma.repository (id) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    oncall_duty_id              uuid            NOT NULL REFERENCES inventory.oncall_duty_teams ( oncall_duty_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableRepositoryOncallProperty"
	idx++

	queryMap["createTableRepositoryServiceProperty"] = `
create table if not exists soma.repository_service_property (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    repository_id               uuid            NOT NULL REFERENCES soma.repository (id) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    service_id                  uuid            NOT NULL,
    team_id                     uuid            NOT NULL REFERENCES inventory.team ( id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    CONSTRAINT __repository_service_property_service_exists FOREIGN KEY ( team_id, service_id ) REFERENCES soma.service_property ( team_id, id ) DEFERRABLE,
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableRepositoryServiceProperty"
	idx++

	queryMap["createTableRepositorySystemProperties"] = `
create table if not exists soma.repository_system_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    repository_id               uuid            NOT NULL REFERENCES soma.repository (id) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    system_property             varchar(64)     NOT NULL REFERENCES soma.system_properties ( system_property ) DEFERRABLE,
    source_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    object_type                 varchar(64)     NOT NULL DEFAULT 'repository' REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    inherited                   boolean         NOT NULL DEFAULT 'yes',
    value                       text            NOT NULL,
    FOREIGN KEY ( system_property, object_type, inherited ) REFERENCES soma.system_property_validity ( system_property, object_type, inherited ) DEFERRABLE,
    CHECK( inherited OR object_type = 'repository' ),
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableRepositorySystemProperties"
	idx++

	queryMap["createTableRepositoryCustomProperty"] = `
create table if not exists soma.repository_custom_properties (
    instance_id                 uuid            NOT NULL REFERENCES soma.property_instances ( instance_id ) DEFERRABLE,
    source_instance_id          uuid            NOT NULL,
    repository_id               uuid            NOT NULL REFERENCES soma.repository (id) DEFERRABLE,
    view                        varchar(64)     NOT NULL DEFAULT 'any' REFERENCES soma.views ( view ) DEFERRABLE,
    custom_property_id          uuid            NOT NULL REFERENCES soma.custom_properties ( custom_property_id ) DEFERRABLE,
    inheritance_enabled         boolean         NOT NULL DEFAULT 'yes',
    children_only               boolean         NOT NULL DEFAULT 'no',
    value                       text            NOT NULL,
    FOREIGN KEY ( repository_id, custom_property_id ) REFERENCES soma.custom_properties ( repository_id, custom_property_id ) DEFERRABLE,
    FOREIGN KEY ( source_instance_id, repository_id ) REFERENCES soma.property_instances ( instance_id, repository_id ) DEFERRABLE
);`
	queries[idx] = "createTableRepositoryCustomProperty"

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
