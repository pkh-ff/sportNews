package service

import "sportNews/internal/model"

// QueryNews
// 新聞列表
func (s *Serv) QueryNews(page, size int) (interface{}, error) {
	start := (page - 1) * size
	news, err := s.Repo.QueryNewsByPage(size, start)
	if err != nil {
		return nil, err
	}

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

	count, err := s.Repo.QueryNewsCount()
	if err != nil {
		return nil, err
	}

	pn := int(count) / size
	if int(count)%size > 0 {
		pn = pn + 1
	}

	return model.BaseResp{
		Data:       data,
		TotalCount: count,
		TotalPage:  pn,
	}, nil
}
