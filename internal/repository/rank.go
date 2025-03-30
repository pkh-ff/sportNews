package repository

import (
	"go.uber.org/zap"
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"sportNews/pkg/log"
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

func (repo *Repository) GetOldestRankDataByType(t enum.RankType) (model.SportRank, error) {
	log.Infof("repository.GetOldestRankDataByType()")
	var data model.SportRank
	_, err := repo.db.Cols("id").
		Where("type = ?", t).
		OrderBy("date ASC").
		Get(&data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func (repo *Repository) GetRankDataCountByType(t enum.RankType) (int64, error) {
	log.Infof("repository.GetRankDataCountByType()")
	return repo.db.
		Where("type = ?", t).
		Count(&model.SportRank{})
}

func (repo *Repository) CheckRankDataExist(date string, t enum.RankType) (bool, error) {
	log.Infof("repository.CheckRankDataExist()")
	return repo.db.
		Where("type = ?", t).
		And("date = ?", date).
		Exist(&model.SportRank{})
}

func (repo *Repository) DeleteRankData(m model.SportRank) error {
	_, err := repo.db.Delete(m)

	return err
}
