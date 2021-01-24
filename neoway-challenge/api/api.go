package api

import (
	"github.com/gin-gonic/gin"
	"neoway-challenge/api/handlers"
	"neoway-challenge/config"
)

func Start() {
	router := gin.Default()

	router.POST("/send-data", handlers.ConvertAndSaveData)
	port := config.PORT
	router.Run(port)
}
