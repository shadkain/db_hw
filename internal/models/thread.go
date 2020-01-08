package models

type NewThread struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Created string `json:"created"`
	Author  string `json:"author"`
}

type ChangeThread struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type Thread struct {
	ID      int    `json:"id"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Created string `json:"created"`
	Votes   int    `json:"votes"`
	Forum   string `json:"forum"`
	Author  string `json:"author"`
}

type Threads []*Thread
