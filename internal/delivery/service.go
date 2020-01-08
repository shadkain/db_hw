package delivery

import (
	"github.com/labstack/echo"
)

func (d *Delivery) takeStatus(ctx echo.Context) error {
	status, err := d.uc.GetStatus()
	if err != nil {
		return err
	}

	if err := ctx.JSON(200, status); err != nil {
		return err
	}

	return nil
}

func (d *Delivery) clearAll(c echo.Context) (Err error) {
	err := d.uc.ClearAll()
	if err != nil {
		return err
	}

	if err := c.JSON(200, nil); err != nil {
		return err
	}

	return nil
}
