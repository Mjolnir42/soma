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

func (n *Node) evalNativeProp(prop string, val string) bool {
	switch prop {
	case msg.NativePropertyEnvironment:
		env := n.Parent.(Bucketeer).GetEnvironment()
		if val == env {
			return true
		}
	case msg.NativePropertyEntity:
		if val == msg.EntityNode {
			return true
		}
	case msg.NativePropertyState:
		if val == n.State {
			return true
		}
	case msg.NativePropertyHardwareNode:
		// XX needs n.ServerName extension of ten
		// if val == n.ServerName { return true }
		return false
	}
	return false
}

func (n *Node) evalSystemProp(prop string, val string, view string) (string, bool, string) {
	for _, v := range n.PropertySystem {
		t := v.(*PropertySystem)
		if t.Key == prop && (t.Value == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Key, true, t.Value
		}
	}
	return "", false, ""
}

func (n *Node) evalOncallProp(prop string, val string, view string) (string, bool) {
	for _, v := range n.PropertyOncall {
		t := v.(*PropertyOncall)
		if "OncallID" == prop && t.ID.String() == val && (t.View == view || t.View == `any`) {
			return t.ID.String(), true
		}
	}
	return "", false
}

func (n *Node) evalCustomProp(prop string, val string, view string) (string, bool, string) {
	for _, v := range n.PropertyCustom {
		t := v.(*PropertyCustom)
		if t.Key == prop && (t.Value == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Key, true, t.Value
		}
	}
	return "", false, ""
}

func (n *Node) evalServiceProp(prop string, val string, view string) (string, bool, string) {
	for _, v := range n.PropertyService {
		t := v.(*PropertyService)
		if ((prop == "name" && (t.ServiceName == val || val == `@defined`)) || (prop == "id" && (t.ServiceID.String() == val))) && (t.View == view || t.View == `any`) {
			return t.ID.String(), true, t.ServiceName
		}
	}
	return "", false, ""
}

func (n *Node) evalAttributeOfService(svcID string, view string, attribute string, value string) (bool, string) {
	t := n.PropertyService[svcID].(*PropertyService)
	for _, a := range t.Attributes {
		if a.Name == attribute && (t.View == view || t.View == `any`) && (a.Value == value || value == `@defined`) {
			return true, a.Value
		}
	}
	return false, ""
}

