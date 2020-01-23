package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shadkain/db_hw/internal/models"
	"github.com/shadkain/db_hw/internal/reqmodels"
	"strconv"
	"strings"
	"time"
)

const (
	SortFlat       = "flat"
	SortTree       = "tree"
	SortParentTree = "parent_tree"

	pathDelim     = "."
	maxIDLength   = 7
	maxTreeLevel  = 5
	postChunkSize = 50
	columnCount   = 8
)

var zeroPath = strings.Repeat("0", maxIDLength)

func (this *repositoryImpl) GetPostByID(id int) (*models.Post, error) {
	return this.getPost("id=$1", id)
}

func (this *repositoryImpl) CreatePosts(posts []*reqmodels.PostCreate, thread *models.Thread) (models.Posts, error) {
	forum, err := this.GetForumSlug(thread.Forum)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	result := make(models.Posts, 0, len(posts))
	for _, chunk := range this.chunkPosts(posts) {
		createdIDs, err := this.createPostsChunk(forum, thread, chunk, now)
		if err != nil {
			return nil, err
		}

		created, err := this.getPostsByIDs(createdIDs)
		if err != nil {
			return nil, err
		}

		result = append(result, created...)
	}

	return result, nil
}

func (this *repositoryImpl) UpdatePostMessage(id int, message string) (*models.Post, error) {
	if message != "" {
		if _, err := this.db.Exec(
			`UPDATE post SET "message"=$1, "isEdited"=true WHERE id=$2 AND "message"<>$1`,
			message, id,
		); err != nil {
			return nil, err
		}
	}

	return this.GetPostByID(id)
}

func (this *repositoryImpl) chunkPosts(posts []*reqmodels.PostCreate) [][]*reqmodels.PostCreate {
	chunks := make([][]*reqmodels.PostCreate, 0)
	for i := 0; i < len(posts); i += postChunkSize {
		end := i + postChunkSize
		if end > len(posts) {
			end = len(posts)
		}

		chunks = append(chunks, posts[i:end])
	}

	return chunks
}

func (this *repositoryImpl) createPostsChunk(forum *models.Forum, thread *models.Thread, posts []*reqmodels.PostCreate, created time.Time) ([]int, error) {
	columns := columnCount
	placeholders := make([]string, 0, len(posts))
	args := make([]interface{}, 0, len(posts)*columns)
	ids := this.postsIDGenerator.Next(len(posts))

	for i, post := range posts {
		id := ids[i]
		path, err := this.getPostPath(id, post.Parent)
		if err != nil {
			return nil, err
		}
		args = append(args, id, thread.ID, thread.Forum, post.Parent, path, post.Author, post.Message, created)
		placeholders = append(placeholders, fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*columns+1, i*columns+2, i*columns+3, i*columns+4, i*columns+5, i*columns+6, i*columns+7, i*columns+8,
		))
	}

	query := fmt.Sprintf(
		"INSERT INTO post (id, thread, forum, parent, path, author, message, created) VALUES %s",
		strings.Join(placeholders, ","),
	)

	var err error
	for {
		_, err = this.db.Exec(query, args...)
		if err == nil || err.Error() != "ERROR: deadlock detected (SQLSTATE 40P01)" {
			break
		}
	}

	return ids, err
}

func (this *repositoryImpl) getPost(filter string, params ...interface{}) (*models.Post, error) {
	return this.getPostFields("*", filter, params...)
}

func (this *repositoryImpl) getPostFields(fields, filter string, params ...interface{}) (*models.Post, error) {
	var post models.Post
	if err := this.db.Get(
		&post,
		"SELECT "+fields+" FROM post WHERE "+filter,
		params...,
	); err != nil {
		return nil, Error(err)
	}

	return &post, nil
}

func (this *repositoryImpl) getPostsByIDs(ids []int) (models.Posts, error) {
	posts := make(models.Posts, 0)

	if query, args, err := sqlx.In(
		`SELECT * FROM post WHERE id IN (?) ORDER BY id`,
		ids,
	); err != nil {
		return nil, err
	} else {
		query = this.db.Rebind(query)
		err = this.db.Select(&posts, query, args...)
		return posts, err
	}
}

func (this *repositoryImpl) getPosts(orderBy []string, limit int, filter string, params ...interface{}) (models.Posts, error) {
	query := fmt.Sprintf(`SELECT * FROM post WHERE %s ORDER BY %s`, filter, strings.Join(orderBy, ","))

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	posts := make(models.Posts, 0, limit)
	err := this.db.Select(&posts, query, params...)

	return posts, err
}

func (this *repositoryImpl) getPostPath(id, parentID int) (string, error) {
	var base string

	if parentID == 0 {
		base = this.getZeroPostPath()
	} else {
		parent, err := this.getPostFields("path", "id=$1", parentID)
		if err != nil {
			return "", err
		}

		base = parent.Path
	}

	path := strings.Replace(base, zeroPath, this.padPostID(id), 1)

	return path, nil
}

func (this *repositoryImpl) getZeroPostPath() string {
	path := zeroPath

	for i := 0; i < maxTreeLevel-1; i++ {
		path += pathDelim + zeroPath
	}

	return path
}

func (this *repositoryImpl) padPostID(id int) string {
	return fmt.Sprintf("%0"+strconv.Itoa(maxIDLength)+"d", id)
}

func (this *repositoryImpl) updatePostPath(tx *sqlx.Tx, id int, path string) error {
	_, err := tx.Exec(`UPDATE post SET path=$1 WHERE id=$2`, path, id)
	return err
}
