package repository

import (
	"bytes"
	"github.com/shadkain/db_hw/internal/models"
	"github.com/shadkain/db_hw/internal/query"
	"errors"
	"github.com/jackc/pgx"
	xsort "sort"
	"strconv"
	"strings"
)

func (r *repositoryImpl) InsertUser(newUser models.NewUser, nickname string) (Err error) {
	var lastID int
	err := r.db.QueryRow(query.InsertUser, newUser.About, newUser.Email,
		newUser.Fullname, nickname).Scan(&lastID)

	if err != nil {
		return err
	}
	return nil
}

func (r *repositoryImpl) SelectUsersByForum(slug, limit, since, desc string) (Users *models.Users, Err error) {
	users := models.Users{}
	var rows *pgx.Rows
	var err error

	if desc == "false" {
		rows, err = r.db.Query(query.SelectUsersByForumSlug, slug)
	} else {
		rows, err = r.db.Query(query.SelectUsersByForumSlugDesc, slug)
	}

	defer rows.Close()
	if err != nil {
		return &users, err
	}

	for rows.Next() {
		scanUser := models.User{}
		err := rows.Scan(&scanUser.About, &scanUser.Email, &scanUser.Fullname,
			&scanUser.Nickname)
		if err != nil {
			return &users, err
		}
		users = append(users, &scanUser)
	}

	resUsers := models.Users{}

	limitDigit, _ := strconv.Atoi(limit)

	if desc == "false" {

		xsort.Slice(users, func(i, j int) bool {
			return bytes.Compare([]byte(strings.ToLower(users[i].Nickname)), []byte(strings.ToLower(users[j].Nickname))) < 0
		})

		if since == "" {
			for i := 0; i < limitDigit && i < len(users); i++ {
				resUsers = append(resUsers, users[i])
			}
		} else {
			j := 0
			for i := 0; j < limitDigit && i < len(users); {
				if bytes.Compare([]byte(strings.ToLower(users[i].Nickname)), []byte(strings.ToLower(since))) > 0 {
					resUsers = append(resUsers, users[i])
					j++
				}
				i++
			}
		}
	} else {

		xsort.Slice(users, func(i, j int) bool {
			return bytes.Compare([]byte(strings.ToLower(users[i].Nickname)), []byte(strings.ToLower(users[j].Nickname))) > 0
		})

		if since == "" {
			for i := 0; i < limitDigit && i < len(users); i++ {
				resUsers = append(resUsers, users[i])
			}
		} else {
			j := 0
			for i := 0; j < limitDigit && i < len(users); {
				if bytes.Compare([]byte(strings.ToLower(users[i].Nickname)), []byte(strings.ToLower(since))) < 0 {
					resUsers = append(resUsers, users[i])
					j++
				}
				i++
			}
		}
	}

	return &resUsers, nil
}

func (r *repositoryImpl) SelectUsersByNicknameOrEmail(email string, nickname string) (Users []models.User, Err error) {
	var users []models.User
	rows, err := r.db.Query(query.SelectUsersByNicknameOrEmail, email, nickname)
	defer rows.Close()
	if err != nil {
		return users, err
	}

	scanUser := models.User{}
	for rows.Next() {
		err := rows.Scan(&scanUser.About, &scanUser.Email, &scanUser.Fullname,
			&scanUser.Nickname)
		if err != nil {
			return users, err
		}
		users = append(users, scanUser)
	}
	return users, nil
}

func (r *repositoryImpl) SelectUserByNickname(nickname string) (user models.User, Err error) {
	var users []models.User
	rows, err := r.db.Query(query.SelectUsersByNickname, nickname)
	defer rows.Close()
	if err != nil {
		return models.User{}, err
	}

	scanUser := models.User{}
	for rows.Next() {
		err := rows.Scan(&scanUser.About, &scanUser.Email, &scanUser.Fullname,
			&scanUser.Nickname)
		if err != nil {
			return models.User{}, err
		}
		users = append(users, scanUser)
	}

	if len(users) == 0 {
		return models.User{}, errors.New("Can't find user by nickname")
	}
	return users[0], nil
}

func (r *repositoryImpl) SelectUsersByEmail(email string) (Users []models.User, Err error) {
	var users []models.User
	rows, err := r.db.Query(query.SelectUsersByEmail, email)
	defer rows.Close()
	if err != nil {
		return users, err
	}

	scanUser := models.User{}
	for rows.Next() {
		err := rows.Scan(&scanUser.About, &scanUser.Email, &scanUser.Fullname,
			&scanUser.Nickname)
		if err != nil {
			return users, err
		}
		users = append(users, scanUser)
	}
	return users, nil
}

func (r *repositoryImpl) UpdateUser(newProfile models.NewUser, nickname string) (Err error) {
	rows, err := r.db.Query(query.UpdateUserByNickname, newProfile.About, newProfile.Email,
		newProfile.Fullname, nickname)
	defer rows.Close()

	if err != nil {
		return err
	}
	return nil
}
