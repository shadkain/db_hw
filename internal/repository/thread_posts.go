package repository

import (
	"fmt"
	"github.com/shadkain/db_hw/internal/models"
	"github.com/shadkain/db_hw/internal/vars"
	"strings"
)

func (this *repositoryImpl) GetThreadPosts(thread, limit int, since *int, sort string, desc bool) (models.Posts, error) {
	switch sort {
	case SortFlat, "":
		return this.getThreadPostsFlat(thread, limit, since, desc)
	case SortTree:
		return this.getThreadPostsTree(thread, limit, since, desc)
	case SortParentTree:
		return this.getThreadPostsParentTree(thread, limit, since, desc)
	}
	return nil, fmt.Errorf("%w: unknown sort method '%s'", vars.ErrNotFound, sort)
}

func (this *repositoryImpl) getThreadPostsFlat(thread, limit int, since *int, desc bool) (models.Posts, error) {
	order := "asc"
	if desc {
		order = "desc"
	}

	orderBy := []string{"created " + order, "id " + order}
	filter := "thread=$1"
	params := []interface{}{thread}

	if since != nil {
		if desc {
			filter += " AND id<$2"
		} else {
			filter += " AND id>$2"
		}

		params = append(params, *since)
	}

	return this.getPosts(orderBy, limit, filter, params...)
}

func (this *repositoryImpl) getThreadPostsTree(thread, limit int, since *int, desc bool) (models.Posts, error) {
	conditions := []string{"thread=$1"}
	params := []interface{}{thread}

	if since != nil {
		sinceCond, err := this.getSinceCondition(since, desc)
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, sinceCond)
	}

	orderBy := []string{"path " + this.getOrder(desc)}
	filter := strings.Join(conditions, " AND ")

	return this.getPosts(orderBy, limit, filter, params...)
}

func (this *repositoryImpl) getThreadPostsParentTree(thread, limit int, since *int, desc bool) (models.Posts, error) {
	conditions := []string{"parent=0", "thread=$1"}

	if since != nil {
		var operator = ">"
		if desc {
			operator = "<"
		}
		sincePost, err := this.getPostFields("path", "id=$1", *since)
		if err != nil {
			return nil, err
		}
		sinceCond := fmt.Sprintf("path %s '%s'", operator, this.getRootPath(sincePost.Path))
		conditions = append(conditions, sinceCond)
	}

	filter := strings.Join(conditions, " AND ")
	var parents models.Posts
	if err := this.db.Select(&parents, fmt.Sprintf(
		`SELECT * FROM post WHERE %s ORDER BY id %s LIMIT %d`, filter, this.getOrder(desc), limit),
		thread,
	); err != nil {
		return nil, err
	}

	posts := make(models.Posts, 0)
	for _, parent := range parents {
		var childs models.Posts
		if err := this.db.Select(
			&childs,
			fmt.Sprintf(
				`SELECT * FROM post WHERE substring(path,1,7) = '%s' AND parent<>0 ORDER BY path`,
				this.padPostID(parent.ID),
			),
		); err != nil {
			return nil, err
		}

		posts = append(posts, parent)
		posts = append(posts, childs...)
	}

	return posts, nil
}

func (this *repositoryImpl) getSinceCondition(since *int, desc bool) (string, error) {
	var operator = ">"
	if desc {
		operator = "<"
	}

	if sincePost, err := this.getPostFields("path", "id=$1", *since); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("path %s '%s'", operator, sincePost.Path), nil
	}
}

func (this *repositoryImpl) getRootPath(path string) string {
	root := strings.Split(path, pathDelim)[0]
	return root + strings.Repeat(pathDelim+zeroPath, maxTreeLevel-1)
}
