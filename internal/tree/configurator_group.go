/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import uuid "github.com/satori/go.uuid"

func (teg *Group) updateCheckInstances() {
	repoName := teg.GetRepositoryName()

	// object may have no checks, but there could be instances to mop up
	if len(teg.Checks) == 0 && len(teg.Instances) == 0 {
		teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, HasChecks=%t",
			repoName,
			`UpdateCheckInstances`,
			`group`,
			teg.ID.String(),
			false,
		)
		// found nothing to do, ensure update flag is unset again
		teg.hasUpdate = false
		return
	}

	// if there are loaded instances, then this is the initial rebuild
	// of the tree
	startupLoad := false
	if len(teg.loadedInstances) > 0 {
		startupLoad = true
	}

	// if this is not the startupLoad and there are no updates, then there
	// is noting to do
	if !startupLoad && !teg.hasUpdate {
		return
	}

	// scan over all current checkinstances if their check still exists.
	// If not the check has been deleted and the spawned instances need
	// a good deletion
	for ck := range teg.CheckInstances {
		if _, ok := teg.Checks[ck]; ok {
			// check still exists
			continue
		}

		// check no longer exists -> cleanup
		inst := teg.CheckInstances[ck]
		for _, i := range inst {
			teg.actionCheckInstanceDelete(teg.Instances[i].MakeAction())
			teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
				repoName,
				`CleanupInstance`,
				`group`,
				teg.ID.String(),
				ck,
				i,
			)
			delete(teg.Instances, i)
		}
		delete(teg.CheckInstances, ck)
	}

	// loop over all checks and test if there is a reason to disable
	// its check instances. And with disable we mean delete.
	for chk := range teg.Checks {
		disableThis := false
		// disable this check if the system property
		// `disable_all_monitoring` is set for the view that the check
		// uses.
		if _, hit, _ := teg.evalSystemProp(
			`disable_all_monitoring`,
			`true`,
			teg.Checks[chk].View,
		); hit {
			disableThis = true
		}
		// disable this check if the system property
		// `disable_check_configuration` is set to the
		// check_configuration that spawned this check
		if _, hit, _ := teg.evalSystemProp(
			`disable_check_configuration`,
			teg.Checks[chk].ConfigID.String(),
			teg.Checks[chk].View,
		); hit {
			disableThis = true
		}
		// if there was a reason to disable this check, all instances
		// are deleted
		if disableThis {
			if instanceArray, ok := teg.CheckInstances[chk]; ok {
				for _, i := range instanceArray {
					teg.actionCheckInstanceDelete(teg.Instances[i].MakeAction())
					teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
						repoName,
						`RemoveDisabledInstance`,
						`group`,
						teg.ID.String(),
						chk,
						i,
					)
					delete(teg.Instances, i)
				}
				delete(teg.CheckInstances, chk)
			}
		}
	}

	// process remaining checks
