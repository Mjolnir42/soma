/*-
 * Copyright (c) 2017-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma // import "github.com/mjolnir42/soma/internal/soma"

import "github.com/mjolnir42/soma/internal/super"

// Start launches all application handlers
func (s *Soma) Start() {
	// grimReaper and supervisor must run first
	s.handlerMap.Add(`grimreaper`, newGrimReaper(1, s))
	s.handlerMap.Register(`grimreaper`, s.dbConnection, s.exportLogger())
	s.handlerMap.Run(`grimreaper`)

	sv := super.New(s.conf)
	sv.RegisterAuditLog(s.auditLog)
	s.handlerMap.Add(`supervisor`, sv)
	s.handlerMap.Register(`supervisor`, s.dbConnection, s.exportLogger())
	s.handlerMap.Run(`supervisor`)

	// start regular handlers
	s.handlerMap.Add(newAttributeRead(s.conf.QueueLen))
	s.handlerMap.Add(newBucketRead(s.conf.QueueLen))
	s.handlerMap.Add(newCapabilityRead(s.conf.QueueLen))
	s.handlerMap.Add(newCheckConfigurationRead(s.conf.QueueLen))
	s.handlerMap.Add(newClusterRead(s.conf.QueueLen))
	s.handlerMap.Add(newDatacenterRead(s.conf.QueueLen))
	s.handlerMap.Add(newEntityRead(s.conf.QueueLen))
	s.handlerMap.Add(newEnvironmentRead(s.conf.QueueLen))
	s.handlerMap.Add(newGroupRead(s.conf.QueueLen))
	s.handlerMap.Add(newHostDeploymentRead(s.conf.QueueLen))
	s.handlerMap.Add(newInstanceRead(s.conf.QueueLen))
	s.handlerMap.Add(newJobRead(s.conf.QueueLen))
	s.handlerMap.Add(newJobResultRead(s.conf.QueueLen))
	s.handlerMap.Add(newJobStatusRead(s.conf.QueueLen))
	s.handlerMap.Add(newJobTypeRead(s.conf.QueueLen))
	s.handlerMap.Add(newLevelRead(s.conf.QueueLen))
	s.handlerMap.Add(newMetricRead(s.conf.QueueLen))
	s.handlerMap.Add(newModeRead(s.conf.QueueLen))
	s.handlerMap.Add(newMonitoringRead(s.conf.QueueLen))
	s.handlerMap.Add(newNodeRead(s.conf.QueueLen))
	s.handlerMap.Add(newOncallRead(s.conf.QueueLen))
	s.handlerMap.Add(newPredicateRead(s.conf.QueueLen))
	s.handlerMap.Add(newPropertyRead(s.conf.QueueLen))
	s.handlerMap.Add(newProviderRead(s.conf.QueueLen))
	s.handlerMap.Add(newRepositoryRead(s.conf.QueueLen))
	s.handlerMap.Add(newServerRead(s.conf.QueueLen))
	s.handlerMap.Add(newStateRead(s.conf.QueueLen))
	s.handlerMap.Add(newStatusRead(s.conf.QueueLen))
	s.handlerMap.Add(newTeamRead(s.conf.QueueLen))
	s.handlerMap.Add(newTreeRead(s.conf.QueueLen))
	s.handlerMap.Add(newUnitRead(s.conf.QueueLen))
	s.handlerMap.Add(newUserRead(s.conf.QueueLen))
	s.handlerMap.Add(newValidityRead(s.conf.QueueLen))
	s.handlerMap.Add(newViewRead(s.conf.QueueLen))
	s.handlerMap.Add(newWorkflowRead(s.conf.QueueLen))

	if !s.conf.ReadOnly {
		s.handlerMap.Add(`forest_custodian`, newForestCustodian(s.conf.QueueLen, s))
		s.handlerMap.Add(`guidepost`, newGuidePost(s.conf.QueueLen, s))
		s.handlerMap.Add(`lifecycle`, newLifeCycle(s))

		if !s.conf.Observer {
			s.handlerMap.Add(newAttributeWrite(s.conf.QueueLen))
			s.handlerMap.Add(newCapabilityWrite(s.conf.QueueLen))
			s.handlerMap.Add(newDatacenterWrite(s.conf.QueueLen))
			s.handlerMap.Add(newDeploymentWrite(s.conf.QueueLen))
			s.handlerMap.Add(newEntityWrite(s.conf.QueueLen))
			s.handlerMap.Add(newEnvironmentWrite(s.conf.QueueLen))
			s.handlerMap.Add(`job_block`, newJobBlock(s.conf.QueueLen))
			s.handlerMap.Add(newJobResultWrite(s.conf.QueueLen))
			s.handlerMap.Add(newJobStatusWrite(s.conf.QueueLen))
			s.handlerMap.Add(newJobTypeWrite(s.conf.QueueLen))
			s.handlerMap.Add(newLevelWrite(s.conf.QueueLen))
			s.handlerMap.Add(newMetricWrite(s.conf.QueueLen))
			s.handlerMap.Add(newModeWrite(s.conf.QueueLen))
			s.handlerMap.Add(newMonitoringWrite(s.conf.QueueLen))
			s.handlerMap.Add(newNodeWrite(s.conf.QueueLen))
			s.handlerMap.Add(newOncallWrite(s.conf.QueueLen))
			s.handlerMap.Add(newPredicateWrite(s.conf.QueueLen))
			s.handlerMap.Add(newPropertyWrite(s.conf.QueueLen))
			s.handlerMap.Add(newProviderWrite(s.conf.QueueLen))
			s.handlerMap.Add(newServerWrite(s.conf.QueueLen))
			s.handlerMap.Add(newStateWrite(s.conf.QueueLen))
			s.handlerMap.Add(newStatusWrite(s.conf.QueueLen))
			s.handlerMap.Add(newTeamWrite(s.conf.QueueLen, s))
			s.handlerMap.Add(newUnitWrite(s.conf.QueueLen))
			s.handlerMap.Add(newUserWrite(s.conf.QueueLen, s))
			s.handlerMap.Add(newValidityWrite(s.conf.QueueLen))
			s.handlerMap.Add(newViewWrite(s.conf.QueueLen))
			s.handlerMap.Add(newWorkflowWrite(s.conf.QueueLen))
		}
	}

	// fully initialize the handlers and fire them up
	for handler := range s.handlerMap.Range() {
		switch handler {
		case `supervisor`, `grimReaper`:
			// already running
			continue
		}
		s.handlerMap.Register(
			handler,
			s.dbConnection,
			s.exportLogger(),
		)
		// starts the handler in a goroutine
		s.handlerMap.Run(handler)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
