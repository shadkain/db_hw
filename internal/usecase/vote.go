package usecase

import (
	"errors"
	"github.com/shadkain/db_hw/internal/models"
	"github.com/jackc/pgx"
	"strconv"
)

func (uc *usecaseImpl) SetVote(newVote models.NewVote, slugOrID string) (Thread models.Thread, Err error) {
	_, err := uc.repo.SelectUserByNickname(newVote.Nickname)
	if err != nil {
		return models.Thread{}, err
	}

	var thread models.Thread
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		threads, err := uc.repo.SelectThreadsBySlug(slugOrID)
		if err != nil {
			return models.Thread{}, err
		}
		if len(*threads) != 1 {
			return models.Thread{}, errors.New("Can't find thread")
		}
		id = (*threads)[0].ID
		thread = *(*threads)[0]
	} else {
		threads, err := uc.repo.SelectThreadsByID(id)
		if err != nil {
			return models.Thread{}, err
		}
		if len(*threads) != 1 {
			return models.Thread{}, errors.New("Can't find thread")
		}
		id = (*threads)[0].ID
		thread = *(*threads)[0]
	}

	err = uc.repo.InsertVote(newVote, id)
	if err != nil {
		pqErr, ok := err.(pgx.PgError)
		if !ok {
			return models.Thread{}, err
		}
		if pqErr.Code == "23505" {
			err = uc.repo.UpdateVote(newVote, id)
		} else {
			return models.Thread{}, err
		}
	}
	if newVote.Voice == 1 {
		thread.Votes = thread.Votes + newVote.Voice
	} else {
		thread.Votes = thread.Votes - 2
	}

	return thread, nil
}
