package models

type Status struct {
	Post   int `json:"post"`
	Thread int `json:"thread"`
	User   int `json:"user"`
	Forum  int `json:"forum"`
}

type Error struct {
	Message string `json:"message"`
}

type Body struct {
	Body interface{}
}
