package models

type GetShowRequest struct {
	Date string `json:"Date"`
}

type AddUserShowRequest struct {
	Id   string `json:"id,omitempty"`
	Date string `json:"date,omitempty"`
}
