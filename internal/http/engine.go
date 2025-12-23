package http

import (
	"sportNews/pkg/log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Engine
// init gin Engine
func engine(debug bool) *gin.Engine {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.Default()

	c := cors.DefaultConfig()
	c.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	c.AllowOrigins = []string{"*"}

	e.Use(cors.New(c))

	log.Info("Gin engine initialized", zap.Bool("debugMode", debug))
	return e
}
