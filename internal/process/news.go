package process

import (
	"go.uber.org/zap"
	"sportNews/conf"
	"sportNews/internal/service/crawler"
	"sportNews/pkg/log"
	"time"
	"xorm.io/xorm"
)

type Process struct {
	DB *xorm.EngineGroup
}

func NDTVProcess(conf conf.App, db *xorm.EngineGroup) {
	serv := crawler.NewNDTVServ(db)
	for {
		log.Info("NDTVProcess: Starting NDTV crawler", zap.Time("startTime", time.Now()))
		serv.Crawler()

		time.Sleep(time.Duration(conf.Process.News) * time.Second)
	}

}

func BCCIProcess(conf conf.App, db *xorm.EngineGroup) {
	serv := crawler.NewBCCIServ(db)
	for {
		log.Info("BCCIProcess: Starting BCCI crawler", zap.Time("startTime", time.Now()))
		serv.Crawler()

		time.Sleep(time.Duration(conf.Process.Ranking) * time.Second)
	}
}
