package soma

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mjolnir42/soma/internal/tree"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (tk *TreeKeeper) startupChecks(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}
	tk.startLog.Printf("TK[%s]: loading checks\n", tk.meta.repoName)

	//
	// load checks for the entire tree, in order from root to leaf.
	// Afterwards, load all check instances. This does not require
	// ordering.
	for _, typ := range []string{`repository`, `bucket`, `group`, `cluster`, `node`} {
		tk.startupScopedChecks(typ, stMap)
	}

	//
	// if this a checks level rebuild, then we need to populate the
	// tree again with the original CheckConfigs since all checks that
	// could be loaded have been deleted
	if tk.status.requiresRebuild && tk.status.rebuildLevel == `checks` {
		tk.startLog.Printf("TK[%s]: starting checks rebuild", tk.meta.repoID)
		for _, typ := range []string{`repository`, `bucket`, `group`, `cluster`, `node`} {
			tk.startupScopedReapplyCheckConfig(typ, stMap)
		}
	}

	if !tk.status.requiresRebuild {
		// recompute instances with preloaded IDs
		tk.tree.ComputeCheckInstances()

		// ensure there are no leftovers
		tk.tree.ClearLoadInfo()

		// this startup drains actions after checks, then suppresses
		// actions for instances that could be matched to loaded
		// information. Leftovers indicate that loaded and computed
		// instances diverge!
		// If requested, print all encountered messages instead of
		// simply bailing out.
		if tk.soma.conf.PrintChannels {
			if len(tk.actions) > 0 {
				tk.status.isBroken = true
				for i := len(tk.actions); i > 0; i-- {
					a := <-tk.actions
					jBxX, _ := json.Marshal(a)
					tk.startLog.Printf("TK[%s], startupChecks(): leftover action in channel: %s", tk.meta.repoName, string(jBxX))
				}
			}
			if len(tk.errors) > 0 {
				tk.status.isBroken = true
				for i := len(tk.errors); i > 0; i-- {
					e := <-tk.errors
					jBxX, _ := json.Marshal(e)
					tk.startLog.Printf("TK[%s], startupChecks(): error in channel: %s", tk.meta.repoName, string(jBxX))
				}
			}
			if tk.status.isBroken {
				return
			}
		}
		// drain the action channel
		if tk.drain(`action`) > 0 {
			tk.status.isBroken = true
			tk.startLog.Printf("TK[%s], startupChecks(): leftovers in actionChannel after drain", tk.meta.repoName)
			return
		}

		// drain the error channel
		if tk.drain(`error`) > 0 {
			tk.status.isBroken = true
			tk.startLog.Printf("TK[%s], startupChecks(): leftovers in errorChannel after drain", tk.meta.repoName)
			return
		}
	}
}

