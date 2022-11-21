package phishnet

import (
	"context"
	"encoding/json"
	"fmt"
	config "github.com/calebtracey/config-yaml"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

//go:generate mockgen -destination=mockService.go -package=phishnet . ServiceI
type ServiceI interface {
	GetShow(ctx context.Context, method string) (PNShowResponse, error)
}

type Service struct {
	Client    *http.Client
	BaseUrl   string
	ApiKeyUri string
	Format    string
}

func InitializePhishNetService(config *config.Config) (*Service, error) {
	phishNetSvc, err := config.GetServiceConfig("phishnet")
	if err != nil {
		return nil, err
	}
	err = godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}
	ApiKey := os.Getenv(phishNetSvc.ApiKeyEnvironmentVariable.Value)

	return &Service{
		Client:    &http.Client{},
		BaseUrl:   phishNetSvc.URL.Value,
		ApiKeyUri: "apikey=" + ApiKey,
		Format:    ".json",
	}, nil
}

func (s *Service) GetShow(ctx context.Context, date string) (PNShowResponse, error) {
	var response PNShowResponse
	req, reqErr := newPhishNetRequest(ctx, s.BaseUrl+"/setlists/showdate/"+date+s.Format+"?"+s.ApiKeyUri)
	if reqErr != nil {
		return response, reqErr
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return response, err
	}
	responseBody, readErr := io.ReadAll(resp.Body)
	if err != nil {
		return response, readErr
	}
	unmarshallErr := json.Unmarshal(responseBody, &response)
	if err != nil {
		return response, unmarshallErr
	}

	if response.Error {
		return response, fmt.Errorf(response.ErrorMessage)
	}

	return response, nil
}

func newPhishNetRequest(ctx context.Context, url string) (*http.Request, error) {
	logrus.Infoln(url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	return req, nil
}
