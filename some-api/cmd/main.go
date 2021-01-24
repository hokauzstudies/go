package main

import (
	"pep-api/api"
	"pep-api/api/tools/router"
	"pep-api/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db.Connect() // TODO criar verificação de conexão
	defer db.CloseConn()

	const mainPath = "pep/"
	rout := gin.Default()

	rout.Use(router.Cors())
	api.Start(rout, mainPath)

	rout.Run((":3002"))
}
