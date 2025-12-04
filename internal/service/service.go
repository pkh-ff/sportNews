package service

import (
	"sportNews/internal/repository"

	"xorm.io/xorm"
)

type Serv struct {
	Repo *repository.Repository
}

func New(db *xorm.EngineGroup) *Serv {
	repo := repository.New(db)

	return &Serv{
		Repo: &repo,
	}
}
