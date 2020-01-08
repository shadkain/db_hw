package main

import (
	"github.com/shadkain/db_hw/internal/delivery"
	"github.com/shadkain/db_hw/internal/storage"
	"github.com/shadkain/db_hw/internal/usecase"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	st := storage.NewStorage()
	if err := st.Open("postgresql://jason:12345678@localhost:2389/dd"); err != nil {
		return
	}

	uc := usecase.NewUsecase(st.Repository())
	handler := delivery.NewDelivery(uc)
	handler.Configure(e)

	if err := e.Start("localhost:5000"); err != nil {
		return
	}
}
