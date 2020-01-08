package models

type NewForum struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	User  string `json:"user"`
}

type Forum struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Posts   int    `json:"posts"`
	Threads int    `json:"threads"`
	User    string `json:"user"`
}
