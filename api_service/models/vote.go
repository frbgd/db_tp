package models

//easyjson:json
type Vote struct {
	Nickname string `json:"nickname,omitempty"`
	Voice    int64  `json:"voice,omitempty"`
}
