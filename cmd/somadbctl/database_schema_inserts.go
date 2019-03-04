package main

import "fmt"

func schemaInserts(printOnly bool, verbose bool) {
	idx := 0
	// map for storing the SQL statements by name
	queryMap := make(map[string]string)
	// slice storing the required statement order so foreign keys can
	// resolve successfully
	queries := make([]string, 100)

	queryMap[`insert__inventory.dictionary=system`] = `
INSERT INTO inventory.dictionary (
            id,
            name,
            created_by
) VALUES (
            '00000000-0000-0000-0000-000000000000'::uuid,
            'system'::varchar,
            '00000000-0000-0000-0000-000000000000'::uuid
);
`
	queries[idx] = `insert__inventory.dictionary=system`
	idx++

	queryMap["insertSystemGroupWheel"] = `
INSERT INTO inventory.team (
            id,
            name,
            ldap_id,
            is_system,
            dictionary_id,
            created_by
) VALUES (
            '00000000-0000-0000-0000-000000000000',
            'wheel',
            0,
            'yes',
            '00000000-0000-0000-0000-000000000000',
            '00000000-0000-0000-0000-000000000000'
);`
	queries[idx] = "insertSystemGroupWheel"
	idx++

	queryMap["insertSystemUserRootAndFriends"] = `
INSERT INTO inventory.user (
            id,
            uid,
            first_name,
            last_name,
            employee_number,
            mail_address,
            is_active,
            is_system,
            is_deleted,
            team_id,
            dictionary_id,
            created_by
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
            '00000000-0000-0000-0000-000000000000',
            '00000000-0000-0000-0000-000000000000',
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
            '00000000-0000-0000-0000-000000000000',
            '00000000-0000-0000-0000-000000000000',
            '00000000-0000-0000-0000-000000000000'
);`
	queries[idx] = "insertSystemUserRootAndFriends"
	idx++

	queryMap[`activate_circular_dependency__1of3`] = `
ALTER TABLE inventory.dictionary ADD CONSTRAINT _dictionary_creator_exists FOREIGN KEY (created_by) REFERENCES inventory.user (id) DEFERRABLE;
`
	queries[idx] = `activate_circular_dependency__1of3`
	idx++

	queryMap[`activate_circular_dependency__2of3`] = `
ALTER TABLE inventory.team ADD CONSTRAINT _team_creator_exists FOREIGN KEY (created_by) REFERENCES inventory.user (id) DEFERRABLE;
`
	queries[idx] = `activate_circular_dependency__2of3`
	idx++

	queryMap[`activate_circular_dependency__3of3`] = `
ALTER TABLE inventory.user ADD CONSTRAINT _user_creator_exists FOREIGN KEY (created_by) REFERENCES inventory.user (id) DEFERRABLE;
`
	queries[idx] = `activate_circular_dependency__3of3`
	idx++

	queryMap["insertCategoryOmnipotence"] = `
INSERT INTO soma.category (
            name,
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
INSERT INTO soma.permission (
            id,
            name,
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
            201811150001,
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
            201811150001,
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
            201811150001,
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
