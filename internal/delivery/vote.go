package delivery

import (
	"encoding/json"
	"github.com/shadkain/db_hw/internal/reqmodels"
	"github.com/valyala/fasthttp"
)

func (this *Handler) voteForThread(c *fasthttp.RequestCtx) {
	var vote reqmodels.Vote
	if err := json.Unmarshal(c.PostBody(), &vote); err != nil {
		BadRequest(c, err)
		return
	}

	if thread, err := this.uc.VoteForThread(PathParam(c, "slug_or_id"), vote); err != nil {
		Error(c, err)
	} else {
		Ok(c, thread)
	}
}
