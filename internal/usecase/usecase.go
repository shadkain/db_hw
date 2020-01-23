package usecase

import (
	"github.com/shadkain/db_hw/internal/models"
	"github.com/shadkain/db_hw/internal/repository"
	"github.com/shadkain/db_hw/internal/reqmodels"
)

type usecaseImpl struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecaseImpl{
		repo: repo,
	}
}

type Usecase interface {
	// Forum
	GetForum(slug string) (*models.Forum, error)
	GetForumThreads(forumSlug, since string, limit int, desc bool) ([]*models.Thread, error)
	GetForumUsers(forum, since string, limit int, desc bool) ([]*models.User, error)
	CreateForum(title, slug, nickname string) (*models.Forum, error)
	// Thread
	GetThread(threadSlugOrID string) (*models.Thread, error)
	GetThreadPosts(threadSlugOrID string, limit int, since *int, sort string, desc bool) ([]*models.Post, error)
	CreateThread(forumSlug string, thread reqmodels.ThreadCreate) (*models.Thread, error)
	UpdateThread(threadSlugOrID string, message, title string) (*models.Thread, error)
	// Post
	GetPostDetails(id int, related []string) (*models.PostDetails, error)
	CreatePosts(threadSlugOrID string, posts []*reqmodels.PostCreate) ([]*models.Post, error)
	UpdatePost(id int, message string) (*models.Post, error)
	// Vote
	VoteForThread(threadSlugOrID string, vote reqmodels.Vote) (*models.Thread, error)
	// User
	GetUserByNickname(nickname string) (*models.User, error)
	CreateUser(nickname, email, fullname, about string) ([]*models.User, error)
	UpdateUser(nickname, email, fullname, about string) (*models.User, error)
	// Service
	GetStatus() (reqmodels.Status, error)
	Clear() error
}
