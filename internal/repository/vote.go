package repository

import (
	"db_hw/internal/models"
	"db_hw/internal/query"
)

func (r *repositoryImpl) InsertVote(newVote models.NewVote, threadID int) error {
	var lastID int
	err := r.db.QueryRow(query.InsertVote, newVote.Nickname, newVote.Voice,
		threadID).Scan(&lastID)

	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) UpdateVote(newVote models.NewVote, threadID int) error {
	rows, err := r.db.Query(query.UpdateVote, newVote.Voice, newVote.Nickname,
		threadID)
	defer rows.Close()

	if err != nil {
		return err
	}
	return nil
}
