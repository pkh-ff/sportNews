package crawler

import (
	"sportNews/internal/repository"
	"xorm.io/xorm"
)

type Serv struct {
	Repo   *repository.Repository
	Source string
}

func New(db *xorm.EngineGroup, source string) *Serv {
	repo := repository.New(db)

	return &Serv{
		Repo:   &repo,
		Source: source,
	}
}
