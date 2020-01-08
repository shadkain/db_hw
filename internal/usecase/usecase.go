package usecase

type usecaseImpl struct {
}

func NewUsecase() Usecase {
	return &usecaseImpl{}
}

type Usecase interface {
}
