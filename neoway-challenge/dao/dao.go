package dao

import (
	"fmt"

	"neoway-challenge/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Initialize() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=America/Sao_Paulo",
		config.DB_HOST,
		config.DB_PORT,
		config.DB_USER,
		config.DB_PASS,
		config.DB_NAME,
	)
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	if !db.HasTable(&ShoppingData{}) {
		db.CreateTable(&ShoppingData{})
	}

	db.LogMode(true)

	return nil
}

func GetDBInstance() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=America/Sao_Paulo",
		config.DB_HOST,
		config.DB_PORT,
		config.DB_USER,
		config.DB_PASS,
		config.DB_NAME,
	)

	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
