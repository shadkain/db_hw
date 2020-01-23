package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/shadkain/db_hw/internal/repository"
	// For using pgx driver
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Storage struct {
	db   *sqlx.DB
	repo repository.Repository
}

func NewStorage() *Storage {
	return &Storage{}
}

func (this *Storage) Open(url string) error {
	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(8)

	if err := db.Ping(); err != nil {
		return err
	}

	this.db = db

	return nil
}

func (this *Storage) Close() {
	if this.db != nil {
		this.db.Close()
	}
}

func (this *Storage) Repository() repository.Repository {
	if this.repo == nil {
		this.repo = repository.NewRepository(this.db)
	}

	return this.repo
}
