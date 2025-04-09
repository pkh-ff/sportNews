package repository

import (
	"go.uber.org/zap"
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"sportNews/pkg/log"
	"xorm.io/xorm"
)

// GetNoCoverNews
// 取得還沒有同步封面新聞列表
func (repo *Repository) GetNoCoverNews() ([]model.News, error) {
	log.Info("GetNoCoverNews: Start fetching news")

	limit := 500
	sess := repo.db.Where("status = ?", enum.Enable).And("cover IS NULL OR cover = ?", "")
	var news = make([]model.News, 0)
	err := sess.OrderBy("pub_date DESC").Limit(limit).Find(&news)
	if err != nil {
		log.Error("GetNoCoverNews: Failed to fetch news", zap.Int("limit", limit), zap.Error(err))
		return nil, err
	}

	log.Info("QueryNewsByPage: Successfully fetched news", zap.Int("newsCount", len(news)))
	return news, nil
}

// QueryNewsByPage
// 分頁取得新聞列表
func (repo *Repository) QueryNewsByPage(limit, start int) ([]model.News, error) {
	log.Info("QueryNewsByPage: Start fetching news", zap.Int("limit", limit), zap.Int("start", start))
	cols := []string{"id", "title", "description", "cover", "cover_source", "pub_date"}
	sess := repo.queryNews(cols)

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
func (repo *Repository) QueryNewsCount() (int64, error) {
	log.Info("QueryNewsCount: Start counting news")
	sess := repo.queryNews([]string{"id"})

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
func (repo *Repository) FindNews(id int) (model.News, error) {
	log.Info("FindNews: Start fetching news details", zap.Int("id", id))
	cols := []string{"title", "description", "cover", "cover_source", "content", "pub_date"}
	sess := repo.db.Cols(cols...)
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
func (repo *Repository) InsertNews(m model.News) error {
	log.Info("InsertNews: Start inserting news", zap.Any("news", m))
	_, err := repo.db.Insert(m)
	if err != nil {
		log.Error("InsertNews: Failed to insert news", zap.Any("news", m), zap.Error(err))
		return err
	}

	log.Info("InsertNews: Successfully inserted news", zap.Any("news", m))
	return nil
}

// UpdateNews
// 更新DB中新聞資料
func (repo *Repository) UpdateNews(m model.News) error {
	log.Info("UpdateNews: Start updating news", zap.Any("news", m))
	_, err := repo.db.ID(m.Id).Update(&m)
	if err != nil {
		log.Error("UpdateNews,  Failed to update news", zap.Int32("id", m.Id), zap.Error(err))
		return err
	}

	return nil
}

func (repo *Repository) GetCountByTitle(title, source string) (int64, error) {
	log.Info("GetCountByTitle: Start counting news by title and source", zap.String("title", title), zap.String("source", source))
	var news = model.News{}

	sess := repo.db.Cols([]string{"id"}...)
	sess.Where("title = ?", title).And("source = ?", source)

	count, err := sess.Count(&news)
	if err != nil {
		log.Error("GetCountByTitle: Failed to count news", zap.Error(err))
		return 0, err
	}

	log.Info("GetCountByTitle: Successfully fetched count", zap.Int64("count", count))
	return count, nil
}

func (repo *Repository) queryNews(cols []string) *xorm.Session {
	if len(cols) == 0 {
		cols = []string{"*"}
	}
	sess := repo.db.Cols(cols...)
	sess.Where("status = ?", enum.Enable)

	return sess
}
