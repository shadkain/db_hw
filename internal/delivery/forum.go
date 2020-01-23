package delivery

import (
	"encoding/json"
	"errors"
	"github.com/shadkain/db_hw/internal/reqmodels"
	"github.com/shadkain/db_hw/internal/vars"
	"github.com/valyala/fasthttp"
	"strconv"
)

func (this *Handler) getForumDetails(c *fasthttp.RequestCtx) {
	if forum, err := this.uc.GetForum(PathParam(c, "slug")); err != nil {
		Error(c, err)
	} else {
		Ok(c, forum)
	}
}

func (this *Handler) getForumThreads(c *fasthttp.RequestCtx) {
	limit, _ := strconv.Atoi(QueryParam(c, "limit"))
	desc, _ := strconv.ParseBool(QueryParam(c, "desc"))

	if threads, err := this.uc.GetForumThreads(PathParam(c, "slug"), QueryParam(c, "since"), limit, desc); err != nil {
		Error(c, err)
	} else {
		Ok(c, threads)
	}
}

func (this *Handler) getForumUsers(c *fasthttp.RequestCtx) {
	limit, _ := strconv.Atoi(QueryParam(c, "limit"))
	desc, _ := strconv.ParseBool(QueryParam(c, "desc"))

	if users, err := this.uc.GetForumUsers(PathParam(c, "slug"), QueryParam(c, "since"), limit, desc); err != nil {
		Error(c, err)
	} else {
		Ok(c, users)
	}
}

func (this *Handler) handleForumCreate(c *fasthttp.RequestCtx) {
	var fc reqmodels.ForumCreate
	if err := json.Unmarshal(c.PostBody(), &fc); err != nil {
		BadRequest(c, err)
		return
	}

	forum, err := this.uc.CreateForum(fc.Title, fc.Slug, fc.User)
	if errors.Is(err, vars.ErrConflict) {
		Conflict(c, forum)
		return
	}

	if err != nil {
		Error(c, err)
	} else {
		Created(c, forum)
	}
}
