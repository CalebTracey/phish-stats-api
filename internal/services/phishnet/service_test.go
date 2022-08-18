package phishnet

import (
	"context"
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"reflect"
	"testing"
)

func TestService_GetShow(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://testUrl200.com/setlists/showdate/11/10/1991.json?apiKey=123Test",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Accept") != "application/json" {
				t.Errorf("Expected Accept: application/json header, got: %s", req.Header.Get("Accept"))
			}
			resp, err := httpmock.NewJsonResponse(200, ShowResponse{
				Error:        false,
				ErrorMessage: "",
				Data: []Data{
					{
						Showdate: "11/10/1991",
						Song:     "song 1",
					},
					{
						Showdate: "11/10/1991",
						Song:     "song 2",
					},
					{
						Showdate: "11/10/1991",
						Song:     "song 3",
					},
				},
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("GET", "http://testUrl500.com/setlists/showdate/11/10/1991.json?apiKey=123Test",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Accept") != "application/json" {
				t.Errorf("Expected Accept: application/json header, got: %s", req.Header.Get("Accept"))
			}
			testErr := fmt.Errorf("test error")
			return &http.Response{}, testErr
		},
	)

	type fields struct {
		Client    *http.Client
		BaseUrl   string
		ApiKeyUri string
		Format    string
	}

	tests := []struct {
		name    string
		fields  fields
		ctx     context.Context
		date    string
		want    ShowResponse
		wantErr bool
	}{
		{
			name: "Happy Path",
			fields: fields{
				Client:    &http.Client{},
				BaseUrl:   "http://testUrl200.com",
				ApiKeyUri: "apiKey=123Test",
				Format:    ".json",
			},
			ctx:  context.Background(),
			date: "11/10/1991",
			want: ShowResponse{
				Error:        false,
				ErrorMessage: "",
				Data: []Data{
					{
						Showdate: "11/10/1991",
						Song:     "song 1",
					},
					{
						Showdate: "11/10/1991",
						Song:     "song 2",
					},
					{
						Showdate: "11/10/1991",
						Song:     "song 3",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Path",
			fields: fields{
				Client:    &http.Client{},
				BaseUrl:   "http://testUrl500.com",
				ApiKeyUri: "apiKey=123Test",
				Format:    ".json",
			},
			ctx:     context.Background(),
			date:    "11/10/1991",
			want:    ShowResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Client:    tt.fields.Client,
				BaseUrl:   tt.fields.BaseUrl,
				ApiKeyUri: tt.fields.ApiKeyUri,
				Format:    tt.fields.Format,
			}
			got, err := s.GetShow(tt.ctx, tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetShow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetShow() got = %v, want %v", got, tt.want)
			}
		})
	}
}
