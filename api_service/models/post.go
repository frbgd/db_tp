package models

import "time"

type Post struct {
	ID       int64     `json:"id,omitempty"`
	Message  string    `json:"message,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Author   string    `json:"author,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Parent   int64     `json:"parent,omitempty"`
	Thread   int64     `json:"thread,omitempty"`
}

type PostUpdate struct {
	Message string `json:"message,omitempty"`
}

type FullPost struct {
	Post   Post   `json:"post,omitempty"`
	Author User   `json:"author,omitempty"`
	Thread Thread `json:"thread,omitempty"`
	Forum  Forum  `json:"forum,omitempty"`
}
