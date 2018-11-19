package main

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/lib/pq"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
)

func connectToDatabase(appLog, errLog *log.Logger) {
	var err error
	var rows *sql.Rows
	var schema string
	var schemaVer int64

	driver := "postgres"

	connect := fmt.Sprintf("dbname='%s' user='%s' password='%s' host='%s' port='%s' sslmode='%s' connect_timeout='%s'",
		SomaCfg.Database.Name,
		SomaCfg.Database.User,
		SomaCfg.Database.Pass,
		SomaCfg.Database.Host,
		SomaCfg.Database.Port,
		SomaCfg.Database.TLSMode,
		SomaCfg.Database.Timeout,
	)

	// enable handling of infinity timestamps
	pq.EnableInfinityTs(msg.NegTimeInf, msg.PosTimeInf)

	conn, err = sql.Open(driver, connect)
	if err != nil {
		log.Fatal(err)
	}
	if err = conn.Ping(); err != nil {
		log.Fatal(err)
	}
	appLog.Print("Connected main pool to database")
	if _, err = conn.Exec(stmt.DatabaseTimezone); err != nil {
		errLog.Fatal(err)
	}
	if _, err = conn.Exec(stmt.DatabaseIsolationLevel); err != nil {
		errLog.Fatal(err)
	}

	// size the connection pool
	conn.SetMaxIdleConns(5)
	conn.SetMaxOpenConns(15)
	conn.SetConnMaxLifetime(12 * time.Hour)

	// required schema versions
	required := map[string]int64{
		"inventory": 201811150001,
		"root":      201605160001,
		`auth`:      201811150001,
		`soma`:      201811150001,
	}

	if rows, err = conn.Query(stmt.DatabaseSchemaVersion); err != nil {
		errLog.Fatal("Query db schema versions: ", err)
	}

	for rows.Next() {
		if err = rows.Scan(
			&schema,
			&schemaVer,
		); err != nil {
			errLog.Fatal("Schema check: ", err)
		}
		if rsv, ok := required[schema]; ok {
			if rsv != schemaVer {
				errLog.Fatalf("Incompatible schema %s: %d != %d", schema, rsv, schemaVer)
			} else {
				appLog.Printf("DB Schema %s, version: %d", schema, schemaVer)
				delete(required, schema)
			}
		} else {
			errLog.Fatal("Unknown schema: ", schema)
		}
	}
	if err = rows.Err(); err != nil {
		errLog.Fatal("Schema check: ", err)
	}
	if len(required) != 0 {
		for s := range required {
			errLog.Printf("Missing schema: %s", s)
		}
		errLog.Fatal("FATAL - database incomplete")
	}
}

func pingDatabase(errLog *log.Logger) {
	ticker := time.NewTicker(time.Second).C

	for {
		<-ticker
		err := conn.Ping()
		if err != nil {
			errLog.Print(err)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
