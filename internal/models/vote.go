package models

type NewVote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}
