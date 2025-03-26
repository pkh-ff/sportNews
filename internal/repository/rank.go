package repository

import (
	"go.uber.org/zap"
	"sportNews/internal/log"
	"sportNews/internal/model"
)

func (repo *Repository) GetRankDate(t string) (model.SportRank, error) {
	log.Infof("repository.GetRankDate()\n")
	cols := []string{"data"}
	sess := repo.db.Cols(cols...)
	sess.Where("type = ?", t)

	var data model.SportRank
	_, err := sess.Limit(1).Get(&data)
	if err != nil {
		log.Error("repository.GetRankDate()", zap.Error(err))
		return data, err
	}

	return data, nil
}
