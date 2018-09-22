/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016-2018, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import (
	"sync"

	"github.com/mjolnir42/soma/internal/msg"
	uuid "github.com/satori/go.uuid"
)

func (tec *Cluster) evalNativeProp(prop string, val string) bool {
	switch prop {
	case msg.NativePropertyEnvironment:
		env := tec.Parent.(Bucketeer).GetEnvironment()
		if val == env {
			return true
		}
	case msg.NativePropertyEntity:
		if val == msg.EntityCluster {
			return true
		}
	case msg.NativePropertyState:
		if val == tec.State {
			return true
		}
	case msg.NativePropertyHardwareNode:
		// cluster != hardware
		return false
	}
	return false
}

func (tec *Cluster) evalSystemProp(prop string, val string, view string) (string, bool, string) {
	for _, v := range tec.PropertySystem {
		t := v.(*PropertySystem)
		if t.Key == prop && (t.Value == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Key, true, t.Value
		}
	}
	return "", false, ""
}

func (tec *Cluster) evalOncallProp(prop string, val string, view string) (string, bool) {
	for _, v := range tec.PropertyOncall {
		t := v.(*PropertyOncall)
		if "OncallID" == prop && t.ID.String() == val && (t.View == view || t.View == `any`) {
			return t.ID.String(), true
		}
	}
	return "", false
}

func (tec *Cluster) evalCustomProp(prop string, val string, view string) (string, bool, string) {
	for _, v := range tec.PropertyCustom {
		t := v.(*PropertyCustom)
		if t.Key == prop && (t.Value == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Key, true, t.Value
		}
	}
	return "", false, ""
}

func (tec *Cluster) evalServiceProp(prop string, val string, view string) (string, bool, string) {
	for _, v := range tec.PropertyService {
		t := v.(*PropertyService)
		if prop == "name" && (t.Service == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.ID.String(), true, t.Service
		}
	}
	return "", false, ""
}

func (tec *Cluster) evalAttributeOfService(svcID string, view string, attribute string, value string) (bool, string) {
	t := tec.PropertyService[svcID].(*PropertyService)
	for _, a := range t.Attributes {
		if a.Name == attribute && (t.View == view || t.View == `any`) && (a.Value == value || value == `@defined`) {
			return true, a.Value
		}
	}
	return false, ""
}

func (tec *Cluster) evalAttributeProp(view string, attr string, value string) (bool, map[string]string) {
	f := map[string]string{}
svcloop:
	for _, v := range tec.PropertyService {
		t := v.(*PropertyService)
		for _, a := range t.Attributes {
			if a.Name == attr && (a.Value == value || value == `@defined`) && (t.View == view || t.View == `any`) {
				f[t.ID.String()] = a.Value
				continue svcloop
			}
		}
	}
	if len(f) > 0 {
		return true, f
	}
	return false, f
}

func (tec *Cluster) getServiceMap(serviceID string) map[string][]string {
	svc := new(PropertyService)
	svc = tec.PropertyService[serviceID].(*PropertyService)

	res := map[string][]string{}
	for _, v := range svc.Attributes {
		res[v.Name] = append(res[v.Name], v.Value)
	}
	return res
}

func (c *Cluster) updateCheckInstances() {
	// object may have no checks, but there could be instances to mop up
	if len(c.Checks) == 0 && len(c.Instances) == 0 {
		c.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, HasChecks=%t",
			c.GetRepositoryName(),
			`UpdateCheckInstances`,
			`cluster`,
			c.ID.String(),
			false,
		)
		// found nothing to do, ensure update flag is unset again
		c.hasUpdate = false
		return
	}

	// if there are loaded instances, then this is the initial rebuild
	// of the tree
	startup := false
	if len(c.loadedInstances) > 0 {
		startup = true
	}

	// if this is not the startupLoad and there are no updates, then there
	// is noting to do
	if !startup && !c.hasUpdate {
		return
	}

	c.deleteOrphanCheckInstances()

	c.removeDisabledCheckInstances()

	c.calculateCheckInstances(startup)
}

