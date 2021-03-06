/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package super // import "github.com/mjolnir42/soma/internal/super"

import (
	"database/sql"
	"fmt"

	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	uuid "github.com/satori/go.uuid"
)

func (s *Supervisor) permissionWrite(q *msg.Request, mr *msg.Result) {
	if s.readonly {
		mr.ReadOnly()
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(mr.Error)
		return
	}

	switch q.Action {
	case msg.ActionAdd:
		switch q.Permission.Category {
		case msg.CategoryGlobal,
			msg.CategoryIdentity,
			msg.CategoryPermission,
			msg.CategoryOperation,
			msg.CategoryRepository,
			msg.CategorySelf,
			msg.CategoryTeam,
			msg.CategoryMonitoring:
			s.permissionAdd(q, mr)
		default:
			// Omnipotence, System, Typo, ...
			mr.ServerError(fmt.Errorf(`Illegal category`))
			mr.Super.Audit.WithField(`Code`, mr.Code).
				Warningln(mr.Error)
			return
		}
	case msg.ActionRemove:
		switch q.Permission.Category {
		case msg.CategoryGlobal,
			msg.CategoryIdentity,
			msg.CategoryPermission,
			msg.CategoryOperation,
			msg.CategoryRepository,
			msg.CategorySelf,
			msg.CategoryTeam,
			msg.CategoryMonitoring:
			s.permissionRemove(q, mr)
		default:
			// Omnipotence, System, Typo, ...
			mr.ServerError(fmt.Errorf(`Illegal category`))
			mr.Super.Audit.WithField(`Code`, mr.Code).
				Warningln(mr.Error)
			return
		}
	case msg.ActionMap:
		s.permissionMap(q, mr)
	case msg.ActionUnmap:
		s.permissionUnmap(q, mr)
	}

	if mr.IsOK() {
		go func() {
			s.Update <- msg.CacheUpdateFromRequest(q)
		}()
	}
}

func (s *Supervisor) permissionAdd(q *msg.Request, mr *msg.Result) {
	var (
		err error
		tx  *sql.Tx
		res sql.Result
	)
	txMap := map[string]*sql.Stmt{}

	// open multi-statement transaction
	if tx, err = s.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// prepare statements for this transaction
	for name, statement := range map[string]string{
		`permission_add_tx_perm`: stmt.PermissionAdd,
		`permission_add_tx_link`: stmt.PermissionLinkGrant,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.PermissionTx.Prepare(%s) error: %s",
				name, err.Error())
			mr.ServerError(err, q.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
			tx.Rollback()
			return
		}
	}

	if res, err = s.permissionAddTx(q, txMap); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		tx.Rollback()
		return
	}
	// sets r.OK()
	if !mr.RowCnt(res.RowsAffected()) {
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		tx.Rollback()
		return
	}

	// close transaction
	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	mr.Permission = append(mr.Permission, q.Permission)
	mr.Super.Audit.WithField(`Code`, mr.Code).
		Infoln(fmt.Sprintf(
			"Successfully added permission %s", q.Permission.Name))
}

func (s *Supervisor) permissionRemove(q *msg.Request, mr *msg.Result) {
	var (
		err error
		tx  *sql.Tx
		res sql.Result
	)
	txMap := map[string]*sql.Stmt{}

	// open multi-statement transaction
	if tx, err = s.conn.Begin(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// prepare statements for this transaction
	for name, statement := range map[string]string{
		`permission_rm_tx_rev_gl`: stmt.PermissionRevokeGlobal,
		`permission_rm_tx_rev_rp`: stmt.PermissionRevokeRepository,
		`permission_rm_tx_rev_tm`: stmt.PermissionRevokeTeam,
		`permission_rm_tx_rev_mn`: stmt.PermissionRevokeMonitoring,
		`permission_rm_tx_lookup`: stmt.PermissionLookupGrantID,
		`permission_rm_tx_unlink`: stmt.PermissionRemoveLink,
		`permission_rm_tx_remove`: stmt.PermissionRemove,
		`permission_rm_tx_unmapa`: stmt.PermissionUnmapAll,
	} {
		if txMap[name], err = tx.Prepare(statement); err != nil {
			err = fmt.Errorf("s.PermissionTx.Prepare(%s) error: %s",
				name, err.Error())
			mr.ServerError(err, q.Section)
			mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
			tx.Rollback()
			return
		}
	}

	if res, err = s.permissionRemoveTx(q, txMap); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		tx.Rollback()
		return
	}
	// sets r.OK()
	if !mr.RowCnt(res.RowsAffected()) {
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		tx.Rollback()
		return
	}

	// close transaction
	if err = tx.Commit(); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	mr.Permission = append(mr.Permission, q.Permission)
	mr.Super.Audit.WithField(`Code`, mr.Code).
		Infoln(fmt.Sprintf(
			"Successfully removed permission %s", q.Permission.ID))
}

