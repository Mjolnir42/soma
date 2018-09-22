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

func (g *Group) evalNativeProp(prop string, val string) bool {
	switch prop {
	case msg.NativePropertyEnvironment:
		env := g.Parent.(Bucketeer).GetEnvironment()
		if val == env {
			return true
		}
	case msg.NativePropertyEntity:
		if val == msg.EntityGroup {
			return true
		}
	case msg.NativePropertyState:
		if val == g.State {
			return true
		}
	case msg.NativePropertyHardwareNode:
		// group != hardware
		return false
	}
	return false
}

func (g *Group) evalSystemProp(prop string, val string, view string) (string, bool, string) {
	for _, v := range g.PropertySystem {
		t := v.(*PropertySystem)
		if t.Key == prop && (t.Value == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Key, true, t.Value
		}
	}
	return "", false, ""
}

func (g *Group) evalOncallProp(prop string, val string, view string) (string, bool) {
	for _, v := range g.PropertyOncall {
		t := v.(*PropertyOncall)
		if "OncallID" == prop && t.ID.String() == val && (t.View == view || t.View == `any`) {
			return t.ID.String(), true
		}
	}
	return "", false
}

func (g *Group) evalCustomProp(prop string, val string, view string) (string, bool, string) {
	for _, v := range g.PropertyCustom {
		t := v.(*PropertyCustom)
		if t.Key == prop && (t.Value == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Key, true, t.Value
		}
	}
	return "", false, ""
}

func (g *Group) evalServiceProp(prop string, val string, view string) (string, bool, string) {
	for _, v := range g.PropertyService {
		t := v.(*PropertyService)
		if prop == "name" && (t.Service == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.ID.String(), true, t.Service
		}
	}
	return "", false, ""
}

func (g *Group) evalAttributeOfService(svcID string, view string, attribute string, value string) (bool, string) {
	t := g.PropertyService[svcID].(*PropertyService)
	for _, a := range t.Attributes {
		if a.Name == attribute && (t.View == view || t.View == `any`) && (a.Value == value || value == `@defined`) {
			return true, a.Value
		}
	}
	return false, ""
}

func (g *Group) evalAttributeProp(view string, attr string, value string) (bool, map[string]string) {
	f := map[string]string{}
svcloop:
	for _, v := range g.PropertyService {
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

func (g *Group) getServiceMap(serviceID string) map[string][]string {
	svc := new(PropertyService)
	svc = g.PropertyService[serviceID].(*PropertyService)

	res := map[string][]string{}
	for _, v := range svc.Attributes {
		res[v.Name] = append(res[v.Name], v.Value)
	}
	return res
}

func (g *Group) updateCheckInstances() {
	// object may have no checks, but there could be instances to mop up
	if len(g.Checks) == 0 && len(g.Instances) == 0 {
		g.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, HasChecks=%t",
			g.GetRepositoryName(),
			`UpdateCheckInstances`,
			`group`,
			g.ID.String(),
			false,
		)
		// found nothing to do, ensure update flag is unset again
		g.hasUpdate = false
		return
	}

	// if there are loaded instances, then this is the initial rebuild
	// of the tree
	startup := false
	if len(g.loadedInstances) > 0 {
		startup = true
	}

	// if this is not the startupLoad and there are no updates, then there
	// is noting to do
	if !startup && !g.hasUpdate {
		return
	}

	g.deleteOrphanCheckInstances()

	g.removeDisabledCheckInstances()

	g.calculateCheckInstances(startup)
}

func (g *Group) deleteOrphanCheckInstances() {
	g.lock.Lock()
	defer g.lock.Unlock()
	// scan over all current checkinstances if their check still exists.
	// If not the check has been deleted and the spawned instances need
	// a good deletion
	for ck := range g.CheckInstances {
		if _, ok := g.Checks[ck]; ok {
			// check still exists
			continue
		}

		// check no longer exists -> cleanup
		inst := g.CheckInstances[ck]
		for _, i := range inst {
			g.actionCheckInstanceDelete(g.Instances[i].MakeAction())
			g.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
				g.GetRepositoryName(),
				`CleanupInstance`,
				`group`,
				g.ID.String(),
				ck,
				i,
			)
			delete(g.Instances, i)
		}
		delete(g.CheckInstances, ck)
	}
}

