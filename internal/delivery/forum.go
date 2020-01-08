package delivery

import (
	"db_hw/internal/models"
	"github.com/labstack/echo"
	"net/http"
)

func (d *Delivery) createForum(c echo.Context) error {
	newForum := models.NewForum{}
	if err := c.Bind(&newForum); err != nil {
		return err
	}

	forums, err := d.uc.GetForumsBySlug(newForum.Slug)
	if err != nil {
		return err
	}
	if len(forums) > 0 {
		if err := c.JSON(http.StatusConflict, forums[0]); err != nil {
			return err
		}
		return nil
	}

	forum, err := d.uc.AddForum(newForum)
	if err != nil {
		if err.Error() == "Can't find user by nickname" {
			if err := c.JSON(http.StatusNotFound, models.Error{"Can't find user"}); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if err := c.JSON(http.StatusCreated, forum); err != nil {
		return err
	}

	return nil
}

func (d *Delivery) takeForum(ctx echo.Context) error {
	forums, err := d.uc.GetForumsBySlug(ctx.Param("slug"))

	if err != nil || len(forums) == 0 {
		if err := ctx.JSON(http.StatusNotFound, models.Error{"Can't find forum by slug"}); err != nil {
			return err
		}
		return nil
	}

	if err := ctx.JSON(http.StatusOK, forums[0]); err != nil {
		return err
	}

	return nil
}
