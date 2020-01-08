package models

type NewPost struct {
	Author  string `json:"author"`
	Message string `json:"message"`
	Parent  int    `json:"parent"`
}

type ChangePost struct {
	Message string `json:"message"`
}

type Post struct {
	Author   string `json:"author"`
	Created  string `json:"created"`
	Forum    string `json:"forum"`
	ID       int    `json:"id"`
	IsEdited bool   `json:"isEdited"`
	Message  string `json:"message"`
	Parent   int    `json:"parent"`
	Thread   int    `json:"thread"`
}

type EditedPost struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	ID      int    `json:"id"`
	Message string `json:"message"`
	Parent  int    `json:"parent"`
	Thread  int    `json:"thread"`
}

type PostDetails struct {
	Forum  interface{} `json:"forum,omitempty"`
	Thread interface{} `json:"thread,omitempty"`
	User   interface{} `json:"author,omitempty"`
	Post   interface{} `json:"post"`
}

type NewPosts []*NewPost

type Posts []*Post
