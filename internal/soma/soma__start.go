/*-
 * Copyright (c) 2017, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package soma

// Start launches all application handlers
func (s *Soma) Start() {
	// grimReaper and supervisor must run first
	s.handlerMap.Add(`grimreaper`, newGrimReaper(1, s))
	// TODO: start supervisor completely

	s.handlerMap.Add(`attribute_r`, newAttributeRead(s.conf.QueueLen))
	s.handlerMap.Add(`bucket_r`, newBucketRead(s.conf.QueueLen))
	s.handlerMap.Add(`capability_r`, newCapabilityRead(s.conf.QueueLen))
	s.handlerMap.Add(`checkconfig_r`, newCheckConfigurationRead(s.conf.QueueLen))
	s.handlerMap.Add(`cluster_r`, newClusterRead(s.conf.QueueLen))
	s.handlerMap.Add(`datacenter_r`, newDatacenterRead(s.conf.QueueLen))
	s.handlerMap.Add(`entity_r`, newEntityRead(s.conf.QueueLen))
	s.handlerMap.Add(`environment_r`, newEnvironmentRead(s.conf.QueueLen))
	s.handlerMap.Add(`group_r`, newGroupRead(s.conf.QueueLen))
	s.handlerMap.Add(`hostdeployment_r`, newHostDeploymentRead(s.conf.QueueLen))
	s.handlerMap.Add(`instance_r`, newInstanceRead(s.conf.QueueLen))
	s.handlerMap.Add(`job_r`, newJobRead(s.conf.QueueLen))
	s.handlerMap.Add(`mode_r`, newModeRead(s.conf.QueueLen))
	s.handlerMap.Add(`monitoring_r`, newMonitoringRead(s.conf.QueueLen))
	s.handlerMap.Add(`node_r`, newNodeRead(s.conf.QueueLen))
	s.handlerMap.Add(`oncall_r`, newOncallRead(s.conf.QueueLen))
	s.handlerMap.Add(`predicate_r`, newPredicateRead(s.conf.QueueLen))
	s.handlerMap.Add(`property_r`, newPropertyRead(s.conf.QueueLen))
	s.handlerMap.Add(`provider_r`, newProviderRead(s.conf.QueueLen))
	s.handlerMap.Add(`repository_r`, newRepositoryRead(s.conf.QueueLen))
	s.handlerMap.Add(`server_r`, newServerRead(s.conf.QueueLen))
	s.handlerMap.Add(`state_r`, newStateRead(s.conf.QueueLen))
	s.handlerMap.Add(`status_r`, newStatusRead(s.conf.QueueLen))
	s.handlerMap.Add(`team_r`, newTeamRead(s.conf.QueueLen))
	s.handlerMap.Add(`tree_r`, newTreeRead(s.conf.QueueLen))
	s.handlerMap.Add(`unit_r`, newUnitRead(s.conf.QueueLen))
	s.handlerMap.Add(`user_r`, newUserRead(s.conf.QueueLen))
	s.handlerMap.Add(`validity_r`, newValidityRead(s.conf.QueueLen))
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
		s.handlerMap.Run(handler)
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
