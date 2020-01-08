package repository

import (
	"github.com/shadkain/db_hw/internal/models"
	"github.com/shadkain/db_hw/internal/query"
)

func (r *repositoryImpl) InsertForum(newForum models.NewForum) (Err error) {
	var lastID int
	err := r.db.QueryRow(query.InsertForum, newForum.Slug, newForum.Title,
		newForum.User).Scan(&lastID)

	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) SelectForumsBySlug(slug string) (Forum []models.Forum, Err error) {
	var forums []models.Forum
	rows, err := r.db.Query(query.SelectForumsBySlug, slug)
	defer rows.Close()
	if err != nil {
		return forums, err
	}

	scanForum := models.Forum{}
	for rows.Next() {
		err := rows.Scan(&scanForum.Posts, &scanForum.Slug, &scanForum.Thread,
			&scanForum.Title, &scanForum.User)
		if err != nil {
			return forums, err
		}
		forums = append(forums, scanForum)
	}

	return forums, nil
}
