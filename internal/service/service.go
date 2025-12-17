package service

import (
	"sportNews/internal/repository"

	"xorm.io/xorm"
)

type Serv struct {
	NewsRepo  repository.NewsRepository
	RankRepo  repository.RankRepository
	VideoRepo repository.VideoRepository
}

func New(db *xorm.EngineGroup) *Serv {
	repo := repository.New(db)

	return &Serv{
		NewsRepo:  repo,
		RankRepo:  repo,
		VideoRepo: repo,
	}
}
