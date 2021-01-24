package user

import (
	"pep-api/api/tools/router"

	"github.com/gin-gonic/gin"
)

// GetEndPoints -
func GetEndPoints(g *gin.RouterGroup) []*router.EndPoint {
	return []*router.EndPoint{
		{
			Name:    "users",
			Method:  "POST",
			Handler: Create,
			Group:   g,
		},
		{
			Name:    "users",
			Method:  "GET",
			Handler: ReadAll, // filter ssoID
			Group:   g,
		},
		// {
		// 	Name:    "users/fromMed",
		// 	Method:  "GET",
		// 	Handler: ReadFromMed,
		// 	Group:   g,
		// },
		{
			Name:    "users/:id",
			Method:  "GET",
			Handler: Read,
			Group:   g,
		},
		{
			Name:    "users/:id",
			Method:  "PUT",
			Handler: Update,
			Group:   g,
		},
		{
			Name:    "users/:id",
			Method:  "DELETE",
			Handler: Delete,
			Group:   g,
		},
	}
}
