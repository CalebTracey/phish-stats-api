package facade

import (
	"context"
	"fmt"
	config "github.com/calebtracey/config-yaml"
	"github.com/calebtracey/phish-stats-api/internal/models"
	"github.com/calebtracey/phish-stats-api/internal/services/auth"
	"github.com/calebtracey/phish-stats-api/internal/services/mongodb"
	"github.com/calebtracey/phish-stats-api/internal/services/phishnet"
	"github.com/go-playground/validator/v10"
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
	MongoService    mongodb.ServiceI
	PhishNetService phishnet.ServiceI
	Validator       *validator.Validate
}

func NewService(appConfig *config.Config) (Service, error) {
	mongoService, err := mongodb.InitializeMongoService(appConfig)
	phishNetService, err := phishnet.InitializePhishNetService(appConfig)
	validate := validator.New()

	if err != nil {
		return Service{}, err
	}

	return Service{
		MongoService:    mongoService,
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
		return response
	}

	userRequest.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	userRequest.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	userRequest.ID = primitive.NewObjectIDFromTimestamp(userRequest.CreatedAt)
	userRequest.UserId = userRequest.ID.Hex()

	token, refreshToken, _ := auth.GenerateAllTokens(userRequest)
	userRequest.Token = &token
	userRequest.RefreshToken = &refreshToken

	user, err := s.MongoService.AddNewUser(ctx, &userRequest)
	if err != nil {
		message.ErrorLog = errorLogs([]error{err}, "New user error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.Message = message
		return response
	}

	response.Message = message
	response.User = user

	return response
}

func (s *Service) LoginUser(ctx context.Context, userRequest models.User) models.UserResponse {
	var message models.Message
	var response models.UserResponse

	foundUser, err := s.MongoService.FindUserByUsername(ctx, &userRequest)
	if err != nil {
		message.ErrorLog = errorLogs([]error{err}, "Find user error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.Message = message
		return response
	}

	passwordIsValid, msg := auth.VerifyPassword(*userRequest.Password, *foundUser.Password)
	if passwordIsValid != true {
		message.ErrorLog = errorLogs([]error{fmt.Errorf(msg)}, "Verification error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.Message = message
		return response
	}

	token, refreshToken, _ := auth.GenerateAllTokens(foundUser)

	err = s.MongoService.UpdateAllTokens(ctx, token, refreshToken, foundUser.UserId)
	if err != nil {
		message.ErrorLog = errorLogs([]error{err}, "Token update error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.Message = message
		return response
	}

	response.Message.Status = strconv.Itoa(http.StatusOK)
	response.User = &foundUser

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