func (c *Cluster) deleteOrphanCheckInstances() {
	c.lock.Lock()
	defer c.lock.Unlock()
	// scan over all current checkinstances if their check still exists.
	// If not the check has been deleted and the spawned instances need
	// a good deletion
	for ck := range c.CheckInstances {
		if _, ok := c.Checks[ck]; ok {
			// check still exists
			continue
		}

		// check no longer exists -> cleanup
		inst := c.CheckInstances[ck]
		for _, i := range inst {
			c.actionCheckInstanceDelete(c.Instances[i].MakeAction())
			c.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
				c.GetRepositoryName(),
				`CleanupInstance`,
				`group`,
				c.ID.String(),
				ck,
				i,
			)
			delete(c.Instances, i)
		}
		delete(c.CheckInstances, ck)
	}
}

func (c *Cluster) removeDisabledCheckInstances() {
	c.lock.Lock()
	defer c.lock.Unlock()
	// loop over all checks and test if there is a reason to disable
	// its check instances. And with disable we mean delete.
	for chk := range c.Checks {
		disableThis := false
		// disable this check if the system property
		// `disable_all_monitoring` is set for the view that the check
		// uses.
		if _, hit, _ := c.evalSystemProp(
			msg.SystemPropertyDisableAllMonitoring,
			`true`,
			c.Checks[chk].View,
		); hit {
			disableThis = true
		}
		// disable this check if the system property
		// `disable_check_configuration` is set to the
		// check_configuration that spawned this check
		if _, hit, _ := c.evalSystemProp(
			msg.SystemPropertyDisableCheckConfiguration,
			c.Checks[chk].ConfigID.String(),
			c.Checks[chk].View,
		); hit {
			disableThis = true
		}
		// if there was a reason to disable this check, all instances
		// are deleted
		if disableThis {
			if instanceArray, ok := c.CheckInstances[chk]; ok {
				for _, i := range instanceArray {
					c.actionCheckInstanceDelete(c.Instances[i].MakeAction())
					c.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
						c.GetRepositoryName(),
						`RemoveDisabledInstance`,
						`group`,
						c.ID.String(),
						chk,
						i,
					)
					delete(c.Instances, i)
				}
				delete(c.CheckInstances, chk)
			}
		}
	}
}

func (c *Cluster) calculateCheckInstances(startup bool) {
	wg := sync.WaitGroup{}
	for i := range c.Checks {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			c.processCheckForUpdates(name, startup)
		}(i)
	}
	wg.Wait()

	// completed the pass, reset update flag
	c.hasUpdate = false
}

func (c *Cluster) processCheckForUpdates(chkName string, startup bool) {
	c.lock.RLock()
	if c.Checks[chkName].Inherited == false && c.Checks[chkName].ChildrenOnly == true {
		// not active here
		c.lock.RUnlock()
		return
	}
	if c.Checks[chkName].View == msg.ViewLocal {
		// groups have no local view
		c.lock.RUnlock()
		return
	}
	if _, hit, _ := c.evalSystemProp(
		// skip check if `disable_all_monitoring` property is set
		msg.SystemPropertyDisableAllMonitoring,
		`true`,
		c.Checks[chkName].View,
	); hit {
		c.lock.RUnlock()
		return
	}
	if _, hit, _ := c.evalSystemProp(
		// skip check if `disable_check_configuration` property is set
		msg.SystemPropertyDisableCheckConfiguration,
		c.Checks[chkName].ConfigID.String(),
		c.Checks[chkName].View,
	); hit {
		c.lock.RUnlock()
		return
	}

	ctx := newCheckContext(chkName, c.Checks[chkName].View, startup)
	c.lock.RUnlock()

	c.constraintCheck(ctx)
	c.log.Printf(
		"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, Match=%t",
		c.GetRepositoryName(), `ConstraintEvaluation`, `group`,
		c.ID.String(), chkName, ctx.brokeConstraint,
	)
	if ctx.brokeConstraint {
		return
	}

	// check triggered, create instances
	switch {
	case !ctx.hasServiceConstraint:
		/* if there are no service constraints, one check instance is
		 * created for this check
		 */
		c.createNoServiceCheckInstance(ctx)
	default:
		/* if service constraints are in effect, then we generate
		 * instances for every service that bound.
		 * Since service attributes can be specified more than once,
		 * but the semantics are unclear what the expected behaviour of
		 * for example a file age check is that is specified against
		 * more than one file path; all possible attribute value
		 * permutations for each service are built and then one check
		 * instance is built for each of these service config
		 * permutations.
		 */
		c.createPerServiceCheckInstances(ctx)
	}

	if ctx.startupBroken {
		return
	}

	// all new check instances have been built, check which
	// existing instances did not get an update and need to be
	// deleted
	if !ctx.startup {
		c.pruneOldCheckInstances(ctx)
		c.dispatchCheckInstanceUpdates(ctx)
	}

	c.createNewCheckInstances(ctx)

	delete(c.CheckInstances, ctx.uuid)
	c.CheckInstances[ctx.uuid] = ctx.newCheckInstances
}

