package reqmodels

type ForumCreate struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	User  string `json:"user"`
}

type ThreadCreate struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
}

type ThreadUpdate struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type PostCreate struct {
	Author  string `json:"author"`
	Message string `json:"message"`
	Parent  int    `json:"parent"`
}

type PostUpdate struct {
	Message string `json:"message"`
}

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}

type UserInput struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
}

type Status struct {
	Forum  int `json:"forum"`
	Post   int `json:"post"`
	Thread int `json:"thread"`
	User   int `json:"user"`
}
