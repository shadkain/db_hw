package delivery

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/shadkain/db_hw/internal/usecase"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	uc     usecase.Usecase
	router *fasthttprouter.Router
}

func NewHandler(usecase usecase.Usecase) *Handler {
	this := &Handler{
		uc:     usecase,
		router: fasthttprouter.New(),
	}

	this.configure()

	return this
}

func (this *Handler) configure() {
	// Forum
	this.router.GET("/api/forum/:slug/details", this.getForumDetails)
	this.router.GET("/api/forum/:slug/threads", this.getForumThreads)
	this.router.GET("/api/forum/:slug/users", this.getForumUsers)
	// Thread
	this.router.GET("/api/thread/:slug_or_id/posts", this.getThreadPosts)
	this.router.GET("/api/thread/:slug_or_id/details", this.getThreadDetails)
	this.router.POST("/api/forum/:slug/create", this.createThread)
	this.router.POST("/api/thread/:slug_or_id/details", this.updateThread)
	// Post
	this.router.GET("/api/post/:id/details", this.getPostDetails)
	this.router.POST("/api/thread/:slug_or_id/create", this.createPost)
	this.router.POST("/api/post/:id/details", this.updatePost)
	// Vote
	this.router.POST("/api/thread/:slug_or_id/vote", this.voteForThread)
	// User
	this.router.GET("/api/user/:nickname/profile", this.getUser)
	this.router.POST("/api/user/:nickname/create", this.createUser)
	this.router.POST("/api/user/:nickname/profile", this.updateUser)
	// Service
	this.router.GET("/api/service/status", this.getStatus)
	this.router.POST("/api/service/clear", this.clear)
}

func (this *Handler) GetHandleFunc() fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		if string(c.Path()) == "/api/forum/create" {
			this.handleForumCreate(c)
		} else {
			this.router.Handler(c)
		}
	}
}
