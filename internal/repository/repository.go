package repository

import (
	"sportNews/pkg/log"

	"go.uber.org/zap"
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
	log.Info("NewDBSession: Creating a new DB session", zap.Int("repositoryIdx", repo.idx))
	return repo.db.NewSession()
}

func (repo Repository) Close() (err error) {
	if repo.db != nil {
		log.Info("NewDBSession: Creating a new DB session", zap.Int("repositoryIdx", repo.idx))
		if err = repo.db.Close(); err != nil {
			log.Error("Close: Failed to close database connection", zap.Int("repositoryIdx", repo.idx), zap.Error(err))
		} else {
			log.Infof("Close: Successfully closed the database connection for Repository(%d)", repo.idx)
		}
	} else {
		log.Warn("Close: Repository database connection is already nil", zap.Int("repositoryIdx", repo.idx))
	}

	return err
}
