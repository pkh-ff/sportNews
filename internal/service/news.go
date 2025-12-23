package service

import (
	"regexp"
	"sportNews/internal/assets"
	"sportNews/internal/model/api"
	"sportNews/pkg/log"

	"go.uber.org/zap"
)

// QueryNews
// 新聞列表
func (s *Serv) QueryNews(page, size int) (api.NewsPageResp, error) {
	log.Info("QueryNews: Start fetching news list", zap.Int("page", page), zap.Int("size", size))
	start := (page - 1) * size
	news, err := s.NewsRepo.QueryNewsByPage(size, start)
	if err != nil {
		log.Error("QueryNews: Failed to fetch news from database", zap.Error(err))
		return api.NewsPageResp{}, err
	}
	log.Info("QueryNews: News fetched successfully", zap.Int("newsCount", len(news)))

	data := make([]api.NewsList, 0)

	for _, v := range news {
		if v.CoverCustom == "" {
			v.CoverCustom = assets.FullAssetsPath(v.Cover)
		} else {
			v.CoverCustom = assets.FullAssetsPath(v.CoverCustom)
		}

		data = append(data, api.NewsList{
			Id:          v.Id,
			Title:       v.Title,
			Description: v.Description,
			Cover:       assets.FullAssetsPath(v.Cover),
			CoverSource: v.CoverSource,
			CoverCustom: v.CoverCustom,
			PubDate:     v.PubDate,
		})
	}
	log.Info("QueryNews: News data processed", zap.Int("newsDataCount", len(data)))

	count, err := s.NewsRepo.QueryNewsCount()
	if err != nil {
		log.Error("QueryNews: Failed to fetch news count from database", zap.Error(err))
		return api.NewsPageResp{}, err
	}

	pn := int(count) / size
	if int(count)%size > 0 {
		pn = pn + 1
	}
	log.Info("QueryNews: Pagination calculated", zap.Int("totalCount", int(count)), zap.Int("totalPages", pn))

	return api.NewsPageResp{
		Records:    data,
		TotalCount: count,
		TotalPage:  pn,
	}, nil
}

// FindNews
// 取得新聞詳情
func (s *Serv) FindNews(id int) (api.NewsDetail, error) {
	log.Info("FindNews: Start fetching news details", zap.Int("id", id))
	data, err := s.NewsRepo.FindNews(id)
	if err != nil {
		log.Error("FindNews: Failed to fetch news details from database", zap.Int("id", id), zap.Error(err))
		return api.NewsDetail{}, err
	}
	log.Info("FindNews: News details fetched successfully", zap.Int("id", id))

	re := regexp.MustCompile(`https?://[^\s]+`) // 匹配任何以 api:// 或 https:// 開頭，後面接著不是空白字元（[^\s]）的一串文字
	content := re.ReplaceAllString(data.Content, "")

	if data.CoverCustom == "" {
		data.CoverCustom = assets.FullAssetsPath(data.Cover)
	} else {
		data.CoverCustom = assets.FullAssetsPath(data.CoverCustom)
	}

	return api.NewsDetail{
		Title:       data.Title,
		Cover:       assets.FullAssetsPath(data.Cover),
		CoverSource: data.CoverSource,
		CoverCustom: data.CoverCustom,
		Description: data.Description,
		Content:     content,
		PubDate:     data.PubDate,
	}, nil
}
