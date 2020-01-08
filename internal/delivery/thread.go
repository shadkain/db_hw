package delivery

import (
	"github.com/shadkain/db_hw/internal/models"
	"github.com/jackc/pgx"
	"github.com/labstack/echo"
	"net/http"
)

func (d *Delivery) createThread(c echo.Context) error {
	newThread := models.NewThread{}

	if err := c.Bind(&newThread); err != nil {
		return err
	}

	forum := c.Param("slug")

	users, err := d.uc.GetUsersByNicknameOrEmail("", newThread.Author)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		newThread.Author = users[0].Nickname
	} else {
		if err := c.JSON(http.StatusNotFound, models.Error{"Can't find user"}); err != nil {
			return err
		}
		return nil
	}
	forums, err := d.uc.GetForumsBySlug(forum)
	if err != nil {
		return err
	}
	if len(forums) > 0 {
		forum = forums[0].Slug
	} else {
		if err := c.JSON(http.StatusNotFound, models.Error{"Can't find forum"}); err != nil {
			return err
		}
		return nil
	}

	thread, err := d.uc.AddThread(newThread, forum)
	if err != nil {
		if err.Error() == "conflict" {
			if err := c.JSON(http.StatusConflict, thread); err != nil {
				return err
			}
			return nil
		}
		pqErr, ok := err.(pgx.PgError)
		if !ok {
			return err
		}
		if pqErr.Code == "23503" {
			if err := c.JSON(http.StatusNotFound, models.Error{"Can't find user"}); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if err := c.JSON(http.StatusCreated, thread); err != nil {
		return err
	}

	return nil
}

func (d *Delivery) takeThread(ctx echo.Context) error {
	slug_or_id := ctx.Param("slug_or_id")

	thread, err := d.uc.GetThreadBySlug(slug_or_id)

	if err != nil {
		if err := ctx.JSON(http.StatusNotFound, models.Error{"Can't thread"}); err != nil {
			return err
		}
		return nil
	}

	if err := ctx.JSON(http.StatusOK, thread); err != nil {
		return err
	}

	return nil
}

func (d *Delivery) takeForumThreads(ctx echo.Context) error {
	limit := ctx.QueryParam("limit")
	since := ctx.QueryParam("since")
	desc := ctx.QueryParam("desc")

	if limit == "" {
		limit = "100"
	}
	if desc == "" {
		desc = "false"
	}

	forums, err := d.uc.GetForumsBySlug(ctx.Param("slug"))

	if err != nil || len(forums) == 0 {
		if err := ctx.JSON(http.StatusNotFound, models.Error{"Can't find forum by slug"}); err != nil {
			return err
		}
		return nil
	}

	threads, err := d.uc.GetThreadsByForum(ctx.Param("slug"), limit, since, desc)
	if err != nil {
		return err
	}

	if err := ctx.JSON(http.StatusOK, threads); err != nil {
		return err
	}

	return nil
}

func (d *Delivery) changeThread(c echo.Context) error {
	changeThread := models.ChangeThread{}

	if err := c.Bind(&changeThread); err != nil {
		return err
	}

	slugOrID := c.Param("slug_or_id")

	thread, err := d.uc.SetThread(changeThread, slugOrID)
	if err != nil {
		if err := c.JSON(http.StatusNotFound, models.Error{"Can't find thread"}); err != nil {
			return err
		}
		return nil
	}

	if err := c.JSON(http.StatusOK, thread); err != nil {
		return err
	}

	return nil
}