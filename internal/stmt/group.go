/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <joerg.pernfuss@1und1.de>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

const (
	GroupStatements = ``

	GroupShow = `
SELECT group_id,
       bucket_id,
       group_name,
       object_state,
       organizational_team_id
FROM   soma.groups
WHERE  group_id = $1::uuid;`

	GroupMemberGroupList = `
SELECT sg.group_id,
       sg.group_name,
       osg.group_name
FROM   soma.group_membership_groups sgmg
JOIN   soma.groups sg
ON     sgmg.child_group_id = sg.group_id
JOIN   soma.groups osg
ON     sgmg.group_id = osg.group_id
WHERE  sgmg.group_id = $1::uuid;`

	GroupMemberClusterList = `
SELECT sc.cluster_id,
       sc.cluster_name,
       sg.group_name
FROM   soma.group_membership_clusters sgmc
JOIN   soma.clusters sc
ON     sgmc.child_cluster_id = sc.cluster_id
JOIN   soma.groups sg
ON     sgmc.group_id = sg.group_id
WHERE  sgmc.group_id = $1::uuid;`

	GroupMemberNodeList = `
SELECT sn.node_id,
       sn.node_name,
       sg.group_name
FROM   soma.group_membership_nodes sgmn
JOIN   soma.nodes sn
ON     sgmn.child_node_id = sn.node_id
JOIN   soma.groups sg
ON     sgmn.group_id = sg.group_id
WHERE  sgmn.group_id = $1::uuid;`

	GroupBucketID = `
SELECT sg.bucket_id
FROM   soma.groups sg
WHERE  sg.group_id = $1;`

	GroupOncProps = `
SELECT op.instance_id,
       op.source_instance_id,
       op.view,
       op.oncall_duty_id,
       iot.name
FROM   soma.group_oncall_properties op
JOIN   inventory.oncall_team iot
  ON   op.oncall_duty_id = iot.id
WHERE  op.group_id = $1::uuid;`

	GroupSvcProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.service_id
FROM   soma.group_service_property sp
WHERE  sp.group_id = $1::uuid;`

	GroupSysProps = `
SELECT sp.instance_id,
       sp.source_instance_id,
       sp.view,
       sp.system_property,
       sp.value
FROM   soma.group_system_properties sp
WHERE  sp.group_id = $1::uuid;`

	GroupCstProps = `
SELECT cp.instance_id,
       cp.source_instance_id,
       cp.view,
       cp.custom_property_id,
       cp.value,
       scp.custom_property
FROM   soma.group_custom_properties cp
JOIN   soma.custom_properties scp
  ON   cp.custom_property_id = scp.custom_property_id
WHERE  cp.group_id = $1::uuid;`

	GroupSystemPropertyForDelete = `
SELECT view,
       system_property,
       value
FROM   soma.group_system_properties
WHERE  source_instance_id = $1::uuid
  AND  source_instance_id = instance_id;`

	GroupCustomPropertyForDelete = `
SELECT sgcp.view,
       sgcp.custom_property_id,
       sgcp.value,
       scp.custom_property
FROM   soma.group_custom_properties sgcp
JOIN   soma.custom_properties scp
  ON   sgcp.repository_id = scp.repository_id
 AND   sgcp.custom_property_id = scp.custom_property_id
WHERE  sgcp.source_instance_id = $1::uuid
  AND  sgcp.source_instance_id = sgcp.instance_id;`

	GroupOncallPropertyForDelete = `
SELECT sgop.view,
       sgop.oncall_duty_id,
       iot.name,
       iot.phone_number
FROM   soma.group_oncall_properties sgop
JOIN   inventory.oncall_team iot
  ON   sgop.oncall_duty_id = iot.id
WHERE  sgop.source_instance_id = $1::uuid
  AND  sgop.source_instance_id = sgop.instance_id;`

	GroupServicePropertyForDelete = `
SELECT sgsp.view,
       sgsp.service_id
FROM   soma.group_service_property sgsp
JOIN   soma.service_property ssp
  ON   sgsp.team_id = ssp.team_id
 AND   sgsp.service_id = ssp.id
WHERE  sgsp.source_instance_id = $1::uuid
  AND  sgsp.source_instance_id = sgsp.instance_id;`
)

func init() {
	m[GroupBucketID] = `GroupBucketID`
	m[GroupCstProps] = `GroupCstProps`
	m[GroupCustomPropertyForDelete] = `GroupCustomPropertyForDelete`
	m[GroupMemberClusterList] = `GroupMemberClusterList`
	m[GroupMemberGroupList] = `GroupMemberGroupList`
	m[GroupMemberNodeList] = `GroupMemberNodeList`
	m[GroupOncProps] = `GroupOncProps`
	m[GroupOncallPropertyForDelete] = `GroupOncallPropertyForDelete`
	m[GroupServicePropertyForDelete] = `GroupServicePropertyForDelete`
	m[GroupShow] = `GroupShow`
	m[GroupSvcProps] = `GroupSvcProps`
	m[GroupSysProps] = `GroupSysProps`
	m[GroupSystemPropertyForDelete] = `GroupSystemPropertyForDelete`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
