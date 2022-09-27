package models

type GetShowRequest struct {
	Date string `json:"Date"`
}

type AddUserShowRequest struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Date     string `json:"date,omitempty"`
}
