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
	GetShow(ctx context.Context, req models.GetShowRequest) models.ShowResponse
	GetUser(ctx context.Context, id string) models.UserResponse
	AddUserShow(ctx context.Context, request models.AddUserShowRequest) models.AddShowResponse
}

type Service struct {
	PsqlService psql.ServiceI
	PNService   phishnet.ServiceI
	AuthService auth.ServiceI
	PNMapper    phishnet.MapperI
	Validator   *validator.Validate
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
		PsqlService: psqlService,
		PNService:   phishNetService,
		AuthService: auth.Service{},
		PNMapper:    phishnet.Mapper{},
		Validator:   validate,
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
		response.User = &models.UserParsedResponse{}
		return response
	}

	u := updateUserRequest(userRequest)
	created := u.CreatedAt.String()
	updated := u.UpdatedAt.String()
	pwHash := s.AuthService.HashPassword(u.Password)

	// if a user was inserted successfully, create auth tokens for response
	token, refreshToken, err := s.updateUserTokens(ctx, userRequest, updated)
	if err != nil {
		message.ErrorLog = errorLogs([]error{err}, "Token update error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.User = &models.UserParsedResponse{}
		response.Message = message
		return response
	}
	var shows []string
	exec := fmt.Sprintf(psql.AddUser, u.ID, u.FullName, u.Email, u.Username, pwHash, token, refreshToken, created, updated, shows)

	_, errs := s.PsqlService.InsertNewUser(ctx, exec)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Warnf("insert new user error: %v", err.Error())
		}
		//TODO probably remove these error logs and just have them log to warnings?
		message.ErrorLog = errorLogs(errs, "New user error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
	}

	message.Status = strconv.Itoa(http.StatusOK)
	message.Count = 1
	response.User = &models.UserParsedResponse{
		FullName:     u.FullName,
		Username:     u.Username,
		RefreshToken: refreshToken,
		Token:        token,
	}
	response.Message = message
	log.Infof("New user %v registered", response.User.Username)
	return response
}

// LoginUser uses auth service to verify user request and psql service to access user data
// TODO test coverage for error cases
func (s *Service) LoginUser(ctx context.Context, userRequest models.User) models.UserResponse {
	var response models.UserResponse
	var message models.Message

	foundUserExec := fmt.Sprintf(psql.FindUserByEmail, userRequest.Email)
	foundUser, errs := s.PsqlService.FindUser(ctx, foundUserExec)

	if len(errs) > 0 {
		message.ErrorLog = errorLogs(errs, "User not found", http.StatusNotFound)
		message.Status = strconv.Itoa(http.StatusNotFound)
		response.User = &models.UserParsedResponse{}
		response.Message = message
		return response
	}

	passwordIsValid, msg := s.AuthService.VerifyPassword(userRequest.Password, foundUser.Password)

	if !passwordIsValid {
		message.ErrorLog = errorLogs([]error{fmt.Errorf(msg)}, fmt.Sprintf("Verification error %v", userRequest.Email), http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.User = &models.UserParsedResponse{}
		response.Message = message
		return response
	}

	updated := time.Now().Format(time.RFC3339)
	userRequest.Email = foundUser.Email
	userRequest.ID = foundUser.ID

	token, refreshToken, err := s.updateUserTokens(ctx, userRequest, updated)
	if err != nil {
		message.ErrorLog = errorLogs([]error{err}, "Token update error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.User = &models.UserParsedResponse{}
		response.Message = message
		return response
	}

	message.Status = strconv.Itoa(http.StatusOK)
	message.Count = 1
	response.Message = message
	foundUser.Token = token
	foundUser.RefreshToken = refreshToken
	response.User = mapUserPublicData(foundUser)
	log.Infof("User %v logged in", response.User.Username)
	return response
}

func (s *Service) GetUser(ctx context.Context, id string) models.UserResponse {
	var response models.UserResponse
	var message models.Message

	foundUserExec := fmt.Sprintf(psql.FindUserById, id)
	foundUser, errs := s.PsqlService.FindUser(ctx, foundUserExec)
	if len(errs) > 0 {
		message.ErrorLog = errorLogs(errs, "User not found", http.StatusNotFound)
		message.Status = strconv.Itoa(http.StatusNotFound)
		response.User = &models.UserParsedResponse{}
		response.Message = message
		return response
	}

	message.Status = strconv.Itoa(http.StatusOK)
	message.Count = 1
	response.Message = message
	response.User = mapUserPublicData(foundUser)
	return response
}

func (s *Service) GetShow(ctx context.Context, req models.GetShowRequest) (response models.ShowResponse) {
	var message models.Message
	pnResponse, err := s.PNService.GetShow(ctx, req.Date)
	if err != nil {
		message.ErrorLog = errorLogs([]error{err}, "Get show error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.Message = message
		return response
	}
	response = s.PNMapper.PhishNetResponseToShowResponse(pnResponse)
	message.Status = strconv.Itoa(http.StatusOK)
	message.Count = 1

	return response
}

func (s *Service) AddUserShow(ctx context.Context, request models.AddUserShowRequest) models.AddShowResponse {
	var response models.AddShowResponse
	var message models.Message
	u := s.GetUser(ctx, request.Id)
	if containsStr(u.User.Shows, request.Date) {
		log.Infof("Show %v not added, already exists", request.Date)
		return response
	}
	exec := fmt.Sprintf(psql.AddUserShow, request.Date, request.Id)
	err := s.PsqlService.InsertOne(ctx, exec)

	if err != nil {
		message.ErrorLog = errorLogs([]error{err}, "Add User Show error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.Message = message
		return response
	}

	response.Date = request.Date
	message.Status = strconv.Itoa(http.StatusOK)
	message.Count = 1
	response.Message = message
	return response
}

func mapUserPublicData(user *models.UserParsedResponse) *models.UserParsedResponse {
	var res models.UserParsedResponse
	res.Username = user.Username
	res.Shows = user.Shows
	res.Email = user.Email
	res.ID = user.ID
	res.Token = user.Token
	res.RefreshToken = user.RefreshToken
	res.FullName = user.FullName
	return &res
}

func (s *Service) updateUserTokens(ctx context.Context, user models.User, updated string) (string, string, error) {
	token, refreshToken, _ := s.AuthService.GenerateAllTokens(user)
	exec := fmt.Sprintf(psql.UpdateTokens, token, refreshToken, updated, user.ID)
	err := s.PsqlService.InsertOne(ctx, exec)
	if err != nil {
		log.Errorf("failed to update tokens for user: %v", user.Username)
		return "", "", err
	}
	return token, refreshToken, nil
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
		log.Errorf("%v: %v", rootCause, err.Error())
		errLogs = append(errLogs, models.ErrorLog{
			RootCause: rootCause,
			Status:    strconv.Itoa(status),
			Trace:     err.Error(),
		})
	}
	return errLogs
}

func containsStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
