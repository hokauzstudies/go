package dao

import (
	gormbulk "github.com/t-tiger/gorm-bulk-insert/v2"
	"log"
	"neoway-challenge/dao/models"
)

type ShoppingData = models.ShoppingData

// automigrate the Data Model on every run to keep  db columns in sync with api model properties
func MigrateModel() error {
	db, err := GetDBInstance()
	if err != nil {
		log.Printf("Error getting DB instance to migrate model: %s", err)
		return err

	}
	defer db.Close()

	result := db.AutoMigrate(&ShoppingData{})
	if result.Error != nil {
		log.Printf("Error getting DB instance to migrate model: %s", result.Error)
		return err

	}
	return nil
}

func Save(data []interface{}) error {

	db, err := GetDBInstance()
	if err != nil {
		log.Printf("Error getting DB instance to save data: %s", err)
		return err

	}
	defer db.Close()

	err = gormbulk.BulkInsert(db, data, 3000)
	if err != nil {
		log.Printf("Error saving data in db: %s", err)
		return err
	}

	return nil
}
