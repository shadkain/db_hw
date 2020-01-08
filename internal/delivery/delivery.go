package delivery

import (
	"db_hw/internal/usecase"
	"github.com/labstack/echo"
)

type Delivery struct {
	uc usecase.Usecase
}

func NewDelivery(usecase usecase.Usecase) *Delivery {
	return &Delivery{
		uc: usecase,
	}
}

func (d *Delivery) Configure(e *echo.Echo) {
	// Forum
	e.POST("/api/forum/create", d.createForum)
	e.GET("/api/forum/:slug/details", d.takeForum)
	// Post
	e.POST("/api/thread/:slug_or_id/create", d.createPosts)
	e.POST("/api/post/:id/details", d.changePost)
	e.GET("/api/post/:id/details", d.takePostByID)
	e.GET("/api/thread/:slug_or_id/posts", d.takePosts)
	// Thread
	e.POST("/api/forum/:slug/create", d.createThread)
	e.POST("/api/thread/:slug_or_id/details", d.changeThread)
	e.GET("/api/thread/:slug_or_id/details", d.takeThread)
	e.GET("/api/forum/:slug/threads", d.takeForumThreads)
	// User
	e.POST("/api/user/:nickname/create", d.createUser)
	e.GET("/api/user/:nickname/profile", d.takeUser)
	e.POST("/api/user/:nickname/profile", d.changeUser)
	e.GET("/api/forum/:slug/users", d.takeUsersByForum)
	// Vote
	e.POST("/api/thread/:slug_or_id/vote", d.createVote)
	// Service
	e.GET("/api/service/status", d.takeStatus)
	e.POST("/api/service/clear", d.clearAll)
}