func (tk *TreeKeeper) startupScopedChecks(typ string, stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	// forward declare variables to allow goto use to dedup exit
	// handling
	var (
		err                                                          error
		checkID, srcCheckID, srcObjType, srcObjID, configID          string
		capabilityID, objID, objType, cfgName, cfgObjID, cfgObjType  string
		externalID, predicate, threshold, levelName, levelShort      string
		cstrType, value1, value2, value3, itemID, itemCfgID          string
		monitoringID, cstrHash, cstrValHash, instSvc, instSvcCfgHash string
		instSvcCfg, errLocation                                      string
		levelNumeric, numVal, interval, version                      int64
		isActive, hasInheritance, isChildrenOnly, isEnabled          bool
		grOrder                                                      map[string][]string
		grWeird                                                      map[string]string
		ckRows, thrRows, cstrRows, itRows, inRows, tckRows           *sql.Rows
		cfgMap                                                       map[string]proto.CheckConfig
		victim                                                       proto.CheckConfig // go/issues/3117 workaround
		ckTree                                                       *tree.Check
		ckItem                                                       tree.CheckItem
		ckOrder                                                      map[string]map[string]tree.Check
		nullBucketID                                                 sql.NullString
	)
	cfgMap = make(map[string]proto.CheckConfig)
	ckOrder = make(map[string]map[string]tree.Check)

	switch typ {
	case "group":
		if grOrder, grWeird, err = tk.orderGroups(stMap); err != nil {
			goto fail
		}
	}

	if ckRows, err = stMap[`LoadChecks`].Query(tk.meta.repoID, typ); err == sql.ErrNoRows {
		// go directly to loading instances since there are no source
		// checks on this type
		goto directinstances
	} else if err != nil {
		goto fail
	}
	defer ckRows.Close()

	// load all checks and start the assembly line
	for ckRows.Next() {
		if err = ckRows.Scan(
			&checkID,
			&nullBucketID,
			&srcCheckID,
			&srcObjType,
			&srcObjID,
			&configID,
			&capabilityID,
			&objID,
			&objType,
		); err != nil {
			goto fail
		}
		// save CheckConfig
		victim = proto.CheckConfig{
			Id:           configID,
			RepositoryId: tk.meta.repoID,
			CapabilityId: capabilityID,
			ObjectId:     objID,
			ObjectType:   objType,
		}
		if nullBucketID.Valid {
			victim.BucketId = nullBucketID.String
		}
		cfgMap[checkID] = victim
	}
	if ckRows.Err() != nil {
		goto fail
	}

	// iterate over the loaded checks and continue assembly with values
	// from the stored checkconfiguration
	for checkID = range cfgMap {
		if err = stMap[`LoadConfig`].QueryRow(cfgMap[checkID].Id, tk.meta.repoID).Scan(
			&nullBucketID,
			&cfgName,
			&cfgObjID,
			&cfgObjType,
			&isActive,
			&hasInheritance,
			&isChildrenOnly,
			&capabilityID,
			&interval,
			&isEnabled,
			&externalID,
		); err != nil {
			// sql.ErrNoRows is fatal here, the check exists - there
			// must be a configuration for it
			goto fail
		}

		victim = cfgMap[checkID]
		victim.Name = cfgName
		victim.Interval = uint64(interval)
		victim.IsActive = isActive
		victim.IsEnabled = isEnabled
		victim.Inheritance = hasInheritance
		victim.ChildrenOnly = isChildrenOnly
		victim.ExternalId = externalID
		cfgMap[checkID] = victim
	}

	// iterate over the loaded checks and continue assembly with values
	// from the stored thresholds
	for checkID = range cfgMap {
		if thrRows, err = stMap[`LoadThreshold`].Query(cfgMap[checkID].Id); err != nil {
			// sql.ErrNoRows is fatal here since a check without
			// thresholds is rather useless
			goto fail
		}

		victim = cfgMap[checkID]
		victim.Thresholds = []proto.CheckConfigThreshold{}

		// iterate over returned thresholds
		for thrRows.Next() {
			if err = thrRows.Scan(
				&predicate,
				&threshold,
				&levelName,
				&levelShort,
				&levelNumeric,
			); err != nil {
				thrRows.Close()
				goto fail
			}
			// ignore error since we converted this into the DB from int64
			numVal, _ = strconv.ParseInt(threshold, 10, 64)

			// save threshold
			victim.Thresholds = append(victim.Thresholds,
				proto.CheckConfigThreshold{
					Predicate: proto.Predicate{
						Symbol: predicate,
					},
					Level: proto.Level{
						Name:      levelName,
						ShortName: levelShort,
						Numeric:   uint16(levelNumeric),
					},
					Value: numVal,
				},
			)
		}
		if err = thrRows.Err(); err != nil {
			goto fail
		}
		cfgMap[checkID] = victim
	}

	// iterate over the loaded checks and continue assembly with values
	// from the stored constraints
	for checkID = range cfgMap {
		victim = cfgMap[checkID]
		victim.Constraints = []proto.CheckConfigConstraint{}
		for _, cstrType = range []string{`custom`, `native`, `oncall`, `attribute`, `service`, `system`} {
			switch cstrType {
			case `custom`:
				cstrRows, err = stMap[`LoadCustomCstr`].Query(cfgMap[checkID].Id)
			case `native`:
				cstrRows, err = stMap[`LoadNativeCstr`].Query(cfgMap[checkID].Id)
			case `oncall`:
				cstrRows, err = stMap[`LoadOncallCstr`].Query(cfgMap[checkID].Id)
			case `attribute`:
				cstrRows, err = stMap[`LoadAttributeCstr`].Query(cfgMap[checkID].Id)
			case `service`:
				cstrRows, err = stMap[`LoadServiceCstr`].Query(cfgMap[checkID].Id)
			case `system`:
				cstrRows, err = stMap[`LoadSystemCstr`].Query(cfgMap[checkID].Id)
			}
			if err != nil {
				goto fail
			}

			// iterate over returned constraints - no rows is valid, as
			// constraints are not mandatory
			for cstrRows.Next() {
				if err = cstrRows.Scan(&value1, &value2, &value3); err != nil {
					cstrRows.Close()
					goto fail
				}
				switch cstrType {
				case `custom`:
					victim.Constraints = append(victim.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Custom: &proto.PropertyCustom{
								Id:           value1,
								Name:         value2,
								RepositoryId: tk.meta.repoID,
								Value:        value3,
							},
						},
					)
				case `native`:
					victim.Constraints = append(victim.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Native: &proto.PropertyNative{
								Name:  value1,
								Value: value2,
							},
						},
					)
				case `oncall`:
					victim.Constraints = append(victim.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Oncall: &proto.PropertyOncall{
								Id:     value1,
								Name:   value2,
								Number: value3,
							},
						},
					)
				case `attribute`:
					victim.Constraints = append(victim.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Attribute: &proto.ServiceAttribute{
								Name:  value1,
								Value: value2,
							},
						},
					)
				case `service`:
					victim.Constraints = append(victim.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Service: &proto.PropertyService{
								Name:   value2,
								TeamId: value1,
							},
						},
					)
				case `system`:
					victim.Constraints = append(victim.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							System: &proto.PropertySystem{
								Name:  value1,
								Value: value2,
							},
						},
					)
				} // switch cstrType
			} // for cstrRows.Next()
			if cstrRows.Err() != nil {
				goto fail
			}
		}
		cfgMap[checkID] = victim
	}

	// iterate over the checks, convert them to tree.Check. Then load
	// the inherited IDs via loadItems and populate tree.Check.Items.
	// Save in datastructure: ckOrder, map[string]map[string]tree.Check
	//		objID -> checkID -> tree.Check
	// this way it is possible to access the checks by objID, which is
	// required to populate groups in the correct order.
	for checkID = range cfgMap {
		victim = cfgMap[checkID]
		if ckOrder[victim.ObjectId] == nil {
			ckOrder[victim.ObjectId] = map[string]tree.Check{}
		}
		if ckTree, err = tk.convertCheck(&victim); err != nil {
			goto fail
		}
		// add source check as well so it gets recreated with the
		// correct UUID
		ckItem = tree.CheckItem{ObjectType: victim.ObjectType}
		ckItem.ObjectId, _ = uuid.FromString(victim.ObjectId)
		ckItem.ItemId, _ = uuid.FromString(checkID)
		ckTree.Items = []tree.CheckItem{ckItem}
		tk.startLog.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, SrcCheckId=%s, CheckId=%s",
			tk.meta.repoName,
			`AssociateCheck`,
			ckItem.ObjectType,
			ckItem.ObjectId,
			checkID,
			ckItem.ItemId,
		)

		if itRows, err = stMap[`LoadItems`].Query(tk.meta.repoID, checkID); err != nil {
			goto fail
		}

		for itRows.Next() {
			if err = itRows.Scan(
				&itemID,
				&objID,
				&objType,
			); err != nil {
				itRows.Close()
				goto fail
			}

			// create new object per iteration
			ckItem := tree.CheckItem{ObjectType: objType}
			ckItem.ObjectId, _ = uuid.FromString(objID)
			ckItem.ItemId, _ = uuid.FromString(itemID)
			ckTree.Items = append(ckTree.Items, ckItem)
			tk.startLog.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, SrcCheckId=%s, CheckId=%s",
				tk.meta.repoName,
				`AssociateCheck`,
				objType,
				objID,
				checkID,
				itemID,
			)
		}
		if err = itRows.Err(); err != nil {
			goto fail
		}
		ckOrder[victim.ObjectId][checkID] = *ckTree
	}

	// apply all tree.Check object to the tree, special case
	// groups due to their ordering requirements
	//
	// grOrder maps from a standalone groupID to an array of child groupIDs
	// ckOrder maps from a groupID to all source checks on it
	// ==> not every group has to have a check
	// ==> groups can have more than one check
	switch typ {
	case "group":
		for objKey := range grOrder {
			// objKey is a standalone groupID. Test if there are
			// checks for it
			if _, ok := ckOrder[objKey]; ok {
				// apply all checks for objKey
				for ck := range ckOrder[objKey] {
					tk.tree.Find(tree.FindRequest{
						ElementType: cfgMap[ck].ObjectType,
						ElementId:   cfgMap[ck].ObjectId,
					}, true).SetCheck(ckOrder[objKey][ck])
					tk.startLog.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s",
						tk.meta.repoName,
						`SetCheck`,
						typ,
						objKey,
						ck,
					)
					if !tk.status.requiresRebuild {
						// drain after each check
						if tk.drain(`action`) != len(ckOrder[objKey][ck].Items) {
							tk.startLog.Printf("TK[%s]: Error=%s, Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s",
								tk.meta.repoName,
								`CheckCountMismatch`,
								`SetCheck`,
								typ,
								objKey,
								ck,
							)
							tk.status.isBroken = true
							return
						}
						if tk.drain(`error`) > 0 {
							goto fail
						}
					}
				}
			}
			// iterate through all childgroups of objKey
			for pos := range grOrder[objKey] {
				// test if there is a check for it
				if _, ok := ckOrder[grOrder[objKey][pos]]; ok {
					// apply all checks for grOrder[objKey][pos]
					for ck := range ckOrder[objKey] {
						tk.tree.Find(tree.FindRequest{
							ElementType: cfgMap[ck].ObjectType,
							ElementId:   cfgMap[ck].ObjectId,
						}, true).SetCheck(ckOrder[objKey][ck])
						tk.startLog.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s",
							tk.meta.repoName,
							`SetCheck`,
							typ,
							objKey,
							ck,
						)
						if !tk.status.requiresRebuild {
							// drain after each check
							if tk.drain(`action`) != len(ckOrder[objKey][ck].Items) {
								if tk.drain(`action`) != len(ckOrder[objKey][ck].Items) {
									tk.startLog.Printf("TK[%s]: Error=%s, Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s",
										tk.meta.repoName,
										`CheckCountMismatch`,
										`SetCheck`,
										typ,
										objKey,
										ck,
									)
									tk.status.isBroken = true
									return
								}
							}
							if tk.drain(`error`) > 0 {
								goto fail
							}
						}
					}
				}
			}
		}
		// iterate through all weird groups as well
		for objKey := range grWeird {
			// Test if there are checks for it
			if _, ok := ckOrder[objKey]; ok {
				// apply all checks for objKey
				for ck := range ckOrder[objKey] {
					tk.tree.Find(tree.FindRequest{
						ElementType: cfgMap[ck].ObjectType,
						ElementId:   cfgMap[ck].ObjectId,
					}, true).SetCheck(ckOrder[objKey][ck])
					tk.startLog.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s",
						tk.meta.repoName,
						`SetCheck`,
						typ,
						objKey,
						ck,
					)
					if !tk.status.requiresRebuild {
						// drain after each check
						if tk.drain(`action`) != len(ckOrder[objKey][ck].Items) {
							tk.startLog.Printf("TK[%s]: Error=%s, Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s",
								tk.meta.repoName,
								`CheckCountMismatch`,
								`SetCheck`,
								typ,
								objKey,
								ck,
							)
							tk.status.isBroken = true
							return
						}
						if tk.drain(`error`) > 0 {
							goto fail
						}
					}
				}
			}
		}
	default:
		for objKey := range ckOrder {
			for ck := range ckOrder[objKey] {
				tk.tree.Find(tree.FindRequest{
					ElementType: cfgMap[ck].ObjectType,
					ElementId:   cfgMap[ck].ObjectId,
				}, true).SetCheck(ckOrder[objKey][ck])
				tk.startLog.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s",
					tk.meta.repoName,
					`SetCheck`,
					typ,
					objKey,
					ck,
				)
				if !tk.status.requiresRebuild {
					// drain after each check
					if tk.drain(`action`) != len(ckOrder[objKey][ck].Items) {
						tk.startLog.Printf("TK[%s]: Error=%s, Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s",
							tk.meta.repoName,
							`CheckCountMismatch`,
							`SetCheck`,
							typ,
							objKey,
							ck,
						)
						tk.status.isBroken = true
						return
					}
					if tk.drain(`error`) > 0 {
						goto fail
					}
				}
			}
		}
	}

