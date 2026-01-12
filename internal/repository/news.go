package repository

//go:generate mockgen -source=news.go -destination=mocks/news_repository_mock.go -package=mocks

import (
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"sportNews/pkg/log"

	"go.uber.org/zap"
	"xorm.io/xorm"
)

type NewsRepository interface {
	QueryNewsByPage(limit, start int) ([]model.News, error)
	QueryNewsCount() (int64, error)
	FindNews(id int) (model.News, error)
	GetLastUpdateCoverCustomNews() (model.News, error)
	GetNoCoverCustomNews() ([]model.News, error)
	InsertNews(m model.News) error
	UpdateNews(m model.News) error
	GetCountByTitle(title, source string) (int64, error)
}

// GetNoCoverNews
// 取得還沒有同步封面新聞列表
func (r *Repository) GetNoCoverNews() ([]model.News, error) {
	log.Info("GetNoCoverNews: Start fetching news")

	limit := 500
	sess := r.exec.Where("status = ?", enum.Enable).And("cover IS NULL OR cover = ?", "")
	var news = make([]model.News, 0)
	err := sess.OrderBy("pub_date DESC").Limit(limit).Find(&news)
	if err != nil {
		log.Error("GetNoCoverNews: Failed to fetch news", zap.Int("limit", limit), zap.Error(err))
		return nil, err
	}

	log.Info("GetNoCoverNews: Successfully fetched news", zap.Int("newsCount", len(news)))
	return news, nil
}

// GetNoCoverCustomNews
// 取得還沒有設置自定義封面新聞列表
func (r *Repository) GetNoCoverCustomNews() ([]model.News, error) {
	log.Info("GetNoCoverCustomNews: Start fetching news")

	limit := 500
	sess := r.exec.Where("status = ?", enum.Enable).And("cover_custom IS NULL OR cover_custom = ?", "")
	var news = make([]model.News, 0)
	err := sess.OrderBy("update_at DESC, id DESC").Limit(limit).Find(&news)
	if err != nil {
		log.Error("GetNoCoverCustomNews: Failed to fetch news", zap.Int("limit", limit), zap.Error(err))
		return nil, err
	}

	log.Info("GetNoCoverCustomNews: Successfully fetched news", zap.Int("newsCount", len(news)))
	return news, nil
}

// GetLastUpdateCoverCustomNews
// 取得最後更新cover_custom欄位資料
func (r *Repository) GetLastUpdateCoverCustomNews() (model.News, error) {
	log.Info("GetLastUpdateCoverCustomNews: Start fetching news")
	var news model.News
	found, err := r.exec.Cols("cover_custom").
		Where("status = ?", enum.Enable).
		And("cover_custom IS NOT NULL").
		And("cover_custom != ?", "").
		OrderBy("update_at DESC, id DESC").Limit(1).Get(&news)
	if err != nil {
		return model.News{}, err
	}

	if !found {
		log.Warn("FindNews: News not found")
		return news, nil
	}

	return news, nil
}

// QueryNewsByPage
// 分頁取得新聞列表
func (r *Repository) QueryNewsByPage(limit, start int) ([]model.News, error) {
	log.Info("QueryNewsByPage: Start fetching news", zap.Int("limit", limit), zap.Int("start", start))
	cols := []string{"id", "title", "description", "cover", "cover_source", "cover_custom", "pub_date"}
	sess := r.queryNews(cols)

	var news = make([]model.News, 0)
	err := sess.OrderBy("pub_date DESC").Limit(limit, start).Find(&news)
	if err != nil {
		log.Error("QueryNewsByPage: Failed to fetch news", zap.Int("limit", limit), zap.Int("start", start), zap.Error(err))
		return nil, err
	}

	log.Info("QueryNewsByPage: Successfully fetched news", zap.Int("newsCount", len(news)))
	return news, nil
}

// QueryNewsCount
// 取得所有有效新聞筆數
func (r *Repository) QueryNewsCount() (int64, error) {
	log.Info("QueryNewsCount: Start counting news")
	sess := r.queryNews([]string{"id"})

	var news = model.News{}
	count, err := sess.Count(&news)
	if err != nil {
		log.Error("QueryNewsCount: Failed to count news", zap.Error(err))
		return 0, err
	}

	log.Info("QueryNewsCount: Successfully fetched news count", zap.Int64("count", count))
	return count, nil
}

// FindNews
// 取得新聞詳情
func (r *Repository) FindNews(id int) (model.News, error) {
	log.Info("FindNews: Start fetching news details", zap.Int("id", id))
	cols := []string{"title", "description", "cover", "cover_source", "cover_custom", "content", "pub_date"}
	sess := r.exec.Cols(cols...)
	sess.Where("id = ?", id).And("status = ?", enum.Enable)

	var data model.News
	found, err := sess.Limit(1).Get(&data)
	if err != nil {
		log.Error("FindNews: Failed to fetch news", zap.Int("id", id), zap.Error(err))
		return data, err
	}

	if !found {
		log.Warn("FindNews: News not found", zap.Int("id", id))
		return data, nil
	}

	log.Info("FindNews: Successfully fetched news", zap.Int("id", id), zap.Any("news", data))
	return data, nil
}

// InsertNews
// 寫入新聞
func (r *Repository) InsertNews(m model.News) error {
	log.Info("InsertNews: Start inserting news", zap.Any("news", m))
	_, err := r.exec.Insert(m)
	if err != nil {
		log.Error("InsertNews: Failed to insert news", zap.Any("news", m), zap.Error(err))
		return err
	}

	log.Info("InsertNews: Successfully inserted news", zap.Any("news", m))
	return nil
}

// UpdateNews
// 更新DB中新聞資料
func (r *Repository) UpdateNews(m model.News) error {
	log.Info("UpdateNews: Start updating news", zap.Any("news", m))
	_, err := r.exec.ID(m.Id).Update(&m)
	if err != nil {
		log.Error("UpdateNews,  Failed to update news", zap.Int32("id", m.Id), zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) GetCountByTitle(title, source string) (int64, error) {
	log.Info("GetCountByTitle: Start counting news by title and source", zap.String("title", title), zap.String("source", source))
	var news = model.News{}

	sess := r.exec.Cols([]string{"id"}...)
	sess.Where("title = ?", title).And("source = ?", source)

	count, err := sess.Count(&news)
	if err != nil {
		log.Error("GetCountByTitle: Failed to count news", zap.Error(err))
		return 0, err
	}

	log.Info("GetCountByTitle: Successfully fetched count", zap.Int64("count", count))
	return count, nil
}

func (r *Repository) queryNews(cols []string) *xorm.Session {
	if len(cols) == 0 {
		cols = []string{"*"}
	}
	sess := r.exec.Cols(cols...)
	sess.Where("status = ?", enum.Enable)

	return sess
}
