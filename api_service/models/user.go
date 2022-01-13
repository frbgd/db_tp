package models

type User struct {
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	About    string `json:"about,omitempty"`
}

type UserUpdate struct {
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	About    string `json:"about,omitempty"`
}
