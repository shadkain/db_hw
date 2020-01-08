package delivery

import (
	"db_hw/internal/models"
	"github.com/labstack/echo"
	"net/http"
)

func (d *Delivery) createVote(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")
	newVote := models.NewVote{}

	if err := c.Bind(&newVote); err != nil {
		return err
	}

	thread, err := d.uc.SetVote(newVote, slugOrID)
	if err != nil {
		if err.Error() == "Can't find thread" {
			if err := c.JSON(http.StatusNotFound, models.Error{"Can't find thread"}); err != nil {
				return err
			}
			return nil
		}
		if err := c.JSON(http.StatusNotFound, models.Error{"Can't find user"}); err != nil {
			return err
		}
		return nil
	}

	if err := c.JSON(http.StatusOK, thread); err != nil {
		return err
	}

	return nil
}