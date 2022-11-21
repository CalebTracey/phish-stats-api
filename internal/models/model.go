package models

import (
	"time"
)

type User struct {
	ID           string    `bson:"id,omitempty" json:"id,omitempty"`
	FullName     string    `bson:"fullName,omitempty" json:"fullName,omitempty"`
	Email        string    `bson:"email,omitempty"  json:"email,omitempty"`
	Username     string    `bson:"username,omitempty"  json:"username,omitempty"`
	Password     string    `bson:"password,omitempty"  json:"password,omitempty"`
	Token        string    `bson:"token,omitempty"  json:"token,omitempty"`
	RefreshToken string    `bson:"refreshToken,omitempty"  json:"refreshToken,omitempty"`
	CreatedAt    time.Time `bson:"createdAt"  json:"createdAt,omitempty"`
	UpdatedAt    time.Time `bson:"updatedAt"  json:"updatedAt,omitempty"`
	Shows        []string  `json:"shows"`
}

type ShowDate struct {
	Date uint8
}

type Show struct {
	Date  string
	Venue string
	Songs []Song
}

type Song struct {
	SongID    string `json:"songid,omitempty"`
	Title     string `json:"title,omitempty"`
	TrackTime string `json:"tracktime,omitempty"`
	Gap       string `json:"gap,omitempty"`
}