func (c *Cluster) constraintCheck(ctx *checkContext) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	// these constaint types must always match for the instance to
	// be valid. defer service and attribute
	for _, cc := range c.Checks[ctx.uuid].Constraints {
		switch cc.Type {
		case msg.ConstraintNative:
			if c.evalNativeProp(cc.Key, cc.Value) {
				ctx.nativeConstr[cc.Key] = cc.Value
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintSystem:
			if id, hit, bind := c.evalSystemProp(cc.Key, cc.Value, ctx.view); hit {
				ctx.systemConstr[id] = bind
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintOncall:
			if id, hit := c.evalOncallProp(cc.Key, cc.Value, ctx.view); hit {
				ctx.oncallConstr = id
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintCustom:
			if id, hit, bind := c.evalCustomProp(cc.Key, cc.Value, ctx.view); hit {
				ctx.customConstr[id] = bind
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintService:
			ctx.hasServiceConstraint = true
			if id, hit, bind := c.evalServiceProp(cc.Key, cc.Value, ctx.view); hit {
				ctx.serviceConstr[id] = bind
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintAttribute:
			ctx.hasAttributeConstraint = true
			ctx.attributes = append(ctx.attributes, cc)
		}
	}

	switch {
	case ctx.hasServiceConstraint && ctx.hasAttributeConstraint:
		/* if the check has both service and attribute constraints,
		* then for the check to hit, the tree element needs to have
		* all the services, and each of them needs to match all
		* attribute constraints
		 */
		for id := range ctx.serviceConstr {
			for _, attr := range ctx.attributes {
				hit, bind := c.evalAttributeOfService(id, ctx.view, attr.Key, attr.Value)
				if hit {
					// attributeC[id] might still be a nil map
					if ctx.attributeConstr[id] == nil {
						ctx.attributeConstr[id] = make(map[string][]string)
					}
					ctx.attributeConstr[id][attr.Key] = append(
						ctx.attributeConstr[id][attr.Key],
						bind,
					)
					continue
				}
				ctx.brokeConstraint = true
				return
			}
		}
	case ctx.hasAttributeConstraint && !ctx.hasServiceConstraint:
		/* if the check has only attribute constraints and no
		* service constraint, then we pull in every service that
		* matches all attribute constraints and generate a check
		* instance for it
		 */
		attrCount := len(ctx.attributes)
		for _, attr := range ctx.attributes {
			if hit, svcIDMap := c.evalAttributeProp(ctx.view, attr.Key, attr.Value); hit {
				for id, bind := range svcIDMap {
					ctx.serviceConstr[id] = svcIDMap[id]
					// attributeC[id] might still be a nil map
					if ctx.attributeConstr[id] == nil {
						ctx.attributeConstr[id] = make(map[string][]string)
					}
					ctx.attributeConstr[id][attr.Key] = append(
						ctx.attributeConstr[id][attr.Key],
						bind,
					)
				}
			}
		}
		// delete all services that did not match all attributes
		//
		// if a check has two attribute constraints on the same
		// attribute, then len(attributeC[id]) != len(attributes)
		for id := range ctx.attributeConstr {
			if countAttributeConstraints(ctx.attributeConstr[id]) != attrCount {
				delete(ctx.serviceConstr, id)
				delete(ctx.attributeConstr, id)
			}
		}
		// declare service constraints in effect if we found a
		// service that bound all attribute constraints
		switch len(ctx.serviceConstr) {
		case 0:
			// found no services that fulfilled all constraints
			ctx.brokeConstraint = true
		default:
			ctx.hasServiceConstraint = true
		}
	}
}

func (c *Cluster) createNoServiceCheckInstance(ctx *checkContext) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	inst := CheckInstance{
		InstanceID: uuid.UUID{},
		CheckID: func(id string) uuid.UUID {
			f, _ := uuid.FromString(id)
			return f
		}(ctx.uuid),
		ConfigID: func(id string) uuid.UUID {
			f, _ := uuid.FromString(c.Checks[id].ConfigID.String())
			return f
		}(ctx.uuid),
		InstanceConfigID:      uuid.Must(uuid.NewV4()),
		ConstraintOncall:      ctx.oncallConstr,
		ConstraintService:     ctx.serviceConstr,
		ConstraintSystem:      ctx.systemConstr,
		ConstraintCustom:      ctx.customConstr,
		ConstraintNative:      ctx.nativeConstr,
		ConstraintAttribute:   ctx.attributeConstr,
		InstanceService:       ``,
		InstanceServiceConfig: nil,
		InstanceSvcCfgHash:    ``,
	}
	inst.calcConstraintHash()
	inst.calcConstraintValHash()

	switch ctx.startup {
	case true:
		// upgrade to writelock
		c.lock.RUnlock()
		c.lock.Lock()
		matched := false

		for loadedID, loadedInst := range c.loadedInstances[ctx.uuid] {
			if loadedInst.InstanceSvcCfgHash != `` {
				continue
			}
			// check if an instance exists bound against the
			// same constraints
			if inst.MatchConstraints(&loadedInst) {
				// found a match
				matched = true

				inst.InstanceID, _ = uuid.FromString(loadedID)
				inst.InstanceConfigID, _ = uuid.FromString(
					loadedInst.InstanceConfigID.String(),
				)
				inst.Version = loadedInst.Version
				delete(c.loadedInstances[ctx.uuid], loadedID)
				c.log.Printf(
					"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, "+
						"InstanceID=%s, ServiceConstrained=%t", c.GetRepositoryName(),
					`ComputeInstance`, `cluster`, c.ID.String(), ctx.uuid,
					loadedID, false,
				)
				break
			}
		}
		if !matched {
			// downgrade to readlock
			c.lock.Unlock()
			c.lock.RLock()

			c.log.Printf("TK[%s]: Failed to match computed instance to loaded instances."+
				" ObjType=%s, ObjId=%s, CheckID=%s", `cluster`, c.ID.String(), ctx.uuid,
				c.GetRepositoryName())
			c.Fault.Error <- &Error{
				Action: `Failed to match a computed instance to loaded data`,
			}
			ctx.startupBroken = true
			return
		}
		// downgrade to readlock
		c.lock.Unlock()
		c.lock.RLock()
	default:
		for _, existingInstanceID := range c.CheckInstances[ctx.uuid] {
			existingInstance := c.Instances[existingInstanceID]

			// ignore instances with service constraints
			if existingInstance.InstanceSvcCfgHash != inst.InstanceSvcCfgHash {
				continue
			}

			// check if an instance exists bound against the same constraints
			if existingInstance.ConstraintHash == inst.ConstraintHash {
				inst.InstanceID, _ = uuid.FromString(
					existingInstance.InstanceID.String(),
				)
				inst.Version = existingInstance.Version + 1
				break
			}
		}
		if uuid.Equal(uuid.Nil, inst.InstanceID) {
			// no match was found during nosvcinstanceloop, this
			// is a new instance
			inst.Version = 0
			inst.InstanceID = uuid.Must(uuid.NewV4())
		}
		c.log.Printf(
			"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, "+
				"InstanceID=%s, ServiceConstrained=%t", c.GetRepositoryName(),
			`ComputeInstance`, `cluster`, c.ID.String(), ctx.uuid, inst.InstanceID.String(),
			false,
		)
	}

	ctx.newInstances[inst.InstanceID.String()] = inst
	ctx.newCheckInstances = append(ctx.newCheckInstances, inst.InstanceID.String())
}

func (c *Cluster) createPerServiceCheckInstances(ctx *checkContext) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for svcID := range ctx.serviceConstr {
		svcCfg := c.getServiceMap(svcID)

		// calculate how many instances this service spawns
		combinations := 1
		for attr := range svcCfg {
			combinations = combinations * len(svcCfg[attr])
		}

		// build all attribute combinations
		results := make([]map[string]string, 0, combinations)
		for attr := range svcCfg {
			if len(results) == 0 {
				for i := range svcCfg[attr] {
					res := map[string]string{}
					res[attr] = svcCfg[attr][i]
					results = append(results, res)
				}
				continue
			}
			ires := make([]map[string]string, 0, combinations)
			for r := range results {
				for j := range svcCfg[attr] {
					res := map[string]string{}
					for k, v := range results[r] {
						res[k] = v
					}
					res[attr] = svcCfg[attr][j]
					ires = append(ires, res)
				}
			}
			results = ires
		}
		// build a CheckInstance for every result
		for _, y := range results {
			// ensure we have a full copy and not a header copy
			cfg := map[string]string{}
			for k, v := range y {
				cfg[k] = v
			}
			inst := CheckInstance{
				InstanceID: uuid.UUID{},
				CheckID: func(id string) uuid.UUID {
					f, _ := uuid.FromString(id)
					return f
				}(ctx.uuid),
				ConfigID: func(id string) uuid.UUID {
					f, _ := uuid.FromString(c.Checks[id].ConfigID.String())
					return f
				}(ctx.uuid),
				InstanceConfigID:      uuid.Must(uuid.NewV4()),
				ConstraintOncall:      ctx.oncallConstr,
				ConstraintService:     ctx.serviceConstr,
				ConstraintSystem:      ctx.systemConstr,
				ConstraintCustom:      ctx.customConstr,
				ConstraintNative:      ctx.nativeConstr,
				ConstraintAttribute:   ctx.attributeConstr,
				InstanceService:       svcID,
				InstanceServiceConfig: cfg,
			}
			inst.calcConstraintHash()
			inst.calcConstraintValHash()
			inst.calcInstanceSvcCfgHash()

			switch ctx.startup {
			case true:
				// upgrade to writelock
				c.lock.RUnlock()
				c.lock.Lock()
				matched := false

				for ldInstID, ldInst := range c.loadedInstances[ctx.uuid] {
					// check for data from loaded instance
					if inst.MatchServiceConstraints(&ldInst) {
						// found a match
						matched = true

						inst.InstanceID, _ = uuid.FromString(ldInstID)
						inst.InstanceConfigID, _ = uuid.FromString(ldInst.InstanceConfigID.String())
						inst.Version = ldInst.Version
						// we can assume InstanceServiceConfig to
						// be equal, since InstanceSvcCfgHash is
						// equal
						delete(c.loadedInstances[ctx.uuid], ldInstID)
						c.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s,"+
							"CheckID=%s, InstanceID=%s, ServiceConstrained=%t",
							c.GetRepositoryName(), `ComputeInstance`, `cluster`,
							c.ID.String(), ctx.uuid, ldInstID, true,
						)
						break
					}
				}
				if !matched {
					// downgrade to readlock
					c.lock.Unlock()
					c.lock.RLock()

					c.log.Printf(
						"TK[%s]: Failed to match computed instance to loaded instances."+
							" ObjType=%s, ObjId=%s, CheckID=%s", `cluster`, c.ID.String(),
						ctx.uuid, c.GetRepositoryName())
					c.Fault.Error <- &Error{
						Action: `Failed to match a computed instance to loaded data`,
					}
					ctx.startupBroken = true
					return
				}
				// downgrade to readlock
				c.lock.Unlock()
				c.lock.RLock()
			default:
				// lookup existing instance ids for check in teg.CheckInstances
				// to determine if this is an update
				for _, exInstID := range c.CheckInstances[ctx.uuid] {
					exInst := c.Instances[exInstID]
					// this existing instance is for the same service
					// configuration -> this is an update
					if exInst.InstanceSvcCfgHash == inst.InstanceSvcCfgHash {
						inst.InstanceID, _ = uuid.FromString(exInst.InstanceID.String())
						inst.Version = exInst.Version + 1
						break
					}
				}
				if uuid.Equal(uuid.Nil, inst.InstanceID) {
					// no match was found during instanceloop, this is
					// a new instance
					inst.Version = 0
					inst.InstanceID = uuid.Must(uuid.NewV4())
				}
				c.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, "+
					"InstanceID=%s, ServiceConstrained=%t", c.GetRepositoryName(),
					`ComputeInstance`, `cluster`, c.ID.String(), ctx.uuid,
					inst.InstanceID.String(), true,
				)
			}

			ctx.newInstances[inst.InstanceID.String()] = inst
			ctx.newCheckInstances = append(
				ctx.newCheckInstances,
				inst.InstanceID.String(),
			)
		}
	}
	if ctx.startup {
		// all instances have been built and matched to
		// loaded instances, but there are loaded
		// instances left. why?
		if len(c.loadedInstances[ctx.uuid]) != 0 {
			c.Fault.Error <- &Error{
				Action: "Leftover matched instances after assignment, " +
					"computed instances missing",
			}
			ctx.startupBroken = true
		}
	}
}

