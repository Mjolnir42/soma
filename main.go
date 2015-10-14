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
	version := "0.0.2"
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
	router.POST("/views", AddView)
	router.DELETE("/views/:view", DeleteView)
	router.PUT("/views/:view", RenameView)

	router.GET("/environments", ListEnvironments)
	router.GET("/environments/:environment", ShowEnvironment)
	router.POST("/environments", AddEnvironment)
	router.DELETE("/environments/:environment", DeleteEnvironment)
	router.PUT("/environments/:environment", RenameEnvironment)

	router.GET("/objstates", ListObjectStates)
	router.GET("/objstates/:state", ShowObjectState)
	router.POST("/objstates", AddObjectState)
	router.DELETE("/objstates/:state", DeleteObjectState)
	router.PUT("/objstates/:state", RenameObjectState)

	log.Fatal(http.ListenAndServe(":8888", router))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
