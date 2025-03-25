package repository

import (
	"go.uber.org/zap"
	"sportNews/internal/enum"
	"sportNews/internal/log"
	"sportNews/internal/model"
	"xorm.io/xorm"
)

func (repo *Repository) QueryNewsByPage(limit, start int) ([]model.News, error) {
	cols := []string{"title", "description", "cover", "cover_source", "link", "pub_date"}
	sess := repo.queryNews(cols)

	var news = make([]model.News, 0)
	err := sess.Limit(limit, start).Find(&news)
	if err != nil {
		log.Error("repository::QueryActorsByPage", zap.Error(err))
		return nil, err
	}

	return news, nil
}

func (repo *Repository) QueryNewsCount() (int64, error) {
	sess := repo.queryNews([]string{"id"})

	var actor = model.News{}
	count, err := sess.Count(actor)
	if err != nil {
		log.Error("repository::QueryActorsCount", zap.Error(err))
		return 0, err
	}

	return count, nil
}

func (repo *Repository) queryNews(cols []string) *xorm.Session {
	sess := repo.db.Cols(cols...)
	sess.Where("status = ?", enum.Enable)

	return sess
}
