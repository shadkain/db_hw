package repository

import (
	"github.com/shadkain/db_hw/internal/models"
	"github.com/shadkain/db_hw/internal/query"
	"errors"
	"github.com/jackc/pgx"
	"strconv"
	"time"
)

func (r *repositoryImpl) InsertPost(newPost models.NewPost, id int, forum string, created time.Time) (LastID int, ThreadID int, Err error) {
	var lastID, threadID int
	err := r.db.QueryRow(query.InsertPost, newPost.Author, newPost.Message,
		newPost.Parent, id, forum).Scan(&lastID, &threadID)

	if err != nil {
		return lastID, threadID, err
	}
	return lastID, threadID, nil
}

func (r *repositoryImpl) SelectPostByID(ID int) (Post models.Post, Err error) {
	var posts []models.Post
	rows, err := r.db.Query(query.SelectPostsByID, ID)
	defer rows.Close()
	if err != nil {
		return models.Post{}, err
	}

	scanPost := models.Post{}
	for rows.Next() {
		var timetz time.Time
		err := rows.Scan(&scanPost.Author, &timetz, &scanPost.Forum,
			&scanPost.ID, &scanPost.IsEdited, &scanPost.Message, &scanPost.Parent,
			&scanPost.Thread)
		if err != nil {
			return models.Post{}, err
		}
		timetz.Format(time.RFC3339Nano)
		posts = append(posts, scanPost)
	}

	if len(posts) == 0 {
		return models.Post{}, errors.New("Can't find user by nickname")
	}
	return posts[0], nil
}

func (r *repositoryImpl) SelectPostByIDThreadID(ID int, threadID int) (Post models.Post, Err error) {
	var posts []models.Post
	rows, err := r.db.Query(query.SelectPostsByIDThreadID, ID, threadID)
	defer rows.Close()
	if err != nil {
		return models.Post{}, err
	}

	scanPost := models.Post{}
	for rows.Next() {
		var timetz time.Time
		err := rows.Scan(&scanPost.Author, &timetz, &scanPost.Forum,
			&scanPost.ID, &scanPost.IsEdited, &scanPost.Message, &scanPost.Parent,
			&scanPost.Thread)
		if err != nil {
			return models.Post{}, err
		}
		scanPost.Created = timetz.Format(time.RFC3339Nano)
		posts = append(posts, scanPost)
	}

	if len(posts) == 0 {
		return models.Post{}, errors.New("Can't find user by nickname")
	}
	return posts[0], nil
}

