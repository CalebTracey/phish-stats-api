package auth

import (
	"fmt"
	"github.com/calebtracey/phish-stats-api/internal/models"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

var SecretKey = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	Email    string
	FullName string
	Uid      string
	jwt.RegisteredClaims
}

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

//VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}

func GenerateAllTokens(user models.User) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:    *user.Email,
		FullName: *user.FullName,
		Uid:      user.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Local().Add(time.Hour * time.Duration(24)),
			},
		},
	}

	refreshClaims := &SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Local().Add(time.Hour * time.Duration(168)),
			},
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SecretKey))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SecretKey))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg := err.Error()
		return nil, fmt.Errorf("the token is invalid; err: %v", msg)
	}

	if claims.ExpiresAt.Time.After(time.Now().Local()) {
		msg := err.Error()
		return nil, fmt.Errorf("token is expired; err: %v", msg)
	}

	return claims, nil
}