func (s *Supervisor) permissionMap(q *msg.Request, mr *msg.Result) {
	var (
		err                 error
		res                 sql.Result
		sectionID, actionID sql.NullString
		mapID               string
	)

	// determine if an entire section is mapped or a specific action
	if q.Permission.Actions != nil {
		sectionID.String = (*q.Permission.Actions)[0].SectionID
		sectionID.Valid = true
		actionID.String = (*q.Permission.Actions)[0].ID
		actionID.Valid = true
	} else if q.Permission.Sections != nil {
		sectionID.String = (*q.Permission.Sections)[0].ID
		sectionID.Valid = true
	} else {
		mr.ServerError(fmt.Errorf(`Nothing to map`), q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}
	mapID = uuid.Must(uuid.NewV4()).String()

	if res, err = s.stmtPermissionMapEntry.Exec(
		mapID,
		q.Permission.Category,
		q.Permission.ID,
		sectionID,
		actionID,
		q.AuthUser,
	); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// sets mr.OK()
	if mr.RowCnt(res.RowsAffected()) {
		mr.Permission = append(mr.Permission, q.Permission)
		mr.Super.Audit.WithField(`Code`, mr.Code).Infoln(
			`Successfully mapped to permission`)
		return
	}
	mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
}

func (s *Supervisor) permissionUnmap(q *msg.Request, mr *msg.Result) {
	var (
		err                 error
		res                 sql.Result
		sectionID, actionID sql.NullString
	)

	// determine if an entire section is unmapped or a specific action
	if q.Permission.Actions != nil {
		sectionID.String = (*q.Permission.Actions)[0].SectionID
		sectionID.Valid = true
		actionID.String = (*q.Permission.Actions)[0].ID
		actionID.Valid = true
	} else if q.Permission.Sections != nil {
		sectionID.String = (*q.Permission.Sections)[0].ID
		sectionID.Valid = true
	} else {
		mr.ServerError(fmt.Errorf(`Nothing to map`), q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	if res, err = s.stmtPermissionUnmapEntry.Exec(
		q.Permission.ID,
		q.Permission.Category,
		sectionID,
		actionID,
	); err != nil {
		mr.ServerError(err, q.Section)
		mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
		return
	}

	// sets mr.OK()
	if mr.RowCnt(res.RowsAffected()) {
		mr.Permission = append(mr.Permission, q.Permission)
		mr.Super.Audit.WithField(`Code`, mr.Code).Infoln(
			`Successfully unmapped from permission`)
		return
	}
	mr.Super.Audit.WithField(`Code`, mr.Code).Warningln(err)
}

func (s *Supervisor) permissionAddTx(q *msg.Request,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err                        error
		res                        sql.Result
		grantPermID, grantCategory string
	)

	q.Permission.ID = uuid.Must(uuid.NewV4()).String()
	grantPermID = uuid.Must(uuid.NewV4()).String()
	switch q.Permission.Category {
	case msg.CategoryGlobal:
		grantCategory = msg.CategoryGrantGlobal
	case msg.CategoryPermission:
		grantCategory = msg.CategoryGrantPermission
	case msg.CategoryOperation:
		grantCategory = msg.CategoryGrantOperation
	case msg.CategoryRepository:
		grantCategory = msg.CategoryGrantRepository
	case msg.CategoryTeam:
		grantCategory = msg.CategoryGrantTeam
	case msg.CategoryMonitoring:
		grantCategory = msg.CategoryGrantMonitoring
	case msg.CategorySelf:
		grantCategory = msg.CategoryGrantSelf
	case msg.CategoryIdentity:
		grantCategory = msg.CategoryGrantIdentity
	}

	if res, err = txMap[`permission_add_tx_perm`].Exec(
		q.Permission.ID,
		q.Permission.Name,
		q.Permission.Category,
		q.AuthUser,
	); err != nil {
		return res, err
	}

	if res, err = txMap[`permission_add_tx_perm`].Exec(
		grantPermID,
		q.Permission.Name,
		grantCategory,
		q.AuthUser,
	); err != nil {
		return res, err
	}

	return txMap[`permission_add_tx_link`].Exec(
		grantCategory,
		grantPermID,
		q.Permission.Category,
		q.Permission.ID,
	)
}

func (s *Supervisor) permissionRemoveTx(q *msg.Request,
	txMap map[string]*sql.Stmt) (sql.Result, error) {
	var (
		err                  error
		res                  sql.Result
		grantingPermissionID string
		revocation           string
	)

	// select correct revocation statement scope
	switch q.Permission.Category {
	case msg.CategoryGlobal,
		msg.CategoryPermission,
		msg.CategorySelf,
		msg.CategoryIdentity,
		msg.CategoryOperation:
		revocation = `permission_rm_tx_rev_gl`
	case msg.CategoryRepository:
		revocation = `permission_rm_tx_rev_rp`
	case msg.CategoryTeam:
		revocation = `permission_rm_tx_rev_tm`
	case msg.CategoryMonitoring:
		revocation = `permission_rm_tx_rev_mn`
	}

	// lookup which permission grants this permission
	if err = txMap[`permission_rm_tx_lookup`].QueryRow(
		q.Permission.ID,
	).Scan(
		&grantingPermissionID,
	); err != nil {
		return res, err
	}

	// revoke all grants of the granting permission
	if res, err = txMap[revocation].Exec(
		grantingPermissionID,
	); err != nil {
		return res, err
	}

	// sever the link between permission and granting permission
	if res, err = txMap[`permission_rm_tx_unlink`].Exec(
		q.Permission.ID,
	); err != nil {
		return res, err
	}

	// remove granting permission
	if res, err = txMap[`permission_rm_tx_remove`].Exec(
		grantingPermissionID,
	); err != nil {
		return res, err
	}

	// revoke all grants of the permission
	if res, err = txMap[revocation].Exec(
		q.Permission.ID,
	); err != nil {
		return res, err
	}

	// unmap all sections & actions from the permission
	if res, err = txMap[`permission_rm_tx_unmapa`].Exec(
		q.Permission.ID,
	); err != nil {
		return res, err
	}

	// remove permission
	return txMap[`permission_rm_tx_remove`].Exec(q.Permission.ID)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
