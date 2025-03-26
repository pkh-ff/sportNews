package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sportNews/internal/enum"
	"sportNews/internal/response"
	"strconv"
)

func (app App) router(e *gin.Engine) {
	e.HandleMethodNotAllowed = true
	e.NoMethod(response.HandleNoAllowMethod)
	e.NoRoute(response.HandleNotFound)

	e.GET("/healthz", response.ErrHandler(app.healthHandler))
	e.GET("/news", response.ErrHandler(app.news))
	e.GET("/video", response.ErrHandler(app.video))
	e.GET("/rank/:type", response.ErrHandler(app.rank))
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
		return response.ErrNoRows
	}

	c.JSON(http.StatusOK, response.Success(news))

	return nil
}

// 取得影片列表
func (app App) video(c *gin.Context) error {
	videos, err := app.Serv.GetVideoList()
	if err != nil {
		return response.ErrNoRows
	}

	c.JSON(http.StatusOK, response.Success(videos))

	return nil
}

// 賽事排行榜
func (app App) rank(c *gin.Context) error {
	t := c.Param("type")
	rankType := enum.RankType(t)
	if !enum.IsRankTypeExist(rankType) {
		return response.ErrParameter
	}

	rank, err := app.Serv.GetRankData(rankType)
	if err != nil {
		return response.ErrNoRows
	}

	c.JSON(http.StatusOK, response.Success(rank))
	return nil
}
