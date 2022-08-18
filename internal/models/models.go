package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FullName     *string            `bson:"fullName,omitempty" json:"fullName,omitempty"`
	Email        *string            `bson:"email,omitempty"  json:"email,omitempty"`
	Username     *string            `bson:"username,omitempty"  json:"username,omitempty"`
	Password     *string            `bson:"password,omitempty"  json:"password,omitempty"`
	Token        *string            `bson:"token,omitempty"  json:"token,omitempty"`
	RefreshToken *string            `bson:"refreshToken,omitempty"  json:"refreshToken,omitempty"`
	CreatedAt    time.Time          `bson:"createdAt"  json:"createdAt,omitempty"`
	UpdatedAt    time.Time          `bson:"updatedAt"  json:"updatedAt,omitempty"`
	UserId       string             `bson:"userId,omitempty"  json:"userId,omitempty"`
}

type LoginUser struct {
	Username     *string `bson:"username,omitempty"  json:"username,omitempty"`
	Email        *string `bson:"email,omitempty"  json:"email,omitempty"`
	Token        *string `bson:"token,omitempty"  json:"token,omitempty"`
	RefreshToken *string `bson:"refreshToken,omitempty"  json:"refreshToken,omitempty"`
	UserId       string  `bson:"userId,omitempty"  json:"userId,omitempty"`
}
