package facade

import (
	"context"
	"fmt"
	"github.com/calebtracey/phish-stats-api/internal/models"
	"github.com/calebtracey/phish-stats-api/internal/services/auth"
	"github.com/calebtracey/phish-stats-api/internal/services/phishnet"
	"github.com/calebtracey/phish-stats-api/internal/services/psql"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestService_LoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPsqlSvc := psql.NewMockServiceI(ctrl)
	mockPhNetSvc := phishnet.NewMockServiceI(ctrl)
	mockAuthSvc := auth.NewMockServiceI(ctrl)

	created := time.Now().Format(time.RFC3339)
	updated := time.Now().Format(time.RFC3339)

	type fields struct {
		PsqlService     psql.ServiceI
		PhishNetService phishnet.ServiceI
		AuthService     auth.ServiceI
		Validator       *validator.Validate
	}

	tests := []struct {
		name          string
		fields        fields
		ctx           context.Context
		userRequest   models.User
		query         string
		exec          string
		want          models.UserResponse
		wantVerify    bool
		wantVerifyMsg string
		wantGenErr    error
		wantUpdateErr error
		mockResponse  *models.UserPsqlResponse
		mockErrs      []error
	}{
		{
			name: "Happy Path",
			fields: fields{
				PsqlService:     mockPsqlSvc,
				PhishNetService: mockPhNetSvc,
				AuthService:     mockAuthSvc,
				Validator:       validator.New(),
			},
			ctx: context.Background(),
			userRequest: models.User{
				Username: "testusername123",
				Password: "password123",
			},
			query: fmt.Sprintf(psql.FindUserByUsername, "testusername123"),
			exec:  fmt.Sprintf(psql.UpdateTokens, "39048567301249586", "01938467501934651", updated, "542113"),
			want: models.UserResponse{
				User: &models.UserPsqlResponse{
					ID:           "542113",
					FullName:     "Test User",
					Username:     "testusername123",
					Token:        "39048567301249586",
					RefreshToken: "01938467501934651",
				},
				Message: models.Message{
					ErrorLog: nil,
					Status:   strconv.Itoa(http.StatusOK),
					Count:    1,
				},
			},
			wantVerify:    true,
			wantVerifyMsg: "",
			wantGenErr:    nil,
			wantUpdateErr: nil,
			mockResponse: &models.UserPsqlResponse{
				ID:           "542113",
				FullName:     "Test User",
				Email:        "test@email.com",
				Username:     "testusername123",
				Password:     "password123",
				Token:        "39048567301249586",
				RefreshToken: "01938467501934651",
				CreatedAt:    created,
				UpdatedAt:    updated,
			},
			mockErrs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				PsqlService:     tt.fields.PsqlService,
				PhishNetService: tt.fields.PhishNetService,
				AuthService:     tt.fields.AuthService,
				Validator:       tt.fields.Validator,
			}
			mockPsqlSvc.EXPECT().FindUserByUsername(tt.ctx, tt.query).Return(tt.mockResponse, tt.mockErrs)
			mockAuthSvc.EXPECT().VerifyPassword(gomock.Any(), gomock.Any()).Return(tt.wantVerify, tt.wantVerifyMsg)
			mockAuthSvc.EXPECT().GenerateAllTokens(gomock.Any()).Return("39048567301249586", "01938467501934651", tt.wantGenErr)
			mockPsqlSvc.EXPECT().UpdateAllTokens(tt.ctx, gomock.Any()).Return(tt.wantUpdateErr)
			if got := s.LoginUser(tt.ctx, tt.userRequest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoginUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPsqlSvc := psql.NewMockServiceI(ctrl)
	mockPhNetSvc := phishnet.NewMockServiceI(ctrl)
	mockAuthSvc := auth.NewMockServiceI(ctrl)

	type fields struct {
		PsqlService     psql.ServiceI
		PhishNetService phishnet.ServiceI
		AuthService     auth.ServiceI
		Validator       *validator.Validate
	}

	tests := []struct {
		name          string
		fields        fields
		ctx           context.Context
		userRequest   models.User
		want          models.UserResponse
		mockRes       *models.NewUserResponse
		mockHash      string
		execErrs      []error
		wantUpdateErr error
		wantGenErr    error
	}{
		{
			name: "Happy Path",
			fields: fields{
				PsqlService:     mockPsqlSvc,
				PhishNetService: mockPhNetSvc,
				AuthService:     mockAuthSvc,
				Validator:       validator.New(),
			},
			ctx: context.Background(),
			userRequest: models.User{
				FullName: "Test User",
				Email:    "test@email.com",
				Username: "testusername123",
				Password: "password123",
			},
			execErrs:      nil,
			wantUpdateErr: nil,
			wantGenErr:    nil,
			want: models.UserResponse{
				User: &models.UserPsqlResponse{
					FullName:     "Test User",
					Username:     "testusername123",
					Token:        "39048567301249586",
					RefreshToken: "01938467501934651",
				},
				Message: models.Message{
					ErrorLog: nil,
					Status:   strconv.Itoa(http.StatusOK),
					Count:    1,
				},
			},
			mockHash: "0123230f8h23f8h023fh8",
			mockRes: &models.NewUserResponse{
				LastInsertedId: int64(8013450345345),
				RowsAffected:   int64(9),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				PsqlService:     tt.fields.PsqlService,
				PhishNetService: tt.fields.PhishNetService,
				AuthService:     tt.fields.AuthService,
				Validator:       tt.fields.Validator,
			}
			mockPsqlSvc.EXPECT().InsertNewUser(tt.ctx, gomock.Any()).Return(tt.mockRes, tt.execErrs)
			mockAuthSvc.EXPECT().HashPassword(tt.userRequest.Password).Return(tt.mockHash)
			mockAuthSvc.EXPECT().GenerateAllTokens(gomock.Any()).Return("39048567301249586", "01938467501934651", tt.wantGenErr)
			mockPsqlSvc.EXPECT().UpdateAllTokens(tt.ctx, gomock.Any()).Return(tt.wantUpdateErr)

			if got := s.RegisterUser(tt.ctx, tt.userRequest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetShow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPsqlSvc := psql.NewMockServiceI(ctrl)
	mockPhNetSvc := phishnet.NewMockServiceI(ctrl)
	mockAuthSvc := auth.NewMockServiceI(ctrl)

	type fields struct {
		PsqlService     psql.ServiceI
		PhishNetService phishnet.ServiceI
		AuthService     auth.ServiceI
		Validator       *validator.Validate
	}

	tests := []struct {
		name         string
		fields       fields
		ctx          context.Context
		req          models.GetShowRequest
		want         models.GetShowResponse
		mockResponse phishnet.ShowResponse
		mockErr      error
	}{
		{
			name: "Happy Path",
			fields: fields{
				PsqlService:     mockPsqlSvc,
				PhishNetService: mockPhNetSvc,
				AuthService:     mockAuthSvc,
				Validator:       validator.New(),
			},
			ctx: context.Background(),
			req: models.GetShowRequest{Date: "11/10/1991"},
			want: models.GetShowResponse{
				Show: models.Show{
					Date: "11/10/1991",
					Songs: []models.Song{
						{
							Title: "song 1",
						},
						{
							Title: "song 2",
						},
						{
							Title: "song 3",
						},
					},
				},
				Message: models.Message{},
			},
			mockResponse: phishnet.ShowResponse{
				Error:        false,
				ErrorMessage: "",
				Data: []phishnet.Data{
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
			mockErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				PsqlService:     tt.fields.PsqlService,
				PhishNetService: tt.fields.PhishNetService,
				AuthService:     tt.fields.AuthService,
				Validator:       tt.fields.Validator,
			}
			mockPhNetSvc.EXPECT().GetShow(tt.ctx, tt.req.Date).Return(tt.mockResponse, tt.mockErr)
			if got := s.GetShow(tt.ctx, tt.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetShow() = %v, want %v", got, tt.want)
			}
		})
	}
}