directinstances:
	// repository and bucket elements can not have check instances,
	// they are essentially metadata
	if typ == "repository" || typ == "bucket" {
		goto done
	}

	// iterate over all checks on this object type and load the check
	// instances they have created
	if tckRows, err = stMap[`LoadChecksForType`].Query(tk.meta.repoID, typ); err != nil {
		errLocation = fmt.Sprintf("Function=%s, repoId=%s, objType=%s", `loadTypeChecks.Query()`, tk.meta.repoID, typ)
		goto fail
	}

	for tckRows.Next() {
		// load check information
		if err = tckRows.Scan(
			&checkID,
			&objID,
		); err != nil {
			tckRows.Close()
			errLocation = `loadTypeChecks.Rows.Scan error`
			goto fail
		}

		// lookup instances for that check
		if inRows, err = stMap[`LoadInstances`].Query(checkID); err != nil {
			tckRows.Close()
			errLocation = fmt.Sprintf("Function=%s, repoId=%s, objType=%s, checkID=%s", `loadInstances.Query()`, tk.meta.repoID, typ, checkID)
			goto fail
		}

		for inRows.Next() {
			if err = inRows.Scan(
				&itemID,
				&configID,
			); err != nil {
				tckRows.Close()
				inRows.Close()
				errLocation = `loadInstances.Rows.Scan error`
				goto fail
			}

			// load configuration for check instance
			if err = stMap[`LoadInstanceCfg`].QueryRow(itemID).Scan(
				&itemCfgID,
				&version,
				&monitoringID,
				&cstrHash,
				&cstrValHash,
				&instSvc,
				&instSvcCfgHash,
				&instSvcCfg,
			); err != nil {
				// sql.ErrNoRows is fatal, an instance must have a
				// configuration
				tckRows.Close()
				inRows.Close()
				errLocation = fmt.Sprintf("Function=%s, repoId=%s, objType=%s, checkID=%s, instanceId=%s", `loadInstConfig.QueryRow()`, tk.meta.repoID, typ, checkID, itemID)
				goto fail
			}

			// fresh object per iteration -> memory safe!
			ckInstance := tree.CheckInstance{
				Version:            uint64(version),
				ConstraintHash:     cstrHash,
				ConstraintValHash:  cstrValHash,
				InstanceService:    instSvc,
				InstanceSvcCfgHash: instSvcCfgHash,
			}
			// if we have a saved service configuration, deserialize it
			if ckInstance.InstanceSvcCfgHash != "" {
				ckInstance.InstanceServiceConfig = make(map[string]string)
				if err = json.Unmarshal([]byte(instSvcCfg), &ckInstance.InstanceServiceConfig); err != nil {
					tckRows.Close()
					inRows.Close()
					errLocation = fmt.Sprintf("Function=%s, repoId=%s, objType=%s, checkID=%s, instanceId=%s, instCfgId=%s", `json.Unmarshal(InstanceServiceConfig)`, tk.meta.repoID, typ, checkID, itemID, itemCfgID)
					goto fail
				}
			}
			ckInstance.InstanceId, _ = uuid.FromString(itemID)
			ckInstance.CheckId, _ = uuid.FromString(checkID)
			ckInstance.ConfigId, _ = uuid.FromString(configID)
			ckInstance.InstanceConfigId, _ = uuid.FromString(itemCfgID)

			// attach instance to tree
			tk.tree.Find(tree.FindRequest{
				ElementType: typ,
				ElementId:   objID,
			}, true).LoadInstance(ckInstance)
			tk.startLog.Printf("TK[%s]: Action=%s, ObjectType=%s, ObjectId=%s, CheckId=%s, InstanceId=%s",
				tk.meta.repoName,
				`LoadInstance`,
				typ,
				objID,
				ckInstance.CheckId.String(),
				ckInstance.InstanceId.String(),
			)
		}
		if err = inRows.Err(); err != nil {
			inRows.Close()
			errLocation = fmt.Sprintf("Function=%s, repoId=%s, objType=%s, checkID=%s", `checkInstanceRows.Iterate.Error`, tk.meta.repoID, typ, checkID)
			goto fail
		}
	}
	if err = tckRows.Err(); err != nil {
		errLocation = fmt.Sprintf("Function=%s, repoId=%s, objType=%s", `checksForType.Iterate.Error()`, tk.meta.repoID, typ)
		goto fail
	}