checksloop:
	for i := range teg.Checks {
		if teg.Checks[i].Inherited == false && teg.Checks[i].ChildrenOnly == true {
			continue checksloop
		}
		if teg.Checks[i].View == "local" {
			continue checksloop
		}
		// skip check if its view has `disable_all_monitoring`
		// property set
		if _, hit, _ := teg.evalSystemProp(
			`disable_all_monitoring`,
			`true`,
			teg.Checks[i].View,
		); hit {
			continue checksloop
		}
		// skip check if there is a matching `disable_check_configuration`
		// property
		if _, hit, _ := teg.evalSystemProp(
			`disable_check_configuration`,
			teg.Checks[i].ConfigID.String(),
			teg.Checks[i].View,
		); hit {
			continue checksloop
		}

		hasBrokenConstraint := false
		hasServiceConstraint := false
		hasAttributeConstraint := false
		view := teg.Checks[i].View

		attributes := []CheckConstraint{}
		oncallC := ""                                  // Id
		systemC := map[string]string{}                 // Id->Value
		nativeC := map[string]string{}                 // Property->Value
		serviceC := map[string]string{}                // Id->Value
		customC := map[string]string{}                 // Id->Value
		attributeC := map[string]map[string][]string{} // svcID->attr->[ value, ... ]

		newInstances := map[string]CheckInstance{}
		newCheckInstances := []string{}

		// these constaint types must always match for the instance to
		// be valid. defer service and attribute
	constraintcheck:
		for _, c := range teg.Checks[i].Constraints {
			switch c.Type {
			case "native":
				if teg.evalNativeProp(c.Key, c.Value) {
					nativeC[c.Key] = c.Value
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "system":
				if id, hit, bind := teg.evalSystemProp(c.Key, c.Value, view); hit {
					systemC[id] = bind
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "oncall":
				if id, hit := teg.evalOncallProp(c.Key, c.Value, view); hit {
					oncallC = id
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "custom":
				if id, hit, bind := teg.evalCustomProp(c.Key, c.Value, view); hit {
					customC[id] = bind
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "service":
				hasServiceConstraint = true
				if id, hit, bind := teg.evalServiceProp(c.Key, c.Value, view); hit {
					serviceC[id] = bind
				} else {
					hasBrokenConstraint = true
					break constraintcheck
				}
			case "attribute":
				hasAttributeConstraint = true
				attributes = append(attributes, c)
			}
		}
		if hasBrokenConstraint {
			teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, Match=%t",
				repoName,
				`ConstraintEvaluation`,
				`group`,
				teg.ID.String(),
				i,
				false,
			)
			continue checksloop
		}

		/* if the check has both service and attribute constraints,
		* then for the check to hit, the tree element needs to have
		* all the services, and each of them needs to match all
		* attribute constraints
		 */
		if hasServiceConstraint && hasAttributeConstraint {
		svcattrloop:
			for id := range serviceC {
				for _, attr := range attributes {
					hit, bind := teg.evalAttributeOfService(id, view, attr.Key, attr.Value)
					if hit {
						if attributeC[id] == nil {
							// attributeC[id] might still be a nil map
							attributeC[id] = map[string][]string{}
						}
						attributeC[id][attr.Key] = append(attributeC[id][attr.Key], bind)
					} else {
						hasBrokenConstraint = true
						break svcattrloop
					}
				}
			}
			/* if the check has only attribute constraints and no
			* service constraint, then we pull in every service that
			* matches all attribute constraints and generate a check
			* instance for it
			 */
		} else if hasAttributeConstraint {
			attrCount := len(attributes)
			for _, attr := range attributes {
				hit, svcIDMap := teg.evalAttributeProp(view, attr.Key, attr.Value)
				if hit {
					for id, bind := range svcIDMap {
						serviceC[id] = svcIDMap[id]
						if attributeC[id] == nil {
							// attributeC[id] might still be a nil map
							attributeC[id] = make(map[string][]string)
						}
						attributeC[id][attr.Key] = append(attributeC[id][attr.Key], bind)
					}
				}
			}
			// delete all services that did not match all attributes
			//
			// if a check has two attribute constraints on the same
			// attribute, then len(attributeC[id]) != len(attributes)
			for id := range attributeC {
				if teg.countAttribC(attributeC[id]) != attrCount {
					delete(serviceC, id)
					delete(attributeC, id)
				}
			}
			// declare service constraints in effect if we found a
			// service that bound all attribute constraints
			if len(serviceC) > 0 {
				hasServiceConstraint = true
			} else {
				// found no services that fulfilled all constraints
				hasBrokenConstraint = true
			}
		}
		if hasBrokenConstraint {
			teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, Match=%t",
				repoName,
				`ConstraintEvaluation`,
				`group`,
				teg.ID.String(),
				i,
				false,
			)
			continue checksloop
		}
		// check triggered, create instances
		teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, Match=%t",
			repoName,
			`ConstraintEvaluation`,
			`group`,
			teg.ID.String(),
			i,
			true,
		)

		/* if there are no service constraints, one check instance is
		* created for this check
		 */
		if !hasServiceConstraint {
			inst := CheckInstance{
				InstanceID: uuid.UUID{},
				CheckID: func(id string) uuid.UUID {
					f, _ := uuid.FromString(id)
					return f
				}(i),
				ConfigID: func(id string) uuid.UUID {
					f, _ := uuid.FromString(teg.Checks[id].ConfigID.String())
					return f
				}(i),
				InstanceConfigID:      uuid.NewV4(),
				ConstraintOncall:      oncallC,
				ConstraintService:     serviceC,
				ConstraintSystem:      systemC,
				ConstraintCustom:      customC,
				ConstraintNative:      nativeC,
				ConstraintAttribute:   attributeC,
				InstanceService:       "",
				InstanceServiceConfig: nil,
				InstanceSvcCfgHash:    "",
			}
			inst.calcConstraintHash()
			inst.calcConstraintValHash()

			if startupLoad {
			nosvcstartinstanceloop:
				for ldInstID, ldInst := range teg.loadedInstances[i] {
					if ldInst.InstanceSvcCfgHash != "" {
						continue nosvcstartinstanceloop
					}
					// check if an instance exists bound against the
					// same constraints
					if ldInst.ConstraintHash == inst.ConstraintHash &&
						uuid.Equal(ldInst.ConfigID, inst.ConfigID) &&
						ldInst.ConstraintValHash == inst.ConstraintValHash {

						// found a match
						inst.InstanceID, _ = uuid.FromString(ldInstID)
						inst.InstanceConfigID, _ = uuid.FromString(ldInst.InstanceConfigID.String())
						inst.Version = ldInst.Version
						delete(teg.loadedInstances[i], ldInstID)
						teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s, ServiceConstrained=%t",
							repoName,
							`ComputeInstance`,
							`group`,
							teg.ID.String(),
							i,
							ldInstID,
							false,
						)
						goto nosvcstartinstancematch
					}
				}
				// if we hit here, then we just computed an instance
				// that we could not match to any loaded instances
				// -> something is wrong
				teg.log.Printf("TK[%s]: Failed to match computed instance to loaded instances."+
					" ObjType=%s, ObjId=%s, CheckID=%s", `group`, teg.ID.String(), i, repoName)
				teg.Fault.Error <- &Error{Action: `Failed to match a computed instance to loaded data`}
				return
			nosvcstartinstancematch:
			} else {
			nosvcinstanceloop:
				for _, exInstID := range teg.CheckInstances[i] {
					exInst := teg.Instances[exInstID]
					// ignore instances with service constraints
					if exInst.InstanceSvcCfgHash != "" {
						continue nosvcinstanceloop
					}
					// check if an instance exists bound against the same
					// constraints
					if exInst.ConstraintHash == inst.ConstraintHash {
						inst.InstanceID, _ = uuid.FromString(exInst.InstanceID.String())
						inst.Version = exInst.Version + 1
						break nosvcinstanceloop
					}
				}
				if uuid.Equal(uuid.Nil, inst.InstanceID) {
					// no match was found during nosvcinstanceloop, this
					// is a new instance
					inst.Version = 0
					inst.InstanceID = uuid.NewV4()
				}
				teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s, ServiceConstrained=%t",
					repoName,
					`ComputeInstance`,
					`group`,
					teg.ID.String(),
					i,
					inst.InstanceID.String(),
					false,
				)
			}
			newInstances[inst.InstanceID.String()] = inst
			newCheckInstances = append(newCheckInstances, inst.InstanceID.String())
		}

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
	serviceconstraintloop:
		for svcID := range serviceC {
			if !hasServiceConstraint {
				break serviceconstraintloop
			}

			svcCfg := teg.getServiceMap(svcID)

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
					}(i),
					ConfigID: func(id string) uuid.UUID {
						f, _ := uuid.FromString(teg.Checks[id].ConfigID.String())
						return f
					}(i),
					InstanceConfigID:      uuid.NewV4(),
					ConstraintOncall:      oncallC,
					ConstraintService:     serviceC,
					ConstraintSystem:      systemC,
					ConstraintCustom:      customC,
					ConstraintNative:      nativeC,
					ConstraintAttribute:   attributeC,
					InstanceService:       svcID,
					InstanceServiceConfig: cfg,
				}
				inst.calcConstraintHash()
				inst.calcConstraintValHash()
				inst.calcInstanceSvcCfgHash()

				if startupLoad {
					for ldInstID, ldInst := range teg.loadedInstances[i] {
						// check for data from loaded instance
						if ldInst.InstanceSvcCfgHash == inst.InstanceSvcCfgHash &&
							ldInst.ConstraintHash == inst.ConstraintHash &&
							ldInst.ConstraintValHash == inst.ConstraintValHash &&
							ldInst.InstanceService == inst.InstanceService &&
							uuid.Equal(ldInst.ConfigID, inst.ConfigID) {

							// found a match
							inst.InstanceID, _ = uuid.FromString(ldInstID)
							inst.InstanceConfigID, _ = uuid.FromString(ldInst.InstanceConfigID.String())
							inst.Version = ldInst.Version
							// we can assume InstanceServiceConfig to
							// be equal, since InstanceSvcCfgHash is
							// equal
							delete(teg.loadedInstances[i], ldInstID)
							teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s, ServiceConstrained=%t",
								repoName,
								`ComputeInstance`,
								`group`,
								teg.ID.String(),
								i,
								ldInstID,
								true,
							)
							goto startinstancematch
						}
					}
					// if we hit here, then just computed an
					// instance that we could not match to any
					// loaded instances -> something is wrong
					teg.log.Printf("TK[%s]: Failed to match computed instance to loaded instances."+
						" ObjType=%s, ObjId=%s, CheckID=%s", `group`, teg.ID.String(), i, repoName)
					teg.Fault.Error <- &Error{Action: `Failed to match a computed instance to loaded data`}
					return
				startinstancematch:
				} else {
					// lookup existing instance ids for check in teg.CheckInstances
					// to determine if this is an update
				instanceloop:
					for _, exInstID := range teg.CheckInstances[i] {
						exInst := teg.Instances[exInstID]
						// this existing instance is for the same service
						// configuration -> this is an update
						if exInst.InstanceSvcCfgHash == inst.InstanceSvcCfgHash {
							inst.InstanceID, _ = uuid.FromString(exInst.InstanceID.String())
							inst.Version = exInst.Version + 1
							break instanceloop
						}
					}
					if uuid.Equal(uuid.Nil, inst.InstanceID) {
						// no match was found during instanceloop, this is
						// a new instance
						inst.Version = 0
						inst.InstanceID = uuid.NewV4()
					}
					teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s, ServiceConstrained=%t",
						repoName,
						`ComputeInstance`,
						`group`,
						teg.ID.String(),
						i,
						inst.InstanceID.String(),
						true,
					)
				}
				newInstances[inst.InstanceID.String()] = inst
				newCheckInstances = append(newCheckInstances, inst.InstanceID.String())
			}
		} // LOOPEND: range serviceC

		// all instances have been built and matched to
		// loaded instances, but there are loaded
		// instances left. why?
		if startupLoad && len(teg.loadedInstances[i]) != 0 {
			teg.Fault.Error <- &Error{Action: `Leftover matched instances after assignment, computed instances missing`}
			return
		}

		// all new check instances have been built, check which
		// existing instances did not get an update and need to be
		// deleted
		for _, oldInstanceID := range teg.CheckInstances[i] {
			if _, ok := newInstances[oldInstanceID]; !ok {
				// there is no new version for this instance id
				teg.actionCheckInstanceDelete(teg.Instances[oldInstanceID].MakeAction())
				teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
					repoName,
					`DeleteInstance`,
					`group`,
					teg.ID.String(),
					i,
					oldInstanceID,
				)
				delete(teg.Instances, oldInstanceID)
				continue
			}
			delete(teg.Instances, oldInstanceID)
			teg.Instances[oldInstanceID] = newInstances[oldInstanceID]
			teg.actionCheckInstanceUpdate(teg.Instances[oldInstanceID].MakeAction())
			teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
				repoName,
				`UpdateInstance`,
				`group`,
				teg.ID.String(),
				i,
				oldInstanceID,
			)
		}
		for _, newInstanceID := range newCheckInstances {
			if _, ok := teg.Instances[newInstanceID]; !ok {
				// this instance is new, not an update
				teg.Instances[newInstanceID] = newInstances[newInstanceID]
				// no need to send a create action during load; the
				// action channel is drained anyway
				if !startupLoad {
					teg.actionCheckInstanceCreate(teg.Instances[newInstanceID].MakeAction())
					teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
						repoName,
						`CreateInstance`,
						`group`,
						teg.ID.String(),
						i,
						newInstanceID,
					)
				} else {
					teg.log.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectID=%s, CheckID=%s, InstanceID=%s",
						repoName,
						`RecreateInstance`,
						`group`,
						teg.ID.String(),
						i,
						newInstanceID,
					)
				}
			}
		}
		delete(teg.CheckInstances, i)
		teg.CheckInstances[i] = newCheckInstances
	} // LOOPEND: range teg.Checks

	// completed the pass, reset update flag
	teg.hasUpdate = false
}

