package main

func createTablesProperties(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 10)

	queryMap["createTableTemplateProperties"] = `
create table if not exists soma.template_property (
	id                          uuid            PRIMARY KEY,
    name                        varchar(128)    NOT NULL
);`
	queries[idx] = "createTableTemplateProperties"
	idx++

	queryMap["createTableAttribute"] = `
create table if not exists soma.attribute (
    attribute                   varchar(128)    PRIMARY KEY,
    cardinality                 varchar(8)      NOT NULL DEFAULT 'multi'
);`
	queries[idx] = "createTableAttribute"
	idx++

	queryMap["createTableTemplatePropertyValues"] = `
create table if not exists soma.template_property_value (
    template_id                 varchar(128)    NOT NULL REFERENCES soma.template_property ( id ) DEFERRABLE,
    attribute                   varchar(128)    NOT NULL REFERENCES soma.attribute ( attribute ) DEFERRABLE,
    value                       varchar(512)    NOT NULL
);`
	queries[idx] = "createTableTemplatePropertyValues"
	idx++

	queryMap["createTableServiceProperty"] = `
create table if not exists soma.service_property (
    id                          uuid            PRIMARY KEY,
    team_id                     uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    name                        varchar(128)    NOT NULL,
    UNIQUE( team_id, id ),
    UNIQUE( name, team_id )
);`
	queries[idx] = "createTableTeamServiceProperty"
	idx++

	queryMap["createTableTeamServicePropertyValues"] = `
create table if not exists soma.service_property_value (
	service_id                  uuid            NOT NULL REFERENCES soma.service_property ( id ) DEFERRABLE,
    attribute                   varchar(128)    NOT NULL REFERENCES soma.attribute ( attribute ) DEFERRABLE,
    team_id                     uuid            NOT NULL REFERENCES inventory.organizational_teams ( organizational_team_id ) DEFERRABLE,
    value                       varchar(512)    NOT NULL,
    FOREIGN KEY( organizational_team_id, service_id ) REFERENCES soma.service_property ( organizational_team_id, id ) DEFERRABLE
);`
	queries[idx] = "createTableTeamServicePropertyValues"
	idx++

	queryMap["createTableSystemProperties"] = `
create table if not exists soma.system_properties (
    system_property             varchar(128)    PRIMARY KEY
);`
	queries[idx] = "createTableSystemProperties"
	idx++

	queryMap["createTableSystemPropertyValidity"] = `
create table if not exists soma.system_property_validity (
    system_property             varchar(128)    NOT NULL REFERENCES soma.system_properties ( system_property ) DEFERRABLE,
    object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    inherited                   boolean         NOT NULL DEFAULT 'yes',
    UNIQUE( system_property, object_type, inherited )
);`
	queries[idx] = "createTableSystemPropertyValidity"
	idx++

	queryMap["createTableNativeProperties"] = `
create table if not exists soma.native_properties (
    native_property             varchar(128)    PRIMARY KEY
);`
	queries[idx] = "createTableNativeProperties"

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

func createTableCustomProperties(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 5)

	queryMap["createTableCustomProperties"] = `
create table if not exists soma.custom_properties (
    custom_property_id          uuid            PRIMARY KEY,
    repository_id               uuid            NOT NULL REFERENCES soma.repositories ( repository_id ) DEFERRABLE,
    custom_property             varchar(128)    NOT NULL,
    UNIQUE( repository_id, custom_property ),
    UNIQUE( repository_id, custom_property_id )
);`
	queries[idx] = "createTableCustomProperties"

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
