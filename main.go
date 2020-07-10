package main

import (
	"github.com/gin-gonic/gin"
	"pf-test-backend/controller"
	"pf-test-backend/models"
)

var(
	engine = gin.Default()
)
func main() {
	models.DataInit()
	engine.Use(Cors())

	// Routes
	engine.GET("/orders", controller.FindAll)

	// Run the server
	engine.Run(":8999")
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Next()
	}
}
