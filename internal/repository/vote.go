package repository

import (
	"database/sql"
	"github.com/shadkain/db_hw/internal/models"
)

func (this *repositoryImpl) AddThreadVote(thread *models.Thread, nickname string, voice int) (newVotes int, err error) {
	oldVoice, err := this.getVoice(nickname, thread.ID)
	if err != nil {
		return
	}
	if oldVoice == voice {
		return thread.Votes, nil
	}

	newVoice := voice - oldVoice
	tx, err := this.db.Beginx()
	if err != nil {
		return
	}
	if err = tx.Get(&newVotes, `UPDATE thread SET votes=votes+$1 WHERE id=$2 RETURNING votes`, newVoice, thread.ID); err != nil {
		tx.Rollback()
		return
	}
	if _, err = tx.Exec(`DELETE FROM vote WHERE thread=$1 AND nickname=$2`, thread.ID, nickname); err != nil {
		tx.Rollback()
		return
	}
	if _, err = tx.Exec(`INSERT INTO vote (thread, nickname, voice) VALUES ($1, $2, $3)`, thread.ID, nickname, voice); err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()

	return
}

func (this *repositoryImpl) getVoice(nickname string, threadID int) (int, error) {
	var voice int
	if err := this.db.Get(
		&voice,
		`SELECT voice FROM vote WHERE nickname=$1 AND thread=$2`,
		nickname, threadID,
	); err == sql.ErrNoRows {
		return 0, nil
	} else {
		return voice, err
	}
}