func (teg *Group) evalNativeProp(prop string, val string) bool {
	switch prop {
	case "environment":
		env := teg.Parent.(Bucketeer).GetEnvironment()
		if val == env {
			return true
		}
	case "object_type":
		if val == "group" {
			return true
		}
	case "object_state":
		if val == teg.State {
			return true
		}
	case "hardware_node":
		// group != hardware
		return false
	}
	return false
}

func (teg *Group) evalSystemProp(prop string, val string, view string) (string, bool, string) {
	for _, v := range teg.PropertySystem {
		t := v.(*PropertySystem)
		if t.Key == prop && (t.Value == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Key, true, t.Value
		}
	}
	return "", false, ""
}

func (teg *Group) evalOncallProp(prop string, val string, view string) (string, bool) {
	for _, v := range teg.PropertyOncall {
		t := v.(*PropertyOncall)
		if "OncallID" == prop && t.ID.String() == val && (t.View == view || t.View == `any`) {
			return t.ID.String(), true
		}
	}
	return "", false
}

func (teg *Group) evalCustomProp(prop string, val string, view string) (string, bool, string) {
	for _, v := range teg.PropertyCustom {
		t := v.(*PropertyCustom)
		if t.Key == prop && (t.Value == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.Key, true, t.Value
		}
	}
	return "", false, ""
}

