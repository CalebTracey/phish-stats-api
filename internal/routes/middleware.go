package routes

import (
	"encoding/json"
	"fmt"
	"github.com/calebtracey/phish-stats-api/internal/models"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (h Handler) Middleware(http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		var response models.AuthResponse
		defer func() {
			response, status := setMiddlewareResponse(response)
			response.Message.TimeTaken = time.Since(startTime).String()
			_ = json.NewEncoder(writeHeader(w, status)).Encode(response)
		}()

		clientToken := r.Header.Get("token")

		if clientToken == "" {
			response.Message.ErrorLog = errorLogs([]error{fmt.Errorf("no authorization header provided")}, "Authentication error", http.StatusInternalServerError)
			return
		}

		claims, err := h.AuthService.ValidateToken(clientToken)
		if err != nil {
			response.Message.ErrorLog = errorLogs([]error{err}, "Validation error", http.StatusForbidden)
			return
		}

		response.UserId = claims.Uid
		response.FullName = claims.FullName
		response.Email = claims.Email
	})

}

func setMiddlewareResponse(res models.AuthResponse) (models.AuthResponse, int) {
	hn, _ := os.Hostname()
	status, _ := strconv.Atoi(res.Message.Status)
	res.Message.HostName = hn
	return res, status
}

//func writeHeader(w http.ResponseWriter, code int) http.ResponseWriter {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(code)
//	return w
//}
//
//func errorLogs(errors []error, rootCause string, status int) []models.ErrorLog {
//	var errLogs []models.ErrorLog
//	for _, err := range errors {
//		errLogs = append(errLogs, models.ErrorLog{
//			RootCause: rootCause,
//			Status:    strconv.Itoa(status),
//			Trace:     err.Error(),
//		})
//	}
//	return errLogs
//}
