package router

import (
	"io"

	"github.com/gin-gonic/gin"
)

type (
	// GroupedEndPoints -
	GroupedEndPoints struct {
		Name        string
		Middlewares []gin.HandlerFunc
	}

	// EndPoint -
	EndPoint struct {
		Name    string
		Method  string
		Handler func(ctx *Context) (int, *Response)
		Group   *gin.RouterGroup
	}

	// Response -
	Response struct {
		Status  string      `json:"status,omitempty"` // OK, Error, Warning
		Data    interface{} `json:"data,omitempty"`
		Error   string      `json:"error,omitempty"`
		Message string      `json:"message,omitempty"`
	}

	// Context -
	Context struct {
		Params    map[string]interface{}
		ExtraBody interface{}
		Body      io.ReadCloser
		Headers   interface{}
		Queries   map[string]interface{}
	}
)

var messages = map[string]string{
	"id-not-found": "Não foi possível identificar solicitante",
}

func new() *Context {
	return &Context{}
}

// EnableEndPoint -
func EnableEndPoint(p *EndPoint) {
	switch p.Method {
	case "GET":
		p.Group.GET(p.Name, func(c *gin.Context) { interceptor(c, p) })
	case "POST":
		p.Group.POST(p.Name, func(c *gin.Context) { interceptor(c, p) })
	case "PUT":
		p.Group.PUT(p.Name, func(c *gin.Context) { interceptor(c, p) })
	case "DELETE":
		p.Group.DELETE(p.Name, func(c *gin.Context) { interceptor(c, p) })
	}
}

func interceptor(c *gin.Context, p *EndPoint) {
	setHeaders(c, p.Method)

	ctx := new()
	ctx.Params = make(map[string]interface{})
	ctx.Queries = make(map[string]interface{})

	for _, v := range c.Params {
		ctx.Params[v.Key] = v.Value
	}

	for k, v := range c.Request.URL.Query() {
		ctx.Queries[k] = v[0]
	}

	id, ok := c.Get("id")
	if ok {
		ctx.Params["operator_id"] = id
	}

	ctx.Body = c.Request.Body

	c.JSON(p.Handler(ctx))
}

func setHeaders(c *gin.Context, m string) { // TODO passar para middleware
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Authorization")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	} else {
		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