func (n *Node) evalAttributeProp(view string, attr string, value string) (bool, map[string]string) {
	f := map[string]string{}
svcloop:
	for _, v := range n.PropertyService {
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

func (n *Node) getServiceMap(serviceID string) map[string][]string {
	svc := new(PropertyService)
	svc = n.PropertyService[serviceID].(*PropertyService)

	res := map[string][]string{}
	for _, v := range svc.Attributes {
		res[v.Name] = append(res[v.Name], v.Value)
	}
	return res
}

func (n *Node) updateCheckInstances() {
	// object may have no checks, but there could be instances to mop up
	if len(n.Checks) == 0 && len(n.Instances) == 0 {
		n.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, HasChecks=%t",
			n.GetRepositoryName(), `UpdateCheckInstances`, `node`,
			n.ID.String(), false,
		)
		// found nothing to do, ensure update flag is unset again
		n.hasUpdate = false
		return
	}

	// if there are loaded instances, then this is the initial rebuild
	// of the tree
	startup := false
	if len(n.loadedInstances) > 0 {
		startup = true
	}

	// if this is not the startupLoad and there are no updates, then there
	// is noting to do
	if !startup && !n.hasUpdate {
		return
	}

	n.deleteOrphanCheckInstances()

	n.removeDisabledCheckInstances()

	n.calculateCheckInstances(startup)
}

func (n *Node) deleteOrphanCheckInstances() {
	n.lock.Lock()
	defer n.lock.Unlock()
	// scan over all current checkinstances if their check still exists.
	// If not the check has been deleted and the spawned instances need
	// a good deletion
	for ck := range n.CheckInstances {
		if _, ok := n.Checks[ck]; ok {
			// check still exists
			continue
		}

		// check no longer exists -> cleanup
		inst := n.CheckInstances[ck]
		for _, i := range inst {
			n.actionCheckInstanceDelete(n.Instances[i].MakeAction())
			n.log.Printf(
				"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
				n.GetRepositoryName(), `CleanupInstance`, `node`, n.ID.String(),
				ck, i,
			)
			delete(n.Instances, i)
		}
		delete(n.CheckInstances, ck)
	}
}

func (n *Node) removeDisabledCheckInstances() {
	n.lock.Lock()
	defer n.lock.Unlock()
	// loop over all checks and test if there is a reason to disable
	// its check instances. And with disable we mean delete.
	for chk := range n.Checks {
		disableThis := false
		// disable this check if the system property
		// `disable_all_monitoring` is set for the view that the check
		// uses
		if _, hit, _ := n.evalSystemProp(
			msg.SystemPropertyDisableAllMonitoring,
			`true`,
			n.Checks[chk].View,
		); hit {
			disableThis = true
		}
		// disable this check if the system property
		// `disable_check_configuration` is set to the
		// check_configuration that spawned this check
		if _, hit, _ := n.evalSystemProp(
			msg.SystemPropertyDisableCheckConfiguration,
			n.Checks[chk].ConfigID.String(),
			n.Checks[chk].View,
		); hit {
			disableThis = true
		}
		// if there was a reason to disable this check, all instances
		// are deleted
		if disableThis {
			if instanceArray, ok := n.CheckInstances[chk]; ok {
				for _, i := range instanceArray {
					n.actionCheckInstanceDelete(n.Instances[i].MakeAction())
					n.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s,"+
						" CheckID=%s, InstanceID=%s", n.GetRepositoryName(),
						`RemoveDisabledInstance`, `node`, n.ID.String(),
						chk, i,
					)
					delete(n.Instances, i)
				}
				delete(n.CheckInstances, chk)
			}
		}
	}
}

func (n *Node) calculateCheckInstances(startup bool) {
	wg := sync.WaitGroup{}
	for i := range n.Checks {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			n.processCheckForUpdates(name, startup)
		}(i)
	}
	wg.Wait()

	// completed the pass, reset update flag
	n.hasUpdate = false
}

func (n *Node) processCheckForUpdates(chkName string, startup bool) {
	n.lock.RLock()
	if n.Checks[chkName].Inherited == false && n.Checks[chkName].ChildrenOnly == true {
		// not active here
		n.lock.RUnlock()
		return
	}
	if _, hit, _ := n.evalSystemProp(
		// skip check if `disable_all_monitoring` property is set
		msg.SystemPropertyDisableAllMonitoring,
		`true`,
		n.Checks[chkName].View,
	); hit {
		n.lock.RUnlock()
		return
	}
	if _, hit, _ := n.evalSystemProp(
		// skip check if `disable_check_configuration` property is set
		msg.SystemPropertyDisableCheckConfiguration,
		n.Checks[chkName].ConfigID.String(),
		n.Checks[chkName].View,
	); hit {
		n.lock.RUnlock()
		return
	}

	ctx := newCheckContext(chkName, n.Checks[chkName].View, startup)
	n.lock.RUnlock()
	n.constraintCheck(ctx)
	n.log.Printf(
		"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, Match=%t",
		n.GetRepositoryName(), `ConstraintEvaluation`, `node`,
		n.ID.String(), chkName, ctx.brokeConstraint,
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
		n.createNoServiceCheckInstance(ctx)
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
		n.createPerServiceCheckInstances(ctx)
	}

	if ctx.startupBroken {
		return
	}

	// all new check instances have been built, check which
	// existing instances did not get an update and need to be
	// deleted
	if !ctx.startup {
		n.pruneOldCheckInstances(ctx)
		n.dispatchCheckInstanceUpdates(ctx)
	}

	n.createNewCheckInstances(ctx)

	delete(n.CheckInstances, ctx.uuid)
	n.CheckInstances[ctx.uuid] = ctx.newCheckInstances
}

