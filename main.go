package main

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// global variables
var conn *sql.DB
var handlerMap = make(map[string]interface{})
var SomaCfg SomaConfig

func main() {
	version := "0.0.16"
	log.Printf("Starting runtime config initialization, SOMA v%s", version)
	err := SomaCfg.readConfigFile("soma.conf")
	if err != nil {
		log.Fatal(err)
	}

	connectToDatabase()
	go pingDatabase()

	startHandlers()

	router := httprouter.New()

	router.GET("/views", ListViews)
	router.GET("/views/:view", ShowView)

	router.GET("/environments", ListEnvironments)
	router.GET("/environments/:environment", ShowEnvironment)

	router.GET("/objstates", ListObjectStates)
	router.GET("/objstates/:state", ShowObjectState)

	router.GET("/objtypes", ListObjectTypes)
	router.GET("/objtypes/:type", ShowObjectType)

	router.GET("/datacenters", ListDatacenters)
	router.GET("/datacenters/:datacenter", ShowDatacenter)

	router.GET("/datacentergroups", ListDatacenterGroups)
	router.GET("/datacentergroups/:datacentergroup", ShowDatacenterGroup)

	router.GET("/levels/", ListLevel)
	router.GET("/levels/:level", ShowLevel)

	router.GET("/predicates/", ListPredicate)
	router.GET("/predicates/:predicate", ShowPredicate)

	router.GET("/status/", ListStatus)
	router.GET("/status/:status", ShowStatus)

	router.GET("/oncall/", ListOncall)
	router.GET("/oncall/:oncall", ShowOncall)
	router.POST("/filter/oncall/", ListOncall)

	router.GET("/teams/", ListTeam)
	router.GET("/teams/:team", ShowTeam)
	router.POST("/filter/teams/", ListTeam)

	router.GET("/nodes/", ListNode)
	router.GET("/nodes/:node", ShowNode)
	router.POST("/filter/nodes/", ListNode)

	router.GET("/servers/", ListServer)
	router.GET("/servers/:server", ShowServer)
	router.POST("/filter/servers/", ListServer)

	if !SomaCfg.ReadOnly {
		router.POST("/views", AddView)
		router.DELETE("/views/:view", DeleteView)
		router.PUT("/views/:view", RenameView)

		router.POST("/environments", AddEnvironment)
		router.DELETE("/environments/:environment", DeleteEnvironment)
		router.PUT("/environments/:environment", RenameEnvironment)

		router.POST("/objstates", AddObjectState)
		router.DELETE("/objstates/:state", DeleteObjectState)
		router.PUT("/objstates/:state", RenameObjectState)

		router.POST("/objtypes", AddObjectType)
		router.DELETE("/objtypes/:type", DeleteObjectType)
		router.PUT("/objtypes/:type", RenameObjectType)

		router.POST("/datacenters", AddDatacenter)
		router.DELETE("/datacenters/:datacenter", DeleteDatacenter)
		router.PUT("/datacenters/:datacenter", RenameDatacenter)

		router.PATCH("/datacentergroups/:datacentergroup", AddDatacenterToGroup)
		router.DELETE("/datacentergroups/:datacentergroup", DeleteDatacenterFromGroup)

		router.POST("/levels/", AddLevel)
		router.DELETE("/levels/:level", DeleteLevel)

		router.POST("/predicates/", AddPredicate)
		router.DELETE("/predicates/:predicate", DeletePredicate)

		router.POST("/status/", AddStatus)
		router.DELETE("/status/:status", DeleteStatus)

		router.POST("/oncall/", AddOncall)
		router.PATCH("/oncall/:oncall", UpdateOncall)
		router.DELETE("/oncall/:oncall", DeleteOncall)

		router.POST("/teams/", AddTeam)
		router.DELETE("/teams/:team", DeleteTeam)

		router.POST("/nodes/", AddNode)
		router.DELETE("/nodes/:node", DeleteNode)

		router.POST("/servers/", AddServer)
		router.DELETE("/servers/:server", DeleteServer)
		router.PUT("/servers/:server", InsertNullServer)
	}

	log.Fatal(http.ListenAndServe(":8888", router))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
