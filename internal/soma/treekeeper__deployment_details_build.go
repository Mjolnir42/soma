package soma

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/mjolnir42/soma/lib/proto"
)

func (tk *TreeKeeper) buildDeploymentDetails() {
	var (
		err                                                 error
		instanceCfgID                                       string
		objID, objType                                      string
		rows, thresh, pkgs, gSysProps, cSysProps, nSysProps *sql.Rows
		gCustProps, cCustProps, nCustProps                  *sql.Rows
		callback                                            sql.NullString
	)

	// TODO:
	// * refactoring switch objType {} block
	// * SQL error handling

	if rows, err = tk.stmtList.Query(tk.meta.repoID); err != nil {
		tk.status.isBroken = true
		return
	}
	defer rows.Close()

deploymentbuilder:
	for rows.Next() {
		detail := proto.Deployment{}

		err = rows.Scan(
			&instanceCfgID,
		)
		if err != nil {
			tk.treeLog.Println(`tk.stmtList.Query().Scan():`, err)
			break deploymentbuilder
		}

		//
		detail.CheckInstance = &proto.CheckInstance{
			InstanceConfigID: instanceCfgID,
		}
		tk.stmtCheckInstance.QueryRow(instanceCfgID).Scan(
			&detail.CheckInstance.Version,
			&detail.CheckInstance.InstanceID,
			&detail.CheckInstance.ConstraintHash,
			&detail.CheckInstance.ConstraintValHash,
			&detail.CheckInstance.InstanceService,
			&detail.CheckInstance.InstanceSvcCfgHash,
			&detail.CheckInstance.InstanceServiceConfig,
			&detail.CheckInstance.CheckID,
			&detail.CheckInstance.ConfigID,
		)

		//
		detail.Check = &proto.Check{
			CheckID: detail.CheckInstance.CheckID,
		}
		tk.stmtCheck.QueryRow(detail.CheckInstance.CheckID).Scan(
			&detail.Check.RepositoryID,
			&detail.Check.SourceCheckID,
			&detail.Check.SourceType,
			&detail.Check.InheritedFrom,
			&detail.Check.CapabilityID,
			&objID,
			&objType,
			&detail.Check.Inheritance,
			&detail.Check.ChildrenOnly,
		)
		detail.ObjectType = objType
		if detail.Check.InheritedFrom != objID {
			detail.Check.IsInherited = true
		}
		detail.Check.CheckConfigID = detail.CheckInstance.ConfigID

		//
		detail.CheckConfig = &proto.CheckConfig{
			ID:           detail.Check.CheckConfigID,
			RepositoryID: detail.Check.RepositoryID,
			BucketID:     detail.Check.BucketID,
			CapabilityID: detail.Check.CapabilityID,
			ObjectID:     objID,
			ObjectType:   objType,
			Inheritance:  detail.Check.Inheritance,
			ChildrenOnly: detail.Check.ChildrenOnly,
		}
		tk.stmtCheckConfig.QueryRow(detail.Check.CheckConfigID).Scan(
			&detail.CheckConfig.Name,
			&detail.CheckConfig.Interval,
			&detail.CheckConfig.IsActive,
			&detail.CheckConfig.IsEnabled,
			&detail.CheckConfig.ExternalID,
		)

		//
		detail.CheckConfig.Thresholds = []proto.CheckConfigThreshold{}
		thresh, err = tk.stmtThreshold.Query(detail.CheckConfig.ID)
		if err != nil {
			// a check config must have 1+ thresholds
			tk.treeLog.Println(`DANGER WILL ROBINSON!`,
				`Failed to get thresholds for:`, detail.CheckConfig.ID)
			continue deploymentbuilder
		}
		defer thresh.Close()

		for thresh.Next() {
			thr := proto.CheckConfigThreshold{
				Predicate: proto.Predicate{},
				Level:     proto.Level{},
			}

			err = thresh.Scan(
				&thr.Predicate.Symbol,
				&thr.Value,
				&thr.Level.Name,
				&thr.Level.ShortName,
				&thr.Level.Numeric,
			)
			if err != nil {
				tk.treeLog.Println(`tk.stmtThreshold.Query().Scan():`, err)
				break deploymentbuilder
			}
			detail.CheckConfig.Thresholds = append(detail.CheckConfig.Thresholds, thr)
		}

		// XXX TODO
		//detail.CheckConfiguration.Constraints = []somaproto.CheckConfigurationConstraint{}
		detail.CheckConfig.Constraints = nil

		//
		detail.Capability = &proto.Capability{
			ID: detail.Check.CapabilityID,
		}
		detail.Monitoring = &proto.Monitoring{}
		detail.Metric = &proto.Metric{}
		detail.Unit = &proto.Unit{}
		tk.stmtCapMonMetric.QueryRow(detail.Capability.ID).Scan(
			&detail.Capability.Metric,
			&detail.Capability.MonitoringID,
			&detail.Capability.View,
			&detail.Capability.Thresholds,
			&detail.Monitoring.Name,
			&detail.Monitoring.Mode,
			&detail.Monitoring.Contact,
			&detail.Monitoring.TeamID,
			&callback,
			&detail.Metric.Unit,
			&detail.Metric.Description,
			&detail.Unit.Name,
		)
		if callback.Valid {
			detail.Monitoring.Callback = callback.String
		} else {
			detail.Monitoring.Callback = ""
		}
		detail.Unit.Unit = detail.Metric.Unit
		detail.Metric.Path = detail.Capability.Metric
		detail.Monitoring.ID = detail.Capability.MonitoringID
		detail.Capability.Name = fmt.Sprintf("%s.%s.%s",
			detail.Monitoring.Name,
			detail.Capability.View,
			detail.Metric.Path,
		)
		detail.View = detail.Capability.View

		//
		detail.Metric.Packages = &[]proto.MetricPackage{}
		pkgs, _ = tk.stmtPkgs.Query(detail.Metric.Path)
		defer pkgs.Close()

		for pkgs.Next() {
			pkg := proto.MetricPackage{}

			err = pkgs.Scan(
				&pkg.Provider,
				&pkg.Name,
			)
			if err != nil {
				tk.treeLog.Println(`tk.stmtPkgs.Query().Scan():`, err)
				break deploymentbuilder
			}
			*detail.Metric.Packages = append(*detail.Metric.Packages, pkg)
		}

		//
		detail.Oncall = &proto.Oncall{}
		detail.Service = &proto.PropertyService{}
		switch objType {
		case "group":
			// fetch the group object
			detail.Group = &proto.Group{
				ID: objID,
			}
			tk.stmtGroup.QueryRow(objID).Scan(
				&detail.Group.BucketID,
				&detail.Group.Name,
				&detail.Group.ObjectState,
				&detail.Group.TeamID,
				&detail.Bucket,
				&detail.Environment,
				&detail.Repository,
			)
			// fetch team information
			detail.Team = &proto.Team{
				ID: detail.Group.TeamID,
			}
			// fetch oncall information if the property is set,
			// otherwise cleanup detail.Oncall
			err = tk.stmtGroupOncall.QueryRow(
				detail.Group.ID,
				detail.View,
			).Scan(
				&detail.Oncall.ID,
				&detail.Oncall.Name,
				&detail.Oncall.Number,
			)
			if err == sql.ErrNoRows {
				detail.Oncall = nil
			} else if err != nil {
				tk.treeLog.Println(`tk.stmtGroupOncall.QueryRow():`, err)
				break deploymentbuilder
			}
			// fetch service name, and attributes if applicable
			if detail.CheckInstance.InstanceService != "" {
				err = tk.stmtGroupService.QueryRow(
					detail.CheckInstance.InstanceService,
					detail.View,
				).Scan(
					&detail.Service.Name,
					&detail.Service.TeamID,
				)
				if err == sql.ErrNoRows {
					detail.Service = nil
				} else if err != nil {
					tk.treeLog.Println(`tk.stmtGroupService.QueryRow():`, err)
					break deploymentbuilder
				} else {
					detail.Service.Attributes = []proto.ServiceAttribute{}
					fm := map[string]string{}
					_ = json.Unmarshal([]byte(detail.CheckInstance.InstanceServiceConfig), &fm)
					for k, v := range fm {
						a := proto.ServiceAttribute{
							Name:  k,
							Value: v,
						}
						detail.Service.Attributes = append(detail.Service.Attributes, a)
					}
				}
			}
			// fetch system properties
			detail.Properties = &[]proto.PropertySystem{}
			gSysProps, _ = tk.stmtGroupSysProp.Query(detail.Group.ID, detail.View)
			defer gSysProps.Close()

			for gSysProps.Next() {
				prop := proto.PropertySystem{}
				err = gSysProps.Scan(
					&prop.Name,
					&prop.Value,
				)
				if err != nil {
					tk.treeLog.Println(`tk.stmtGroupSysProp.Query().Scan():`, err)
					break deploymentbuilder
				}
				*detail.Properties = append(*detail.Properties, prop)
				if prop.Name == "group_datacenter" {
					detail.Datacenter = prop.Value
				}
			}
			if len(*detail.Properties) == 0 {
				detail.Properties = nil
			}
			// fetch custom properties
			detail.CustomProperties = &[]proto.PropertyCustom{}
			gCustProps, _ = tk.stmtGroupCustProp.Query(detail.Group.ID, detail.View)
			defer gCustProps.Close()

			for gCustProps.Next() {
				prop := proto.PropertyCustom{}
				err = gCustProps.Scan(
					&prop.ID,
					&prop.Name,
					&prop.Value,
				)
				if err != nil {
					tk.treeLog.Println(`tk.stmtGroupCustProp.Query().Scan():`, err)
					break deploymentbuilder
				}
				*detail.CustomProperties = append(*detail.CustomProperties, prop)
			}
			if len(*detail.CustomProperties) == 0 {
				detail.CustomProperties = nil
			}
		case "cluster":
			// fetch the cluster object
			detail.Cluster = &proto.Cluster{
				ID: objID,
			}
			tk.stmtCluster.QueryRow(objID).Scan(
				&detail.Cluster.Name,
				&detail.Cluster.BucketID,
				&detail.Cluster.ObjectState,
				&detail.Cluster.TeamID,
				&detail.Bucket,
				&detail.Environment,
				&detail.Repository,
			)
			// fetch team information
			detail.Team = &proto.Team{
				ID: detail.Cluster.TeamID,
			}
			// fetch oncall information if the property is set,
			// otherwise cleanup detail.Oncall
			err = tk.stmtClusterOncall.QueryRow(detail.Cluster.ID, detail.View).Scan(
				&detail.Oncall.ID,
				&detail.Oncall.Name,
				&detail.Oncall.Number,
			)
			if err != nil {
				detail.Oncall = nil
			}
			// fetch the service name, and attributes if applicable
			if detail.CheckInstance.InstanceService != "" {
				err = tk.stmtClusterService.QueryRow(
					detail.CheckInstance.InstanceService,
					detail.View,
				).Scan(
					&detail.Service.Name,
					&detail.Service.TeamID,
				)
				if err != nil {
					detail.Service = nil
				} else {
					detail.Service.Attributes = []proto.ServiceAttribute{}
					fm := map[string]string{}
					_ = json.Unmarshal([]byte(detail.CheckInstance.InstanceServiceConfig), &fm)
					for k, v := range fm {
						a := proto.ServiceAttribute{
							Name:  k,
							Value: v,
						}
						detail.Service.Attributes = append(detail.Service.Attributes, a)
					}
				}
			}
			// fetch system properties
			detail.Properties = &[]proto.PropertySystem{}
			cSysProps, _ = tk.stmtClusterSysProp.Query(detail.Cluster.ID, detail.View)
			defer cSysProps.Close()

			for cSysProps.Next() {
				prop := proto.PropertySystem{}
				err = cSysProps.Scan(
					&prop.Name,
					&prop.Value,
				)
				if err != nil {
					tk.treeLog.Println(`tk.stmtClusterSysProp.Query().Scan():`, err)
					break deploymentbuilder
				}
				*detail.Properties = append(*detail.Properties, prop)
				if prop.Name == "cluster_datacenter" {
					detail.Datacenter = prop.Value
				}
			}
			if len(*detail.Properties) == 0 {
				detail.Properties = nil
			}
			// fetch custom properties
			detail.CustomProperties = &[]proto.PropertyCustom{}
			cCustProps, _ = tk.stmtClusterCustProp.Query(detail.Cluster.ID, detail.View)
			defer cCustProps.Close()

			for cCustProps.Next() {
				prop := proto.PropertyCustom{}
				cCustProps.Scan(
					&prop.ID,
					&prop.Name,
					&prop.Value,
				)
				*detail.CustomProperties = append(*detail.CustomProperties, prop)
			}
			if len(*detail.CustomProperties) == 0 {
				detail.CustomProperties = nil
			}
		case "node":
			// fetch the node object
			detail.Server = &proto.Server{}
			detail.Node = &proto.Node{
				ID: objID,
			}
			tk.stmtNode.QueryRow(objID).Scan(
				&detail.Node.AssetID,
				&detail.Node.Name,
				&detail.Node.TeamID,
				&detail.Node.ServerID,
				&detail.Node.State,
				&detail.Node.IsOnline,
				&detail.Node.IsDeleted,
				&detail.Bucket,
				&detail.Environment,
				&detail.Repository,
				&detail.Server.AssetID,
				&detail.Server.Datacenter,
				&detail.Server.Location,
				&detail.Server.Name,
				&detail.Server.IsOnline,
				&detail.Server.IsDeleted,
			)
			detail.Server.ID = detail.Node.ServerID
			detail.Datacenter = detail.Server.Datacenter
			// fetch team information
			detail.Team = &proto.Team{
				ID: detail.Node.TeamID,
			}
			// fetch oncall information if the property is set,
			// otherwise cleanup detail.Oncall
			err = tk.stmtNodeOncall.QueryRow(detail.Node.ID, detail.View).Scan(
				&detail.Oncall.ID,
				&detail.Oncall.Name,
				&detail.Oncall.Number,
			)
			if err != nil {
				detail.Oncall = nil
			}
			// fetch the service name, and attributes if applicable
			if detail.CheckInstance.InstanceService != "" {
				err = tk.stmtNodeService.QueryRow(
					detail.CheckInstance.InstanceService,
					detail.View,
				).Scan(
					&detail.Service.Name,
					&detail.Service.TeamID,
				)
				if err != nil {
					detail.Service = nil
				} else {
					detail.Service.Attributes = []proto.ServiceAttribute{}
					fm := map[string]string{}
					_ = json.Unmarshal([]byte(detail.CheckInstance.InstanceServiceConfig), &fm)
					for k, v := range fm {
						a := proto.ServiceAttribute{
							Name:  k,
							Value: v,
						}
						detail.Service.Attributes = append(detail.Service.Attributes, a)
					}
				}
			}
			// fetch system properties
			detail.Properties = &[]proto.PropertySystem{}
			nSysProps, _ = tk.stmtNodeSysProp.Query(detail.Node.ID, detail.View)
			defer nSysProps.Close()

			for nSysProps.Next() {
				prop := proto.PropertySystem{}
				err = nSysProps.Scan(
					&prop.Name,
					&prop.Value,
				)
				if err != nil {
					tk.treeLog.Println(`tk.stmtNodeSysProp.Query().Scan():`, err)
					break deploymentbuilder
				}
				*detail.Properties = append(*detail.Properties, prop)
			}
			if len(*detail.Properties) == 0 {
				detail.Properties = nil
			}
			// fetch custom properties
			detail.CustomProperties = &[]proto.PropertyCustom{}
			nCustProps, _ = tk.stmtNodeCustProp.Query(detail.Node.ID, detail.View)
			defer nCustProps.Close()

			for nCustProps.Next() {
				prop := proto.PropertyCustom{}
				nCustProps.Scan(
					&prop.ID,
					&prop.Name,
					&prop.Value,
				)
				*detail.CustomProperties = append(*detail.CustomProperties, prop)
			}
			if len(*detail.CustomProperties) == 0 {
				detail.CustomProperties = nil
			}
		}

		tk.stmtTeam.QueryRow(detail.Team.ID).Scan(
			&detail.Team.Name,
			&detail.Team.LdapID,
		)

		// if no datacenter information was gathered, use the default DC
		if detail.Datacenter == "" {
			tk.stmtDefaultDC.QueryRow().Scan(&detail.Datacenter)
		}

		// build JSON of DeploymentDetails
		var detailJSON []byte
		if detailJSON, err = json.Marshal(&detail); err != nil {
			tk.treeLog.Println(`Failed to JSON marshal deployment details:`,
				detail.CheckInstance.InstanceConfigID, err)
			break deploymentbuilder
		}
		if _, err = tk.stmtUpdate.Exec(
			detailJSON,
			detail.Monitoring.ID,
			detail.CheckInstance.InstanceConfigID,
		); err != nil {
			tk.treeLog.Println(`Failed to save DeploymentDetails.JSON:`,
				detail.CheckInstance.InstanceConfigID, err)
			break deploymentbuilder
		}
	}
	// mark the tree as broken to prevent further data processing
	if err != nil {
		tk.status.isBroken = true
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
