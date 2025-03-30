package process

import (
	"context"
	"sportNews/internal/service/crawler"
	"xorm.io/xorm"
)

type Process struct {
	DB *xorm.EngineGroup
}

func NDTVProcess(ctx context.Context, db *xorm.EngineGroup) {
	serv := crawler.NewNDTVServ(db)
	serv.Crawler()
	//serv := crawler.NewEspncricinfoServ(db)
	//serv.Crawler()

	//serv := crawler.NewBCCIServ(db)
	//serv.Crawler()
}
