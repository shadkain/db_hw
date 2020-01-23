package delivery

import (
	"encoding/json"
	"errors"
	"github.com/shadkain/db_hw/internal/reqmodels"
	"github.com/shadkain/db_hw/internal/vars"
	"github.com/valyala/fasthttp"
)

func (this *Handler) createUser(c *fasthttp.RequestCtx) {
	var ui reqmodels.UserInput
	if err := json.Unmarshal(c.PostBody(), &ui); err != nil {
		BadRequest(c, err)
		return
	}

	users, err := this.uc.CreateUser(PathParam(c, "nickname"), ui.Email, ui.Fullname, ui.About)
	if errors.Is(err, vars.ErrConflict) {
		Conflict(c, users)
		return
	}

	if err != nil {
		Error(c, err)
	} else {
		Created(c, users[0])
	}
}

func (this *Handler) getUser(c *fasthttp.RequestCtx) {
	if user, err := this.uc.GetUserByNickname(PathParam(c, "nickname")); err != nil {
		Error(c, err)
	} else {
		Ok(c, user)
	}
}

func (this *Handler) updateUser(c *fasthttp.RequestCtx) {
	var ui reqmodels.UserInput
	if err := json.Unmarshal(c.PostBody(), &ui); err != nil {
		BadRequest(c, err)
		return
	}

	user, err := this.uc.UpdateUser(PathParam(c, "nickname"), ui.Email, ui.Fullname, ui.About)
	if errors.Is(err, vars.ErrConflict) {
		ConflictWithMessage(c, err)
		return
	}

	if err != nil {
		Error(c, err)
	} else {
		Ok(c, user)
	}
}
