package main

import "fmt"

func schemaInserts(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 100)

	queryMap["insertSystemGroupWheel"] = `
INSERT INTO inventory.organizational_teams (
            organizational_team_id,
            organizational_team_name,
            organizational_team_ldap_id,
            organizational_team_system
) VALUES (
            '00000000-0000-0000-0000-000000000000',
            'wheel',
            0,
            'yes'
);`
	queries[idx] = "insertSystemGroupWheel"
	idx++

	queryMap["insertSystemUserRootAndFriends"] = `
INSERT INTO inventory.users (
            user_id,
            user_uid,
            user_first_name,
            user_last_name,
            user_employee_number,
            user_mail_address,
            user_is_active,
            user_is_system,
            user_is_deleted,
            organizational_team_id
) VALUES (
            '00000000-0000-0000-0000-000000000000',
            'root',
            'Charlie',
            'Root',
            0,
            'devnull@example.com',
            'yes',
            'yes',
            'no',
            '00000000-0000-0000-0000-000000000000'
),
(
            'ffffffff-ffff-ffff-ffff-ffffffffffff',
            'AnonymousCoward',
            'Anonymous',
            'Coward',
            9999999999999999,
            'devzero@example.com',
            'yes',
            'yes',
            'no',
            '00000000-0000-0000-0000-000000000000'
);`
	queries[idx] = "insertSystemUserRootAndFriends"
	idx++

	queryMap["insertCategoryOmnipotence"] = `
INSERT INTO soma.categories (
            category,
            created_by
) VALUES (
            'omnipotence',
            '00000000-0000-0000-0000-000000000000'
),
(
            'system',
            '00000000-0000-0000-0000-000000000000'
);`
	queries[idx] = "insertCategoryOmnipotence"
	idx++

	queryMap["insertPermissionOmnipotence"] = `
INSERT INTO soma.permissions (
            permission_id,
            permission_name,
            category,
            created_by
) VALUES (
            '00000000-0000-0000-0000-000000000000',
            'omnipotence',
            'omnipotence',
            '00000000-0000-0000-0000-000000000000'
);`
	queries[idx] = "insertPermissionOmnipotence"
	idx++

	queryMap["grantOmnipotence"] = `
INSERT INTO soma.authorizations_global (
            grant_id,
            user_id,
            permission_id,
            category,
            created_by
) VALUES (
            '00000000-0000-0000-0000-000000000000',
            '00000000-0000-0000-0000-000000000000',
            '00000000-0000-0000-0000-000000000000',
            'omnipotence',
            '00000000-0000-0000-0000-000000000000'
);`
	queries[idx] = "grantOmnipotence"
	idx++

	queryMap["insertJobStatus"] = `
INSERT INTO soma.job_status (
            job_status
) VALUES
            ( 'queued' ),
            ( 'in_progress' ),
            ( 'processed' )
;`
	queries[idx] = "insertJobStatus"
	idx++

	queryMap["insertJobResults"] = `
INSERT INTO soma.job_results (
            job_result
) VALUES
            ( 'pending' ),
            ( 'success' ),
            ( 'failed' )
;`
	queries[idx] = "insertJobResults"
	idx++

	queryMap["insertJobTypes"] = `
INSERT INTO soma.job_types (
            job_type
) VALUES
            ( 'bucket::property-create' ),
            ( 'bucket::property-destroy' ),
            ( 'bucket::property-update' ),
            ( 'bucket::create' ),
            ( 'check-config::create' ),
            ( 'check-config::destroy' ),
            ( 'group::create' ),
            ( 'group::destroy' ),
            ( 'group::member-assign' ),
            ( 'group::property-create' ),
            ( 'group::property-destroy' ),
            ( 'group::property-update' ),
            ( 'cluster::create' ),
            ( 'cluster::destroy' ),
            ( 'cluster::property-create' ),
            ( 'cluster::property-destroy' ),
            ( 'cluster::property-update' ),
            ( 'cluster::member-assign' ),
            ( 'node-config::assign' ),
            ( 'node-config::unassign' ),
            ( 'node-config::property-create' ),
            ( 'node-config::property-destroy' ),
            ( 'node-config::property-update' ),
            ( 'repository-config::property-create' ),
            ( 'repository-config::property-destroy' ),
            ( 'repository-config::property-update' ),
            ( 'repository::rename' ),
            ( 'repository::destroy' ),
            ( 'bucket::rename' ),
            ( 'bucket::destroy' )
;`
	queries[idx] = "insertJobTypes"
	idx++

	queryMap["insertRootRestricted"] = `
INSERT INTO root.flags (
            flag,
            status
) VALUES
            ( 'restricted', false ),
            ( 'disabled', false )
;`
	queries[idx] = "insertRootRestricted"
	idx++

	performDatabaseTask(printOnly, verbose, queries[:idx], queryMap)
}

func schemaVersionInserts(printOnly bool, verbose bool, version string) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 100)

	invString := fmt.Sprintf(`
INSERT INTO public.schema_versions (
            schema,
            version,
            description )
VALUES (
            'inventory',
            201605060001,
            'Initial create - somadbctl %s'
);`, version)
	queryMap["insertInventorySchemaVersion"] = invString
	queries[idx] = "insertInventorySchemaVersion"
	idx++

	authString := fmt.Sprintf(`
INSERT INTO public.schema_versions (
            schema,
            version,
            description
) VALUES (
            'auth',
            201711080001,
            'Initial create - somadbctl %s'
);`, version)
	queryMap["insertAuthSchemaVersion"] = authString
	queries[idx] = "insertAuthSchemaVersion"
	idx++

	somaString := fmt.Sprintf(`
INSERT INTO public.schema_versions (
            schema,
            version,
            description
) VALUES (
            'soma',
            201811120001,
            'Initial create - somadbctl %s'
);`, version)
	queryMap["insertSomaSchemaVersion"] = somaString
	queries[idx] = "insertSomaSchemaVersion"
	idx++

	rootString := fmt.Sprintf(`
INSERT INTO public.schema_versions (
            schema,
            version,
            description
) VALUES (
            'root',
            201605160001,
            'Initial create - somadbctl %s'
);`, version)
	queryMap["insertRootSchemaVersion"] = rootString
	queries[idx] = "insertRootSchemaVersion"
	idx++

	performDatabaseTask(printOnly, verbose, queries[:idx], queryMap)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
