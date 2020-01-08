package usecase

import (
	"errors"
	"db_hw/internal/models"
	"strconv"
)

func (uc *usecaseImpl) AddThread(newThread models.NewThread, forum string) (Thread models.Thread, Err error) {
	threads, err := uc.repo.SelectThreadsBySlug(newThread.Slug)
	if len(*threads) > 0 {
		return *(*threads)[0], errors.New("conflict")
	}
	lastID, err := uc.repo.InsertThread(newThread, forum)
	if err != nil {
		return models.Thread{}, err
	}

	thread := models.Thread{
		Author:  newThread.Author,
		Created: newThread.Created,
		Forum:   forum,
		ID:      lastID,
		Message: newThread.Message,
		Slug:    newThread.Slug,
		Title:   newThread.Title,
		Votes:   0,
	}

	return thread, nil
}

func (uc *usecaseImpl) GetThreadBySlug(slugOrID string) (Thread models.Thread, Err error) {
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
		thread = *(*threads)[0]
	} else {
		threads, err := uc.repo.SelectThreadsByID(id)
		if err != nil {
			return models.Thread{}, err
		}
		if len(*threads) != 1 {
			return models.Thread{}, errors.New("Can't find thread")
		}
		thread = *(*threads)[0]
	}

	return thread, nil
}

func (uc *usecaseImpl) GetThreadsByForum(forum string, limit string, since string, desc string) (Threads *models.Threads, Err error) {
	threads, err := uc.repo.SelectThreadsByForum(forum, limit, since, desc)
	if err != nil {
		return threads, err
	}

	return threads, nil
}

func (uc *usecaseImpl) SetThread(changeThread models.ChangeThread, slugOrID string) (Thread models.Thread, Err error) {
	thread := models.Thread{}
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		threads, err := uc.repo.SelectThreadsBySlug(slugOrID)
		if err != nil {
			return models.Thread{}, err
		}
		if len(*threads) != 1 {
			return models.Thread{}, errors.New("Can't find thread")
		}
		thread = *(*threads)[0]
		id = (*threads)[0].ID
	} else {
		threads, err := uc.repo.SelectThreadsByID(id)
		if err != nil {
			return models.Thread{}, err
		}
		if len(*threads) != 1 {
			return models.Thread{}, errors.New("Can't find thread")
		}
		thread = *(*threads)[0]
		id = (*threads)[0].ID
	}

	if changeThread.Message == "" {
		changeThread.Message = thread.Message
	} else {
		thread.Message = changeThread.Message
	}
	if changeThread.Title == "" {
		changeThread.Title = thread.Title
	} else {
		thread.Title = changeThread.Title
	}

	if err := uc.repo.UpdateThread(changeThread, id); err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}
