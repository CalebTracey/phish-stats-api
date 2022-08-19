package routes

import (
	"encoding/json"
	"github.com/calebtracey/phish-stats-api/internal/facade"
	"github.com/calebtracey/phish-stats-api/internal/models"
	"github.com/calebtracey/phish-stats-api/internal/services/auth"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Handler struct {
	Service     facade.ServiceI
	AuthService auth.ServiceI
}

func (h Handler) InitializeRoutes() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	secure := r.PathPrefix("/api").Subrouter()
	secure.Use(h.AuthService.Middleware)

	// Health check
	r.Handle("/health", h.HealthCheck()).Methods(http.MethodGet)

	// User
	r.Handle("/auth/register", h.RegistrationHandler()).Methods(http.MethodPost)
	r.Handle("/auth/login", h.LoginHandler()).Methods(http.MethodPost)

	r.Handle("/show", h.GetShowHandler()).Methods(http.MethodGet)

	return r
}

func (h Handler) RegistrationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response models.UserResponse
		var request models.User
		startTime := time.Now()

		defer func() {
			response.Message.TimeTaken = time.Since(startTime).String()
			response, status := setUserResponse(response)
			_ = json.NewEncoder(writeHeader(w, status)).Encode(response)
		}()

		requestBody, readErr := ioutil.ReadAll(r.Body)

		if readErr != nil {
			response.Message.ErrorLog = errorLogs([]error{readErr}, "Unable to read request body", http.StatusBadRequest)
			return
		}
		err := json.Unmarshal(requestBody, &request)
		if err != nil {
			response.Message.ErrorLog = errorLogs([]error{err}, "Unable to parse request", http.StatusBadRequest)
			return
		}

		response = h.Service.RegisterUser(r.Context(), request)
	}
}

func (h Handler) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response models.UserResponse
		var request models.User
		startTime := time.Now()

		defer func() {
			response.Message.TimeTaken = time.Since(startTime).String()
			response, status := setUserResponse(response)
			_ = json.NewEncoder(writeHeader(w, status)).Encode(response)
		}()

		requestBody, readErr := ioutil.ReadAll(r.Body)

		if readErr != nil {
			response.Message.ErrorLog = errorLogs([]error{readErr}, "Unable to read request body", http.StatusBadRequest)
			return
		}
		err := json.Unmarshal(requestBody, &request)
		if err != nil {
			response.Message.ErrorLog = errorLogs([]error{err}, "Unable to parse request", http.StatusBadRequest)
			return
		}

		response = h.Service.LoginUser(r.Context(), request)
	}
}

func (h Handler) GetShowHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		var response models.GetShowResponse

		defer func() {
			//response, status := setUserResponse(response)
			status, _ := strconv.Atoi(response.Message.Status)
			response.Message.TimeTaken = time.Since(startTime).String()
			_ = json.NewEncoder(writeHeader(w, status)).Encode(response)
		}()
		var request models.GetShowRequest
		requestBody, readErr := ioutil.ReadAll(r.Body)

		if readErr != nil {
			response.Message.ErrorLog = errorLogs([]error{readErr}, "Unable to read request body", http.StatusBadRequest)
			return
		}
		err := json.Unmarshal(requestBody, &request)
		if err != nil {
			response.Message.ErrorLog = errorLogs([]error{err}, "Unable to parse request", http.StatusBadRequest)
			return
		}

		response = h.Service.GetShow(r.Context(), request)
	}
}

func (h Handler) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			logrus.Errorln(err.Error())
			return
		}
	}
}

func setUserResponse(res models.UserResponse) (models.LoginResponse, int) {
	status, _ := strconv.Atoi(res.Message.Status)
	hn, _ := os.Hostname()
	return models.LoginResponse{
		Username:     res.User.Username,
		Email:        res.User.Email,
		Token:        res.User.Token,
		RefreshToken: res.User.RefreshToken,
		Message: models.Message{
			ErrorLog: res.Message.ErrorLog,
			HostName: hn,
			Status:   res.Message.Status,
		},
	}, status
}

func writeHeader(w http.ResponseWriter, code int) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	return w
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
