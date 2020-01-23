package usecase

import (
	"fmt"
	"github.com/shadkain/db_hw/internal/models"
	"github.com/shadkain/db_hw/internal/reqmodels"
	"github.com/shadkain/db_hw/internal/vars"
)

func (this *usecaseImpl) GetPostDetails(id int, related []string) (*models.PostDetails, error) {
	post, err := this.repo.GetPostByID(id)
	if err != nil {
		return nil, err
	}

	details := models.PostDetails{
		Post: post,
	}

	for _, r := range related {
		switch r {
		case "user":
			details.Author, err = this.repo.GetUserByNickname(post.Author)
		case "forum":
			details.Forum, err = this.repo.GetForumBySlug(post.Forum)
		case "thread":
			details.Thread, err = this.repo.GetThreadByID(post.Thread)
		}

		if err != nil {
			return nil, err
		}
	}

	return &details, nil
}

func (this *usecaseImpl) CreatePosts(threadSlugOrID string, posts []*reqmodels.PostCreate) (models.Posts, error) {
	thread, err := this.repo.GetThreadFieldsBySlugOrID("id, forum", threadSlugOrID)
	if err != nil {
		return nil, err
	}

	if err := this.checkPostsCreate(posts, thread.ID); err != nil {
		return nil, err
	}

	return this.repo.CreatePosts(posts, thread)
}

func (this *usecaseImpl) UpdatePost(id int, message string) (*models.Post, error) {
	return this.repo.UpdatePostMessage(id, message)
}

func (this *usecaseImpl) checkPostsCreate(posts []*reqmodels.PostCreate, threadID int) error {
	for _, post := range posts {
		if err := this.checkPostCreate(post, threadID); err != nil {
			return err
		}
	}

	return nil
}

func (this *usecaseImpl) checkPostCreate(post *reqmodels.PostCreate, threadID int) error {
	if _, err := this.repo.GetUserNickname(post.Author); err != nil {
		return err
	}

	if post.Parent != 0 {
		parent, err := this.repo.GetPostByID(post.Parent)
		if err == vars.ErrNotFound {
			return fmt.Errorf("%w: post parent do not exists", vars.ErrConflict)
		}
		if err != nil {
			return err
		}
		if parent.Thread != threadID {
			return fmt.Errorf("%w: parent post was created in another thread", vars.ErrConflict)
		}
	}

	return nil
}
