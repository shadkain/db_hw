package repository

type repositoryImpl struct {
}

func NewRepository() Repository {
	return &repositoryImpl{}
}

type Repository interface {
}
