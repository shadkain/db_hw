package models

type NewUser struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
}

type User struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
}

type Users []*User
