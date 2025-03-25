package service

import "sportNews/internal/model"

// GetVideoList
// 取得指定筆數影片列表
func (s *Serv) GetVideoList() (interface{}, error) {
	videos, err := s.Repo.VideoList(10)
	if err != nil {
		return nil, err
	}

	data := make([]model.VideoResp, 0)
	for _, v := range videos {
		data = append(data, model.VideoResp{
			Title:       v.Title,
			Description: v.Description,
			Cover:       v.Cover,
			Link:        v.Link,
		})
	}

	return data, nil
}
