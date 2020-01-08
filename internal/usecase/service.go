package usecase

import (
	"db_hw/internal/models"
)

func (uc *usecaseImpl) GetStatus() (Status models.Status, Err error) {
	status, err := uc.repo.SelectStatus()

	if err != nil {
		return models.Status{}, err
	}

	return status, nil
}

func (uc *usecaseImpl) ClearAll() (Err error) {
	if err := uc.repo.ClearAll(); err != nil {
		return err
	}

	return nil
}
