package process

import (
	"log"
	"sportNews/conf"
	"sportNews/internal/service/crawler"
	"time"
	"xorm.io/xorm"
)

type Process struct {
	DB *xorm.EngineGroup
}

func NDTVProcess(conf conf.App, db *xorm.EngineGroup) {
	serv := crawler.NewNDTVServ(db)
	for {
		log.Printf("run NDTVProcess on %v", time.Now().String())
		serv.Crawler()

		time.Sleep(time.Duration(conf.Process.News) * time.Second)
	}

}

func BCCIProcess(conf conf.App, db *xorm.EngineGroup) {
	serv := crawler.NewBCCIServ(db)
	for {
		log.Printf("run BCCIProcess on %v", time.Now().String())
		serv.Crawler()

		time.Sleep(time.Duration(conf.Process.Ranking) * time.Second)
	}
}
