/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

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
	s.handlerMap.Add(`view_r`, newViewRead(s.conf.QueueLen))
	s.handlerMap.Add(`workflow_r`, newWorkflowRead(s.conf.QueueLen))

	if !s.conf.ReadOnly {
		s.handlerMap.Add(`forest_custodian`, newForestCustodian(s.conf.QueueLen, s))
		s.handlerMap.Add(`guidepost`, newGuidePost(s.conf.QueueLen, s))
		s.handlerMap.Add(`lifecycle`, newLifeCycle(s))

		if !s.conf.Observer {
			s.handlerMap.Add(`attribute_w`, newAttributeWrite(s.conf.QueueLen))
			s.handlerMap.Add(`capability_w`, newCapabilityWrite(s.conf.QueueLen))
			s.handlerMap.Add(`datacenter_w`, newDatacenterWrite(s.conf.QueueLen))
			s.handlerMap.Add(`deployment_w`, newDeploymentWrite(s.conf.QueueLen))
			s.handlerMap.Add(`entity_w`, newEntityWrite(s.conf.QueueLen))
			s.handlerMap.Add(`environment_w`, newEnvironmentWrite(s.conf.QueueLen))
			s.handlerMap.Add(`job_block`, newJobBlock(s.conf.QueueLen))
			s.handlerMap.Add(`mode_w`, newModeWrite(s.conf.QueueLen))
			s.handlerMap.Add(`monitoring_w`, newMonitoringWrite(s.conf.QueueLen))
			s.handlerMap.Add(`node_w`, newNodeWrite(s.conf.QueueLen))
			s.handlerMap.Add(`oncall_w`, newOncallWrite(s.conf.QueueLen))
			s.handlerMap.Add(`predicate_w`, newPredicateWrite(s.conf.QueueLen))
			s.handlerMap.Add(`property_w`, newPropertyWrite(s.conf.QueueLen))
			s.handlerMap.Add(`provider_w`, newProviderWrite(s.conf.QueueLen))
			s.handlerMap.Add(`server_w`, newServerWrite(s.conf.QueueLen))
			s.handlerMap.Add(`state_w`, newStateWrite(s.conf.QueueLen))
			s.handlerMap.Add(`status_w`, newStatusWrite(s.conf.QueueLen))
			s.handlerMap.Add(`team_w`, newTeamWrite(s.conf.QueueLen, s))
			s.handlerMap.Add(`unit_w`, newUnitWrite(s.conf.QueueLen))
			s.handlerMap.Add(`user_w`, newUserWrite(s.conf.QueueLen, s))
			s.handlerMap.Add(`validity_w`, newValidityWrite(s.conf.QueueLen))
			s.handlerMap.Add(`view_w`, newViewWrite(s.conf.QueueLen))
			s.handlerMap.Add(`workflow_w`, newWorkflowWrite(s.conf.QueueLen))
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
