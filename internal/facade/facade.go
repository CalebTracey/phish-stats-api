package facade

import (
	"context"
	"fmt"
	config "github.com/calebtracey/config-yaml"
	"github.com/calebtracey/phish-stats-api/internal/models"
	"github.com/calebtracey/phish-stats-api/internal/services/auth"
	"github.com/calebtracey/phish-stats-api/internal/services/phishnet"
	"github.com/calebtracey/phish-stats-api/internal/services/psql"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
)

//go:generate mockgen -destination=mockFacade.go -package=facade . ServiceI
type ServiceI interface {
	LoginUser(ctx context.Context, userRequest models.User) models.UserResponse
	RegisterUser(ctx context.Context, userRequest models.User) models.UserResponse
	GetShow(ctx context.Context, req models.GetShowRequest) models.GetShowResponse
}

type Service struct {
	PsqlService     psql.ServiceI
	PhishNetService phishnet.ServiceI
	Validator       *validator.Validate
}

func NewService(appConfig *config.Config) (Service, error) {
	psqlConfig, psqlErr := appConfig.GetDatabaseConfig("PSTGQL")
	if psqlErr != nil {
		return Service{}, psqlErr
	}
	psqlService := psql.InitializePsqlService(psqlConfig)
	phishNetService, err := phishnet.InitializePhishNetService(appConfig)
	validate := validator.New()

	if err != nil {
		return Service{}, err
	}

	return Service{
		PsqlService:     psqlService,
		PhishNetService: phishNetService,
		Validator:       validate,
	}, nil
}

func (s *Service) RegisterUser(ctx context.Context, userRequest models.User) models.UserResponse {
	var response models.UserResponse
	var message models.Message

	validationErr := s.Validator.Struct(userRequest)

	if validationErr != nil {
		message.ErrorLog = errorLogs([]error{validationErr}, "Validation error", http.StatusBadRequest)
		message.Status = strconv.Itoa(http.StatusBadRequest)
		response.Message = message
		response.User = &models.UserPsqlResponse{}
		return response
	}

	u := updateUserRequest(userRequest)
	created := u.CreatedAt.String()
	updated := u.UpdatedAt.String()
	pwHash := auth.HashPassword(u.Password)

	exec := fmt.Sprintf(psql.AddUser, u.ID, u.FullName, u.Email, u.Username, pwHash, u.Token, u.RefreshToken, created, updated)

	_, errs := s.PsqlService.InsertNewUser(ctx, exec)

	if errs != nil && len(errs) > 0 {
		for _, err := range errs {
			log.Warnf("insert new user error: %v", err.Error())
		}
		//TODO probably remove these error logs and just have them log to warnings?
		message.ErrorLog = errorLogs(errs, "New user error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
	}

	return models.UserResponse{
		User: &models.UserPsqlResponse{
			ID:           u.ID,
			FullName:     u.FullName,
			Email:        u.Email,
			Username:     u.Username,
			Password:     u.Password,
			Token:        u.Token,
			RefreshToken: u.RefreshToken,
			CreatedAt:    created,
			UpdatedAt:    updated,
		},
		Message: message,
	}
}

func (s *Service) LoginUser(ctx context.Context, userRequest models.User) models.UserResponse {
	var message models.Message
	var response models.UserResponse

	foundUserExec := fmt.Sprintf(psql.FindUserByUsername, userRequest.Username)
	foundUser, errs := s.PsqlService.FindUserByUsername(ctx, foundUserExec)

	if errs != nil && len(errs) > 0 {
		message.ErrorLog = errorLogs(errs, "User not found", http.StatusNotFound)
		message.Status = strconv.Itoa(http.StatusNotFound)
		response.User = &models.UserPsqlResponse{}
		response.Message = message
		return response
	}

	passwordIsValid, msg := auth.VerifyPassword(userRequest.Password, foundUser.Password)

	if passwordIsValid != true {
		message.ErrorLog = errorLogs([]error{fmt.Errorf(msg)}, "Verification error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.User = &models.UserPsqlResponse{}
		response.Message = message
		return response
	}
	updated, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	token, refreshToken, _ := auth.GenerateAllTokens(*foundUser)
	exec := fmt.Sprintf(psql.UpdateTokens, token, refreshToken, updated, foundUser.ID)
	err := s.PsqlService.UpdateAllTokens(ctx, exec)

	if err != nil {
		message.ErrorLog = errorLogs([]error{err}, "Token update error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.User = &models.UserPsqlResponse{}
		response.Message = message
		return response
	}
	response.Message.Status = strconv.Itoa(http.StatusOK)
	response.User = foundUser
	response.User.RefreshToken = refreshToken
	response.User.Token = token
	response.User.UpdatedAt = updated.String()

	log.Infof("User %v logged in", response.User.Username)
	return response
}

func (s *Service) GetShow(ctx context.Context, req models.GetShowRequest) models.GetShowResponse {
	var response models.GetShowResponse
	var songs []models.Song
	var message models.Message

	showData, err := s.PhishNetService.GetShow(ctx, req.Date)

	if err != nil {
		message.ErrorLog = errorLogs([]error{err}, "Get show error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.Message = message
		return response
	}

	for _, data := range showData.Data {
		songs = append(songs, models.Song{
			Title:     data.Song,
			TrackTime: data.Tracktime,
		})
	}

	response.Show = models.Show{
		Date:  showData.Data[0].Showdate,
		Venue: showData.Data[0].Venue,
		Songs: songs,
	}

	return response
}

func updateUserRequest(userRequest models.User) models.User {
	userRequest.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	userRequest.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	userRequest.ID = primitive.NewObjectIDFromTimestamp(userRequest.CreatedAt).Hex()

	return userRequest
}

func errorLogs(errors []error, rootCause string, status int) []models.ErrorLog {
	var errLogs []models.ErrorLog
	for _, err := range errors {
		errLogs = append(errLogs, models.ErrorLog{
			RootCause: rootCause,
			Status:    strconv.Itoa(status),
			Trace:     err.Error(),
		})
	}
	return errLogs
}
