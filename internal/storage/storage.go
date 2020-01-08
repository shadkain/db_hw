package storage

import (
	"github.com/shadkain/db_hw/internal/repository"
	"github.com/jackc/pgx"
)

type Storage struct {
	repo repository.Repository
	db   *pgx.ConnPool
}

func NewStorage() *Storage {
	return &Storage{}
}

const maxConn = 2000

func (s *Storage) Open(psqURI string) error {
	if s.db != nil {
		return nil
	}
	config, err := pgx.ParseURI(psqURI)
	if err != nil {
		return err
	}

	s.db, err = pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     config,
			MaxConnections: maxConn,
		},
	)

	return err
}

func (s *Storage) Repository() repository.Repository {
	if s.repo == nil {
		s.repo = repository.NewRepository(s.db)
	}

	return s.repo
}
