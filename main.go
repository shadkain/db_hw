package main

import (
	"fmt"
	"github.com/shadkain/db_hw/internal/delivery"
	"github.com/shadkain/db_hw/internal/storage"
	"github.com/shadkain/db_hw/internal/usecase"
	"github.com/valyala/fasthttp"
	"log"
	"os"
)

const PORT = 5000

func main() {
	storage := storage.NewStorage()
	if err := storage.Open(getDatabaseUrl("dev")); err != nil {
		log.Fatal(err)
	}

	usecase := usecase.NewUsecase(storage.Repository())
	handler := delivery.NewHandler(usecase)

	fmt.Printf("→ Started listening port: %d\n", PORT)
	log.Fatal(fasthttp.ListenAndServe(fmt.Sprintf(":%d", PORT), handler.GetHandleFunc()))
}

func getDatabaseUrl(mode string) string {
	switch mode {
	case "dev":
		return "postgres://jason:@localhost:2389/subd"
	case "prod":
		return os.Getenv("DATABASE_URL")
	default:
		panic("unknown mode")
	}
}
