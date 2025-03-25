package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sportNews/internal/httpError"
)

func (app App) router(e *gin.Engine) {
	e.HandleMethodNotAllowed = true
	e.NoMethod(httpError.HandleNoAllowMethod)
	e.NoRoute(httpError.HandleNotFound)

	e.GET("/healthz", httpError.ErrHandler(app.healthHandler))
}

func (app App) healthHandler(c *gin.Context) error {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

	return nil
}
