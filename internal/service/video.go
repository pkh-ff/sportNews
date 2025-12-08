package service

import (
	"sportNews/internal/assets"
	"sportNews/internal/model/api"
	"sportNews/pkg/log"

	"go.uber.org/zap"
)

// GetVideoList
// 取得指定筆數影片列表
func (s *Serv) GetVideoList() (interface{}, error) {
	log.Info("GetVideoList: Start fetching video list")
	videos, err := s.VideoRepo.VideoList(10)
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
			Cover:       assets.FullAssetsPath(v.Cover),
			Link:        assets.FullAssetsPath(v.Link),
		})
	}
	log.Info("GetVideoList: Video data transformed into API response", zap.Int("videoDataCount", len(data)))

	return data, nil
}
