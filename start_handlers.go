package main

func startHandlers() {
	spawnViewReadHandler()
	spawnViewWriteHandler()

	spawnEnvironmentReadHandler()
	spawnEnvironmentWriteHandler()

	spawnObjectStateReadHandler()
	spawnObjectStateWriteHandler()

	spawnObjectTypeReadHandler()
	spawnObjectTypeWriteHandler()

	spawnDatacenterReadHandler()
	spawnDatacenterWriteHandler()

	spawnLevelReadHandler()
}

func spawnViewReadHandler() {
	var viewReadHandler somaViewReadHandler
	viewReadHandler.input = make(chan somaViewRequest)
	viewReadHandler.shutdown = make(chan bool)
	viewReadHandler.conn = conn
	handlerMap["viewReadHandler"] = viewReadHandler
	go viewReadHandler.run()
}

func spawnViewWriteHandler() {
	var viewWriteHandler somaViewWriteHandler
	viewWriteHandler.input = make(chan somaViewRequest, 64)
	viewWriteHandler.shutdown = make(chan bool)
	viewWriteHandler.conn = conn
	handlerMap["viewWriteHandler"] = viewWriteHandler
	go viewWriteHandler.run()
}

func spawnEnvironmentReadHandler() {
	var environmentReadHandler somaEnvironmentReadHandler
	environmentReadHandler.input = make(chan somaEnvironmentRequest)
	environmentReadHandler.shutdown = make(chan bool)
	environmentReadHandler.conn = conn
	handlerMap["environmentReadHandler"] = environmentReadHandler
	go environmentReadHandler.run()
}

func spawnEnvironmentWriteHandler() {
	var environmentWriteHandler somaEnvironmentWriteHandler
	environmentWriteHandler.input = make(chan somaEnvironmentRequest, 64)
	environmentWriteHandler.shutdown = make(chan bool)
	environmentWriteHandler.conn = conn
	handlerMap["environmentWriteHandler"] = environmentWriteHandler
	go environmentWriteHandler.run()
}

func spawnObjectStateReadHandler() {
	var objectStateReadHandler somaObjectStateReadHandler
	objectStateReadHandler.input = make(chan somaObjectStateRequest)
	objectStateReadHandler.shutdown = make(chan bool)
	objectStateReadHandler.conn = conn
	handlerMap["objectStateReadHandler"] = objectStateReadHandler
	go objectStateReadHandler.run()
}

func spawnObjectStateWriteHandler() {
	var objectStateWriteHandler somaObjectStateWriteHandler
	objectStateWriteHandler.input = make(chan somaObjectStateRequest, 64)
	objectStateWriteHandler.shutdown = make(chan bool)
	objectStateWriteHandler.conn = conn
	handlerMap["objectStateWriteHandler"] = objectStateWriteHandler
	go objectStateWriteHandler.run()
}

func spawnObjectTypeReadHandler() {
	var objectTypeReadHandler somaObjectTypeReadHandler
	objectTypeReadHandler.input = make(chan somaObjectTypeRequest)
	objectTypeReadHandler.shutdown = make(chan bool)
	objectTypeReadHandler.conn = conn
	handlerMap["objectTypeReadHandler"] = objectTypeReadHandler
	go objectTypeReadHandler.run()
}

func spawnObjectTypeWriteHandler() {
	var objectTypeWriteHandler somaObjectTypeWriteHandler
	objectTypeWriteHandler.input = make(chan somaObjectTypeRequest, 64)
	objectTypeWriteHandler.shutdown = make(chan bool)
	objectTypeWriteHandler.conn = conn
	handlerMap["objectTypeWriteHandler"] = objectTypeWriteHandler
	go objectTypeWriteHandler.run()
}

func spawnDatacenterReadHandler() {
	var datacenterReadHandler somaDatacenterReadHandler
	datacenterReadHandler.input = make(chan somaDatacenterRequest)
	datacenterReadHandler.shutdown = make(chan bool)
	datacenterReadHandler.conn = conn
	handlerMap["datacenterReadHandler"] = datacenterReadHandler
	go datacenterReadHandler.run()
}

func spawnDatacenterWriteHandler() {
	var datacenterWriteHandler somaDatacenterWriteHandler
	datacenterWriteHandler.input = make(chan somaDatacenterRequest, 64)
	datacenterWriteHandler.shutdown = make(chan bool)
	datacenterWriteHandler.conn = conn
	handlerMap["datacenterWriteHandler"] = datacenterWriteHandler
	go datacenterWriteHandler.run()
}

func spawnLevelReadHandler() {
	var levelReadHandler somaLevelReadHandler
	levelReadHandler.input = make(chan somaLevelRequest, 64)
	levelReadHandler.shutdown = make(chan bool)
	levelReadHandler.conn = conn
	handlerMap["levelReadHandler"] = levelReadHandler
	go levelReadHandler.run()
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
