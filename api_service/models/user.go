package models

//easyjson:json
type User struct {
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	About    string `json:"about,omitempty"`
}

//easyjson:json
type UserUpdate struct {
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	About    string `json:"about,omitempty"`
}
