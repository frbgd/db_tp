package models

import "time"

//easyjson:json
type Thread struct {
	Id      int64     `json:"id,omitempty"`
	Slug    string    `json:"slug,omitempty"`
	Title   string    `json:"title,omitempty"`
	Message string    `json:"message,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Author  string    `json:"author,omitempty"`
	Forum   string    `json:"forum,omitempty"`
	Votes   int       `json:"votes,omitempty"`
}

//easyjson:json
type ThreadUpdate struct {
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}
