/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/stmt"
)

func (s *Soma) newDatabaseConn() (*sql.DB, error) {
	driver := `postgres`

	connect := fmt.Sprintf("dbname='%s' user='%s' password='%s' host='%s' port='%s' sslmode='%s' connect_timeout='%s'",
		s.conf.Database.Name,
		s.conf.Database.User,
		s.conf.Database.Pass,
		s.conf.Database.Host,
		s.conf.Database.Port,
		s.conf.Database.TLSMode,
		s.conf.Database.Timeout,
	)

	dbcon, err := sql.Open(driver, connect)
	if err != nil {
		return nil, err
	}
	if err = dbcon.Ping(); err != nil {
		return nil, err
	}

	if _, err = dbcon.Exec(
		stmt.DatabaseTimezone,
	); err != nil {
		return nil, err
	}

	if _, err = dbcon.Exec(
		stmt.DatabaseIsolationLevel,
	); err != nil {
		log.Fatal(err)
	}
	dbcon.SetMaxIdleConns(1)
	dbcon.SetMaxOpenConns(5)
	dbcon.SetConnMaxLifetime(12 * time.Hour)
	s.appLog.Infoln(`Connected new secondary pool to database`)
	return dbcon, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
