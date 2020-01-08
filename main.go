package main

import (
	"github.com/labstack/echo"
	"github.com/shadkain/db_hw/internal/delivery"
	"github.com/shadkain/db_hw/internal/storage"
	"github.com/shadkain/db_hw/internal/usecase"
)

func main() {
	e := echo.New()

	st := storage.NewStorage()
	if err := st.Open("postgresql://forum_user:forum_pass@localhost:5432/forum_db"); err != nil {
		return
	}

	uc := usecase.NewUsecase(st.Repository())
	handler := delivery.NewDelivery(uc)
	handler.Configure(e)

	if err := e.Start("localhost:5000"); err != nil {
		return
	}
}
