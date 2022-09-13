package models

type AuthResponse struct {
	FullName string `bson:"fullName,omitempty" json:"fullName,omitempty"`
	Email    string `bson:"email,omitempty"  json:"email,omitempty"`
	UserId   string `bson:"userId,omitempty"  json:"userId,omitempty"`
	Message  Message
}

type GetShowResponse struct {
	Show    Show
	Message Message
}

type UserResponse struct {
	User    *UserPsqlResponse
	Message Message
}

type AllUsersResponse struct {
	Users   []*User
	Message Message
}

type LoginResponse struct {
	Username     string  `bson:"username,omitempty"  json:"username,omitempty"`
	Email        string  `bson:"email,omitempty"  json:"email,omitempty"`
	Token        string  `bson:"token,omitempty"  json:"token,omitempty"`
	RefreshToken string  `bson:"refreshToken,omitempty"  json:"refreshToken,omitempty"`
	UserId       string  `bson:"userId,omitempty"  json:"userId,omitempty"`
	Message      Message `json:"message,omitempty"`
}

type RegistrationResponse struct {
	AccessToken string `json:"accessToken,omitempty"`
	Email       string `json:"email,omitempty"`
	Id          string `json:"id,omitempty"`
	Message     Message
}

type UserPsqlResponse struct {
	ID           string `bson:"id,omitempty" json:"id,omitempty"`
	FullName     string `bson:"fullName,omitempty" json:"fullName,omitempty"`
	Email        string `bson:"email,omitempty"  json:"email,omitempty"`
	Username     string `bson:"username,omitempty"  json:"username,omitempty"`
	Password     string `bson:"password,omitempty"  json:"password,omitempty"`
	Token        string `bson:"token,omitempty"  json:"token,omitempty"`
	RefreshToken string `bson:"refreshToken,omitempty"  json:"refreshToken,omitempty"`
	CreatedAt    string `bson:"createdAt"  json:"createdAt,omitempty"`
	UpdatedAt    string `bson:"updatedAt"  json:"updatedAt,omitempty"`
}

type NewUserResponse struct {
	LastInsertedId int64
	RowsAffected   int64
}

type LoginUser struct {
	Username     string `bson:"username,omitempty"  json:"username,omitempty"`
	Email        string `bson:"email,omitempty"  json:"email,omitempty"`
	Token        string `bson:"token,omitempty"  json:"token,omitempty"`
	RefreshToken string `bson:"refreshToken,omitempty"  json:"refreshToken,omitempty"`
	UserId       string `bson:"userId,omitempty"  json:"userId,omitempty"`
}

type Message struct {
	ErrorLog  []ErrorLog `json:"errorLog,omitempty"`
	HostName  string     `json:"hostName,omitempty"`
	Status    string     `json:"status,omitempty"`
	TimeTaken string     `json:"timeTaken,omitempty"`
	Count     int        `json:"count,omitempty"`
}

type ErrorLog struct {
	Status    string `json:"status,omitempty"`
	Trace     string `json:"trace,omitempty"`
	RootCause string `json:"rootCause,omitempty"`
}
