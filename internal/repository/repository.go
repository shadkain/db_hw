package repository

import (
	"database/sql"
)

type repositoryImpl struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repositoryImpl{
		db: db,
	}
}

type Repository interface {
}
