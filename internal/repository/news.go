package repository

import (
	"go.uber.org/zap"
	"sportNews/internal/enum"
	"sportNews/internal/log"
	"sportNews/internal/model"
	"xorm.io/xorm"
)

// QueryNewsByPage
// 分頁取得新聞列表
func (repo *Repository) QueryNewsByPage(limit, start int) ([]model.News, error) {
	log.Infof("repository.QueryNewsByPage(), limit:%v, start:%v", limit, start)
	cols := []string{"id", "title", "description", "cover", "cover_source", "pub_date"}
	sess := repo.queryNews(cols)

	var news = make([]model.News, 0)
	err := sess.OrderBy("pub_date DESC").Limit(limit, start).Find(&news)
	if err != nil {
		log.Error("repository.QueryNewsByPage()", zap.Error(err))
		return nil, err
	}

	return news, nil
}

func (repo *Repository) QueryNewsCount() (int64, error) {
	log.Infof("repository.QueryNewsCount()")
	sess := repo.queryNews([]string{"id"})

	var news = model.News{}
	count, err := sess.Count(news)
	if err != nil {
		log.Error("repository.QueryNewsCount()", zap.Error(err))
		return 0, err
	}

	return count, nil
}

// FindNews
// 取得新聞詳情
func (repo *Repository) FindNews(id int) (model.News, error) {
	log.Infof("repository.FindNews()")
	cols := []string{"title", "description", "cover", "cover_source", "content", "pub_date"}
	sess := repo.db.Cols(cols...)
	sess.Where("id = ?", id).And("status = ?", enum.Enable)

	var data model.News
	_, err := sess.Limit(1).Get(&data)
	if err != nil {
		log.Error("repository.GetRankDate()", zap.Error(err))
		return data, err
	}

	return data, nil
}

func (repo *Repository) InsertNews(m model.News) error {
	log.Infof("repository.InsertNews()")
	_, err := repo.db.Insert(m)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) GetCountByTitle(title, source string) (int64, error) {
	log.Infof("repository.GetCountByTitle()")
	var news = model.News{}

	sess := repo.db.Cols([]string{"id"}...)
	sess.Where("title = ?", title).And("source = ?", source)
	count, err := sess.Count(news)
	if err != nil {
		log.Error("repository.GetCountByTitle()", zap.Error(err))
		return 0, err
	}

	return count, nil
}

func (repo *Repository) queryNews(cols []string) *xorm.Session {
	sess := repo.db.Cols(cols...)
	sess.Where("status = ?", enum.Enable)

	return sess
}
