package main

func createTablesInventoryAssets(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 15)

	queryMap["createTableDatacenters"] = `
create table if not exists inventory.datacenters (
    datacenter                  varchar(32)     PRIMARY KEY
);`
	queries[idx] = "createTableDatacenters"
	idx++

	queryMap["createTableServers"] = `
create table if not exists inventory.servers (
    server_id                   uuid            PRIMARY KEY,
    server_asset_id             numeric(16,0)   UNIQUE NOT NULL,
    server_datacenter_name      varchar(32)     NOT NULL REFERENCES inventory.datacenters ( datacenter ) DEFERRABLE,
    server_datacenter_location  varchar(256)    NOT NULL,
    server_name                 varchar(256)    NOT NULL,
    server_online               boolean         NOT NULL DEFAULT 'yes',
    server_deleted              boolean         NOT NULL DEFAULT 'no',
    CHECK( NOT (server_online AND server_deleted) )
);`
	queries[idx] = "createTableServers"
	idx++

	queryMap["createIndexUniqueServersOnline"] = `
create unique index _unique_server_online
    on inventory.servers ( server_name )
    where server_online
;`
	queries[idx] = "createIndexUniqueServersOnline"

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

func createTablesInventoryAccounts(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 20)

	queryMap[`create__inventory.dictionary`] = `
create table if not exists inventory.dictionary (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    name                        varchar(128)    NOT NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT now(),
    CONSTRAINT _dictionary_primary_key PRIMARY KEY( id ),
    CONSTRAINT _dictionary_unique_name UNIQUE ( name ) DEFERRABLE,
    CONSTRAINT _dictionary_timezone_utc CHECK( EXTRACT( TIMEZONE FROM created_at ) = '0' )
);`
	queries[idx] = `create__inventory.dictionary`
	idx++

	queryMap[`create__inventory.team`] = `
create table if not exists inventory.team (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionary_id               uuid            NOT NULL,
    name                        varchar(384)    NOT NULL,
    ldap_id                     numeric(16,0)   NOT NULL,
    is_system                   boolean         NOT NULL DEFAULT 'no'::boolean,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT now(),
    CONSTRAINT _team_primary_key PRIMARY KEY( id ),
    CONSTRAINT _team_from_dictionary_for_fk UNIQUE (dictionary_id, id),
    CONSTRAINT _team_unique_ldap_id_per_dictionary UNIQUE (dictionary_id, ldap_id),
    CONSTRAINT _team_unique_name UNIQUE (name),
    CONSTRAINT _team_timezone_utc CHECK( EXTRACT( TIMEZONE FROM created_at ) = '0' ),
    CONSTRAINT _team_dictionary_exists FOREIGN KEY (dictionary_id) REFERENCES inventory.dictionary (id) DEFERRABLE
);`
	queries[idx] = `create__inventory.team`
	idx++

	queryMap[`create__inventory.user`] = `
create table if not exists inventory.user (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionary_id               uuid            NOT NULL,
    uid                         varchar(256)    NOT NULL,
    first_name                  varchar(256)    NOT NULL,
    last_name                   varchar(256)    NOT NULL,
    employee_number             numeric(16,0)   NOT NULL,
    mail_address                text            NOT NULL,
    is_active                   boolean         NOT NULL DEFAULT 'yes',
    is_system                   boolean         NOT NULL DEFAULT 'no',
    is_deleted                  boolean         NOT NULL DEFAULT 'no',
    team_id                     uuid            NOT NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT now(),
    CONSTRAINT _user_primary_key PRIMARY KEY (id),
    CONSTRAINT _user_unique_employee_number_per_dictionary UNIQUE (dictionary_id, employee_number),
    CONSTRAINT _user_unique_name UNIQUE (uid),
    CONSTRAINT _user_timezone_utc CHECK( EXTRACT( TIMEZONE FROM created_at ) = '0' ),
    CONSTRAINT _user_dictionary_exists FOREIGN KEY (dictionary_id) REFERENCES inventory.dictionary (id) DEFERRABLE,
    CONSTRAINT _user_from_same_dictionary_as_team FOREIGN KEY (dictionary_id, team_id) REFERENCES inventory.team (dictionary_id, id) DEFERRABLE,
    CONSTRAINT _user_team_exists FOREIGN KEY (team_id) REFERENCES inventory.team (id) DEFERRABLE
);`
	queries[idx] = `create__inventory.user`
	idx++

	queryMap["createIndexUsersDeleted"] = `
create index _user_is_deleted on inventory.user ( is_deleted, id )
    where is_deleted
;`
	queries[idx] = "createIndexUsersDeleted"
	idx++

	queryMap["createIndexUsersSystem"] = `
create index _user_is_system on inventory.user ( is_system, id )
    where is_system
;`
	queries[idx] = "createIndexUsersSystem"
	idx++

	queryMap["createIndexUsersDeactivated"] = `
create index _user_is_inactive on inventory.user ( is_active, id )
    where is_active = 'no'::boolean
;`
	queries[idx] = "createIndexUsersDeactivated"
	idx++

	queryMap[`create__inventory.oncall_team`] = `
create table if not exists inventory.oncall_team (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    dictionary_id               uuid            NOT NULL,
    name                        varchar(256)    NOT NULL,
    phone_number                numeric(5,0)    NOT NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT now(),
    CONSTRAINT _oncall_id_primary_key PRIMARY KEY (id),
    CONSTRAINT _oncall_team_creator_exists FOREIGN KEY(created_by) REFERENCES inventory.user(id) DEFERRABLE,
    CONSTRAINT _oncall_team_dictionary_exists FOREIGN KEY (dictionary_id) REFERENCES inventory.dictionary (id) DEFERRABLE,
    CONSTRAINT _oncall_team_from_dictionary_for_fk UNIQUE (dictionary_id, id),
    CONSTRAINT _oncall_team_timezone_utc CHECK( EXTRACT( TIMEZONE FROM created_at ) = '0' ),
    CONSTRAINT _oncall_team_unique_name UNIQUE ( name ),
    CONSTRAINT _oncall_team_unique_phone_number UNIQUE ( phone_number )
);`
	queries[idx] = `create__inventory.oncall_team`
	idx++

	queryMap[`create__inventory.oncall_membership`] = `
create table if not exists inventory.oncall_membership (
    user_id                     uuid            NOT NULL,
    oncall_id                   uuid            NOT NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT now(),
    CONSTRAINT _oncall_membership_creator_exists FOREIGN KEY(created_by) REFERENCES inventory.user(id) DEFERRABLE,
    CONSTRAINT _oncall_membership_oncall_exists FOREIGN KEY (oncall_id) REFERENCES inventory.oncall_team (id) ON DELETE CASCADE DEFERRABLE,
    CONSTRAINT _oncall_membership_only_once UNIQUE (user_id, oncall_id),
    CONSTRAINT _oncall_membership_timezone_utc CHECK( EXTRACT( TIMEZONE FROM created_at ) = '0' ),
    CONSTRAINT _oncall_membership_user_exists FOREIGN KEY (user_id) REFERENCES inventory.user (id) ON DELETE CASCADE DEFERRABLE
);`
	queries[idx] = `create__inventory.oncall_membership`

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
