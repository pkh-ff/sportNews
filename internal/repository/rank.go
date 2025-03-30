package repository

import (
	"go.uber.org/zap"
	"sportNews/internal/enum"
	"sportNews/internal/log"
	"sportNews/internal/model"
)

// GetRankDate
// 取得指定類型最新排行榜資料
func (repo *Repository) GetRankDate(t enum.RankType) (model.SportRank, error) {
	log.Infof("repository.GetRankDate()")
	cols := []string{"data"}
	sess := repo.db.Cols(cols...)
	sess.Where("type = ?", t)

	var data model.SportRank
	_, err := sess.OrderBy("date DESC").Limit(1).Get(&data)
	if err != nil {
		log.Error("repository.GetRankDate()", zap.Error(err))
		return data, err
	}

	return data, nil
}

// InsertRank
// 寫入排行榜資料
func (repo *Repository) InsertRank(m model.SportRank) error {
	log.Infof("repository.InsertRank()")
	_, err := repo.db.Insert(m)
	if err != nil {
		return err
	}

	return nil
}
