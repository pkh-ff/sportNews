package api

import (
	"errors"
	"net/http"
	"sportNews/internal/enum"
	"sportNews/internal/model/api"
	"sportNews/internal/response"
	"sportNews/pkg/log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func (app App) router(e *gin.Engine) {
	e.HandleMethodNotAllowed = true
	e.NoMethod(response.HandleNoAllowMethod)
	e.NoRoute(response.HandleNotFound)

	e.GET("/healthz", response.ErrHandler(app.healthHandler))
	e.GET("/news", response.ErrHandler(app.news))
	e.GET("/news/:id", response.ErrHandler(app.newsDetail))
	e.GET("/video", response.ErrHandler(app.video))
	e.GET("/rank/:type", response.ErrHandler(app.rank))
	e.POST("/feedback", response.ErrHandler(app.feedback))
}

func (app App) healthHandler(c *gin.Context) error {
	log.Info("health check endpoint accessed")
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

	return nil
}

// 分頁方式取得新聞列表
func (app App) news(c *gin.Context) error {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "15"))
	if err != nil || size < 1 {
		size = 15
	}

	news, err := app.Serv.QueryNews(page, size)
	if err != nil {
		log.Error("Failed to fetch news", zap.Error(err))
		return response.ErrNoRows
	}

	c.JSON(http.StatusOK, response.Success(news))

	return nil
}

// 新聞詳情
func (app App) newsDetail(c *gin.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		return response.ErrParameter
	}

	data, err := app.Serv.FindNews(id)
	if err != nil {
		log.Error("Failed to fetch news detail", zap.Int("id", id), zap.Error(err))
		return response.ErrNoRows
	}

	c.JSON(http.StatusOK, response.Success(data))
	return nil
}

// 取得影片列表
func (app App) video(c *gin.Context) error {
	videos, err := app.Serv.GetVideoList()
	if err != nil {
		log.Error("Failed to fetch video list", zap.Error(err))
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
		log.Error("Failed to fetch rank data", zap.String("type", string(rankType)), zap.Error(err))
		return response.ErrNoRows
	}

	c.JSON(http.StatusOK, response.Success(rank))
	return nil
}

func (app App) feedback(c *gin.Context) error {
	log.Info("health check endpoint accessed")

	var req api.FeedbackReq

	if err := c.ShouldBindJSON(&req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errMsgs := make([]string, 0)
			// 檢查哪些欄位錯誤
			for _, fe := range ve {
				field := fe.Field()
				errMsgs = append(errMsgs, field+" is invalid")
			}

			// 返回錯誤訊息
			return response.ParameterError(errMsgs)
		}
	}

	c.JSON(http.StatusOK, response.Success(""))

	return nil
}
