package main

func createTablesJobs(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 10)

	queryMap[`createTableJobStatus`] = `
create table if not exists soma.job_status (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    name                        varchar(32)     NOT NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW()::timestamptz(3),
    CONSTRAINT _job_status_primary_key          PRIMARY KEY ( id ),
    CONSTRAINT _job_status_unique_name          UNIQUE ( name ),
    CONSTRAINT _job_status_user_exists          FOREIGN KEY ( created_by ) REFERENCES inventory.user ( id ) DEFERRABLE,
    CONSTRAINT _job_status_timezone_utc         CHECK( EXTRACT( TIMEZONE FROM created_at )  = '0' )
);`
	queries[idx] = `createTableJobStatus`
	idx++

	queryMap[`createTableJobResult`] = `
create table if not exists soma.job_result (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    name                        varchar(32)     NOT NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW()::timestamptz(3),
    CONSTRAINT _job_result_primary_key          PRIMARY KEY (id),
    CONSTRAINT _job_result_unique_name          UNIQUE (name),
    CONSTRAINT _job_result_user_exists          FOREIGN KEY ( created_by ) REFERENCES inventory.user ( id ) DEFERRABLE,
    CONSTRAINT _job_result_timezone_utc         CHECK( EXTRACT( TIMEZONE FROM created_at )  = '0' )
);`
	queries[idx] = `createTableJobResult`
	idx++

	queryMap[`createTableJobType`] = `
create table if not exists soma.job_type (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    name                        varchar(128)    NOT NULL,
    created_by                  uuid            NOT NULL,
    created_at                  timestamptz(3)  NOT NULL DEFAULT NOW()::timestamptz(3),
    CONSTRAINT _job_type_primary_key            PRIMARY KEY (id),
    CONSTRAINT _job_type_unique_name            UNIQUE (name),
    CONSTRAINT _job_type_user_exists            FOREIGN KEY ( created_by ) REFERENCES inventory.user ( id ) DEFERRABLE,
    CONSTRAINT _job_type_timezone_utc           CHECK( EXTRACT( TIMEZONE FROM created_at )  = '0' )
);`
	queries[idx] = `createTableJobType`
	idx++

	queryMap[`createTableJob`] = `
create table if not exists soma.job (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    status                      varchar(32)     NOT NULL,
    result                      varchar(32)     NOT NULL,
    type                        varchar(128)    NOT NULL,
    serial                      bigserial       NOT NULL,
    repository_id               uuid            NOT NULL,
    user_id                     uuid            NOT NULL,
    team_id                     uuid            NOT NULL,
    error                       text            NOT NULL DEFAULT '',
    queued_at                   timestamptz(3)  NOT NULL DEFAULT NOW()::timestamptz(3),
    started_at                  timestamptz(3),
    finished_at                 timestamptz(3),
    job                         jsonb           NOT NULL,
    CONSTRAINT _job_primary_key                 PRIMARY KEY (id),
    CONSTRAINT _job_status_exists               FOREIGN KEY ( status ) REFERENCES soma.job_status ( name ) DEFERRABLE,
    CONSTRAINT _job_result_exists               FOREIGN KEY ( result ) REFERENCES soma.job_result ( name ) DEFERRABLE,
    CONSTRAINT _job_type_exists                 FOREIGN KEY ( type ) REFERENCES soma.job_type ( name ) DEFERRABLE,
    CONSTRAINT _job_repository_exists           FOREIGN KEY ( repository_id ) REFERENCES soma.repository (id) DEFERRABLE,
    CONSTRAINT _job_user_exists                 FOREIGN KEY ( user_id ) REFERENCES inventory.user ( id ) DEFERRABLE,
    CONSTRAINT _job_team_exists                 FOREIGN KEY ( team_id ) REFERENCES inventory.team ( id ) DEFERRABLE
);`
	queries[idx] = `createTableJob`
	idx++

	queryMap[`createIndexJobStatus`] = `
create index _job_not_processed
    on soma.job ( team_id, user_id, status, id )
    where status != 'processed'
;`
	queries[idx] = `createIndexJobStatus`
	idx++

	queryMap[`createIndexRepoJob`] = `
create index _job_by_repository
    on soma.job ( repository_id, serial, id, status )
;`
	queries[idx] = `createIndexRepoJob`

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
