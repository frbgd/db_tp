package models

//easyjson:json
type Status struct {
	Forum  int64 `json:"forum,omitempty"`
	Post   int64 `json:"post,omitempty"`
	Thread int64 `json:"thread,omitempty"`
	User   int64 `json:"user,omitempty"`
}
