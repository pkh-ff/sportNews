package service

import (
	"sportNews/internal/model/api"
	"sportNews/pkg/log"
)

// GetVideoList
// 取得指定筆數影片列表
func (s *Serv) GetVideoList() (interface{}, error) {
	log.Info("service.GetVideoList()")
	videos, err := s.Repo.VideoList(10)
	if err != nil {
		log.Errorf("service.GetVideoList(), get data error: %v", err)
		return nil, err
	}
	log.Infof("service.GetVideoList(), videos:%v", videos)

	data := make([]api.VideoResp, 0)
	for _, v := range videos {
		data = append(data, api.VideoResp{
			Title:       v.Title,
			Description: v.Description,
			Cover:       v.Cover,
			Link:        v.Link,
		})
	}
	log.Infof("service.GetVideoList(), data:%v", data)

	return data, nil
}
