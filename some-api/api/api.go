package api

import (
	"pep-api/api/modules/user"
	"pep-api/api/tools/router"

	"github.com/gin-gonic/gin"
)

// Start -
func Start(r *gin.Engine, mainPath string) {
	innerPath := mainPath + "/api"
	private := r.Group(innerPath)
	// private.Use(middleware.AuthRequired())

	endpoints := []*router.EndPoint{}
	endpoints = append(endpoints, user.GetEndPoints(private)...)

	for _, point := range endpoints {
		router.EnableEndPoint(point)
	}
}
