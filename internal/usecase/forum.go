package usecase

import (
	"fmt"
	"github.com/shadkain/db_hw/internal/models"
	"github.com/shadkain/db_hw/internal/vars"
)

func (this *usecaseImpl) GetForum(slug string) (*models.Forum, error) {
	return this.repo.GetForumBySlug(slug)
}

func (this *usecaseImpl) GetForumThreads(forumSlug, since string, limit int, desc bool) ([]*models.Thread, error) {
	forum, err := this.repo.GetForumSlug(forumSlug)
	if err != nil {
		return nil, err
	}

	var threads models.Threads
	if since == "" {
		threads, err = this.repo.GetForumThreads(forum.Slug, limit, desc)
	} else {
		threads, err = this.repo.GetForumThreadsSince(forum.Slug, since, limit, desc)
	}
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (this *usecaseImpl) GetForumUsers(forum, since string, limit int, desc bool) ([]*models.User, error) {
	return this.repo.GetForumUsers(forum, since, limit, desc)
}

func (this *usecaseImpl) CreateForum(title, slug, nickname string) (*models.Forum, error) {
	trueNickname, err := this.repo.GetUserNickname(nickname)
	if err != nil {
		return nil, err
	}

	if exForum, err := this.repo.GetForumBySlug(slug); err != nil && err != vars.ErrNotFound {
		return nil, err
	} else if exForum != nil {
		return exForum, fmt.Errorf("%w: forum with this slug already exists", vars.ErrConflict)
	}

	return this.repo.CreateForum(title, slug, trueNickname)
}
