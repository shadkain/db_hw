package repository

import (
	"fmt"
	"github.com/shadkain/db_hw/internal/models"
	"github.com/shadkain/db_hw/internal/reqmodels"
	"strconv"
)

func (this *repositoryImpl) GetThreadByID(id int) (*models.Thread, error) {
	return this.getThread("*", "id=$1", id)
}

func (this *repositoryImpl) GetThreadBySlug(slug string) (*models.Thread, error) {
	return this.getThread("*", "slug=$1", slug)
}

func (this *repositoryImpl) GetThreadBySlugOrID(slugOrID string) (*models.Thread, error) {
	return this.GetThreadFieldsBySlugOrID("*", slugOrID)
}

func (this *repositoryImpl) GetForumThreads(forum string, limit int, desc bool) ([]*models.Thread, error) {
	query := fmt.Sprintf(
		"SELECT * FROM thread WHERE forum=$1 ORDER BY created %s LIMIT $2",
		this.getOrder(desc),
	)

	var threads []*models.Thread
	err := this.db.Select(&threads, query, forum, limit)

	return threads, err
}

func (this *repositoryImpl) GetForumThreadsSince(forum, since string, limit int, desc bool) ([]*models.Thread, error) {
	createdCond := ">="
	if desc {
		createdCond = "<="
	}

	query := fmt.Sprintf(
		"SELECT * FROM thread WHERE forum=$1 AND created %s $2 ORDER BY created %s LIMIT $3",
		createdCond, this.getOrder(desc),
	)

	threads := make([]*models.Thread, 0)
	err := this.db.Select(&threads, query, forum, since, limit)

	return threads, err
}

func (this *repositoryImpl) GetThreadFieldsBySlugOrID(fields, slugOrID string) (*models.Thread, error) {
	if id, err := strconv.Atoi(slugOrID); err != nil {
		return this.getThread(fields, "slug=$1", slugOrID)
	} else {
		return this.getThread(fields, "id=$1", id)
	}
}

func (this *repositoryImpl) CreateThread(forum *models.Forum, thread reqmodels.ThreadCreate) (*models.Thread, error) {
	var id int
	err := this.db.
		QueryRow(
			`INSERT INTO thread (title, author, forum, message, slug, created) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			thread.Title, thread.Author, forum.Slug, thread.Message, thread.Slug, thread.Created,
		).
		Scan(&id)
	if err != nil {
		return nil, err
	}

	return this.GetThreadByID(id)
}

func (this *repositoryImpl) UpdateThread(threadSlugOrID string, message, title string) (*models.Thread, error) {
	thread, err := this.GetThreadBySlugOrID(threadSlugOrID)
	if err != nil {
		return nil, err
	}

	if message != "" {
		thread.Message = message
	}
	if title != "" {
		thread.Title = title
	}

	_, err = this.db.Exec(
		`UPDATE thread SET "message"=$1, title=$2 WHERE id=$3`,
		thread.Message, thread.Title, thread.ID,
	)

	return thread, err
}

func (this *repositoryImpl) getThread(fields, filter string, params ...interface{}) (*models.Thread, error) {
	t := models.Thread{}
	err := this.db.Get(&t, "SELECT "+fields+" FROM thread WHERE "+filter, params...)
	if err != nil {
		return nil, Error(err)
	}
	return &t, nil
}
