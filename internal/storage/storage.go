package storage

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/shadkain/db_hw/internal/repository"
	"io/ioutil"
)

type Storage struct {
	repo repository.Repository
	db   *pgx.ConnPool
}

func NewStorage() *Storage {
	return &Storage{}
}

const dbSchema = "assets/db/db.sql"
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

	if err := s.LoadSchemaSQL(); err != nil {
		return err
	}

	return err
}

func (s *Storage) LoadSchemaSQL() error {
	if s.db == nil {
		return pgx.ErrDeadConn
	}

	content, err := ioutil.ReadFile(dbSchema)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(string(content)); err != nil {
		return err
	}
	tx.Commit()

	fmt.Println("sql schema loaded")

	return nil
}

func (s *Storage) Disconn() {
	if s.db != nil {
		s.db.Close()
		s.db = nil
	}
}

func (RS *Storage) Query(sql string, args ...interface{}) (*pgx.Rows, error) {
	if RS.db == nil {
		return nil, pgx.ErrDeadConn
	}
	return RS.db.Query(sql, args...)
}

func (RS *Storage) QueryRow(sql string, args ...interface{}) *pgx.Row {
	if RS.db == nil {
		return nil
	}
	return RS.db.QueryRow(sql, args...)
}

func (s *Storage) Exec(sql string, args ...interface{}) (pgx.CommandTag, error) {
	if s.db == nil {
		return "", pgx.ErrDeadConn
	}

	tx, err := s.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	tag, err := tx.Exec(sql, args...)
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return tag, nil
}

func (s *Storage) Repository() repository.Repository {
	if s.repo == nil {
		s.repo = repository.NewRepository(s.db)
	}

	return s.repo
}
