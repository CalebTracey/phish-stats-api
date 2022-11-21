package psql

import (
	"fmt"
	"github.com/calebtracey/phish-stats-api/internal/models"
)

//go:generate mockgen -destination=mockMapper.go -package=psql . MapperI
type MapperI interface {
	CreatePSQLUserExec(u models.User, pwHash, token, refreshToken string) string
}

type Mapper struct{}

func (m Mapper) CreatePSQLUserExec(u models.User, pwHash, token, refreshToken string) string {
	var shows []string
	created := u.CreatedAt.String()
	updated := u.UpdatedAt.String()
	return fmt.Sprintf(AddUser, u.ID, u.FullName, u.Email, u.Username, pwHash, token, refreshToken, created, updated, shows)
}
