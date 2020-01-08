package delivery

import (
	"db_hw/internal/models"
	"github.com/jackc/pgx"
	"github.com/labstack/echo"
	"strconv"
)

func (d *Delivery) createPosts(c echo.Context) (Err error) {
	newPosts := models.NewPosts{}

	if err := c.Bind(&newPosts); err != nil {
		return err
	}

	posts, err := d.uc.AddPosts(newPosts, c.Param("slug_or_id"))
	if err != nil {
		_, ok := err.(pgx.PgError)
		if !ok {
			if err.Error() == "Can't find thread" {
				if err := c.JSON(404, models.Error{"Can't find thread"}); err != nil {
					return err
				}
				return nil
			}
			if err := c.JSON(409, models.Error{err.Error()}); err != nil {
				return err
			}
			return nil
		}
		if err := c.JSON(404, models.Error{"Can't find user"}); err != nil {
			return err
		}
		return nil
	}

	if err := c.JSON(201, posts); err != nil {
		return err
	}

	return nil
}

func (d *Delivery) takePostByID(ctx echo.Context) error {
	ID := ctx.Param("id")
	postID, err := strconv.Atoi(ID)
	if err != nil {
		return nil
	}

	related := ctx.QueryParam("related")
	println(related)

	postDetails, err := d.uc.GetPostByID(postID, related)
	if err != nil {
		if err := ctx.JSON(404, models.Error{"Can't find post"}); err != nil {
			return err
		}
		return nil
	}

	if err := ctx.JSON(200, postDetails); err != nil {
		return err
	}

	return nil
}

func (d *Delivery) takePosts(ctx echo.Context) error {
	slugOrID := ctx.Param("slug_or_id")

	limit := ctx.QueryParam("limit")
	since := ctx.QueryParam("since")
	sort := ctx.QueryParam("sort")
	desc := ctx.QueryParam("desc")

	if limit == "" {
		limit = "100"
	}

	if sort == "" {
		sort = "flat"
	}
	if desc == "" {
		desc = "false"
	}
	if since == "" {
		if desc == "false" {
			since = "0"
		} else {
			since = "999999999"
		}
	}

	posts, err := d.uc.GetPosts(slugOrID, limit, since, sort, desc)
	if err != nil {
		if err := ctx.JSON(404, models.Error{"Can't find thread"}); err != nil {
			return err
		}
		return nil
	}

	if err := ctx.JSON(200, posts); err != nil {
		return err
	}

	return nil
}

func (d *Delivery) changePost(c echo.Context) (Err error) {
	changePost := models.ChangePost{}

	if err := c.Bind(&changePost); err != nil {
		return err
	}

	ID := c.Param("id")
	postID, err := strconv.Atoi(ID)
	if err != nil {
		return err
	}

	post, err := d.uc.SetPost(changePost, postID)
	if err != nil {
		if err := c.JSON(404, models.Error{"Can't find post"}); err != nil {
			return err
		}
		return nil
	}

	if err := c.JSON(200, post); err != nil {
		return err
	}

	return nil
}
