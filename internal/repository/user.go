package repository

import (
	"database/sql"
	"github.com/shadkain/db_hw/internal/models"
	"github.com/shadkain/db_hw/internal/vars"
)

func (this *repositoryImpl) GetUserByNickname(nickname string) (*models.User, error) {
	var user models.User
	if err := this.db.Get(
		&user,
		`SELECT * FROM "user" WHERE nickname=$1`,
		nickname,
	); err != nil {
		return nil, Error(err)
	}

	return &user, nil
}

func (this *repositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := this.db.Get(
		&user,
		`SELECT * FROM "user" WHERE email=$1`,
		email,
	); err != nil {
		return nil, Error(err)
	}

	return &user, nil
}

func (this *repositoryImpl) GetUserNickname(nickname string) (string, error) {
	if realNickname, err := this.userCache.GetNicknameCaseInsensitive(nickname); err == nil {
		return realNickname, nil
	}

	var user models.User
	if err := this.db.Get(
		&user,
		`SELECT id, nickname FROM "user" WHERE nickname=$1`,
		nickname,
	); err != nil {
		return "", Error(err)
	}

	this.userCache.Set(user.ID, user.Nickname)

	return user.Nickname, nil
}

func (this *repositoryImpl) GetUsersByNicknameOrEmail(nickname, email string) ([]*models.User, error) {
	var users []*models.User
	if err := this.db.Select(
		&users,
		`SELECT * FROM "user" WHERE nickname=$1 OR email=$2`,
		nickname, email,
	); err != nil {
		return nil, Error(err)
	}

	return users, nil
}

func (this *repositoryImpl) CreateUser(nickname, email, fullname, about string) (*models.User, error) {
	var id int
	if err := this.db.QueryRow(
		`INSERT INTO "user" (nickname, email, fullname, about) VALUES ($1, $2, $3, $4) RETURNING id`,
		nickname, email, fullname, about,
	).Scan(&id); err != nil {
		return nil, err
	}

	this.userCache.Set(id, nickname)

	return this.getUserByID(id)
}

func (this *repositoryImpl) UpdateUserByNickname(nickname, email, fullname, about string) error {
	if result, err := this.db.Exec(
		`UPDATE "user" SET email=$2, fullname=$3, about=$4 WHERE nickname=$1`,
		nickname, email, fullname, about,
	); err != nil {
		return Error(err)
	} else if rows, _ := result.RowsAffected(); rows == 0 {
		return vars.ErrNotFound
	}

	return nil
}

func (this *repositoryImpl) getUserByID(userID int) (*models.User, error) {
	var user models.User
	if err := this.db.Get(
		&user,
		`SELECT * FROM "user" WHERE id=$1`,
		userID,
	); err == sql.ErrNoRows {
		return nil, vars.ErrNotFound
	}

	return &user, nil
}