func (g *Group) removeDisabledCheckInstances() {
	g.lock.Lock()
	defer g.lock.Unlock()
	// loop over all checks and test if there is a reason to disable
	// its check instances. And with disable we mean delete.
	for chk := range g.Checks {
		disableThis := false
		// disable this check if the system property
		// `disable_all_monitoring` is set for the view that the check
		// uses.
		if _, hit, _ := g.evalSystemProp(
			msg.SystemPropertyDisableAllMonitoring,
			`true`,
			g.Checks[chk].View,
		); hit {
			disableThis = true
		}
		// disable this check if the system property
		// `disable_check_configuration` is set to the
		// check_configuration that spawned this check
		if _, hit, _ := g.evalSystemProp(
			msg.SystemPropertyDisableCheckConfiguration,
			g.Checks[chk].ConfigID.String(),
			g.Checks[chk].View,
		); hit {
			disableThis = true
		}
		// if there was a reason to disable this check, all instances
		// are deleted
		if disableThis {
			if instanceArray, ok := g.CheckInstances[chk]; ok {
				for _, i := range instanceArray {
					g.actionCheckInstanceDelete(g.Instances[i].MakeAction())
					g.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
						g.GetRepositoryName(),
						`RemoveDisabledInstance`,
						`group`,
						g.ID.String(),
						chk,
						i,
					)
					delete(g.Instances, i)
				}
				delete(g.CheckInstances, chk)
			}
		}
	}
}

func (g *Group) calculateCheckInstances(startup bool) {
	wg := sync.WaitGroup{}
	for i := range g.Checks {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			g.processCheckForUpdates(name, startup)
		}(i)
	}
	wg.Wait()

	// completed the pass, reset update flag
	g.hasUpdate = false
}

func (g *Group) processCheckForUpdates(chkName string, startup bool) {
	g.lock.RLock()
	if g.Checks[chkName].Inherited == false && g.Checks[chkName].ChildrenOnly == true {
		// not active here
		g.lock.RUnlock()
		return
	}
	if g.Checks[chkName].View == msg.ViewLocal {
		// groups have no local view
		g.lock.RUnlock()
		return
	}
	if _, hit, _ := g.evalSystemProp(
		// skip check if `disable_all_monitoring` property is set
		msg.SystemPropertyDisableAllMonitoring,
		`true`,
		g.Checks[chkName].View,
	); hit {
		g.lock.RUnlock()
		return
	}
	if _, hit, _ := g.evalSystemProp(
		// skip check if `disable_check_configuration` property is set
		msg.SystemPropertyDisableCheckConfiguration,
		g.Checks[chkName].ConfigID.String(),
		g.Checks[chkName].View,
	); hit {
		g.lock.RUnlock()
		return
	}

	ctx := newCheckContext(chkName, g.Checks[chkName].View, startup)
	g.lock.RUnlock()

	g.constraintCheck(ctx)
	g.log.Printf(
		"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, Match=%t",
		g.GetRepositoryName(), `ConstraintEvaluation`, `group`,
		g.ID.String(), chkName, ctx.brokeConstraint,
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
		g.createNoServiceCheckInstance(ctx)
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
		g.createPerServiceCheckInstances(ctx)
	}

	if ctx.startupBroken {
		return
	}

	// all new check instances have been built, check which
	// existing instances did not get an update and need to be
	// deleted
	if !ctx.startup {
		g.pruneOldCheckInstances(ctx)
		g.dispatchCheckInstanceUpdates(ctx)
	}

	g.createNewCheckInstances(ctx)

	delete(g.CheckInstances, ctx.uuid)
	g.CheckInstances[ctx.uuid] = ctx.newCheckInstances
}

