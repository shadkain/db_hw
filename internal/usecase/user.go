package usecase

import (
	"db_hw/internal/models"
)

func (uc *usecaseImpl) AddUser(newUser models.NewUser, nickname string) (User models.User, Err error) {
	if err := uc.repo.InsertUser(newUser, nickname); err != nil {
		return models.User{}, err
	}

	user := models.User{
		About:    newUser.About,
		Email:    newUser.Email,
		Fullname: newUser.Fullname,
		Nickname: nickname,
	}

	return user, nil
}

func (uc *usecaseImpl) GetUsersByForum(slug, limit, since, desc string) (Users *models.Users, Err error) {
	users, err := uc.repo.SelectUsersByForum(slug, limit, since, desc)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (uc *usecaseImpl) GetUsersByNicknameOrEmail(email string, nickname string) (User []models.User, Err error) {
	users, err := uc.repo.SelectUsersByNicknameOrEmail(email, nickname)

	if err != nil {
		return users, err
	}

	return users, nil
}

func (uc *usecaseImpl) GetUserByNickname(nickname string) (user models.User, Err error) {
	user, err := uc.repo.SelectUserByNickname(nickname)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (uc *usecaseImpl) GetUsersByEmail(email string) (User []models.User, Err error) {
	users, err := uc.repo.SelectUsersByEmail(email)

	if err != nil {
		return users, err
	}

	return users, nil
}

func (uc *usecaseImpl) SetUser(newProfile models.NewUser, nickname string) (User models.User, Err error) {
	curentUser, err := uc.repo.SelectUserByNickname(nickname)
	if err != nil {
		return models.User{}, err
	}

	if newProfile.Email == "" {
		newProfile.Email = curentUser.Email
	}
	if newProfile.About == "" {
		newProfile.About = curentUser.About
	}
	if newProfile.Fullname == "" {
		newProfile.Fullname = curentUser.Fullname
	}

	if err := uc.repo.UpdateUser(newProfile, nickname); err != nil {
		return models.User{}, err
	}

	user := models.User{
		About:    newProfile.About,
		Email:    newProfile.Email,
		Fullname: newProfile.Fullname,
		Nickname: nickname,
	}

	return user, nil
}
