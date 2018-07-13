/*-
 * Copyright (c) 2016-2017, Jörg Pernfuß
 * Copyright (c) 2016, 1&1 Internet SE
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/mjolnir42/soma/internal/handler"
	"github.com/mjolnir42/soma/internal/msg"
	"github.com/mjolnir42/soma/internal/stmt"
	"github.com/mjolnir42/soma/lib/proto"
)

// CheckConfigurationRead handles read requests for
// check configurations
type CheckConfigurationRead struct {
	Input                       chan msg.Request
	Shutdown                    chan struct{}
	handlerName                 string
	conn                        *sql.DB
	stmtList                    *sql.Stmt
	stmtShow                    *sql.Stmt
	stmtShowThreshold           *sql.Stmt
	stmtShowConstraintCustom    *sql.Stmt
	stmtShowConstraintSystem    *sql.Stmt
	stmtShowConstraintNative    *sql.Stmt
	stmtShowConstraintService   *sql.Stmt
	stmtShowConstraintAttribute *sql.Stmt
	stmtShowConstraintOncall    *sql.Stmt
	stmtShowInstanceInfo        *sql.Stmt
	appLog                      *logrus.Logger
	reqLog                      *logrus.Logger
	errLog                      *logrus.Logger
}

// newCheckConfigurationRead returns a new
// CheckConfigurationRead handler with input
// buffer of length
func newCheckConfigurationRead(length int) (string, *CheckConfigurationRead) {
	r := &CheckConfigurationRead{}
	r.handlerName = generateHandlerName() + `_r`
	r.Input = make(chan msg.Request, length)
	r.Shutdown = make(chan struct{})
	return r.handlerName, r
}

// Register initializes resources provided by the Soma app
func (r *CheckConfigurationRead) Register(c *sql.DB, l ...*logrus.Logger) {
	r.conn = c
	r.appLog = l[0]
	r.reqLog = l[1]
	r.errLog = l[2]
}

// RegisterRequests links the handler inside the handlermap to the requests
// it processes
func (r *CheckConfigurationRead) RegisterRequests(hmap *handler.Map) {
	for _, action := range []string{
		msg.ActionList,
		msg.ActionShow,
		msg.ActionSearch,
	} {
		hmap.Request(msg.SectionCheckConfig, action, r.handlerName)
	}
}

// Intake exposes the Input channel as part of the handler interface
func (r *CheckConfigurationRead) Intake() chan msg.Request {
	return r.Input
}

// PriorityIntake aliases Intake as part of the handler interface
func (r *CheckConfigurationRead) PriorityIntake() chan msg.Request {
	return r.Intake()
}

// Run is the event loop for CheckConfigurationRead
func (r *CheckConfigurationRead) Run() {
	var err error

	for statement, prepStmt := range map[string]*sql.Stmt{
		stmt.CheckConfigList:                r.stmtList,
		stmt.CheckConfigShowBase:            r.stmtShow,
		stmt.CheckConfigShowThreshold:       r.stmtShowThreshold,
		stmt.CheckConfigShowConstrCustom:    r.stmtShowConstraintCustom,
		stmt.CheckConfigShowConstrSystem:    r.stmtShowConstraintSystem,
		stmt.CheckConfigShowConstrNative:    r.stmtShowConstraintNative,
		stmt.CheckConfigShowConstrService:   r.stmtShowConstraintService,
		stmt.CheckConfigShowConstrAttribute: r.stmtShowConstraintAttribute,
		stmt.CheckConfigShowConstrOncall:    r.stmtShowConstraintOncall,
		stmt.CheckConfigInstanceInfo:        r.stmtShowInstanceInfo,
	} {
		if prepStmt, err = r.conn.Prepare(statement); err != nil {
			r.errLog.Fatal(`checkconfig`, err, stmt.Name(statement))
		}
		defer prepStmt.Close()
	}

runloop:
	for {
		select {
		case <-r.Shutdown:
			break runloop
		case req := <-r.Input:
			go func() {
				r.process(&req)
			}()
		}
	}
}

// process is the request dispatcher
func (r *CheckConfigurationRead) process(q *msg.Request) {
	result := msg.FromRequest(q)
	msgRequest(r.reqLog, q)

	switch q.Action {
	case msg.ActionList:
		r.list(q, &result)
	case msg.ActionShow:
		r.show(q, &result)
	case msg.ActionSearch:
		// XXX BUG x.search(q, &result)
	default:
		result.UnknownRequest(q)
	}
	q.Reply <- result
}

// list returns all check configurations
func (r *CheckConfigurationRead) list(q *msg.Request, mr *msg.Result) {
	var (
		configID, repoID, configName string
		bucketNULL                   sql.NullString
		bucketID                     string
		rows                         *sql.Rows
		err                          error
	)

	if rows, err = r.stmtList.Query(
		q.CheckConfig.RepositoryID,
	); err != nil {
		mr.ServerError(err, q.Section)
		return
	}

	for rows.Next() {
		if err = rows.Scan(
			&configID,
			&repoID,
			&bucketNULL,
			&configName,
		); err != nil {
			rows.Close()
			mr.ServerError(err, q.Section)
			return
		}
		if bucketNULL.Valid {
			// bucketNULL is null if the check configuration
			// is on a repository
			bucketID = bucketNULL.String
		}
		mr.CheckConfig = append(mr.CheckConfig,
			proto.CheckConfig{
				ID:           configID,
				RepositoryID: repoID,
				BucketID:     bucketID,
				Name:         configName,
			},
		)
	}
	if err = rows.Err(); err != nil {
		mr.ServerError(err, q.Section)
		return
	}
	mr.OK()
}

// show returns details for a check configuration
func (r *CheckConfigurationRead) show(q *msg.Request, mr *msg.Result) {
	var (
		configID, repoID, configName   string
		bucketID, objectID, objectType string
		capabilityID, externalID       string
		bucketNULL                     sql.NullString
		err                            error
		interval                       int64
		isActive, hasInheritance       bool
		isChildrenOnly, isEnabled      bool
		checkConfig                    proto.CheckConfig
	)

	if err = r.stmtShow.QueryRow(
		q.CheckConfig.ID,
	).Scan(
		&configID,
		&repoID,
		&bucketNULL,
		&configName,
		&objectID,
		&objectType,
		&isActive,
		&hasInheritance,
		&isChildrenOnly,
		&capabilityID,
		&interval,
		&isEnabled,
		&externalID,
	); err == sql.ErrNoRows {
		mr.NotFound(err, q.Section)
		return
	} else if err != nil {
		goto fail
	}

	if bucketNULL.Valid {
		bucketID = bucketNULL.String
	}

	checkConfig = proto.CheckConfig{
		ID:           configID,
		Name:         configName,
		Interval:     uint64(interval),
		RepositoryID: repoID,
		BucketID:     bucketID,
		CapabilityID: capabilityID,
		ObjectID:     objectID,
		ObjectType:   objectType,
		IsActive:     isActive,
		IsEnabled:    isEnabled,
		Inheritance:  hasInheritance,
		ChildrenOnly: isChildrenOnly,
		ExternalID:   externalID,
	}

	// retrieve check configuration thresholds
	if err = r.thresholds(&checkConfig); err != nil {
		goto fail
	}

	if err = r.constraints(&checkConfig); err != nil {
		goto fail
	}

	if err = r.instances(&checkConfig); err != nil {
		goto fail
	}

	mr.CheckConfig = append(mr.CheckConfig, checkConfig)
	mr.OK()
	return

fail:
	mr.ServerError(err, q.Section)
}

// thresholds add thresholds to a check configuration
func (r *CheckConfigurationRead) thresholds(cnf *proto.CheckConfig) error {
	var (
		predicate, threshold, lvlName, lvlShort string
		configID                                string
		lvlNumeric, value                       int64
		err                                     error
		rows                                    *sql.Rows
	)

	if rows, err = r.stmtShowThreshold.Query(
		cnf.ID,
	); err != nil {
		return err
	}
	defer rows.Close()

	cnf.Thresholds = make(
		[]proto.CheckConfigThreshold,
		0,
	)

	for rows.Next() {
		if err = rows.Scan(
			&configID,
			&predicate,
			&threshold,
			&lvlName,
			&lvlShort,
			&lvlNumeric,
		); err != nil {
			return err
		}

		if value, err = strconv.ParseInt(
			threshold, 10, 64,
		); err != nil {
			return err
		}

		thr := proto.CheckConfigThreshold{
			Predicate: proto.Predicate{
				Symbol: predicate,
			},
			Level: proto.Level{
				Name:      lvlName,
				ShortName: lvlShort,
				Numeric:   uint16(lvlNumeric),
			},
			Value: value,
		}

		cnf.Thresholds = append(cnf.Thresholds, thr)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	// check configurations must have at least one threshold
	if len(cnf.Thresholds) == 0 {
		return fmt.Errorf(`CheckConfiguration has no` +
			`thresholds defined`)
	}
	return nil
}

// constraints adds constraints to a check configuration
func (r *CheckConfigurationRead) constraints(cnf *proto.CheckConfig) error {
	var err error
	cnf.Constraints = make([]proto.CheckConfigConstraint, 0)

	for _, typ := range []string{
		`custom`,
		`system`,
		`native`,
		`service`,
		`attribute`,
		`oncall`,
	} {
		switch typ {
		case `custom`:
			if err = r.constraintCustom(cnf); err != nil {
				return err
			}
		case `system`:
			if err = r.constraintSystem(cnf); err != nil {
				return err
			}
		case `native`:
			if err = r.constraintNative(cnf); err != nil {
				return err
			}
		case `service`:
			if err = r.constraintService(cnf); err != nil {
				return err
			}
		case `attribute`:
			if err = r.constraintAttribute(cnf); err != nil {
				return err
			}
		case `oncall`:
			if err = r.constraintOncall(cnf); err != nil {
				return err
			}
		}
	}
	return nil
}

// constraintCustom adds constraints on custom properties to
// a check configuration
func (r *CheckConfigurationRead) constraintCustom(cnf *proto.CheckConfig) error {
	var (
		configID, propertyID, repoID, property, value string
		rows                                          *sql.Rows
		err                                           error
	)

	if rows, err = r.stmtShowConstraintCustom.Query(
		cnf.ID,
	); err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&configID,
			&propertyID,
			&repoID,
			&value,
			&property,
		); err != nil {
			return err
		}

		constraint := proto.CheckConfigConstraint{
			ConstraintType: `custom`,
			Custom: &proto.PropertyCustom{
				ID:           propertyID,
				RepositoryID: repoID,
				Name:         property,
				Value:        value,
			},
		}
		cnf.Constraints = append(cnf.Constraints, constraint)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// constraintSystem adds constraints on system properties to
// a check configuration
func (r *CheckConfigurationRead) constraintSystem(cnf *proto.CheckConfig) error {
	var (
		configID, property, value string
		rows                      *sql.Rows
		err                       error
	)

	if rows, err = r.stmtShowConstraintSystem.Query(
		cnf.ID,
	); err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&configID,
			&property,
			&value,
		); err != nil {
			return err
		}

		constraint := proto.CheckConfigConstraint{
			ConstraintType: `system`,
			System: &proto.PropertySystem{
				Name:  property,
				Value: value,
			},
		}
		cnf.Constraints = append(cnf.Constraints, constraint)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// constraintNative adds constraints on native properties to
// a check configuration
func (r *CheckConfigurationRead) constraintNative(cnf *proto.CheckConfig) error {
	var (
		configID, property, value string
		rows                      *sql.Rows
		err                       error
	)

	if rows, err = r.stmtShowConstraintNative.Query(
		cnf.ID,
	); err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&configID,
			&property,
			&value,
		); err != nil {
			return err
		}

		constraint := proto.CheckConfigConstraint{
			ConstraintType: `native`,
			Native: &proto.PropertyNative{
				Name:  property,
				Value: value,
			},
		}

		cnf.Constraints = append(cnf.Constraints, constraint)
	}
	return nil
}

// constraintService adds constraints on service properties to
// a check configuration
func (r *CheckConfigurationRead) constraintService(cnf *proto.CheckConfig) error {
	var (
		configID, teamID, svcName string
		rows                      *sql.Rows
		err                       error
	)

	if rows, err = r.stmtShowConstraintService.Query(
		cnf.ID,
	); err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&configID,
			&teamID,
			&svcName,
		); err != nil {
			return err
		}

		constraint := proto.CheckConfigConstraint{
			ConstraintType: `service`,
			Service: &proto.PropertyService{
				Name:   svcName,
				TeamID: teamID,
			},
		}
		cnf.Constraints = append(cnf.Constraints, constraint)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// constraintAttribute adds constraints on service attributes
// to a check configuration
func (r *CheckConfigurationRead) constraintAttribute(cnf *proto.CheckConfig) error {
	var (
		configID, attribute, value string
		rows                       *sql.Rows
		err                        error
	)

	if rows, err = r.stmtShowConstraintAttribute.Query(
		cnf.ID,
	); err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&configID,
			&attribute,
			&value,
		); err != nil {
			return err
		}

		constraint := proto.CheckConfigConstraint{
			ConstraintType: `attribute`,
			Attribute: &proto.ServiceAttribute{
				Name:  attribute,
				Value: value,
			},
		}
		cnf.Constraints = append(cnf.Constraints, constraint)
	}
	return nil
}

// constraintOncall adds constraints on oncall properties
// to a check configuration
func (r *CheckConfigurationRead) constraintOncall(cnf *proto.CheckConfig) error {
	var (
		configID, oncallID, oncallName, oncallNumber string
		rows                                         *sql.Rows
		err                                          error
	)

	if rows, err = r.stmtShowConstraintOncall.Query(
		cnf.ID,
	); err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(
			&configID,
			&oncallID,
			&oncallName,
			&oncallNumber,
		); err != nil {
			return err
		}

		constraint := proto.CheckConfigConstraint{
			ConstraintType: `oncall`,
			Oncall: &proto.PropertyOncall{
				ID:     oncallID,
				Name:   oncallName,
				Number: oncallNumber,
			},
		}
		cnf.Constraints = append(cnf.Constraints, constraint)
	}
	return nil
}

// instances adds information about spawned instances
// to a check configuration
func (r *CheckConfigurationRead) instances(cnf *proto.CheckConfig) error {
	var (
		rows                             *sql.Rows
		instanceID, objectID, objectType string
		currentStatus, nextStatus        string
		err                              error
	)

	if rows, err = r.stmtShowInstanceInfo.Query(
		cnf.ID,
	); err != nil {
		return err
	}
	defer rows.Close()

	instances := []proto.CheckInstanceInfo{}

	for rows.Next() {
		if err = rows.Scan(
			&instanceID,
			&objectID,
			&objectType,
			&currentStatus,
			&nextStatus,
		); err != nil {
			return err
		}

		instance := proto.CheckInstanceInfo{
			ID:            instanceID,
			ObjectID:      objectID,
			ObjectType:    objectType,
			CurrentStatus: currentStatus,
			NextStatus:    nextStatus,
		}
		instances = append(instances, instance)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	if len(instances) > 0 {
		cnf.Details = &proto.CheckConfigDetails{
			Instances: instances,
		}
	}
	return nil
}

// ShutdownNow signals the handler to shut down
func (r *CheckConfigurationRead) ShutdownNow() {
	close(r.Shutdown)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
