package soma

import (
	"database/sql"

	"github.com/mjolnir42/soma/internal/tree"
	uuid "github.com/satori/go.uuid"
)

func (tk *TreeKeeper) startupOncallProperties(stMap map[string]*sql.Stmt) {
	if tk.status.isBroken {
		return
	}

	var (
		err                                                           error
		instanceID, srcInstanceID, objectID, view, oncallID           string
		inInstanceID, inObjectType, inObjID, oncallName, oncallNumber string
		inheritance, childrenOnly                                     bool
		rows, instanceRows                                            *sql.Rows
	)

	for loopType, loopStmt := range map[string]string{
		`repository`: `LoadPropRepoOncall`,
		`bucket`:     `LoadPropBuckOncall`,
		`group`:      `LoadPropGrpOncall`,
		`cluster`:    `LoadPropClrOncall`,
		`node`:       `LoadPropNodeOncall`,
	} {

		tk.startLog.Printf("TK[%s]: loading %s oncall properties\n", tk.meta.repoName, loopType)
		rows, err = stMap[loopStmt].Query(tk.meta.repoID)
		if err != nil {
			tk.startLog.Printf("TK[%s] Error loading %s oncall properties: %s", tk.meta.repoName, loopType, err.Error())
			tk.status.isBroken = true
			return
		}
		defer rows.Close()

	oncallloop:
		// load all oncall properties defined directly on objects
		for rows.Next() {
			err = rows.Scan(
				&instanceID,
				&srcInstanceID,
				&objectID,
				&view,
				&oncallID,
				&inheritance,
				&childrenOnly,
				&oncallName,
				&oncallNumber,
			)
			if err != nil {
				if err == sql.ErrNoRows {
					break oncallloop
				}
				tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
				tk.status.isBroken = true
				return
			}

			// build the property
			prop := tree.PropertyOncall{
				Inheritance:  inheritance,
				ChildrenOnly: childrenOnly,
				View:         view,
				Name:         oncallName,
				Number:       oncallNumber,
			}
			prop.Id, _ = uuid.FromString(instanceID)
			prop.OncallId, _ = uuid.FromString(oncallID)
			prop.Instances = make([]tree.PropertyInstance, 0)

			instanceRows, err = stMap[`LoadPropOncallInstance`].Query(
				tk.meta.repoID,
				srcInstanceID,
			)
			if err != nil {
				tk.startLog.Printf("TK[%s] Error loading %s oncall properties: %s", tk.meta.repoName, loopType, err.Error())
				tk.status.isBroken = true
				return
			}
			defer instanceRows.Close()

		inproploop:
			// load all all ids for properties that were inherited from the
			// current oncall property so the IDs can be set correctly
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
					tk.startLog.Printf("TK[%s] Error: %s\n", tk.meta.repoName, err.Error())
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
					ObjectId:   propObjectID,
					ObjectType: inObjectType,
					InstanceId: propInstanceID,
				}
				prop.Instances = append(prop.Instances, pi)
			}

			// lookup the object and set the prepared property
			tk.tree.Find(tree.FindRequest{
				ElementType: loopType,
				ElementId:   objectID,
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
