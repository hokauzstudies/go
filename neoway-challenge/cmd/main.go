package main

import (
	"fmt"
	"log"

	"neoway-challenge/api"
	"neoway-challenge/dao"
)

func main() {

	// initialize database to make sure its reachable
	err := dao.Initialize()
	if err != nil {
		errMsg := fmt.Sprintf("Error initialising database: %s", err)
		log.Fatal(errMsg)
	}
	fmt.Println("Connected to database")

	// auto migrate database models on start
	err = dao.MigrateModel()
	if err != nil {
		errMsg := fmt.Sprintf("Error migrating database models: %s", err)
		log.Fatal(errMsg)
	}

	api.Start()

}
