package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/shadkain/db_hw/internal/cache"
	"github.com/shadkain/db_hw/internal/generator"
	"github.com/shadkain/db_hw/internal/models"
	"github.com/shadkain/db_hw/internal/reqmodels"
)

type repositoryImpl struct {
	db               *sqlx.DB
	userCache        *cache.UserCache
	postsIDGenerator *generator.Generator
}

func NewRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{
		db:               db,
		userCache:        cache.NewUserCache(),
		postsIDGenerator: generator.NewGenerator(),
	}
}

type Repository interface {
	// Forum
	GetForumByID(id int) (*models.Forum, error)
	GetForumBySlug(slug string) (*models.Forum, error)
	GetForumSlug(slug string) (*models.Forum, error)
	CreateForum(title, slug, user string) (*models.Forum, error)
	GetForumUsers(forumSlug, since string, limit int, desc bool) ([]*models.User, error)
	// Thread
	GetThreadByID(id int) (*models.Thread, error)
	GetThreadBySlug(slug string) (*models.Thread, error)
	GetThreadBySlugOrID(slugOrID string) (*models.Thread, error)
	GetForumThreads(forum string, limit int, desc bool) ([]*models.Thread, error)
	GetForumThreadsSince(forum, since string, limit int, desc bool) ([]*models.Thread, error)
	GetThreadFieldsBySlugOrID(fields, slugOrID string) (*models.Thread, error)
	CreateThread(forum *models.Forum, thread reqmodels.ThreadCreate) (*models.Thread, error)
	UpdateThread(threadSlugOrID string, message, title string) (*models.Thread, error)
	GetThreadPosts(thread, limit int, since *int, sort string, desc bool) ([]*models.Post, error)
	// Post
	GetPostByID(id int) (*models.Post, error)
	CreatePosts(posts []*reqmodels.PostCreate, thread *models.Thread) ([]*models.Post, error)
	UpdatePostMessage(id int, message string) (*models.Post, error)
	// Vote
	AddThreadVote(thread *models.Thread, nickname string, voice int) (int, error)
	// User
	GetUserByNickname(nickname string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserNickname(nickname string) (string, error)
	GetUsersByNicknameOrEmail(nickname, email string) ([]*models.User, error)
	CreateUser(nickname, email, fullname, about string) (*models.User, error)
	UpdateUserByNickname(nickname, email, fullname, about string) error
	// Service
	CountForums() (int, error)
	CountThreads() (int, error)
	CountPosts() (int, error)
	CountUsers() (int, error)
	Clear() error
}
