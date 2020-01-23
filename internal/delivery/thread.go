package delivery

import (
	"encoding/json"
	"errors"
	"github.com/shadkain/db_hw/internal/reqmodels"
	"github.com/shadkain/db_hw/internal/vars"
	"github.com/valyala/fasthttp"
	"strconv"
)

func (this *Handler) createThread(c *fasthttp.RequestCtx) {
	var tc reqmodels.ThreadCreate
	if err := json.Unmarshal(c.PostBody(), &tc); err != nil {
		BadRequest(c, err)
		return
	}

	thread, err := this.uc.CreateThread(PathParam(c, "slug"), tc)
	if errors.Is(err, vars.ErrConflict) {
		Conflict(c, thread)
		return
	}

	if err != nil {
		Error(c, err)
	} else {
		Created(c, thread)
	}
}

func (this *Handler) getThreadDetails(c *fasthttp.RequestCtx) {
	if thread, err := this.uc.GetThread(PathParam(c, "slug_or_id")); err != nil {
		Error(c, err)
	} else {
		Ok(c, thread)
	}
}

func (this *Handler) updateThread(c *fasthttp.RequestCtx) {
	var tu reqmodels.ThreadUpdate
	if err := json.Unmarshal(c.PostBody(), &tu); err != nil {
		BadRequest(c, err)
		return
	}

	if thread, err := this.uc.UpdateThread(PathParam(c, "slug_or_id"), tu.Message, tu.Title); err != nil {
		Error(c, err)
	} else {
		Ok(c, thread)
	}
}

func (this *Handler) getThreadPosts(c *fasthttp.RequestCtx) {
	sinceParam := QueryParam(c, "since")
	var since *int = nil
	if sinceParam != "" {
		n, _ := strconv.Atoi(sinceParam)
		since = &n
	}

	limit, _ := strconv.Atoi(QueryParam(c, "limit"))
	desc, _ := strconv.ParseBool(QueryParam(c, "desc"))
	if posts, err := this.uc.GetThreadPosts(
		PathParam(c, "slug_or_id"),
		limit,
		since,
		QueryParam(c, "sort"),
		desc,
	); err != nil {
		Error(c, err)
	} else {
		Ok(c, posts)
	}
}
