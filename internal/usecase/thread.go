package usecase

import (
	"fmt"
	"github.com/shadkain/db_hw/internal/reqmodels"
	"github.com/shadkain/db_hw/internal/vars"
	"github.com/shadkain/db_hw/internal/models"
	"time"
)

func (this *usecaseImpl) GetThread(threadSlugOrID string) (*models.Thread, error) {
	return this.repo.GetThreadBySlugOrID(threadSlugOrID)
}

func (this *usecaseImpl) GetThreadPosts(threadSlugOrID string, limit int, since *int, sort string, desc bool) (models.Posts, error) {
	thread, err := this.repo.GetThreadFieldsBySlugOrID("id", threadSlugOrID)
	if err != nil {
		return nil, err
	}

	return this.repo.GetThreadPosts(thread.ID, limit, since, sort, desc)
}

func (this *usecaseImpl) CreateThread(forumSlug string, thread reqmodels.ThreadCreate) (*models.Thread, error) {
	if _, err := this.repo.GetUserNickname(thread.Author); err != nil {
		return nil, err
	}

	forum, err := this.repo.GetForumSlug(forumSlug)
	if err != nil {
		return nil, err
	}

	if thread.Slug != "" {
		existing, err := this.repo.GetThreadBySlug(thread.Slug)
		if err != nil && err != vars.ErrNotFound {
			return nil, err
		}
		if existing != nil {
			return existing, fmt.Errorf("%w: thread with this slug already exists", vars.ErrConflict)
		}
	}

	if thread.Created == "" {
		thread.Created = time.Now().Format(time.RFC3339)
	}

	return this.repo.CreateThread(forum, thread)
}

func (this *usecaseImpl) UpdateThread(threadSlugOrID string, message, title string) (*models.Thread, error) {
	return this.repo.UpdateThread(threadSlugOrID, message, title)
}
