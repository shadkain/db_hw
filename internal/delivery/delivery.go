package delivery

import (
	"db_hw/internal/usecase"
	"github.com/labstack/echo"
)

type Delivery struct {
	uc usecase.Usecase
}

func NewDelivery(uc usecase.Usecase) *Delivery {
	return &Delivery{
		uc: uc,
	}
}

func (d *Delivery) Configure(e *echo.Echo) {
}
