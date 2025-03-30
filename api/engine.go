package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Engine
// init gin Engine
func engine(debug bool) *gin.Engine {
	e := &gin.Engine{}

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e = gin.Default()

	c := cors.DefaultConfig()
	c.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	c.AllowAllOrigins = true

	e.Use(cors.New(c))

	return e
}
