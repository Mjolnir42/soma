package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func connectToDatabase() {
	var err error
	driver := "postgres"

	connect := fmt.Sprintf("dbname='%s' user='%s' password='%s' host='%s' port='%s' sslmode='%s' connect_timeout='%s'",
		Eye.Database.Name,
		Eye.Database.User,
		Eye.Database.Pass,
		Eye.Database.Host,
		Eye.Database.Port,
		Eye.TlsMode,
		Eye.Timeout,
	)

	Eye.run.conn, err = sql.Open(driver, connect)
	if err != nil {
		log.Fatal(err)
	}
	if err = Eye.run.conn.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Print("Connected to database")
	if _, err = Eye.run.conn.Exec(`SET TIME ZONE 'UTC';`); err != nil {
		log.Fatal(err)
	}
}

func pingDatabase() {
	ticker := time.NewTicker(time.Second).C

	for {
		<-ticker
		err := Eye.run.conn.Ping()
		if err != nil {
			log.Print(err)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
