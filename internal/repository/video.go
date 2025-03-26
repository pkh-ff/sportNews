package repository

import (
	"go.uber.org/zap"
	"sportNews/internal/enum"
	"sportNews/internal/log"
	"sportNews/internal/model"
)

// VideoList
// 取得指定筆數影片列表
func (repo *Repository) VideoList(limit int) ([]model.Video, error) {
	cols := []string{"title", "description", "cover", "link"}
	sess := repo.db.Cols(cols...)
	sess.Where("status = ?", enum.Enable)

	var videos = make([]model.Video, 0)
	err := sess.Limit(limit).Find(&videos)
	if err != nil {
		log.Error("repository::VideoList()", zap.Error(err))
		return nil, err
	}

	return videos, nil
}
