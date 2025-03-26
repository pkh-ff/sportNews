package service

import (
	"sportNews/internal/log"
	"sportNews/internal/model"
)

// GetVideoList
// 取得指定筆數影片列表
func (s *Serv) GetVideoList() (interface{}, error) {
	log.Info("service.GetVideoList()")
	videos, err := s.Repo.VideoList(10)
	if err != nil {
		log.Errorf("service.GetVideoList(), get data error: %v\n", err)
		return nil, err
	}
	log.Infof("service.GetVideoList(), videos:%v\n", videos)

	data := make([]model.VideoResp, 0)
	for _, v := range videos {
		data = append(data, model.VideoResp{
			Title:       v.Title,
			Description: v.Description,
			Cover:       v.Cover,
			Link:        v.Link,
		})
	}
	log.Infof("service.GetVideoList(), data:%v\n", data)

	return data, nil
}