func (n *Node) constraintCheck(ctx *checkContext) {
	n.lock.RLock()
	defer n.lock.RUnlock()
	// these constaint types must always match for the instance to
	// be valid. defer service and attribute
	for _, cc := range n.Checks[ctx.uuid].Constraints {
		switch cc.Type {
		case msg.ConstraintNative:
			if n.evalNativeProp(cc.Key, cc.Value) {
				ctx.nativeConstr[cc.Key] = cc.Value
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintSystem:
			if id, hit, bind := n.evalSystemProp(cc.Key, cc.Value, ctx.view); hit {
				ctx.systemConstr[id] = bind
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintOncall:
			if id, hit := n.evalOncallProp(cc.Key, cc.Value, ctx.view); hit {
				ctx.oncallConstr = id
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintCustom:
			if id, hit, bind := n.evalCustomProp(cc.Key, cc.Value, ctx.view); hit {
				ctx.customConstr[id] = bind
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintService:

			ctx.hasServiceConstraint = true
			if id, hit, bind := n.evalServiceProp(cc.Key, cc.Value, ctx.view); hit {
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
				hit, bind := n.evalAttributeOfService(id, ctx.view, attr.Key, attr.Value)
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
			if hit, svcIDMap := n.evalAttributeProp(ctx.view, attr.Key, attr.Value); hit {
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

func (n *Node) createNoServiceCheckInstance(ctx *checkContext) {
	n.lock.RLock()
	defer n.lock.RUnlock()
	inst := CheckInstance{
		InstanceID: uuid.UUID{},
		CheckID: func(id string) uuid.UUID {
			f, _ := uuid.FromString(id)
			return f
		}(ctx.uuid),
		ConfigID: func(id string) uuid.UUID {
			f, _ := uuid.FromString(n.Checks[id].ConfigID.String())
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
		n.lock.RUnlock()
		n.lock.Lock()
		matched := false
		for loadedID, loadedInst := range n.loadedInstances[ctx.uuid] {
			if loadedInst.InstanceSvcCfgHash != `` {
				continue
			}
			// check if an instance exists bound against the same
			// constraints
			if inst.MatchConstraints(&loadedInst) {
				// found a match
				matched = true

				inst.InstanceID, _ = uuid.FromString(loadedID)
				inst.InstanceConfigID, _ = uuid.FromString(
					loadedInst.InstanceConfigID.String(),
				)
				inst.Version = loadedInst.Version
				delete(n.loadedInstances[ctx.uuid], loadedID)
				n.log.Printf(
					"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, "+
						"InstanceID=%s, ServiceConstrained=%t", n.GetRepositoryName(),
					`ComputeInstance`, `node`, n.ID.String(), ctx.uuid,
					loadedID, false,
				)
				break
			}
		}
		if !matched {
			// downgrade to readlock
			n.lock.Unlock()
			n.lock.RLock()

			n.log.Printf("TK[%s]: Failed to match computed instance to loaded instances."+
				" ObjType=%s, ObjId=%s, CheckID=%s", `node`, n.ID.String(), ctx.uuid,
				n.GetRepositoryName())
			n.Fault.Error <- &Error{
				Action: `Failed to match a computed instance to loaded data`,
			}
			ctx.startupBroken = true
			return
		}
		// downgrade to readlock
		n.lock.Unlock()
		n.lock.RLock()
	default:
		for _, existingInstanceID := range n.CheckInstances[ctx.uuid] {
			existingInstance := n.Instances[existingInstanceID]

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
		n.log.Printf(
			"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, "+
				"InstanceID=%s, ServiceConstrained=%t", n.GetRepositoryName(),
			`ComputeInstance`, `node`, n.ID.String(), ctx.uuid, inst.InstanceID.String(),
			false,
		)
	}

	ctx.newInstances[inst.InstanceID.String()] = inst
	ctx.newCheckInstances = append(ctx.newCheckInstances, inst.InstanceID.String())
}

func (n *Node) createPerServiceCheckInstances(ctx *checkContext) {
	n.lock.RLock()
	defer n.lock.RUnlock()
	for svcID := range ctx.serviceConstr {
		svcCfg := n.getServiceMap(svcID)

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
					f, _ := uuid.FromString(n.Checks[id].ConfigID.String())
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
				n.lock.RUnlock()
				n.lock.Lock()
				matched := false

				for ldInstID, ldInst := range n.loadedInstances[ctx.uuid] {
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
						delete(n.loadedInstances[ctx.uuid], ldInstID)
						n.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, "+
							"CheckID=%s, InstanceID=%s, ServiceConstrained=%t",
							n.GetRepositoryName(), `ComputeInstance`, `node`,
							n.ID.String(), ctx.uuid, ldInstID, true,
						)
						break
					}
				}
				if !matched {
					// downgrade to readlock
					n.lock.Unlock()
					n.lock.RLock()

					n.log.Printf(
						"TK[%s]: Failed to match computed instance to loaded instances."+
							" ObjType=%s, ObjId=%s, CheckID=%s", `node`, n.ID.String(),
						ctx.uuid, n.GetRepositoryName())
					n.Fault.Error <- &Error{
						Action: `Failed to match a computed instance to loaded data`,
					}
					ctx.startupBroken = true
					return
				}
				// downgrade to readlock
				n.lock.Unlock()
				n.lock.RLock()
			default:
				// lookup existing instance ids for check in n.CheckInstances
				// to determine if this is an update
				for _, exInstID := range n.CheckInstances[ctx.uuid] {
					exInst := n.Instances[exInstID]
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
				n.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, "+
					"InstanceID=%s, ServiceConstrained=%t", n.GetRepositoryName(),
					`ComputeInstance`, `node`, n.ID.String(), ctx.uuid,
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
		if len(n.loadedInstances[ctx.uuid]) != 0 {
			n.Fault.Error <- &Error{
				Action: `Leftover matched instances after assignment, ` +
					`computed instances missing`,
			}
			ctx.startupBroken = true
		}
	}
}

func (n *Node) pruneOldCheckInstances(ctx *checkContext) {
	n.lock.Lock()
	defer n.lock.Unlock()
	Instances := n.CheckInstances[ctx.uuid]
	for _, oldInstanceID := range Instances {
		if _, ok := ctx.newInstances[oldInstanceID]; !ok {
			// there is no new version for this instance id
			n.actionCheckInstanceDelete(n.Instances[oldInstanceID].MakeAction())
			n.log.Printf(
				"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
				n.GetRepositoryName(), `DeleteInstance`, `node`, n.ID.String(),
				ctx.uuid, oldInstanceID,
			)
			n.CheckInstances[ctx.uuid] = removeFromArray(n.CheckInstances[ctx.uuid], oldInstanceID)
			delete(n.Instances, oldInstanceID)
		}
	}
}

func (n *Node) dispatchCheckInstanceUpdates(ctx *checkContext) {
	n.lock.Lock()
	defer n.lock.Unlock()

	for _, oldInstanceID := range n.CheckInstances[ctx.uuid] {
		delete(n.Instances, oldInstanceID)
		n.Instances[oldInstanceID] = ctx.newInstances[oldInstanceID]
		n.actionCheckInstanceUpdate(n.Instances[oldInstanceID].MakeAction())
		n.log.Printf(
			"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
			n.GetRepositoryName(), `UpdateInstance`, `node`, n.ID.String(),
			ctx.uuid, oldInstanceID,
		)
	}
}

func (n *Node) createNewCheckInstances(ctx *checkContext) {
	n.lock.Lock()
	defer n.lock.Unlock()

	for _, newInstanceID := range ctx.newCheckInstances {
		if _, ok := n.Instances[newInstanceID]; !ok {
			// this instance is new, not an update
			n.Instances[newInstanceID] = ctx.newInstances[newInstanceID]

			action := `CreateInstance`
			switch ctx.startup {
			case true:
				action = `RecreateInstance`
			default:
				n.actionCheckInstanceCreate(n.Instances[newInstanceID].MakeAction())
			}
			n.log.Printf(
				"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
				n.GetRepositoryName(), action, `node`, n.ID.String(),
				ctx.uuid, newInstanceID,
			)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
