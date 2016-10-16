package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1and1/soma/lib/proto"
	log "github.com/Sirupsen/logrus"
)

// Message structs
type somaViewRequest struct {
	action string
	name   string
	View   proto.View
	reply  chan somaResult
}

type somaViewResult struct {
	ResultError error
	View        proto.View
}

func (a *somaViewResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Views = append(r.Views, somaViewResult{ResultError: err})
	}
}

func (a *somaViewResult) SomaAppendResult(r *somaResult) {
	r.Views = append(r.Views, *a)
}

/*  Read Access
 */
type somaViewReadHandler struct {
	input     chan somaViewRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
	appLog    *log.Logger
	reqLog    *log.Logger
	errLog    *log.Logger
}

func (r *somaViewReadHandler) run() {
	var err error

	r.list_stmt, err = r.conn.Prepare(`
SELECT view
FROM   soma.views;`)
	if err != nil {
		log.Fatal("view/list: ", err)
	}
	defer r.list_stmt.Close()

	r.show_stmt, err = r.conn.Prepare(`
SELECT view
FROM   soma.views
WHERE  view = $1::varchar;`)
	if err != nil {
		log.Fatal("view/show: ", err)
	}
	defer r.show_stmt.Close()

	for {
		select {
		case <-r.shutdown:
			break
		case req := <-r.input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

func (r *somaViewReadHandler) process(q *somaViewRequest) {
	var (
		view string
		rows *sql.Rows
		err  error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: view/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err = rows.Scan(&view)
			result.Append(err, &somaViewResult{
				View: proto.View{
					Name: view,
				},
			})
		}
		if err = rows.Err(); err != nil {
			result.Append(err, &somaViewResult{})
			err = nil
		}
	case "show":
		log.Printf("R: view/show for %s", q.View.Name)
		err = r.show_stmt.QueryRow(q.View.Name).Scan(
			&view,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				result.SetNotFound()
			} else {
				_ = result.SetRequestError(err)
			}
			q.reply <- result
			return
		}

		result.Append(err, &somaViewResult{
			View: proto.View{
				Name: view,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */

type somaViewWriteHandler struct {
	input    chan somaViewRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
	ren_stmt *sql.Stmt
	appLog   *log.Logger
	reqLog   *log.Logger
	errLog   *log.Logger
}

func (w *somaViewWriteHandler) run() {
	var err error

	w.add_stmt, err = w.conn.Prepare(`
INSERT INTO soma.views (
	view)
SELECT $1::varchar WHERE NOT EXISTS (
    SELECT view
	FROM   soma.views
	WHERE  view = $1::varchar);`)
	if err != nil {
		log.Fatal("view/add: ", err)
	}
	defer w.add_stmt.Close()

	w.del_stmt, err = w.conn.Prepare(`
DELETE FROM soma.views
WHERE  view = $1::varchar;`)
	if err != nil {
		log.Fatal("view/delete: ", err)
	}
	defer w.del_stmt.Close()

	w.ren_stmt, err = w.conn.Prepare(`
UPDATE soma.views
SET    view = $1::varchar
WHERE  view = $2::varchar;`)
	if err != nil {
		log.Fatal("view/rename: ", err)
	}
	defer w.ren_stmt.Close()

runloop:
	for {
		select {
		case <-w.shutdown:
			break runloop
		case req := <-w.input:
			w.process(&req)
		}
	}
}

func (w *somaViewWriteHandler) process(q *somaViewRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: view/add for %s", q.View.Name)
		res, err = w.add_stmt.Exec(
			q.View.Name,
		)
	case "delete":
		log.Printf("R: view/delete for %s", q.View.Name)
		res, err = w.del_stmt.Exec(
			q.View.Name,
		)
	case "rename":
		log.Printf("R: view/rename for %s", q.name)
		res, err = w.ren_stmt.Exec(
			q.View.Name,
			q.name,
		)
	default:
		log.Printf("R: unimplemented levels/%s", q.action)
		result.SetNotImplemented()
		q.reply <- result
		return
	}
	if result.SetRequestError(err) {
		q.reply <- result
		return
	}

	rowCnt, _ := res.RowsAffected()
	switch {
	case rowCnt == 0:
		result.Append(errors.New("No rows affected"), &somaViewResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaViewResult{})
	default:
		result.Append(nil, &somaViewResult{
			View: q.View,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaViewReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaViewWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
