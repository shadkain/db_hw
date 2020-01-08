package main

import (
	"db_hw/internal/delivery"
	"db_hw/internal/storage"
	"db_hw/internal/usecase"
	"github.com/labstack/echo"
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
