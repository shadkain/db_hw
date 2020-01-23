package repository

import (
	"fmt"
	"github.com/shadkain/db_hw/internal/models"
	"strings"
)

func (this *repositoryImpl) GetForumByID(id int) (*models.Forum, error) {
	return this.getForum("*", "id=$1", id)
}

func (this *repositoryImpl) GetForumBySlug(slug string) (*models.Forum, error) {
	return this.getForum("*", "slug=$1", slug)
}

func (this *repositoryImpl) GetForumSlug(slug string) (*models.Forum, error) {
	return this.getForum("slug", "slug=$1", slug)
}

func (this *repositoryImpl) CreateForum(title, slug, user string) (*models.Forum, error) {
	var id int
	if err := this.db.QueryRow(
		`INSERT INTO forum (title, slug, "user") VALUES ($1, $2, $3) RETURNING id`,
		title, slug, user,
	).Scan(&id); err != nil {
		return nil, err
	}

	return this.GetForumByID(id)
}

func (this *repositoryImpl) GetForumUsers(forumSlug, since string, limit int, desc bool) ([]*models.User, error) {
	forum, err := this.GetForumSlug(forumSlug)
	if err != nil {
		return nil, err
	}

	sinceFilter := ""
	if since != "" {
		if desc {
			sinceFilter = "AND nickname<$2"
		} else {
			sinceFilter = "AND nickname>$2"
		}
	}

	query := fmt.Sprintf(
		`SELECT "user".* FROM "user"
         		JOIN forum_user ON nickname = forum_user.user
				WHERE forum = $1 %s ORDER BY nickname %s %s`,
		sinceFilter, this.getOrder(desc), this.getLimit(limit),
	)

	users := make(models.Users, 0)
	if since == "" {
		err = this.db.Select(&users, query, forum.Slug)
	} else {
		err = this.db.Select(&users, query, forum.Slug, since)
	}

	return users, err
}

func (this *repositoryImpl) getForum(fields, filter string, params ...interface{}) (*models.Forum, error) {
	var forum models.Forum
	if err := this.db.Get(
		&forum,
		`SELECT `+fields+` FROM forum WHERE `+filter,
		params...,
	); err != nil {
		return nil, Error(err)
	}

	var err error
	var includePosts = fields == "*" || strings.Contains(fields, "posts")

	if forum.Posts != 0 || !includePosts {
		return &forum, nil
	}
	if forum.Posts, err = this.countForumPosts(forum.Slug); err != nil {
		return nil, err
	}
	if err := this.updateForumPostsCount(forum.ID, forum.Posts); err != nil {
		return nil, err
	}

	return &forum, nil
}

func (this *repositoryImpl) countForumPosts(forumSlug string) (int, error) {
	var count int
	if err := this.db.Get(
		&count,
		`SELECT COUNT(*) FROM post WHERE forum=$1`,
		forumSlug,
	); err != nil {
		return 0, err
	}

	return count, nil
}

func (this *repositoryImpl) updateForumPostsCount(id, posts int) error {
	_, err := this.db.Exec(
		`UPDATE forum SET posts=$1 WHERE id=$2`,
		posts, id,
	)
	return err
}
