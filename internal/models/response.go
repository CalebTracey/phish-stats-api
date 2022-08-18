package models

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

type Show struct {
	Date  string
	Venue string
	Songs []Song
}

type Song struct {
	Title     string
	TrackTime string
}

type UserResponse struct {
	User    *User
	Message Message
}

type AllUsersResponse struct {
	Users   []*User
	Message Message
}

type LoginResponse struct {
	Username     *string `bson:"username,omitempty"  json:"username,omitempty"`
	Email        *string `bson:"email,omitempty"  json:"email,omitempty"`
	Token        *string `bson:"token,omitempty"  json:"token,omitempty"`
	RefreshToken *string `bson:"refreshToken,omitempty"  json:"refreshToken,omitempty"`
	UserId       string  `bson:"userId,omitempty"  json:"userId,omitempty"`
	Message      Message `json:"message,omitempty"`
}

type RegistrationResponse struct {
	AccessToken *string `json:"accessToken,omitempty"`
	Email       *string `json:"email,omitempty"`
	Id          string  `json:"id,omitempty"`
	Message     Message
}
