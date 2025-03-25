package repository

import (
	"sportNews/internal/log"
	"xorm.io/xorm"
)

type Repository struct {
	idx int
	db  *xorm.EngineGroup
}

func New(db *xorm.EngineGroup) Repository {
	return Repository{
		db: db,
	}
}

func (repo Repository) NewDBSession() *xorm.Session {
	return repo.db.NewSession()
}

func (repo Repository) Close() (err error) {
	if repo.db != nil {
		if err = repo.db.Close(); err != nil {
			log.Errorf("repository::Close", "Repository(%d) failed to close database connection, err = %v\n", repo.idx, err)
		}
		log.Infof("Repository(%d) closed the db connection.\n", repo.idx)
	}

	return
}
