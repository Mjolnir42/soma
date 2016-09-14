package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/1and1/soma/lib/proto"
)

type somaAttributeRequest struct {
	action    string
	Attribute proto.Attribute
	reply     chan somaResult
}

type somaAttributeResult struct {
	ResultError error
	Attribute   proto.Attribute
}

func (a *somaAttributeResult) SomaAppendError(r *somaResult, err error) {
	if err != nil {
		r.Attributes = append(r.Attributes,
			somaAttributeResult{ResultError: err})
	}
}

func (a *somaAttributeResult) SomaAppendResult(r *somaResult) {
	r.Attributes = append(r.Attributes, *a)
}

/* Read Access
 */
type somaAttributeReadHandler struct {
	input     chan somaAttributeRequest
	shutdown  chan bool
	conn      *sql.DB
	list_stmt *sql.Stmt
	show_stmt *sql.Stmt
}

func (r *somaAttributeReadHandler) run() {
	var err error

	if r.list_stmt, err = r.conn.Prepare(stmtAttributeList); err != nil {
		log.Fatal("attribute/list: ", err)
	}
	defer r.list_stmt.Close()

	if r.show_stmt, err = r.conn.Prepare(stmtAttributeShow); err != nil {
		log.Fatal("attribute/show: ", err)
	}
	defer r.show_stmt.Close()

runloop:
	for {
		select {
		case <-r.shutdown:
			break runloop
		case req := <-r.input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

func (r *somaAttributeReadHandler) process(q *somaAttributeRequest) {
	var (
		attribute, cardinality string
		rows                   *sql.Rows
		err                    error
	)
	result := somaResult{}

	switch q.action {
	case "list":
		log.Printf("R: attributes/list")
		rows, err = r.list_stmt.Query()
		defer rows.Close()
		if result.SetRequestError(err) {
			q.reply <- result
			return
		}

		for rows.Next() {
			err := rows.Scan(&attribute, &cardinality)
			result.Append(err, &somaAttributeResult{
				Attribute: proto.Attribute{
					Name:        attribute,
					Cardinality: cardinality,
				},
			})
		}
	case "show":
		log.Printf("R: attribute/show for %s", q.Attribute.Name)
		err = r.show_stmt.QueryRow(q.Attribute.Name).Scan(
			&attribute,
			&cardinality,
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

		result.Append(err, &somaAttributeResult{
			Attribute: proto.Attribute{
				Name:        attribute,
				Cardinality: cardinality,
			},
		})
	default:
		result.SetNotImplemented()
	}
	q.reply <- result
}

/* Write Access
 */
type somaAttributeWriteHandler struct {
	input    chan somaAttributeRequest
	shutdown chan bool
	conn     *sql.DB
	add_stmt *sql.Stmt
	del_stmt *sql.Stmt
}

func (w *somaAttributeWriteHandler) run() {
	var err error

	if w.add_stmt, err = w.conn.Prepare(stmtAttributeAdd); err != nil {
		log.Fatal("attribute/add: ", err)
	}
	defer w.add_stmt.Close()

	if w.del_stmt, err = w.conn.Prepare(stmtAttributeDelete); err != nil {
		log.Fatal("attribute/delete: ", err)
	}
	defer w.del_stmt.Close()

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

func (w *somaAttributeWriteHandler) process(q *somaAttributeRequest) {
	var (
		res sql.Result
		err error
	)
	result := somaResult{}

	switch q.action {
	case "add":
		log.Printf("R: attributes/add for %s", q.Attribute.Name)
		res, err = w.add_stmt.Exec(
			q.Attribute.Name,
			q.Attribute.Cardinality,
		)
	case "delete":
		log.Printf("R: attributes/del for %s", q.Attribute.Name)
		res, err = w.del_stmt.Exec(
			q.Attribute.Name,
		)
	default:
		log.Printf("R: unimplemented attributes/%s", q.action)
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
		result.Append(errors.New("No rows affected"), &somaAttributeResult{})
	case rowCnt > 1:
		result.Append(fmt.Errorf("Too many rows affected: %d", rowCnt),
			&somaAttributeResult{})
	default:
		result.Append(nil, &somaAttributeResult{
			Attribute: q.Attribute,
		})
	}
	q.reply <- result
}

/* Ops Access
 */
func (r *somaAttributeReadHandler) shutdownNow() {
	r.shutdown <- true
}

func (w *somaAttributeWriteHandler) shutdownNow() {
	w.shutdown <- true
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
