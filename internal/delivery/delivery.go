package delivery

import "github.com/labstack/echo"

type Delivery struct {
}

func NewDelivery() *Delivery {
	return &Delivery{}
}

func (d *Delivery) Configure(e *echo.Echo) {
}
