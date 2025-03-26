package service

import (
	"encoding/json"
	"sportNews/internal/enum"
	"sportNews/internal/log"
	"sportNews/internal/model"
)

// GetRankData
// 取得指定類型排行榜資料
func (s *Serv) GetRankData(t enum.RankType) (interface{}, error) {
	log.Infof("service.GetRankData()\n")
	data, err := s.Repo.GetRankDate(t)
	if err != nil {
		log.Errorf("service.GetRankData(), get data error: %v, type:%v\n", err, t)
		return nil, err
	}
	log.Infof("service.GetRankData(), type:%v, rank data:%v\n", t, data)

	if len(data.Data) < 1 {
		return []model.RankDetail{}, nil
	}

	var r []model.RankDetail

	err = json.Unmarshal([]byte(data.Data), &r)
	if err != nil {
		log.Errorf("service.GetRankData(), parser json error: %v, type:%v, data:%v\n", err, data.Data)
		return nil, err
	}
	log.Infof("service.GetRankData(), type:%v, r:%v\n", t, r)

	return r, nil
}
