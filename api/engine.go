package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Engine
// init gin Engine
func engine() *gin.Engine {
	e := &gin.Engine{}

	gin.SetMode(gin.DebugMode) // TODO
	e = gin.Default()

	c := cors.DefaultConfig()
	c.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	c.AllowAllOrigins = true

	e.Use(cors.New(c))

	return e
}
