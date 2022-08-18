package phishnet

import (
	"context"
	"encoding/json"
	config "github.com/calebtracey/config-yaml"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

type ServiceI interface {
	GetShow(ctx context.Context, method string) (ShowResponse, error)
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

func (s *Service) GetShow(ctx context.Context, date string) (ShowResponse, error) {
	var response ShowResponse

	apiUrl := s.BaseUrl + "/setlists/showdate/" + date + s.Format + "?" + s.ApiKeyUri
	logrus.Infoln(apiUrl)
	req, _ := http.NewRequestWithContext(ctx, "GET", apiUrl, nil)

	resp, err := s.Client.Do(req)

	if err != nil {
		return response, err
	}
	responseBody, readErr := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, readErr
	}
	unmarshallErr := json.Unmarshal(responseBody, &response)
	if err != nil {
		return response, unmarshallErr
	}

	return response, nil
}
