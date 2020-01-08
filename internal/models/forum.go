package models

type NewForum struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	User  string `json:"user"`
}

type Forum struct {
	Posts  int    `json:"posts"`
	Slug   string `json:"slug"`
	Thread int    `json:"threads"`
	Title  string `json:"title"`
	User   string `json:"user"`
}
