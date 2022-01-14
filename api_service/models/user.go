package models

//easyjson:json
type UserUpdate struct {
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	About    string `json:"about,omitempty"`
}
