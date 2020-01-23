package models

type Forum struct {
	ID      int    `db:"id" json:"-"`
	Title   string `db:"title" json:"title"`
	User    string `db:"user" json:"user"`
	Slug    string `db:"slug" json:"slug"`
	Posts   int    `db:"posts" json:"posts"`
	Threads int    `db:"threads" json:"threads"`
}

type Thread struct {
	ID      int    `db:"id" json:"id"`
	Title   string `db:"title" json:"title"`
	Author  string `db:"author" json:"author"`
	Forum   string `db:"forum" json:"forum"`
	Message string `db:"message" json:"message"`
	Votes   int    `db:"votes" json:"votes"`
	Slug    string `db:"slug" json:"slug"`
	Created string `db:"created" json:"created"`
}

type Post struct {
	ID       int    `db:"id" json:"id"`
	Parent   int    `db:"parent" json:"parent"`
	Path     string `db:"path" json:"-"`
	Author   string `db:"author" json:"author"`
	Forum    string `db:"forum" json:"forum"`
	Thread   int    `db:"thread" json:"thread"`
	Message  string `db:"message" json:"message"`
	IsEdited bool   `db:"isEdited" json:"isEdited"`
	Created  string `db:"created" json:"created"`
}

type PostDetails struct {
	Post   *Post
	Author *User
	Forum  *Forum
	Thread *Thread
}

type Vote struct {
	ID       int    `db:"id" json:"id"`
	Thread   int    `db:"thread" json:"thread"`
	Nickname string `db:"nickname" json:"nickname"`
	Voice    int    `db:"voice" json:"voice"`
}

type User struct {
	ID       int    `db:"id" json:"-"`
	Nickname string `db:"nickname" json:"nickname"`
	Fullname string `db:"fullname" json:"fullname"`
	About    string `db:"about" json:"about"`
	Email    string `db:"email" json:"email"`
}

type Users = []*User
type Threads = []*Thread
type Posts = []*Post
