package storage

import (
	"database/sql"
	"db_hw/internal/repository"
	_ "github.com/lib/pq"
)

type Storage struct {
	db   *sql.DB
	repo repository.Repository
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Open(url string) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Storage) Close() {
	if s.db != nil {
		s.db.Close()
	}

	s.db = nil
}

func (s *Storage) Repository() repository.Repository {
	if s.repo == nil {
		s.repo = repository.NewRepository(s.db)
	}

	return s.repo
}
