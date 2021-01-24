// +build prod

package config

import "os"

var (
	// API RELATED
	PORT = os.Getenv("PORT")

	// DB RELATED
	DB_HOST = os.Getenv("DB_HOST")
	DB_PORT = os.Getenv("DB_PORT")
	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")
	DB_NAME = os.Getenv("DB_NAME")
)
