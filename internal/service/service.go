package service

import (
	"sportNews/internal/repository"

	"xorm.io/xorm"
)

type Serv struct {
	Repo      *repository.Repository
	VideoRepo repository.VideoRepository
}

func New(db *xorm.EngineGroup) *Serv {
	repo := repository.New(db)

	return &Serv{
		Repo:      &repo,
		VideoRepo: &repo,
	}
}