func (teg *Group) evalServiceProp(prop string, val string, view string) (string, bool, string) {
	for _, v := range teg.PropertyService {
		t := v.(*PropertyService)
		if prop == "name" && (t.Service == val || val == `@defined`) && (t.View == view || t.View == `any`) {
			return t.ID.String(), true, t.Service
		}
	}
	return "", false, ""
}

func (teg *Group) evalAttributeOfService(svcID string, view string, attribute string, value string) (bool, string) {
	t := teg.PropertyService[svcID].(*PropertyService)
	for _, a := range t.Attributes {
		if a.Name == attribute && (t.View == view || t.View == `any`) && (a.Value == value || value == `@defined`) {
			return true, a.Value
		}
	}
	return false, ""
}

func (teg *Group) evalAttributeProp(view string, attr string, value string) (bool, map[string]string) {
	f := map[string]string{}
svcloop:
	for _, v := range teg.PropertyService {
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

func (teg *Group) getServiceMap(serviceID string) map[string][]string {
	svc := new(PropertyService)
	svc = teg.PropertyService[serviceID].(*PropertyService)

	res := map[string][]string{}
	for _, v := range svc.Attributes {
		res[v.Name] = append(res[v.Name], v.Value)
	}
	return res
}

func (teg *Group) countAttribC(attributeC map[string][]string) int {
	var count int
	for key := range attributeC {
		count = count + len(attributeC[key])
	}
	return count
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
