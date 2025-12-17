package repository

import (
	"xorm.io/xorm"
)

type Repository struct {
	idx  int
	eg   *xorm.EngineGroup
	exec xorm.Interface
}

func New(db *xorm.EngineGroup) *Repository {
	return &Repository{
		eg:   db,
		exec: db,
	}
}

func (r *Repository) withSession(s *xorm.Session) *Repository {
	rr := *r
	rr.exec = s
	return &rr
}

func (r *Repository) InTx(fn func(tx *Repository) error) error {
	s := r.eg.NewSession()
	defer s.Close()

	if err := s.Begin(); err != nil {
		return err
	}

	txRepo := r.withSession(s)

	if err := fn(txRepo); err != nil {
		_ = s.Rollback()
		return err
	}
	return s.Commit()
}
