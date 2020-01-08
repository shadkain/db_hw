package repository

import (
	"db_hw/internal/models"
	"db_hw/internal/query"
	"github.com/jackc/pgx"
	"time"
)

func (r *repositoryImpl) InsertThread(newThread models.NewThread, forum string) (LastID int, Err error) {
	var lastID int
	if newThread.Slug == "" {
		if newThread.Created == "" {
			Err = r.db.QueryRow(query.InsertThreadWithoutCreated, newThread.Author,
				newThread.Message, newThread.Title, forum).Scan(&lastID)
		} else {
			Err = r.db.QueryRow(query.InsertThread, newThread.Author, newThread.Created,
				newThread.Message, newThread.Title, forum).Scan(&lastID)
		}
	} else {
		if newThread.Created == "" {
			Err = r.db.QueryRow(query.InsertThreadWithSlugWithoutCreated, newThread.Author,
				newThread.Message, newThread.Title, forum, newThread.Slug).Scan(&lastID)
		} else {
			Err = r.db.QueryRow(query.InsertThreadWithSlug, newThread.Author, newThread.Created,
				newThread.Message, newThread.Title, forum, newThread.Slug).Scan(&lastID)
		}
	}
	if Err != nil {
		return lastID, Err
	}
	return lastID, nil
}

func (r *repositoryImpl) SelectThreadsBySlug(slug string) (Threads *models.Threads, Err error) {
	threads := models.Threads{}

	rows, err := r.db.Query(query.SelectThreadsBySlug, slug)

	defer rows.Close()
	if err != nil {
		return &threads, err
	}

	for rows.Next() {
		scanThread := models.Thread{}
		var timetz time.Time
		err := rows.Scan(&scanThread.Author, &timetz, &scanThread.Forum,
			&scanThread.ID, &scanThread.Message, &scanThread.Slug, &scanThread.Title,
			&scanThread.Votes)
		if err != nil {
			return &threads, err
		}
		scanThread.Created = timetz.Format(time.RFC3339Nano)
		threads = append(threads, &scanThread)
	}
	return &threads, nil
}

func (r *repositoryImpl) SelectThreadsByID(id int) (Threads *models.Threads, Err error) {
	threads := models.Threads{}

	rows, err := r.db.Query(query.SelectThreadsByID, id)

	defer rows.Close()
	if err != nil {
		return &threads, err
	}

	for rows.Next() {
		scanThread := models.Thread{}
		var timetz time.Time
		err := rows.Scan(&scanThread.Author, &timetz, &scanThread.Forum,
			&scanThread.ID, &scanThread.Message, &scanThread.Slug, &scanThread.Title,
			&scanThread.Votes)
		if err != nil {
			return &threads, err
		}
		scanThread.Created = timetz.Format(time.RFC3339Nano)
		threads = append(threads, &scanThread)
	}
	return &threads, nil
}

func (r *repositoryImpl) SelectThreadsByForum(forum string, limit string, since string, desc string) (Threads *models.Threads, Err error) {
	threads := models.Threads{}
	var rows *pgx.Rows
	var err error
	if since == "" && desc == "false" {
		rows, err = r.db.Query(query.SelectThreadsByForum, forum, limit)
	} else if since != "" && desc == "false" {
		rows, err = r.db.Query(query.SelectThreadsByForumSince, forum, limit, since)
	} else if since == "" && desc == "true" {
		rows, err = r.db.Query(query.SelectThreadsByForumDesc, forum, limit)
	} else {
		rows, err = r.db.Query(query.SelectThreadsByForumSinceDesc, forum, limit, since)
	}
	defer rows.Close()
	if err != nil {
		return &threads, err
	}

	for rows.Next() {
		scanThread := models.Thread{}
		var timetz time.Time
		err := rows.Scan(&scanThread.Author, &timetz, &scanThread.Forum,
			&scanThread.ID, &scanThread.Message, &scanThread.Slug, &scanThread.Title,
			&scanThread.Votes)
		if err != nil {
			return &threads, err
		}
		scanThread.Created = timetz.Format(time.RFC3339Nano)
		threads = append(threads, &scanThread)
	}
	return &threads, nil
}

func (r *repositoryImpl) UpdateThread(changeThread models.ChangeThread, id int) (Err error) {
	rows, err := r.db.Query(query.UpdateThreadByID, changeThread.Message,
		changeThread.Title, id)
	defer rows.Close()

	if err != nil {
		return err
	}
	return nil
}
