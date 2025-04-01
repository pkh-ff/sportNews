package service

import (
	"go.uber.org/zap"
	"sportNews/internal/model/api"
	"sportNews/pkg/log"
)

// GetVideoList
// 取得指定筆數影片列表
func (s *Serv) GetVideoList() (interface{}, error) {
	log.Info("GetVideoList: Start fetching video list")
	videos, err := s.Repo.VideoList(10)
	if err != nil {
		log.Error("GetVideoList: Failed to fetch video list from repository", zap.Error(err))
		return nil, err
	}
	log.Info("GetVideoList: Video list fetched successfully", zap.Int("videoCount", len(videos)))

	data := make([]api.VideoResp, 0)
	for _, v := range videos {
		data = append(data, api.VideoResp{
			Title:       v.Title,
			Description: v.Description,
			Cover:       v.Cover,
			Link:        v.Link,
		})
	}
	log.Info("GetVideoList: Video data transformed into API response", zap.Int("videoDataCount", len(data)))

	return data, nil
}
