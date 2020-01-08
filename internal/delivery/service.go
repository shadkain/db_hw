package delivery

import (
	"github.com/labstack/echo"
	"net/http"
)

func (d *Delivery) takeStatus(ctx echo.Context) error {
	status, err := d.uc.GetStatus()
	if err != nil {
		return err
	}

	if err := ctx.JSON(http.StatusOK, status); err != nil {
		return err
	}

	return nil
}

func (d *Delivery) clearAll(c echo.Context) error {
	err := d.uc.ClearAll()
	if err != nil {
		return err
	}

	if err := c.JSON(http.StatusOK, nil); err != nil {
		return err
	}

	return nil
}