package repository

import (
	"db_hw/internal/models"
	"db_hw/internal/query"
)

func (r *repositoryImpl) SelectStatus() (Status models.Status, Err error) {
	err := r.db.QueryRow(query.SelectStatus).Scan(&Status.Post, &Status.Thread, &Status.User, &Status.Forum)
	if err != nil {
		return Status, err
	}
	return
}

func (r *repositoryImpl) ClearAll() (Err error) {
	rows, err := r.db.Query(query.ClearAll)
	defer rows.Close()
	if err != nil {
		return err
	}

	return nil
}
