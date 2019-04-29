/*-
 * Copyright (c) 2016,2018, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß <code.jpe@gmail.com>
 * All rights reserved
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package stmt

import (
	"github.com/mjolnir42/soma/lib/proto"
)

const (
	LifecycleStatements = ``

	LifecycleActiveUnblockCondition = `
SELECT  scicd.blocked_instance_config_id,
        scicd.blocking_instance_config_id,
        scicd.unblocking_state,
        p.status,
        p.next_status,
        p.check_instance_id
FROM    soma.check_instance_configuration_dependencies scicd
JOIN    soma.check_instance_configurations scic
ON      scicd.blocking_instance_config_id = scic.check_instance_config_id
AND     scicd.unblocking_state = scic.status
JOIN    soma.check_instance_configurations p
ON      scicd.blocked_instance_config_id = p.check_instance_config_id
JOIN    soma.check_instances sci
ON      p.check_instance_id = sci.check_instance_id
AND     scicd.blocking_instance_config_id = sci.current_instance_config_id;`

	LifecycleUpdateInstance = `
UPDATE  soma.check_instances
SET     update_available = $1::boolean,
        current_instance_config_id = $2::uuid
WHERE   check_instance_id = $3::uuid;`

	LifecycleUpdateConfig = `
UPDATE  soma.check_instance_configurations
SET     status = $1::varchar,
        next_status = $2::varchar,
        awaiting_deletion = $3::boolean,
        status_last_updated_at = NOW()::timestamptz
WHERE   check_instance_config_id = $4::uuid;`

	LifecycleDeleteDependency = `
DELETE FROM soma.check_instance_configuration_dependencies
WHERE       blocked_instance_config_id = $1::uuid
AND         blocking_instance_config_id = $2::uuid
AND         unblocking_state = $3::varchar;`

	LifecycleReadyDeployments = `
SELECT scic.check_instance_id,
       scic.monitoring_id,
       sms.monitoring_callback_uri
FROM   soma.check_instance_configurations scic
JOIN   soma.monitoring_systems sms
ON     scic.monitoring_id = sms.monitoring_id
JOIN   soma.check_instances sci
ON     scic.check_instance_id = sci.check_instance_id
AND    scic.check_instance_config_id = sci.current_instance_config_id
WHERE  (  scic.status = '` + proto.DeploymentAwaitingRollout + `'::varchar
       OR scic.status = '` + proto.DeploymentAwaitingDeprovision + `'::varchar)
AND    sms.monitoring_callback_uri IS NOT NULL
AND    sci.update_available;`

	LifecycleRescheduleDeployments = `
SELECT scic.check_instance_id,
       scic.monitoring_id,
       sms.monitoring_callback_uri
FROM   soma.check_instance_configurations scic
JOIN   soma.monitoring_systems sms
ON     scic.monitoring_id = sms.monitoring_id
JOIN   soma.check_instances sci
ON     scic.check_instance_id = sci.check_instance_id
AND    scic.check_instance_config_id = sci.current_instance_config_id
WHERE  (  scic.status = '` + proto.DeploymentRolloutInProgress + `'::varchar
       OR scic.status = '` + proto.DeploymentAwaitingRollout + `'::varchar
       OR scic.status = '` + proto.DeploymentAwaitingDeprovision + `'::varchar
       OR scic.status = '` + proto.DeploymentDeprovisionInProgress + `'::varchar)
AND    sms.monitoring_callback_uri IS NOT NULL
AND    scic.status_last_updated_at IS NOT NULL
AND    scic.notified_at IS NOT NULL
AND    NOW() > (scic.status_last_updated_at + '5 minute'::interval)
AND    NOW() > (scic.notified_at + '5 minute'::interval)
AND    NOT sci.update_available;`

	LifecycleSetNotified = `
UPDATE soma.check_instance_configurations scic
SET    notified_at = NOW()::timestamptz
FROM   soma.check_instances sci
WHERE  sci.current_instance_config_id = scic.check_instance_config_id
  AND  sci.check_instance_id = $1::uuid;`

	LifecycleClearUpdateFlag = `
UPDATE soma.check_instances
SET    update_available = 'false'::boolean
WHERE  check_instance_id = $1::uuid;`

	LifecycleBlockedConfigsForDeletedInstance = `
SELECT scicd.blocked_instance_config_id,
       scicd.blocking_instance_config_id,
       scicd.unblocking_state
FROM   soma.check_instances sci
JOIN   soma.check_instance_configurations scic
  ON   sci.check_instance_id = scic.check_instance_id
JOIN   soma.check_instance_configuration_dependencies scicd
  ON   scic.check_instance_config_id = scicd.blocked_instance_config_id
WHERE  sci.deleted
  AND  scic.status = '` + proto.DeploymentBlocked + `'::varchar;`

	LifecycleConfigAwaitingDeletion = `
UPDATE soma.check_instance_configurations
SET    status = '` + proto.DeploymentAwaitingDeletion + `'::varchar,
       next_status = '` + proto.DeploymentNone + `'::varchar,
       awaiting_deletion = 'yes'::boolean
WHERE  check_instance_config_id = $1::uuid;`

	LifecycleDeleteGhosts = `
UPDATE soma.check_instance_configurations scic
SET    status = '` + proto.DeploymentAwaitingDeletion + `'::varchar,
       next_status = '` + proto.DeploymentNone + `'::varchar,
       awaiting_deletion = 'yes'::boolean
FROM   soma.check_instances sci
WHERE  scic.check_instance_id = sci.check_instance_id
AND    scic.status = '` + proto.DeploymentAwaitingRollout + `'::varchar
AND    sci.deleted
AND    sci.update_available;`

	LifecycleDeleteFailedRollouts = `
UPDATE soma.check_instance_configurations scic
SET    status = '` + proto.DeploymentAwaitingDeletion + `'::varchar,
       next_status = '` + proto.DeploymentNone + `'::varchar,
       awaiting_deletion = 'yes'::boolean
FROM   soma.check_instances sci
WHERE  scic.check_instance_id = sci.check_instance_id
AND    sci.deleted
AND    scic.status = '` + proto.DeploymentRolloutFailed + `'::varchar;`

	LifecycleDeleteDeprovisioned = `
UPDATE soma.check_instance_configurations scic
SET    status = '` + proto.DeploymentAwaitingDeletion + `'::varchar,
       next_status = '` + proto.DeploymentNone + `'::varchar,
       awaiting_deletion = 'yes'::boolean
FROM   soma.check_instances sci
WHERE  scic.check_instance_id = sci.check_instance_id
AND    sci.deleted
AND    scic.status = '` + proto.DeploymentDeprovisioned + `'::varchar
AND    scic.next_status = '` + proto.DeploymentNone + `'::varchar;`

	LifecycleDeleteOrphanCheckInstances = `
UPDATE soma.check_instances sci
SET    deleted = 'yes'::boolean
FROM   soma.check_configurations sck
WHERE  sci.check_configuration_id = sck.configuration_id
AND NOT sci.deleted
AND sck.deleted;`

	LifecycleDeprovisionDeletedActive = `
SELECT scic.check_instance_config_id,
       sci.check_instance_id
FROM   soma.check_instance_configurations scic
JOIN   soma.check_instances sci
  ON   scic.check_instance_id = sci.check_instance_id
WHERE  sci.deleted
  AND  scic.status = '` + proto.DeploymentActive + `'::varchar
  AND  scic.next_status = '` + proto.DeploymentNone + `'::varchar;`

	LifecycleDeprovisionConfiguration = `
UPDATE soma.check_instance_configurations
SET    status = '` + proto.DeploymentAwaitingDeprovision + `'::varchar,
       next_status = '` + proto.DeploymentDeprovisionInProgress + `'::varchar
WHERE  check_instance_config_id = $1::uuid;`

	LifecycleDeadLockResolver = `
SELECT ci.check_instance_id,
       ci.current_instance_config_id
FROM   check_instances ci
JOIN   check_instance_configurations cic
  ON   ci.check_instance_id = cic.check_instance_id
 AND   ci.current_instance_config_id = cic.check_instance_config_id
JOIN   check_instance_configuration_dependencies cicd
  ON   ci.current_instance_config_id = cicd.blocking_instance_config_id
WHERE  cic.status = '` + proto.DeploymentActive + `'::varchar;`
)

func init() {
	m[LifecycleActiveUnblockCondition] = `LifecycleActiveUnblockCondition`
	m[LifecycleBlockedConfigsForDeletedInstance] = `LifecycleBlockedConfigsForDeletedInstance`
	m[LifecycleClearUpdateFlag] = `LifecycleClearUpdateFlag`
	m[LifecycleConfigAwaitingDeletion] = `LifecycleConfigAwaitingDeletion`
	m[LifecycleDeadLockResolver] = `LifecycleDeadLockResolver`
	m[LifecycleDeleteDependency] = `LifecycleDeleteDependency`
	m[LifecycleDeleteDeprovisioned] = `LifecycleDeleteDeprovisioned`
	m[LifecycleDeleteOrphanCheckInstances] = `LifecycleDeleteOrphanCheckInstances`
	m[LifecycleDeleteFailedRollouts] = `LifecycleDeleteFailedRollouts`
	m[LifecycleDeleteGhosts] = `LifecycleDeleteGhosts`
	m[LifecycleDeprovisionConfiguration] = `LifecycleDeprovisionConfiguration`
	m[LifecycleDeprovisionDeletedActive] = `LifecycleDeprovisionDeletedActive`
	m[LifecycleReadyDeployments] = `LifecycleReadyDeployments`
	m[LifecycleRescheduleDeployments] = `LifecycleRescheduleDeployments`
	m[LifecycleSetNotified] = `LifecycleSetNotified`
	m[LifecycleUpdateConfig] = `LifecycleUpdateConfig`
	m[LifecycleUpdateInstance] = `LifecycleUpdateInstance`
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
