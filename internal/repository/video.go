package repository

import (
	"go.uber.org/zap"
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"sportNews/pkg/log"
)

// VideoList
// 隨機取得指定筆數影片列表
func (repo *Repository) VideoList(limit int) ([]model.Video, error) {
	log.Info("VideoList: Fetching video list", zap.Int("repositoryIdx", repo.idx), zap.Int("limit", limit))
	cols := []string{"title", "description", "cover", "link"}
	sess := repo.db.Cols(cols...)
	sess.Where("status = ?", enum.Enable)

	var videos = make([]model.Video, 0)
	err := sess.OrderBy("RAND()").Limit(limit).Find(&videos)
	if err != nil {
		log.Error("VideoList: Error fetching video list", zap.Int("repositoryIdx", repo.idx), zap.Error(err))
		return nil, err
	}

	log.Info("VideoList: Successfully fetched video list", zap.Int("repositoryIdx", repo.idx), zap.Int("count", len(videos)))
	return videos, nil
}
