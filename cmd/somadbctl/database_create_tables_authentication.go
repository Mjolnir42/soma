package main

func createTablesAuthentication(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 25)

	queryMap["createTableUserAuthentication"] = `
create table if not exists auth.user_authentication (
    user_id                     uuid            NOT NULL REFERENCES inventory.user ( id ) ON DELETE CASCADE DEFERRABLE,
    crypt                       text            NOT NULL,
    reset_pending               boolean         NOT NULL DEFAULT 'no',
    valid_from                  timestamptz(3)  NOT NULL,
    valid_until                 timestamptz(3)  NOT NULL,
    CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
    CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
);`
	queries[idx] = "createTableUserAuthentication"
	idx++

	queryMap["createTableTokenAuthentication"] = `
create table if not exists auth.tokens (
    token                       varchar(256)    UNIQUE NOT NULL,
    salt                        varchar(256)    NOT NULL,
    valid_from                  timestamptz(3)  NOT NULL,
    valid_until                 timestamptz(3)  NOT NULL,
    CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
    CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
);`
	queries[idx] = "createTableTokenAuthentication"
	idx++

	queryMap[`createTableTokenRevocation`] = `
create table if not exists auth.token_revocations (
    user_id                     uuid            NOT NULL REFERENCES inventory.user ( id ) ON DELETE CASCADE DEFERRABLE,
    revoked_at                  timestamptz(3)  NOT NULL,
    CHECK( EXTRACT( TIMEZONE FROM revoked_at ) = '0' )
);`
	queries[idx] = `createTableTokenRevocation`
	idx++

	queryMap[`createIndexTokenRevocationUserID`] = `
create index _token_revocations_user_id
    on auth.token_revocations ( user_id )
;`
	queries[idx] = `createIndexTokenRevocationUserID`
	idx++

	queryMap[`createIndexTokenRevocationRevokedAt`] = `
create index _token_revocations_revoked_at
    on auth.token_revocations ( revoked_at DESC )
;`
	queries[idx] = `createIndexTokenRevocationRevokedAt`
	idx++

	queryMap["createTableUserKeys"] = `
create table if not exists auth.user_keys (
    user_id                     uuid            NOT NULL REFERENCES inventory.user ( id ) ON DELETE CASCADE DEFERRABLE,
    user_key_fingerprint        varchar(128)    NOT NULL,
    user_key_public             text            NOT NULL,
    user_key_active             boolean         NOT NULL DEFAULT 'yes'
);`
	queries[idx] = "createTableUserKeys"
	idx++

	queryMap["createIndexUniqueActiveUserKey"] = `
create unique index _unique_active_user_key
    on auth.user_keys ( user_id, user_key_active )
    where user_key_active
;`
	queries[idx] = "createIndexUniqueActiveUserKey"
	idx++

	queryMap["createTableUserClientCertificates"] = `
create table if not exists auth.user_client_certificates (
    user_id                     uuid            NOT NULL REFERENCES inventory.user ( id ) ON DELETE CASCADE DEFERRABLE,
    user_cert_fingerprint       varchar(128)    NOT NULL,
    user_cert_active            boolean         NOT NULL DEFAULT 'yes'
);`
	queries[idx] = "createTableUserClientCertificates"
	idx++

	queryMap["createIndexUniqueActiveUserCert"] = `
create unique index _unique_active_user_cert
    on auth.user_client_certificates ( user_id, user_cert_active )
    where user_cert_active
;`
	queries[idx] = "createIndexUniqueActiveUserCert"
	idx++

	queryMap[`create__auth.admin`] = `
create table if not exists auth.admin (
    id                          uuid            NOT NULL DEFAULT public.gen_random_uuid(),
    uid                         varchar(256)    NOT NULL,
    user_uid                    varchar(256)    NOT NULL,
    is_active                   boolean         NOT NULL DEFAULT 'yes',
    CONSTRAINT _admin_primary_key PRIMARY KEY (id),
    CONSTRAINT _admin_unique_name UNIQUE (uid),
    CONSTRAINT _admin_check_uid_prefix CHECK( left( uid, 6 ) = 'admin_' ),
    CONSTRAINT _admin_uid_contains_user_uid CHECK( position( user_uid in uid ) != 0 ),
    CONSTRAINT _admin_user_exists FOREIGN KEY (user_uid) REFERENCES inventory.user (uid) ON DELETE CASCADE DEFERRABLE,
);`
	queries[idx] = `create__auth.admin`
	idx++

	queryMap["createTableAdminAuthentication"] = `
create table if not exists auth.admin_authentication (
    admin_id                    uuid            NOT NULL REFERENCES auth.admin (id) ON DELETE CASCADE DEFERRABLE,
    crypt                       text            NOT NULL,
    reset_pending               boolean         NOT NULL DEFAULT 'no',
    valid_from                  timestamptz(3)  NOT NULL,
    valid_until                 timestamptz(3)  NOT NULL,
    CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
    CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
);`
	queries[idx] = "createTableAdminAuthentication"
	idx++

	queryMap["createTableAdminKeys"] = `
create table if not exists auth.admin_keys (
    admin_id                    uuid            NOT NULL REFERENCES auth.admin (id) ON DELETE CASCADE DEFERRABLE,
    admin_key_fingerprint       varchar(128)    NOT NULL,
    admin_key_public            text            NOT NULL,
    admin_key_active            boolean         NOT NULL DEFAULT 'yes'
);`
	queries[idx] = "createTableAdminKeys"
	idx++

	queryMap["createIndexUniqueActiveAdminKey"] = `
create index _unique_active_admin_key
    on auth.admin_keys ( admin_id, admin_key_active )
    where admin_key_active
;`
	queries[idx] = "createIndexUniqueActiveAdminKey"
	idx++

	queryMap["createTableAdminClientCertificates"] = `
create table if not exists auth.admin_client_certificates (
    admin_id                    uuid            NOT NULL REFERENCES auth.admin (id) ON DELETE CASCADE DEFERRABLE,
    admin_cert_fingerprint      varchar(128)    NOT NULL,
    admin_cert_active           boolean         NOT NULL DEFAULT 'yes'
);`
	queries[idx] = "createTableAdminClientCertificates"
	idx++

	queryMap["createIndexUniqueActiveAdminCert"] = `
create unique index _unique_active_admin_cert
    on auth.admin_client_certificates ( admin_id, admin_cert_active )
    where admin_cert_active
;`
	queries[idx] = "createIndexUniqueActiveAdminCert"
	idx++

	queryMap["createTableTools"] = `
create table if not exists auth.tools (
    tool_id                     uuid            PRIMARY KEY,
    tool_name                   varchar(256)    UNIQUE NOT NULL,
    tool_owner_id               uuid            NOT NULL REFERENCES inventory.user ( id ) ON DELETE RESTRICT DEFERRABLE,
    created                     timestamptz(3)  NOT NULL DEFAULT NOW()::timestamptz(3),
    CHECK( EXTRACT( TIMEZONE FROM created ) = '0' ),
    CHECK( left( tool_name, 5 ) = 'tool_' )
);`
	queries[idx] = "createTableTools"
	idx++

	queryMap["createTableToolAuthentication"] = `
create table if not exists auth.tool_authentication (
    tool_id                     uuid            NOT NULL REFERENCES auth.tools ( tool_id ) ON DELETE CASCADE DEFERRABLE,
    crypt                       text            NOT NULL,
    reset_pending               boolean         NOT NULL DEFAULT 'no',
    valid_from                  timestamptz(3)  NOT NULL,
    valid_until                 timestamptz(3)  NOT NULL,
    CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
    CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' )
);`
	queries[idx] = "createTableToolAuthentication"
	idx++

	queryMap["createTableToolKeys"] = `
create table if not exists auth.tool_keys (
    tool_id                     uuid            NOT NULL REFERENCES auth.tools ( tool_id ) ON DELETE CASCADE DEFERRABLE,
    tool_key_fingerprint        varchar(128)    NOT NULL,
    tool_key_public             text            NOT NULL,
    tool_key_active             boolean         NOT NULL DEFAULT 'yes'
);`
	queries[idx] = "createTableToolKeys"
	idx++

	queryMap["createIndexUniqueActiveToolKey"] = `
create unique index _unique_active_tool_key
    on auth.tool_keys ( tool_id, tool_key_active )
    where tool_key_active
;`
	queries[idx] = "createIndexUniqueActiveToolKey"
	idx++

	queryMap["createTableToolClientCertificates"] = `
create table if not exists auth.tool_client_certificates (
    tool_id                     uuid            NOT NULL REFERENCES auth.tools ( tool_id ) ON DELETE CASCADE DEFERRABLE,
    tool_cert_fingerprint       varchar(128)    NOT NULL,
    tool_cert_active            boolean         NOT NULL DEFAULT 'yes'
);`
	queries[idx] = "createTableToolClientCertificates"
	idx++

	queryMap["createIndexUniqueActiveToolCert"] = `
create unique index _unique_active_tool_cert
    on auth.tool_client_certificates ( tool_id, tool_cert_active )
    where tool_cert_active
;`
	queries[idx] = "createIndexUniqueActiveToolCert"
	idx++

	queryMap["createTablePasswordReset"] = `
create table if not exists auth.password_reset (
    user_id                     uuid            NULL REFERENCES inventory.user ( id ) ON DELETE CASCADE DEFERRABLE,
    admin_id                    uuid            NULL REFERENCES auth.admin (id) ON DELETE CASCADE DEFERRABLE,
    tool_id                     uuid            NULL REFERENCES auth.tools ( tool_id ) ON DELETE CASCADE DEFERRABLE,
    token                       varchar(256)    UNIQUE NOT NULL,
    valid_from                  timestamptz(3)  NOT NULL,
    valid_until                 timestamptz(3)  NOT NULL,
    token_used                  boolean         NOT NULL DEFAULT 'no',
    token_invalidated           boolean         NOT NULL DEFAULT 'no',
    CHECK( EXTRACT( TIMEZONE FROM valid_from )  = '0' ),
    CHECK( EXTRACT( TIMEZONE FROM valid_until ) = '0' ),
    CHECK(    ( user_id IS NOT NULL AND admin_id IS     NULL AND tool_id IS     NULL )
           OR ( user_id IS     NULL AND admin_id IS NOT NULL AND tool_id IS     NULL )
           OR ( user_id IS     NULL AND admin_id IS     NULL AND tool_id IS NOT NULL ) )
);`
	queries[idx] = "createTablePasswordReset"

	performDatabaseTask(printOnly, verbose, queries, queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
