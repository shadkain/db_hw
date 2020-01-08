package usecase

import (
	"github.com/shadkain/db_hw/internal/models"
)

func (uc *usecaseImpl) AddForum(newForum models.NewForum) (Forum models.Forum, Err error) {
	user, err := uc.repo.SelectUserByNickname(newForum.User)
	if err != nil {
		return models.Forum{}, err
	}

	newForum.User = user.Nickname

	if err := uc.repo.InsertForum(newForum); err != nil {
		return models.Forum{}, err
	}

	forum := models.Forum{
		Posts:  0,
		Slug:   newForum.Slug,
		Thread: 0,
		Title:  newForum.Title,
		User:   newForum.User,
	}

	return forum, nil
}

func (uc *usecaseImpl) GetForumsBySlug(slug string) (Forum []models.Forum, Err error) {
	forums, err := uc.repo.SelectForumsBySlug(slug)
	if err != nil {
		return forums, err
	}

	return forums, nil
}
