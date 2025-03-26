package service

import (
	"sportNews/internal/log"
	"sportNews/internal/model"
)

// QueryNews
// 新聞列表
func (s *Serv) QueryNews(page, size int) (interface{}, error) {
	log.Info("service.QueryNews()")
	start := (page - 1) * size
	news, err := s.Repo.QueryNewsByPage(size, start)
	if err != nil {
		log.Errorf("service.QueryNews(), get news error: %v\n", err)
		return nil, err
	}
	log.Infof("service.QueryNews(), news:%v\n", news)

	data := make([]model.NewsResp, 0)
	for _, v := range news {
		data = append(data, model.NewsResp{
			Title:       v.Title,
			Description: v.Description,
			Cover:       v.Cover,
			CoverSource: v.CoverSource,
			Link:        v.Link,
			PubDate:     v.PubDate,
		})
	}
	log.Infof("service.QueryNews(), data:%v\n", data)

	count, err := s.Repo.QueryNewsCount()
	if err != nil {
		log.Errorf("service.QueryNews(), get news data count error: %v\n", err)
		return nil, err
	}

	pn := int(count) / size
	if int(count)%size > 0 {
		pn = pn + 1
	}
	log.Infof("service.QueryNews(), TotalCount:%v, TotalPage:%v\n", count, pn)

	return model.BaseResp{
		Data:       data,
		TotalCount: count,
		TotalPage:  pn,
	}, nil
}