func (c *Cluster) pruneOldCheckInstances(ctx *checkContext) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, oldInstanceID := range c.CheckInstances[ctx.uuid] {
		if _, ok := ctx.newInstances[oldInstanceID]; !ok {
			// there is no new version for oldInstanceID
			c.actionCheckInstanceDelete(c.Instances[oldInstanceID].MakeAction())
			c.log.Printf(
				"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
				c.GetRepositoryName(), `DeleteInstance`, `cluster`, c.ID.String(),
				ctx.uuid, oldInstanceID,
			)
			delete(c.Instances, oldInstanceID)
		}
	}
}

func (c *Cluster) dispatchCheckInstanceUpdates(ctx *checkContext) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, oldInstanceID := range c.CheckInstances[ctx.uuid] {
		delete(c.Instances, oldInstanceID)
		c.Instances[oldInstanceID] = ctx.newInstances[oldInstanceID]
		c.actionCheckInstanceUpdate(c.Instances[oldInstanceID].MakeAction())
		c.log.Printf(
			"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
			c.GetRepositoryName(), `UpdateInstance`, `cluster`, c.ID.String(),
			ctx.uuid, oldInstanceID,
		)
	}
}

func (c *Cluster) createNewCheckInstances(ctx *checkContext) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, newInstanceID := range ctx.newCheckInstances {
		if _, ok := c.Instances[newInstanceID]; !ok {
			// this instance is new, not an update
			c.Instances[newInstanceID] = ctx.newInstances[newInstanceID]

			action := `CreateInstance`
			switch ctx.startup {
			case true:
				action = `RecreateInstance`
			default:
				c.actionCheckInstanceCreate(c.Instances[newInstanceID].MakeAction())
			}
			c.log.Printf(
				"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
				c.GetRepositoryName(), action, `cluster`, c.ID.String(),
				ctx.uuid, newInstanceID,
			)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