done:
	return

fail:
	tk.status.isBroken = true
	if err != nil {
		tk.startLog.Println(`BROKEN REPOSITORY ERROR: `, errLocation, err)
	}
	return
}

func (tk *TreeKeeper) startupScopedReapplyCheckConfig(typ string, stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err                                         error
		configRows, threshRows, cstrRows            *sql.Rows
		predicate, threshold, levelName, levelShort string
		cstrType, value1, value2, value3            string
		levelNumeric, numVal                        int64
		treeCheck                                   *tree.Check
		nullBucketID                                sql.NullString
	)

	// 1. identify check configurations to load:
	//    select configuration_id from soma.check_configurations where
	//    configuration_object_type = typ and repository_id =
	//    tk.meta.repoID and not deleted
	if configRows, err = stMap[`LoadAllConfigsForType`].Query(typ, tk.meta.repoID); err != nil {
		goto fail
	}
	defer configRows.Close()

	for configRows.Next() {
		conf := proto.CheckConfig{}
		conf.Thresholds = []proto.CheckConfigThreshold{}
		conf.Constraints = []proto.CheckConfigConstraint{}
		conf.RepositoryId = tk.meta.repoID
		conf.ObjectType = typ
		if err = configRows.Scan(
			&conf.Id,
			&nullBucketID,
			&conf.Name,
			&conf.ObjectId,
			&conf.Inheritance,
			&conf.ChildrenOnly,
			&conf.CapabilityId,
			&conf.Interval,
			&conf.IsEnabled,
			&conf.ExternalId,
		); err != nil {
			goto fail
		}
		if nullBucketID.Valid {
			conf.BucketId = nullBucketID.String
		}
		tk.startLog.Printf("TK[%s]: rebuild processing check configuration %s", tk.meta.repoName, conf.Id)
		// 2. assemble proto.CheckConfig object:
		//    + thresholds
		if threshRows, err = stMap[`LoadThreshold`].Query(conf.Id); err != nil {
			goto fail
		}

		for threshRows.Next() {
			if err = threshRows.Scan(
				&predicate,
				&threshold,
				&levelName,
				&levelShort,
				&levelNumeric,
			); err != nil {
				threshRows.Close()
				goto fail
			}
			// ignore errors, we converted into the DB from int64
			numVal, _ = strconv.ParseInt(threshold, 10, 64)

			// add threshold to config
			conf.Thresholds = append(conf.Thresholds,
				proto.CheckConfigThreshold{
					Predicate: proto.Predicate{
						Symbol: predicate,
					},
					Level: proto.Level{
						Name:      levelName,
						ShortName: levelShort,
						Numeric:   uint16(levelNumeric),
					},
					Value: numVal,
				},
			)
		}
		if err = threshRows.Err(); err != nil {
			goto fail
		}

		//    + constraints
		for _, cstrType = range []string{`custom`, `native`, `oncall`, `attribute`, `service`, `system`} {
			switch cstrType {
			case `custom`:
				cstrRows, err = stMap[`LoadCustomCstr`].Query(conf.Id)
			case `native`:
				cstrRows, err = stMap[`LoadNativeCstr`].Query(conf.Id)
			case `oncall`:
				cstrRows, err = stMap[`LoadOncallCstr`].Query(conf.Id)
			case `attribute`:
				cstrRows, err = stMap[`LoadAttributeCstr`].Query(conf.Id)
			case `service`:
				cstrRows, err = stMap[`LoadServiceCstr`].Query(conf.Id)
			case `system`:
				cstrRows, err = stMap[`LoadSystemCstr`].Query(conf.Id)
			}
			if err != nil {
				goto fail
			}

			for cstrRows.Next() {
				if err = cstrRows.Scan(&value1, &value2, &value3); err != nil {
					cstrRows.Close()
					goto fail
				}
				switch cstrType {
				case `custom`:
					conf.Constraints = append(conf.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Custom: &proto.PropertyCustom{
								Id:           value1,
								Name:         value2,
								RepositoryId: tk.meta.repoID,
								Value:        value3,
							},
						},
					)
				case `native`:
					conf.Constraints = append(conf.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Native: &proto.PropertyNative{
								Name:  value1,
								Value: value2,
							},
						},
					)
				case `oncall`:
					conf.Constraints = append(conf.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Oncall: &proto.PropertyOncall{
								Id:     value1,
								Name:   value2,
								Number: value3,
							},
						},
					)
				case `attribute`:
					conf.Constraints = append(conf.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Attribute: &proto.ServiceAttribute{
								Name:  value1,
								Value: value2,
							},
						},
					)
				case `service`:
					conf.Constraints = append(conf.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							Service: &proto.PropertyService{
								Name:   value2,
								TeamId: value1,
							},
						},
					)
				case `system`:
					conf.Constraints = append(conf.Constraints,
						proto.CheckConfigConstraint{
							ConstraintType: cstrType,
							System: &proto.PropertySystem{
								Name:  value1,
								Value: value2,
							},
						},
					)
				} // switch cstrType
			} // for cstrRows.Next()
			if cstrRows.Err() != nil {
				goto fail
			}
		}
		// 3. convert to treecheck
		if treeCheck, err = tk.convertCheck(&conf); err == nil {
			// 4. apply config
			tk.tree.Find(tree.FindRequest{
				ElementType: conf.ObjectType,
				ElementId:   conf.ObjectId,
			}, true).SetCheck(*treeCheck)
		} else {
			tk.startLog.Printf("Rebuild error during check conversion: %s", err)
		}
	}
	return

