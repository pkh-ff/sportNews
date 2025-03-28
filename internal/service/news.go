package service

import (
	"sportNews/internal/assets"
	"sportNews/internal/log"
	"sportNews/internal/model/api"
)

// QueryNews
// 新聞列表
func (s *Serv) QueryNews(page, size int) (api.NewsPageResp, error) {
	log.Info("service.QueryNews()")
	start := (page - 1) * size
	news, err := s.Repo.QueryNewsByPage(size, start)
	if err != nil {
		log.Errorf("service.QueryNews(), get news error: %v", err)
		return api.NewsPageResp{}, err
	}
	log.Infof("service.QueryNews(), news:%v", news)

	data := make([]api.NewsList, 0)
	for _, v := range news {
		data = append(data, api.NewsList{
			Id:          v.Id,
			Title:       v.Title,
			Description: v.Description,
			Cover:       assets.FullPath(v.Cover),
			CoverSource: v.CoverSource,
			PubDate:     v.PubDate,
		})
	}
	log.Infof("service.QueryNews(), data:%v", data)

	count, err := s.Repo.QueryNewsCount()
	if err != nil {
		log.Errorf("service.QueryNews(), get news data count error: %v", err)
		return api.NewsPageResp{}, err
	}

	pn := int(count) / size
	if int(count)%size > 0 {
		pn = pn + 1
	}
	log.Infof("service.QueryNews(), TotalCount:%v, TotalPage:%v", count, pn)

	return api.NewsPageResp{
		Records:    data,
		TotalCount: count,
		TotalPage:  pn,
	}, nil
}

// FindNews
// 取得新聞詳情
func (s *Serv) FindNews(id int) (api.NewsDetail, error) {
	log.Infof("service.FindNews(), id:", id)
	data, err := s.Repo.FindNews(id)
	if err != nil {
		return api.NewsDetail{}, err
	}

	return api.NewsDetail{
		Title:       data.Title,
		Cover:       assets.FullPath(data.Cover),
		CoverSource: data.CoverSource,
		Description: data.Description,
		Content:     data.Content,
		PubDate:     data.PubDate,
	}, nil
}
