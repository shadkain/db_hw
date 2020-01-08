package usecase

import (
	"errors"
	"db_hw/internal/models"
	"strconv"
	"strings"
	"time"
)

const date = "1970-01-01T00:00:00.000Z"

func (uc *usecaseImpl) AddPosts(newPosts models.NewPosts, slug_or_id string) (Posts models.Posts, Err error) {
	var forum string
	id, err := strconv.Atoi(slug_or_id)
	if err != nil {
		threads, err := uc.repo.SelectThreadsBySlug(slug_or_id)
		if err != nil {
			return models.Posts{}, err
		}
		if len(*threads) != 1 {
			return models.Posts{}, errors.New("Can't find thread")
		}
		forum = (*threads)[0].Forum
		id = (*threads)[0].ID
	} else {
		threads, err := uc.repo.SelectThreadsByID(id)
		if err != nil {
			return models.Posts{}, err
		}
		if len(*threads) != 1 {
			return models.Posts{}, errors.New("Can't find thread")
		}
		forum = (*threads)[0].Forum
		id = (*threads)[0].ID
	}

	posts := models.Posts{}
	created := time.Now()

	for _, newPost := range newPosts {
		if newPost.Parent != 0 {
			_, err := uc.repo.SelectPostByIDThreadID(newPost.Parent, id)
			if err != nil {
				return models.Posts{}, err
			}
		}

		lastID, threadID, err := uc.repo.InsertPost(*newPost, id, forum, created)
		if err != nil {
			return models.Posts{}, err
		}
		post := models.Post{
			Author:   newPost.Author,
			Created:  "",
			Forum:    forum,
			ID:       lastID,
			IsEdited: false,
			Message:  newPost.Message,
			Parent:   newPost.Parent,
			Thread:   threadID,
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (uc *usecaseImpl) GetPostByID(ID int, related string) (Post models.PostDetails, Err error) {
	var postDetails models.PostDetails

	post, err := uc.repo.SelectPostByID(ID)
	if err != nil {
		return postDetails, err
	}
	post.Created = date
	postDetails.Post = post

	var user models.User
	if strings.Contains(related, "user") {
		user, err = uc.repo.SelectUserByNickname(post.Author)
		if err != nil {
			return postDetails, nil
		}
		postDetails.User = user
	}

	if strings.Contains(related, "thread") {
		threads, err := uc.repo.SelectThreadsByID(post.Thread)
		if err != nil || len(*threads) != 1 {
			return postDetails, nil
		}

		postDetails.Thread = (*threads)[0]
	}

	if strings.Contains(related, "forum") {
		forums, err := uc.repo.SelectForumsBySlug(post.Forum)
		if err != nil || len(forums) != 1 {
			return postDetails, nil
		}

		postDetails.Forum = forums[0]
	}

	return postDetails, nil
}

func (uc *usecaseImpl) GetPosts(slugOrID, limit, since, sort, desc string) (Posts *models.Posts, Err error) {
	var thread models.Thread
	id, err := strconv.Atoi(slugOrID)
	if err != nil {
		threads, err := uc.repo.SelectThreadsBySlug(slugOrID)
		if err != nil {
			return &models.Posts{}, err
		}
		if len(*threads) != 1 {
			return &models.Posts{}, errors.New("Can't find thread")
		}
		thread = *(*threads)[0]
	} else {
		threads, err := uc.repo.SelectThreadsByID(id)
		if err != nil {
			return &models.Posts{}, err
		}
		if len(*threads) != 1 {
			return &models.Posts{}, errors.New("Can't find thread")
		}
		thread = *(*threads)[0]
	}

	posts, err := uc.repo.SelectPosts(thread.ID, limit, since, sort, desc)
	if err != nil {
		return posts, err
	}

	return posts, err
}

func (uc *usecaseImpl) SetPost(changePost models.ChangePost, postID int) (Post models.Post, Err error) {
	post, err := uc.repo.SelectPostByID(postID)
	if err != nil {
		return post, err
	}

	post.Created = date
	if changePost.Message == "" {
		changePost.Message = post.Message
		post.IsEdited = false
		return post, nil
	} else if changePost.Message == post.Message {
		post.IsEdited = false
		return post, nil
	} else {
		post.Message = changePost.Message
	}
	post.IsEdited = true

	if err := uc.repo.UpdatePost(changePost, postID); err != nil {
		return models.Post{}, err
	}

	return post, nil
}
