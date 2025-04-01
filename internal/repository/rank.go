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
	log.Info("GetRankDate: Fetching latest rank data", zap.Any("type", t))
	cols := []string{"data"}
	sess := repo.db.Cols(cols...)
	sess.Where("type = ?", t)

	var data model.SportRank
	_, err := sess.OrderBy("date DESC").Limit(1).Get(&data)
	if err != nil {
		log.Error("GetRankDate: Error fetching rank data", zap.Error(err), zap.Any("type", t))
		return data, err
	}
	log.Info("GetRankDate: Successfully fetched latest rank data", zap.Any("rankData", data))

	return data, nil
}

// InsertRank
// 寫入排行榜資料
func (repo *Repository) InsertRank(m model.SportRank) error {
	log.Info("InsertRank: Inserting rank data", zap.Any("rankData", m))
	_, err := repo.db.Insert(m)
	if err != nil {
		log.Error("InsertRank: Error inserting rank data", zap.Error(err), zap.Any("rankData", m))
		return err
	}

	log.Info("InsertRank: Successfully inserted rank data", zap.Any("rankData", m))
	return nil
}

// GetOldestRankDataByType
// 以type為條件取得最舊一筆資料
func (repo *Repository) GetOldestRankDataByType(t enum.RankType) (model.SportRank, error) {
	log.Info("GetOldestRankDataByType: Fetching oldest rank data", zap.Any("type", t))
	var data model.SportRank
	_, err := repo.db.Cols("id").
		Where("type = ?", t).
		OrderBy("date ASC").
		Get(&data)
	if err != nil {
		log.Error("GetOldestRankDataByType: Error fetching oldest rank data", zap.Error(err), zap.Any("type", t))
		return data, err
	}
	log.Info("GetOldestRankDataByType: Successfully fetched oldest rank data", zap.Any("rankData", data))

	return data, nil
}

// GetRankDataCountByType
// 以type作為條件取得資料量
func (repo *Repository) GetRankDataCountByType(t enum.RankType) (int64, error) {
	log.Info("GetRankDataCountByType: Fetching rank data count", zap.Any("type", t))
	count, err := repo.db.
		Where("type = ?", t).
		Count(&model.SportRank{})
	if err != nil {
		log.Error("GetRankDataCountByType: Error fetching rank data count", zap.Error(err), zap.Any("type", t))
		return 0, err
	}

	log.Info("GetRankDataCountByType: Successfully fetched rank data count", zap.Int64("count", count))
	return count, nil
}

// CheckRankDataExist
// 檢查資料是否存在
func (repo *Repository) CheckRankDataExist(date string, t enum.RankType) (bool, error) {
	log.Info("CheckRankDataExist: Checking if rank data exists", zap.String("date", date), zap.Any("type", t))
	exist, err := repo.db.
		Where("type = ?", t).
		And("date = ?", date).
		Exist(&model.SportRank{})
	if err != nil {
		log.Error("CheckRankDataExist: Error checking rank data existence", zap.Error(err), zap.String("date", date), zap.Any("type", t))
		return false, err
	}

	log.Info("CheckRankDataExist: Rank data existence check", zap.Bool("exists", exist))
	return exist, nil
}

// DeleteRankData
// 刪除排行榜資料
func (repo *Repository) DeleteRankData(m model.SportRank) (bool, error) {
	log.Info("DeleteRankData: Deleting rank data", zap.Any("rankData", m))
	_, err := repo.db.Delete(m)
	if err != nil {
		log.Error("DeleteRankData: Error deleting rank data", zap.Error(err), zap.Any("rankData", m))
		return false, err
	}

	log.Info("DeleteRankData: Successfully deleted rank data", zap.Any("rankData", m))
	return true, nil
}