func (g *Group) constraintCheck(ctx *checkContext) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	// these constaint types must always match for the instance to
	// be valid. defer service and attribute
	for _, c := range g.Checks[ctx.uuid].Constraints {
		switch c.Type {
		case msg.ConstraintNative:
			if g.evalNativeProp(c.Key, c.Value) {
				ctx.nativeConstr[c.Key] = c.Value
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintSystem:
			if id, hit, bind := g.evalSystemProp(c.Key, c.Value, ctx.view); hit {
				ctx.systemConstr[id] = bind
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintOncall:
			if id, hit := g.evalOncallProp(c.Key, c.Value, ctx.view); hit {
				ctx.oncallConstr = id
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintCustom:
			if id, hit, bind := g.evalCustomProp(c.Key, c.Value, ctx.view); hit {
				ctx.customConstr[id] = bind
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintService:
			ctx.hasServiceConstraint = true
			if id, hit, bind := g.evalServiceProp(c.Key, c.Value, ctx.view); hit {
				ctx.serviceConstr[id] = bind
				continue
			}
			ctx.brokeConstraint = true
			return
		case msg.ConstraintAttribute:
			ctx.hasAttributeConstraint = true
			ctx.attributes = append(ctx.attributes, c)
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
				hit, bind := g.evalAttributeOfService(id, ctx.view, attr.Key, attr.Value)
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
			if hit, svcIDMap := g.evalAttributeProp(ctx.view, attr.Key, attr.Value); hit {
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

func (g *Group) createNoServiceCheckInstance(ctx *checkContext) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	inst := CheckInstance{
		InstanceID: uuid.UUID{},
		CheckID: func(id string) uuid.UUID {
			f, _ := uuid.FromString(id)
			return f
		}(ctx.uuid),
		ConfigID: func(id string) uuid.UUID {
			f, _ := uuid.FromString(g.Checks[id].ConfigID.String())
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
		g.lock.RUnlock()
		g.lock.Lock()
		matched := false

		for loadedID, loadedInst := range g.loadedInstances[ctx.uuid] {
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
				delete(g.loadedInstances[ctx.uuid], loadedID)
				g.log.Printf(
					"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, "+
						"InstanceID=%s, ServiceConstrained=%t", g.GetRepositoryName(),
					`ComputeInstance`, `group`, g.ID.String(), ctx.uuid,
					loadedID, false,
				)
				break
			}
		}
		if !matched {
			// downgrade to readlock
			g.lock.Unlock()
			g.lock.RLock()

			g.log.Printf("TK[%s]: Failed to match computed instance to loaded instances."+
				" ObjType=%s, ObjId=%s, CheckID=%s", `group`, g.ID.String(), ctx.uuid,
				g.GetRepositoryName())
			g.Fault.Error <- &Error{
				Action: `Failed to match a computed instance to loaded data`,
			}
			ctx.startupBroken = true
			return
		}
		// downgrade to readlock
		g.lock.Unlock()
		g.lock.RLock()
	default:
		for _, existingInstanceID := range g.CheckInstances[ctx.uuid] {
			existingInstance := g.Instances[existingInstanceID]

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
		g.log.Printf(
			"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, "+
				"InstanceID=%s, ServiceConstrained=%t", g.GetRepositoryName(),
			`ComputeInstance`, `group`, g.ID.String(), ctx.uuid, inst.InstanceID.String(),
			false,
		)
	}

	ctx.newInstances[inst.InstanceID.String()] = inst
	ctx.newCheckInstances = append(ctx.newCheckInstances, inst.InstanceID.String())
}

func (g *Group) createPerServiceCheckInstances(ctx *checkContext) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	for svcID := range ctx.serviceConstr {
		svcCfg := g.getServiceMap(svcID)

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
					f, _ := uuid.FromString(g.Checks[id].ConfigID.String())
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
				g.lock.RUnlock()
				g.lock.Lock()
				matched := false

				for ldInstID, ldInst := range g.loadedInstances[ctx.uuid] {
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
						delete(g.loadedInstances[ctx.uuid], ldInstID)
						g.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s,"+
							"CheckID=%s, InstanceID=%s, ServiceConstrained=%t",
							g.GetRepositoryName(), `ComputeInstance`, `group`,
							g.ID.String(), ctx.uuid, ldInstID, true,
						)
						break
					}
				}
				if !matched {
					// downgrade to readlock
					g.lock.Unlock()
					g.lock.RLock()

					g.log.Printf(
						"TK[%s]: Failed to match computed instance to loaded instances."+
							" ObjType=%s, ObjId=%s, CheckID=%s", `group`, g.ID.String(),
						ctx.uuid, g.GetRepositoryName())
					g.Fault.Error <- &Error{
						Action: `Failed to match a computed instance to loaded data`,
					}
					ctx.startupBroken = true
					return
				}
				// downgrade to readlock
				g.lock.Unlock()
				g.lock.RLock()
			default:
				// lookup existing instance ids for check in g.CheckInstances
				// to determine if this is an update
				for _, exInstID := range g.CheckInstances[ctx.uuid] {
					exInst := g.Instances[exInstID]
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
				g.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, "+
					"InstanceID=%s, ServiceConstrained=%t", g.GetRepositoryName(),
					`ComputeInstance`, `group`, g.ID.String(), ctx.uuid,
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
		if len(g.loadedInstances[ctx.uuid]) != 0 {
			g.Fault.Error <- &Error{
				Action: "Leftover matched instances after assignment, " +
					"computed instances missing",
			}
			ctx.startupBroken = true
		}
	}
}

func (g *Group) pruneOldCheckInstances(ctx *checkContext) {
	g.lock.Lock()
	defer g.lock.Unlock()

	for _, oldInstanceID := range g.CheckInstances[ctx.uuid] {
		if _, ok := ctx.newInstances[oldInstanceID]; !ok {
			// there is no new version for oldInstanceID
			g.actionCheckInstanceDelete(g.Instances[oldInstanceID].MakeAction())
			g.log.Printf(
				"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
				g.GetRepositoryName(), `DeleteInstance`, `group`, g.ID.String(),
				ctx.uuid, oldInstanceID,
			)
			delete(g.Instances, oldInstanceID)
		}
	}
}

func (g *Group) dispatchCheckInstanceUpdates(ctx *checkContext) {
	g.lock.Lock()
	defer g.lock.Unlock()

	for _, oldInstanceID := range g.CheckInstances[ctx.uuid] {
		delete(g.Instances, oldInstanceID)
		g.Instances[oldInstanceID] = ctx.newInstances[oldInstanceID]
		g.actionCheckInstanceUpdate(g.Instances[oldInstanceID].MakeAction())
		g.log.Printf(
			"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
			g.GetRepositoryName(), `UpdateInstance`, `group`, g.ID.String(),
			ctx.uuid, oldInstanceID,
		)
	}
}

func (g *Group) createNewCheckInstances(ctx *checkContext) {
	g.lock.Lock()
	defer g.lock.Unlock()

	for _, newInstanceID := range ctx.newCheckInstances {
		if _, ok := g.Instances[newInstanceID]; !ok {
			// this instance is new, not an update
			g.Instances[newInstanceID] = ctx.newInstances[newInstanceID]

			action := `CreateInstance`
			switch ctx.startup {
			case true:
				action = `RecreateInstance`
			default:
				g.actionCheckInstanceCreate(g.Instances[newInstanceID].MakeAction())
			}
			g.log.Printf(
				"TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
				g.GetRepositoryName(), action, `group`, g.ID.String(),
				ctx.uuid, newInstanceID,
			)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
