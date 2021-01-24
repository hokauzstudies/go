package middleware

import (
	"fmt"
	"net/http"
	"pep-api/api/tools/jwt"
	"pep-api/api/tools/router"
	"strings"

	"github.com/gin-gonic/gin"
)

var messages = map[string]string{
	"missing-token":     "Necessário token de autenticação",
	"expored-token":     "Faça login novamente",
	"invalid-token":     "Token inválido",
	"missing-api-token": "Necessário token de autenticação",
	"invalid-api-token": "Token inválido",
}

// ClientAuthRequired -
func ClientAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			auth  = c.GetHeader("Authorization")
			token string
		)

		if !strings.Contains(auth, "Bearer ") || auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &router.Response{Status: "error", Message: messages["missing-token"], Error: "missing-token"})
		}

		token = strings.Replace(auth, "Bearer ", "", 1)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &router.Response{Status: "error", Message: messages["invalid-token"], Error: "invalid-token"})
		}

		check, claims, e := jwt.Verify("", true, token, true)
		if !check {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &router.Response{Status: "error", Message: messages["expired-token"], Error: "expired-token"})
		}
		if e != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &router.Response{Status: "error", Message: messages["invalid-token"], Error: "invalid-token"})
		}

		id := fmt.Sprintf("%v", claims["sub"])
		c.Set("id", id)
		c.Next()
	}
}

// ApiAuthRequired -
func ApiAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			auth  = c.GetHeader("Authorization")
			token string
		)

		if !strings.Contains(auth, "Bearer ") || auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &router.Response{Status: "error", Message: messages["missing-api-token"], Error: "missing-api-token"}) // TODO add router.Reponse
			return
		}

		token = strings.Replace(auth, "Bearer ", "", 1)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &router.Response{Status: "error", Message: messages["invalid-api-token"], Error: "invalid-api-token"})
			return
		}
		c.Next()
	}
}
