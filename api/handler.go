package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sportNews/internal/httpError"
	"strconv"
)

func (app App) router(e *gin.Engine) {
	e.HandleMethodNotAllowed = true
	e.NoMethod(httpError.HandleNoAllowMethod)
	e.NoRoute(httpError.HandleNotFound)

	e.GET("/healthz", httpError.ErrHandler(app.healthHandler))
	e.GET("/news", httpError.ErrHandler(app.news))
	e.GET("/video", httpError.ErrHandler(app.video))
}

func (app App) healthHandler(c *gin.Context) error {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

	return nil
}

// 分頁方式取得新聞列表
func (app App) news(c *gin.Context) error {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	sizeStr := c.DefaultQuery("size", "15")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		page = 1
	}

	news, err := app.Serv.QueryNews(page, size)
	if err != nil {
		return httpError.ErrNoRows
	}

	c.JSON(http.StatusOK, news)

	return nil
}

// 取得影片列表
func (app App) video(c *gin.Context) error {
	videos, err := app.Serv.GetVideoList()
	if err != nil {
		return httpError.ErrNoRows
	}

	c.JSON(http.StatusOK, videos)

	return nil
}
