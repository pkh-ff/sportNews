package repository

//go:generate mockgen -source=video.go -destination=mocks/video_repository_mock.go -package=mocks

import (
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"sportNews/pkg/log"

	"go.uber.org/zap"
)

type VideoRepository interface {
	VideoList(limit int) ([]model.Video, error)
}

// VideoList
// 隨機取得指定筆數影片列表
func (r *Repository) VideoList(limit int) ([]model.Video, error) {
	log.Info("VideoList: Fetching video list", zap.Int("repositoryIdx", r.idx), zap.Int("limit", limit))
	cols := []string{"title", "description", "cover", "link"}
	sess := r.exec.Cols(cols...)
	sess.Where("status = ?", enum.Enable)

	var videos = make([]model.Video, 0)
	err := sess.OrderBy("RAND()").Limit(limit).Find(&videos)
	if err != nil {
		log.Error("VideoList: Error fetching video list", zap.Int("repositoryIdx", r.idx), zap.Error(err))
		return nil, err
	}

	log.Info("VideoList: Successfully fetched video list", zap.Int("repositoryIdx", r.idx), zap.Int("count", len(videos)))
	return videos, nil
}
