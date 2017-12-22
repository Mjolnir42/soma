/*
 * Copyright (c) 2016, 1&1 Internet SE
 * Written by Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved.
 */

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func connectToDatabase() {
	var err error
	driver := "postgres"

	connect := fmt.Sprintf(
		"%s='%s' %s='%s' %s='%s' %s='%s' %s='%s' %s='%s' %s='%s'",
		"dbname",
		Eye.Database.Name,
		"user",
		Eye.Database.User,
		"password",
		Eye.Database.Pass,
		"host",
		Eye.Database.Host,
		"port",
		Eye.Database.Port,
		"sslmode",
		Eye.Database.TLSMode,
		"connect_timeout",
		Eye.Database.Timeout,
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

func prepareStatements() {
	var err error

	Eye.run.checkItem, err = Eye.run.conn.Prepare(stmtCheckItemExists)
	log.Println("Preparing: check_item")
	abortOnError(err)

	Eye.run.checkLookup, err = Eye.run.conn.Prepare(stmtCheckLookupExists)
	log.Println("Preparing: check_lookup")
	abortOnError(err)

	Eye.run.deleteItem, err = Eye.run.conn.Prepare(stmtDeleteConfigurationItem)
	log.Println("Preparing: delete_item")
	abortOnError(err)

	Eye.run.deleteLookup, err = Eye.run.conn.Prepare(stmtDeleteLookupID)
	log.Println("Preparing: delete_lookup")
	abortOnError(err)

	Eye.run.getConfig, err = Eye.run.conn.Prepare(stmtGetSingleConfiguration)
	log.Println("Preparing: get_config")
	abortOnError(err)

	Eye.run.getItems, err = Eye.run.conn.Prepare(stmtGetConfigurationItemIds)
	log.Println("Preparing: get_items")
	abortOnError(err)

	Eye.run.getLookup, err = Eye.run.conn.Prepare(stmtGetLookupIDForItem)
	log.Println("Preparing: get_lookup")
	abortOnError(err)

	Eye.run.insertItem, err = Eye.run.conn.Prepare(stmtInsertConfigurationItem)
	log.Println("Preparing: insert_item")
	abortOnError(err)

	Eye.run.insertLookup, err = Eye.run.conn.Prepare(stmtInsertLookupInformation)
	log.Println("Preparing: insert_lookup")
	abortOnError(err)

	Eye.run.itemCount, err = Eye.run.conn.Prepare(stmtGetItemCountForLookupID)
	log.Println("Preparing: item_count")
	abortOnError(err)

	Eye.run.retrieve, err = Eye.run.conn.Prepare(stmtRetrieveConfigurationsByLookup)
	log.Println("Preparing: retrieve")
	abortOnError(err)

	Eye.run.updateItem, err = Eye.run.conn.Prepare(stmtUpdateConfigurationItem)
	log.Println("Preparing: update_item")
	abortOnError(err)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
