package main

func createTablesPermissions(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 25)

	queryMap[`create__soma.category`] = `
create table if not exists soma.category (
    name                        varchar(32)     NOT NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT _category_primary_key PRIMARY KEY (name),
    CONSTRAINT _category_creator_exists FOREIGN KEY (created_by) REFERENCES inventory.user (id) DEFERRABLE
);`
	queries[idx] = `create__soma.category`
	idx++

	queryMap[`create__soma.section`] = `
create table if not exists soma.section (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    name                        varchar(64)     NOT NULL,
    category                    varchar(32)     NOT NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT _section_primary_key PRIMARY KEY (id),
    CONSTRAINT _section_category_exists FOREIGN KEY (category) REFERENCES soma.category (name) DEFERRABLE,
    CONSTRAINT _section_creator_exists FOREIGN KEY (created_by) REFERENCES inventory.user (id) DEFERRABLE,
    CONSTRAINT _section_from_category_for_fk UNIQUE (id, category),
    CONSTRAINT _section_unique_name UNIQUE (name)
);`
	queries[idx] = `create__soma.section`
	idx++

	queryMap[`create__soma.action`] = `
create table if not exists soma.action (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    name                        varchar(64)     NOT NULL,
    section_id                  uuid            NOT NULL,
    category                    varchar(32)     NOT NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT _action_primary_key PRIMARY KEY (id),
    CONSTRAINT _action_creator_exists FOREIGN KEY (created_by) REFERENCES inventory.user (id) DEFERRABLE,
    CONSTRAINT _action_section_from_category FOREIGN KEY (section_id, category) REFERENCES soma.section (id, category) DEFERRABLE,
    CONSTRAINT _action_unique_per_section UNIQUE (section_id,name),
    CONSTRAINT _action_from_section_for_fk UNIQUE (section_id, id)
);`
	queries[idx] = `create__soma.action`
	idx++

	queryMap[`create__soma.permission`] = `
create table if not exists soma.permission (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    name                        varchar(128)    NOT NULL,
    category                    varchar(32)     NOT NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT _permission_primary_key PRIMARY KEY (id),
    CONSTRAINT _permission_unique_name_per_category UNIQUE (name, category),
    CONSTRAINT _permission_from_category_for_fk UNIQUE ( id, category ),
    CONSTRAINT _permission_creator_exists FOREIGN KEY (created_by) REFERENCES inventory.user (id) DEFERRABLE,
    CONSTRAINT _permission_category_exists FOREIGN KEY (category) REFERENCES soma.category (name) DEFERRABLE,
    CONSTRAINT _permission_timezone_utc CHECK( EXTRACT( TIMEZONE FROM created_at ) = '0' ),
    -- only omnipotence is category omnipotence
    CONSTRAINT _permission_validate_omnipotence CHECK (category != 'omnipotence' OR name = 'omnipotence')
);`
	queries[idx] = `create__soma.permission`
	idx++

	queryMap[`create__soma.permission_map`] = `
create table if not exists soma.permission_map (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    category                    varchar(32)     NOT NULL,
    permission_id               uuid            NOT NULL,
    section_id                  uuid            NOT NULL,
    action_id                   uuid            NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    CONSTRAINT _permission_map_primary_key PRIMARY KEY (id),
    CONSTRAINT _permission_map_action_in_section FOREIGN KEY (section_id, action_id) REFERENCES soma.action(section_id, id) DEFERRABLE,
    CONSTRAINT _permission_map_creator_exists FOREIGN KEY (created_by) REFERENCES inventory.user (id) DEFERRABLE,
    CONSTRAINT _permission_map_permission_in_category FOREIGN KEY ( permission_id, category ) REFERENCES soma.permission (id, category) DEFERRABLE,
    CONSTRAINT _permission_map_section_in_category FOREIGN KEY (section_id, category) REFERENCES soma.section(id, category) DEFERRABLE,
    CONSTRAINT _permission_map_timezone_utc CHECK( EXTRACT( TIMEZONE FROM created_at ) = '0' )
);`
	queries[idx] = `create__soma.permission_map`
	idx++

	queryMap[`create__soma.permission_grant_map`] = `
create table if not exists soma.permission_grant_map (
    category                    varchar(32)     NOT NULL,
    permission_id               uuid            NOT NULL,
    granted_category            varchar(32)     NOT NULL,
    granted_permission_id       uuid            NOT NULL,
    CONSTRAINT _permission_grant_map_check_grant_category_correlation CHECK ( granted_category = substring(category from '^([^:]+):')),
    CONSTRAINT _permission_grant_map_check_is_grant_category CHECK ( category ~ ':grant$' ),
    CONSTRAINT _permission_grant_map_check_no_self_grant CHECK ( permission_id != granted_permission_id ),
    CONSTRAINT _permission_grant_map_granted_only_once UNIQUE (granted_permission_id),
    CONSTRAINT _permission_grant_map_granted_permission_exists FOREIGN KEY (granted_permission_id, granted_category) REFERENCES soma.permission (id, category) DEFERRABLE,
    CONSTRAINT _permission_grant_map_granting_only_once UNIQUE (permission_id),
    CONSTRAINT _permission_grant_map_granting_permission_exists FOREIGN KEY (permission_id, category) REFERENCES soma.permission (id, category) DEFERRABLE
);`
	queries[idx] = `create__soma.permission_grant_map`
	idx++

	queryMap["createTableGlobalAuthorizations"] = `
create table if not exists soma.authorizations_global (
    grant_id                    uuid            PRIMARY KEY,
    admin_id                    uuid            REFERENCES auth.admin ( id ) DEFERRABLE,
    user_id                     uuid            REFERENCES inventory.user ( id ) DEFERRABLE,
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ) DEFERRABLE,
    team_id                     uuid            REFERENCES inventory.team ( id ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permission (id) DEFERRABLE,
    category                    varchar(32)     NOT NULL REFERENCES soma.category (name) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.user ( id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, category ) REFERENCES soma.permission ( id, category ) DEFERRABLE,
    CHECK (   ( admin_id IS NOT NULL AND user_id IS     NULL AND tool_id IS     NULL AND team_id IS     NULL )
           OR ( admin_id IS     NULL AND user_id IS NOT NULL AND tool_id IS     NULL AND team_id IS     NULL )
           OR ( admin_id IS     NULL AND user_id IS     NULL AND tool_id IS NOT NULL AND team_id IS     NULL )
           OR ( admin_id IS     NULL AND user_id IS     NULL AND tool_id IS     NULL AND team_id IS NOT NULL ) ),
    CONSTRAINT authorizations_global_category_check CHECK ( category IN ( 'omnipotence','system','global','global:grant','permission','permission:grant','operations','operations:grant','self','self:grant','identity','identity:grant' ) ),
    -- if system, then it has to be an admin account
    CHECK ( category != 'system' OR admin_id IS NOT NULL ),
    -- admins can only have system
    CHECK ( admin_id IS NULL OR category = 'system' ),
    -- only root can have omnipotence
    CHECK ( permission_id != '00000000-0000-0000-0000-000000000000'::uuid OR user_id = '00000000-0000-0000-0000-000000000000'::uuid ),
    UNIQUE( admin_id, user_id, tool_id, team_id, category, permission_id )
);`
	queries[idx] = "createTableGlobalAuthorizations"
	idx++

	queryMap[`createUniqueIndexAdminGlobalAuthorization`] = `
create unique index _unique_admin_global_authoriz
    on soma.authorizations_global ( admin_id, permission_id )
    where admin_id IS NOT NULL;`
	queries[idx] = `createUniqueIndexAdminGlobalAuthorization`
	idx++

	queryMap[`createUniqueIndexUserGlobalAuthorization`] = `
create unique index _unique_user_global_authoriz
    on soma.authorizations_global ( user_id, permission_id )
    where user_id IS NOT NULL;`
	queries[idx] = `createUniqueIndexUserGlobalAuthorization`
	idx++

	queryMap[`createUniqueIndexToolGlobalAuthorization`] = `
create unique index _unique_tool_global_authoriz
    on soma.authorizations_global ( tool_id, permission_id )
    where tool_id IS NOT NULL;`
	queries[idx] = `createUniqueIndexToolGlobalAuthorization`
	idx++

	queryMap["createTableRepoAuthorizations"] = `
create table if not exists soma.authorizations_repository (
    grant_id                    uuid            PRIMARY KEY,
    user_id                     uuid            REFERENCES inventory.user ( id ) DEFERRABLE,
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ) DEFERRABLE,
    admin_id                    uuid            REFERENCES auth.admin ( id ) DEFERRABLE,
    team_id                     uuid            REFERENCES inventory.team ( id ) DEFERRABLE,
    object_type                 varchar(64)     NOT NULL REFERENCES soma.object_types ( object_type ) DEFERRABLE,
    repository_id               uuid            REFERENCES soma.repository (id) DEFERRABLE,
    bucket_id                   uuid            REFERENCES soma.buckets ( bucket_id ) DEFERRABLE,
    group_id                    uuid            REFERENCES soma.groups ( group_id ) DEFERRABLE,
    cluster_id                  uuid            REFERENCES soma.clusters ( cluster_id ) DEFERRABLE,
    node_id                     uuid            REFERENCES soma.nodes ( node_id ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permission (id) DEFERRABLE,
    category                    varchar(32)     NOT NULL REFERENCES soma.category (name) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.user ( id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, category ) REFERENCES soma.permission (id, category) DEFERRABLE,
    FOREIGN KEY ( bucket_id, repository_id ) REFERENCES soma.buckets ( bucket_id, repository_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, group_id ) REFERENCES soma.groups ( bucket_id, group_id ) DEFERRABLE,
    FOREIGN KEY ( bucket_id, cluster_id ) REFERENCES soma.clusters ( bucket_id, cluster_id ) DEFERRABLE,
    FOREIGN KEY ( node_id, bucket_id ) REFERENCES soma.node_bucket_assignment ( node_id, bucket_id ) DEFERRABLE,
    CONSTRAINT check_single_subject_id_set
    CHECK ( ( user_id IS NOT NULL AND tool_id IS     NULL AND team_id IS     NULL AND admin_id IS     NULL )
         OR ( user_id IS     NULL AND tool_id IS NOT NULL AND team_id IS     NULL AND admin_id IS     NULL )
         OR ( user_id IS     NULL AND tool_id IS     NULL AND team_id IS NOT NULL AND admin_id IS     NULL )
         OR ( user_id IS     NULL AND tool_id IS     NULL AND team_id IS     NULL AND admin_id IS NOT NULL ) ),
    CONSTRAINT check_category CHECK ( category IN ( 'repository', 'repository:grant' ) ),
    CONSTRAINT check_object_types CHECK ( object_type IN ( 'repository', 'bucket', 'group', 'cluster', 'node' )),
    CONSTRAINT check_type_repository_id CHECK ( object_type != 'repository' OR repository_id IS NOT NULL ),
    CONSTRAINT check_type_bucket_id CHECK ( object_type != 'bucket' OR bucket_id IS NOT NULL ),
    CONSTRAINT check_type_group_id CHECK ( object_type != 'group' OR group_id IS NOT NULL ),
    CONSTRAINT check_type_cluster_id CHECK ( object_type != 'cluster' OR cluster_id IS NOT NULL ),
    CONSTRAINT check_type_node_id CHECK ( object_type != 'node' OR node_id IS NOT NULL ),
    CONSTRAINT check_single_object_id_set
    CHECK ( ( repository_id IS NOT NULL AND bucket_id IS     NULL AND group_id IS     NULL AND cluster_id IS     NULL AND node_id IS     NULL )
         OR ( repository_id IS NOT NULL AND bucket_id IS NOT NULL AND group_id IS     NULL AND cluster_id IS     NULL AND node_id IS     NULL )
         OR ( repository_id IS NOT NULL AND bucket_id IS NOT NULL AND group_id IS NOT NULL AND cluster_id IS     NULL AND node_id IS     NULL )
         OR ( repository_id IS NOT NULL AND bucket_id IS NOT NULL AND group_id IS     NULL AND cluster_id IS NOT NULL AND node_id IS     NULL )
         OR ( repository_id IS NOT NULL AND bucket_id IS NOT NULL AND group_id IS     NULL AND cluster_id IS     NULL AND node_id IS NOT NULL )),
    UNIQUE ( user_id, tool_id, team_id, category, permission_id, object_type, repository_id, bucket_id, group_id, cluster_id, node_id )
);`
	queries[idx] = "createTableRepoAuthorizations"
	idx++

	queryMap["createTableMonitoringAuthorizations"] = `
create table if not exists soma.authorizations_monitoring (
    grant_id                    uuid            PRIMARY KEY,
    user_id                     uuid            REFERENCES inventory.user ( id ) DEFERRABLE,
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ) DEFERRABLE,
    team_id                     uuid            REFERENCES inventory.team ( id ) DEFERRABLE,
    monitoring_id               uuid            NOT NULL REFERENCES soma.monitoring_systems ( monitoring_id ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permission (id) DEFERRABLE,
    category                    varchar(32)     NOT NULL REFERENCES soma.category (name) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.user ( id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, category ) REFERENCES soma.permission (id, category) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND team_id IS NOT NULL ) ),
    CHECK ( category IN ( 'monitoring', 'monitoring:grant' ) ),
    UNIQUE ( user_id, tool_id, team_id, category, permission_id, monitoring_id )
);`
	queries[idx] = "createTableMonitoringAuthorizations"
	idx++

	queryMap["createTableTeamAuthorizations"] = `
create table if not exists soma.authorizations_team (
    grant_id                    uuid            PRIMARY KEY,
    user_id                     uuid            REFERENCES inventory.user ( id ) DEFERRABLE,
    tool_id                     uuid            REFERENCES auth.tools ( tool_id ) DEFERRABLE,
    team_id                     uuid            REFERENCES inventory.team ( id ) DEFERRABLE,
    authorized_team_id          uuid            NOT NULL REFERENCES inventory.team ( id ) DEFERRABLE,
    permission_id               uuid            NOT NULL REFERENCES soma.permission (id) DEFERRABLE,
    category                    varchar(32)     NOT NULL REFERENCES soma.category (name) DEFERRABLE,
    created_by                  uuid            NOT NULL REFERENCES inventory.user ( id ) DEFERRABLE,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW(),
    FOREIGN KEY ( permission_id, category ) REFERENCES soma.permission (id, category) DEFERRABLE,
    CHECK (   ( user_id IS NOT NULL AND tool_id IS     NULL AND team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS NOT NULL AND team_id IS     NULL )
           OR ( user_id IS     NULL AND tool_id IS     NULL AND team_id IS NOT NULL ) ),
    CHECK ( category IN ( 'team', 'team:grant' ) ),
    UNIQUE ( user_id, tool_id, team_id, category, permission_id, authorized_team_id )
);`
	queries[idx] = "createTableTeamAuthorizations"

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
