package usecase

import (
	"fmt"
	"github.com/shadkain/db_hw/internal/vars"
	"github.com/shadkain/db_hw/internal/models"
)

func (this *usecaseImpl) GetUserByNickname(nickname string) (*models.User, error) {
	return this.repo.GetUserByNickname(nickname)
}

func (this *usecaseImpl) CreateUser(nickname, email, fullname, about string) ([]*models.User, error) {
	if exUsers, err := this.repo.GetUsersByNicknameOrEmail(nickname, email); err != nil && err != vars.ErrNotFound {
		return nil, err
	} else if exUsers != nil {
		return exUsers, vars.ErrConflict
	}

	user, err := this.repo.CreateUser(nickname, email, fullname, about)

	return []*models.User{user}, err
}

func (this *usecaseImpl) UpdateUser(nickname, email, fullname, about string) (*models.User, error) {
	exUser, err := this.repo.GetUserByNickname(nickname)
	if err != nil {
		return nil, err
	}

	if email == "" {
		email = exUser.Email
	} else {
		if userByEmail, _ := this.repo.GetUserByEmail(email); userByEmail != nil && userByEmail.Nickname != nickname {
			return nil, fmt.Errorf("%w: user with this email already exists", vars.ErrConflict)
		}
	}
	if fullname == "" {
		fullname = exUser.Fullname
	}
	if about == "" {
		about = exUser.About
	}

	if err := this.repo.UpdateUserByNickname(nickname, email, fullname, about); err != nil {
		return nil, err
	}

	return this.repo.GetUserByNickname(nickname)
}
