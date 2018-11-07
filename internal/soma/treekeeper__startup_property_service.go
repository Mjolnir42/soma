package soma

import (
	"database/sql"

	"github.com/mjolnir42/soma/internal/tree"
	"github.com/mjolnir42/soma/lib/proto"
	uuid "github.com/satori/go.uuid"
)

func (tk *TreeKeeper) startupServiceProperties(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err                                          error
		instanceID, srcInstanceID, objectID, view    string
		inInstanceID, inObjectType, inObjID, attrKey string
		serviceProperty, teamID, attrValue           string
		inheritance, childrenOnly                    bool
		rows, attributeRows, instanceRows            *sql.Rows
	)

	for loopType, loopStmt := range map[string][2]string{
		`repository`: [2]string{
			`LoadPropRepoService`,
			`LoadPropRepoSvcAttr`},
		`bucket`: [2]string{
			`LoadPropBuckService`,
			`LoadPropBuckSvcAttr`},
		`group`: [2]string{
			`LoadPropGrpService`,
			`LoadPropGrpSvcAttr`},
		`cluster`: [2]string{
			`LoadPropClrService`,
			`LoadPropClrSvcAttr`},
		`node`: [2]string{
			`LoadPropNodeService`,
			`LoadPropNodeSvcAttr`},
	} {

		tk.startLog.Printf("TK[%s]: loading %s service properties", tk.meta.repoName, loopType)
		rows, err = stMap[loopStmt[0]].Query(tk.meta.repoID)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading %s service properties: %s", tk.meta.repoName, loopType, err.Error())
			tk.status.isBroken = true
			return
		}
		defer rows.Close()

	serviceloop:
		// load all service properties defined directly on the object
		for rows.Next() {
			err = rows.Scan(
				&instanceID,
				&srcInstanceID,
				&objectID,
				&view,
				&serviceProperty,
				&teamID,
				&inheritance,
				&childrenOnly,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break serviceloop
				}
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
				tk.status.isBroken = true
				return
			}

			// build the property
			prop := tree.PropertyService{
				Inheritance:  inheritance,
				ChildrenOnly: childrenOnly,
				View:         view,
				Service:      serviceProperty,
			}
			prop.ID, _ = uuid.FromString(instanceID)
			prop.Attributes = make([]proto.ServiceAttribute, 0)
			prop.Instances = make([]tree.PropertyInstance, 0)

			attributeRows, err = stMap[loopStmt[1]].Query(
				teamID,
				serviceProperty,
			)
			if err != nil {
				tk.startLog.Printf("TK[%s] Error loading %s service properties: %s", tk.meta.repoName, loopType, err.Error())
				tk.status.isBroken = true
				return
			}
			defer attributeRows.Close()

		attributeloop:
			// load service attributes
			for attributeRows.Next() {
				err = attributeRows.Scan(
					&attrKey,
					&attrValue,
				)
				if err != nil {
					if err == sql.ErrNoRows {
						break attributeloop
					}
					tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
					tk.status.isBroken = true
					return
				}

				pa := proto.ServiceAttribute{
					Name:  attrKey,
					Value: attrValue,
				}
				prop.Attributes = append(prop.Attributes, pa)
			}

			instanceRows, err = stMap[`LoadPropSvcInstance`].Query(
				tk.meta.repoID,
				srcInstanceID,
			)
			if err != nil {
				tk.startLog.Printf("TK[%s] Error loading %s service properties: %s", tk.meta.repoName, loopType, err.Error())
				tk.status.isBroken = true
				return
			}
			defer instanceRows.Close()

		inproploop:
			// load all all ids for properties that were inherited from the
			// current service property so the IDs can be set correctly
			for instanceRows.Next() {
				err = instanceRows.Scan(
					&inInstanceID,
					&inObjectType,
					&inObjID,
				)
				if err != nil {
					if err == sql.ErrNoRows {
						break inproploop
					}
					tk.startLog.Printf("TK[%s] Error: %s", tk.meta.repoName, err.Error())
					tk.status.isBroken = true
					return
				}

				var propObjectID, propInstanceID uuid.UUID
				if propObjectID, err = uuid.FromString(inObjID); err != nil {
					tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
					tk.status.isBroken = true
					return
				}
				if propInstanceID, err = uuid.FromString(inInstanceID); err != nil {
					tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
					tk.status.isBroken = true
					return
				}
				if uuid.Equal(uuid.Nil, propObjectID) || uuid.Equal(uuid.Nil, propInstanceID) {
					continue inproploop
				}
				if inObjectType == "MAGIC_NO_RESULT_VALUE" {
					continue inproploop
				}

				pi := tree.PropertyInstance{
					ObjectID:   propObjectID,
					ObjectType: inObjectType,
					InstanceID: propInstanceID,
				}
				prop.Instances = append(prop.Instances, pi)
			}

			// lookup the object and set the prepared property
			tk.tree.Find(tree.FindRequest{
				ElementType: loopType,
				ElementID:   objectID,
			}, true).SetProperty(&prop)

			// throw away all generated actions, we do this for every
			// property since with inheritance this can create a lot of
			// actions
			tk.drain(`action`)
			tk.drain(`error`)
		}
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
