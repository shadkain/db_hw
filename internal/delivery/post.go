package delivery

import (
	"encoding/json"
	"github.com/shadkain/db_hw/internal/reqmodels"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

func (this *Handler) getPostDetails(c *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(PathParam(c, "id"))
	related := strings.Split(QueryParam(c, "related"), ",")
	details, err := this.uc.GetPostDetails(id, related)
	if err != nil {
		Error(c, err)
		return
	}

	result := map[string]interface{}{
		"post": details.Post,
	}

	for _, r := range related {
		switch r {
		case "user":
			result["author"] = details.Author
		case "forum":
			result["forum"] = details.Forum
		case "thread":
			result["thread"] = details.Thread
		}
	}

	Ok(c, result)
}

func (this *Handler) createPost(c *fasthttp.RequestCtx) {
	var posts []*reqmodels.PostCreate
	if err := json.Unmarshal(c.PostBody(), &posts); err != nil {
		BadRequest(c, err)
		return
	}

	if result, err := this.uc.CreatePosts(PathParam(c, "slug_or_id"), posts); err != nil {
		Error(c, err)
	} else {
		Created(c, result)
	}
}

func (this *Handler) updatePost(c *fasthttp.RequestCtx) {
	var pu reqmodels.PostUpdate
	if err := json.Unmarshal(c.PostBody(), &pu); err != nil {
		BadRequest(c, err)
		return
	}

	id, _ := strconv.Atoi(PathParam(c, "id"))

	if thread, err := this.uc.UpdatePost(id, pu.Message); err != nil {
		Error(c, err)
	} else {
		Ok(c, thread)
	}
}
