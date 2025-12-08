package service

import (
	"encoding/json"
	"sportNews/internal/assets"
	"sportNews/internal/enum"
	"sportNews/internal/model"
	"sportNews/pkg/log"

	"go.uber.org/zap"
)

// GetRankData
// 取得指定類型排行榜資料
func (s *Serv) GetRankData(t enum.RankType) (interface{}, error) {
	log.Info("GetRankData: Start fetching rank data", zap.String("rankType", string(t)))

	data, err := s.RankRepo.GetRankDate(t)
	if err != nil {
		log.Error("GetRankData: Failed to fetch rank data from repository", zap.Error(err), zap.String("rankType", string(t)))
		return nil, err
	}
	log.Info("GetRankData: Rank data fetched successfully", zap.String("rankType", string(t)), zap.Any("rankData", data))

	if len(data.Data) < 1 {
		log.Info("GetRankData: No rank data available for the given type", zap.String("rankType", string(t)))
		return []model.RankDetail{}, nil
	}

	var r []model.RankDetail

	err = json.Unmarshal([]byte(data.Data), &r)
	if err != nil {
		log.Error("GetRankData: Failed to parse rank data JSON", zap.Error(err), zap.String("rankType", string(t)), zap.String("rawData", data.Data))
		return nil, err
	}
	log.Info("GetRankData: Rank data successfully parsed", zap.String("rankType", string(t)), zap.Int("rankDataCount", len(r)))

	for i := range r {
		r[i].Icon = assets.FullAssetsPath(r[i].Icon)

		// 暫時作法
		if r[i].Team == "New Zealand" {
			r[i].Icon = assets.FullAssetsPath("n8/temp/New_Zealand.png")
		} else if r[i].Team == "England" {
			r[i].Icon = assets.FullAssetsPath("n8/temp/England.png")
		}
	}
	return r, nil
}
