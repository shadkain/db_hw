package delivery

import (
	"db_hw/internal/models"
	"github.com/labstack/echo"
)

func (d *Delivery) createUser(c echo.Context) error {
	newUser := models.NewUser{}

	if err := c.Bind(&newUser); err != nil {
		return err
	}

	users, err := d.uc.GetUsersByNicknameOrEmail(newUser.Email, c.Param("nickname"))
	if err != nil {
		return err
	}
	if len(users) > 0 {
		if err := c.JSON(409, users); err != nil {
			return err
		}
		return nil
	}

	user, err := d.uc.AddUser(newUser, c.Param("nickname"))
	if err != nil {
		return err
	}

	if err := c.JSON(201, user); err != nil {
		return err
	}

	return nil
}

func (d *Delivery) takeUser(ctx echo.Context) error {
	user, err := d.uc.GetUserByNickname(ctx.Param("nickname"))

	if err != nil {
		if err.Error() == "Can't find user by nickname" {
			if err := ctx.JSON(404, models.Error{err.Error()}); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if err := ctx.JSON(200, user); err != nil {
		return err
	}

	return nil
}

func (d *Delivery) takeUsersByForum(ctx echo.Context) error {
	slug := ctx.Param("slug")

	forums, err := d.uc.GetForumsBySlug(slug)
	if len(forums) != 1 || err != nil {
		if err := ctx.JSON(404, models.Error{"Can't find forum by slug"}); err != nil {
			return err
		}
		return nil
	}

	limit := ctx.QueryParam("limit")
	since := ctx.QueryParam("since")
	desc := ctx.QueryParam("desc")

	if limit == "" {
		limit = "100"
	}
	if desc == "" {
		desc = "false"
	}

	users, err := d.uc.GetUsersByForum(slug, limit, since, desc)
	if err != nil {
		return err
	}

	if err := ctx.JSON(200, users); err != nil {
		return err
	}

	return nil
}

func (d *Delivery) changeUser(c echo.Context) (Err error) {
	newProfile := models.NewUser{}

	if err := c.Bind(&newProfile); err != nil {
		return err
	}

	users, err := d.uc.GetUsersByEmail(newProfile.Email)
	if err != nil {
		return err
	}
	if len(users) > 0 && !(len(users) == 1 &&
		users[0].Nickname == c.Param("nickname")) {
		if err := c.JSON(409, models.Error{"Conflict"}); err != nil {
			return err
		}

		return nil
	}

	user, err := d.uc.SetUser(newProfile, c.Param("nickname"))
	if err != nil {
		if err.Error() == "Can't find user by nickname" {
			if err := c.JSON(404, models.Error{err.Error()}); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if err := c.JSON(200, user); err != nil {
		return err
	}

	return nil
}