fail:
	tk.status.isBroken = true
	if err != nil {
		tk.startLog.Printf("Error during rebuild loading of checks: %s", err)
	}
}

// orderGroups orders the groups in a repository so they can be
// processed from root to leaf
func (tk *TreeKeeper) orderGroups(stMap map[string]*sql.Stmt) (map[string][]string, map[string]string, error) {
	if tk.status.isBroken {
		return nil, nil, fmt.Errorf("Broken tree detected")
	}

	var (
		err                                                 error
		groupID, groupState, parentID, childID, candidateID string
		stRows, rlRows                                      *sql.Rows
		ok                                                  bool
		grStateMap, grRelMap, grWeirdMap                    map[string]string
		grOrder                                             map[string][]string
		children                                            []string
		oldLen, sameCount                                   int
	)

	grStateMap = map[string]string{}
	grRelMap = map[string]string{}
	grWeirdMap = map[string]string{}
	grOrder = map[string][]string{}

	// load groups in this repository
	if stRows, err = stMap[`LoadGroupState`].Query(tk.meta.repoID); err != nil {
		tk.status.isBroken = true
		return nil, nil, err
	}
	defer stRows.Close()

	for stRows.Next() {
		if err = stRows.Scan(&groupID, &groupState); err != nil {
			// error loading group state
			tk.status.isBroken = true
			return nil, nil, err
		}
		grStateMap[groupID] = groupState
	}
	if err = stRows.Err(); err != nil {
		tk.status.isBroken = true
		return nil, nil, err
	}
	if len(grStateMap) == 0 {
		// repository has no groups, return empty handed
		return grOrder, grWeirdMap, nil
	}

	// load relations between groups in this repository
	if rlRows, err = stMap[`LoadGroupRelations`].Query(tk.meta.repoID); err != nil {
		tk.status.isBroken = true
		return nil, nil, err
	}
	defer rlRows.Close()

	for rlRows.Next() {
		if err = rlRows.Scan(&groupID, &childID); err != nil {
			// error loading relations
			tk.status.isBroken = true
			return nil, nil, err
		}
		// every group can only be child to one parent group, but
		// a parent group can have multiple child groups => data
		// needs to be stored this way
		grRelMap[childID] = groupID
	}
	if err = rlRows.Err(); err != nil {
		tk.status.isBroken = true
		return nil, nil, err
	}

	// iterate over all groups and identify state standalone,
	// attached to the bucket. once ordered, remove them.
	for groupID, groupState = range grStateMap {
		switch groupState {
		case "grouped":
			continue
		case "standalone":
			grOrder[groupID] = []string{}
		default:
			// groups should really not be in a different state
			grWeirdMap[groupID] = groupState
		}
		delete(grStateMap, groupID)
	}

	// simplified first ordering round, only sort children of
	// standalone groups
	for childID, groupID = range grRelMap {
		if _, ok = grOrder[groupID]; !ok {
			// groupID is not standalone
			continue
		}
		// groupID is standalone
		grOrder[groupID] = append(grOrder[groupID], childID)
		delete(grRelMap, childID)
		delete(grStateMap, childID)
	}

	// breaker switch variables to short-circuit an unbounded loop
	oldLen = 0
	sameCount = 0

	// full ordering, grStateMap contains all grouped groups left
orderloop:
	for len(grStateMap) > 0 {
		// install a breaker switch in case the groups can not be
		// ordered. if no elements were deleted from grStateMap
		// for three full iterations => give up
		// XXX 3 was chosen via dice roll
		if len(grStateMap) == oldLen {
			sameCount++
		} else {
			oldLen = len(grStateMap)
			sameCount = 0
		}
		if sameCount >= 3 {
			break orderloop
		}

		// iterate through all unordered children
	childloop:
		for childID, parentID = range grRelMap {
			// since childID was not ordered during the first
			// pass, its parentID is a grouped group itself
			for groupID, children = range grOrder {
				// iterate through all children
				for _, candidateID = range children {
					if candidateID == parentID {
						// this candidateID is the searched parent.
						// since candidateID is a child of
						// groupID, childID is appended there
						grOrder[groupID] = append(grOrder[groupID], childID)
						delete(grStateMap, childID)
						delete(grRelMap, childID)
						continue childloop
					}
				}
			}
		}
	}
	if sameCount >= 3 || len(grRelMap) != 0 {
		// breaker went off or we have unordered grRelMap left
		tk.status.isBroken = true
		return nil, nil, fmt.Errorf("Failed to order groups")
	}

	// return order
	return grOrder, grWeirdMap, nil
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
