package models

import "time"

//easyjson:json
type ErrMsg struct {
	Message string `json:"message,omitempty"`
}

//easyjson:json
type Forum struct {
	Slug    string `json:"slug,omitempty"`
	Title   string `json:"title,omitempty"`
	User    string `json:"user,omitempty"`
	Posts   int    `json:"posts,omitempty"`
	Threads int    `json:"threads,omitempty"`
}

//easyjson:json
type Post struct {
	Id       int       `json:"id,omitempty"`
	Message  string    `json:"message,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Author   string    `json:"author,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Parent   int       `json:"parent,omitempty"`
	Thread   int       `json:"thread,omitempty"`
}

//easyjson:json
type Posts []Post

//easyjson:json
type FullPost struct {
	Post   *Post   `json:"post,omitempty"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}

//easyjson:json
type PostUpdate struct {
	Message string `json:"message,omitempty"`
}

//easyjson:json
type Status struct {
	Forum  int `json:"forum,omitempty"`
	Post   int `json:"post,omitempty"`
	Thread int `json:"thread,omitempty"`
	User   int `json:"user,omitempty"`
}

//easyjson:json
type Thread struct {
	Id      int       `json:"id,omitempty"`
	Slug    string    `json:"slug,omitempty"`
	Title   string    `json:"title,omitempty"`
	Message string    `json:"message,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Author  string    `json:"author,omitempty"`
	Forum   string    `json:"forum,omitempty"`
	Votes   int       `json:"votes,omitempty"`
}

//easyjson:json
type Threads []Thread

//easyjson:json
type ThreadUpdate struct {
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

//easyjson:json
type User struct {
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	About    string `json:"about,omitempty"`
}

//easyjson:json
type Users []User

//easyjson:json
type UserUpdate struct {
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	About    string `json:"about,omitempty"`
}

//easyjson:json
type Vote struct {
	Nickname string `json:"nickname,omitempty"`
	Voice    int    `json:"voice,omitempty"`
}
