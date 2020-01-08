package usecase

import (
	"db_hw/internal/repository"
)

type usecaseImpl struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecaseImpl{
		repo: repo,
	}
}

type Usecase interface {
}