func (r *repositoryImpl) SelectPosts(threadID int, limit, since, sort, desc string) (Posts *models.Posts, Err error) {
	posts := models.Posts{}

	var rows *pgx.Rows
	var err error
	if sort == "flat" {
		if desc == "false" {
			rows, err = r.db.Query(query.SelectPostsFlat, threadID, limit, since)
		} else {
			rows, err = r.db.Query(query.SelectPostsFlatDesc, threadID, limit, since)
		}

	} else if sort == "tree" {
		if desc == "false" {
			if since != "0" && since != "999999999" {
				rows, err = r.db.Query(query.SelectPostsTree, threadID, 100000)
			} else {
				rows, err = r.db.Query(query.SelectPostsTree, threadID, limit)
			}
		} else {
			if since != "0" && since != "999999999" {
				rows, err = r.db.Query(query.SelectPostsTreeSinceDesc, threadID)
			} else {
				rows, err = r.db.Query(query.SelectPostsTreeDesc, threadID, limit, 1000000)
			}
		}
	} else if sort == "parent_tree" {
		if desc == "false" {
			rows, err = r.db.Query(query.SelectPostsParentTree, threadID)
		} else {
			rows, err = r.db.Query(query.SelectPostsParentTreeDesc, threadID)
		}
	}

	if sort != "parent_tree" {
		defer rows.Close()
		if err != nil {
			return &posts, err
		}

		for rows.Next() {
			scanPost := models.Post{}
			var timetz time.Time
			err := rows.Scan(&scanPost.Author, &timetz, &scanPost.Forum,
				&scanPost.ID, &scanPost.IsEdited, &scanPost.Message, &scanPost.Parent,
				&scanPost.Thread)
			if err != nil {
				return &posts, err
			}
			scanPost.Created = timetz.Format(time.RFC3339Nano)
			posts = append(posts, &scanPost)
		}
	} else {
		if err != nil {
			rows.Close()
			return &posts, err
		}

		count := 0
		limitDigit, _ := strconv.Atoi(limit)

		for rows.Next() {
			scanPost := models.Post{}
			var timetz time.Time
			err := rows.Scan(&scanPost.Author, &timetz, &scanPost.Forum,
				&scanPost.ID, &scanPost.IsEdited, &scanPost.Message, &scanPost.Parent,
				&scanPost.Thread)
			if err != nil {
				return &posts, err
			}

			if scanPost.Parent == 0 {
				count = count + 1
			}
			if count > limitDigit && (since == "0" || since == "999999999") {
				break
			} else {
				scanPost.Created = timetz.Format(time.RFC3339Nano)
				posts = append(posts, &scanPost)
			}

		}
		rows.Close()
	}

	if since != "0" && since != "999999999" && sort == "tree" {
		limitDigit, _ := strconv.Atoi(limit)
		sinceDigit, _ := strconv.Atoi(since)
		sincePosts := models.Posts{}
		counter := 0

		if desc == "false" {
			startIndex := 1000000000
			//postMinStartIndex
			minValue := 100000000000
			for i := 0; i < len(posts); i++ {
				if posts[i].ID == sinceDigit {
					startIndex = i + 1
					break
				}
				if (posts[i].ID > sinceDigit) && (posts[i].ID < minValue) {
					startIndex = i
					minValue = posts[i].ID
				}
			}
			sincePostsCount := 0
			counter = startIndex
			for sincePostsCount < limitDigit && counter < len(posts) {
				scanPost := models.Post{}
				scanPost = *posts[counter]
				sincePosts = append(sincePosts, &scanPost)
				if sort == "tree" {
					sincePostsCount++
				} else {
					if scanPost.Parent == 0 {
						sincePostsCount++
					}
				}
				counter++
			}
		} else {
			startIndex := -1000000000
			//postMinStartIndex
			maxValue := 0
			for i := len(posts) - 1; i >= 0; i-- {
				if posts[i].ID == sinceDigit {
					startIndex = i - 1
					break
				}
				if (posts[i].ID < sinceDigit) && (posts[i].ID > maxValue) {
					startIndex = i
					maxValue = posts[i].ID
				}
			}

			sincePostsCount := 0
			counter = startIndex
			for sincePostsCount < limitDigit && counter >= 0 {
				scanPost := models.Post{}
				scanPost = *posts[counter]
				sincePosts = append(sincePosts, &scanPost)
				if sort == "tree" {
					sincePostsCount++
				} else {
					if scanPost.Parent == 0 {
						sincePostsCount++
					}
				}
				counter--
			}
		}
		return &sincePosts, nil
	}

	if since != "0" && since != "999999999" && sort == "parent_tree" {
		limitDigit, _ := strconv.Atoi(limit)
		sinceDigit, _ := strconv.Atoi(since)
		sincePosts := models.Posts{}
		counter := 0
		if desc == "false" {
			startIndex := 1000000000
			minValue := 100000000000
			for i := 0; i < len(posts); i++ {
				if posts[i].ID == sinceDigit {
					startIndex = i + 1
					break
				}
				if (posts[i].ID > sinceDigit) && (posts[i].ID < minValue) {
					startIndex = i
					minValue = posts[i].ID
				}
			}
			sincePostsCount := 0
			counter = startIndex
			for sincePostsCount < limitDigit && counter < len(posts) {
				scanPost := models.Post{}
				scanPost = *posts[counter]
				sincePosts = append(sincePosts, &scanPost)
				sincePostsCount++
				counter++
			}
		} else {
			startIndex := -1000000000
			//postMinStartIndex
			maxValue := 100000000000
			for i := len(posts) - 1; i >= 0; i-- {
				if posts[i].ID == sinceDigit {
					startIndex = i + 1
					break
				}
				if (posts[i].ID < sinceDigit) && (posts[i].ID < maxValue) {
					startIndex = i
					maxValue = posts[i].ID
				}
			}

			sincePostsCount := 0
			counter = startIndex
			for sincePostsCount < limitDigit && counter < len(posts) {
				scanPost := models.Post{}
				scanPost = *posts[counter]
				sincePosts = append(sincePosts, &scanPost)
				if sort == "tree" {
					sincePostsCount++
				} else {
					if scanPost.Parent == 0 {
						sincePostsCount++
					}
				}
				counter++
			}
		}
		return &sincePosts, nil
	}

	return &posts, nil
}

func (r *repositoryImpl) UpdatePost(changePost models.ChangePost, postID int) (Err error) {
	rows, err := r.db.Query(query.UpdatePostByID, changePost.Message, true, postID)
	defer rows.Close()

	if err != nil {
		return err
	}
	return nil
}
