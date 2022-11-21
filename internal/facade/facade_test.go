package facade

import (
	"context"
	"fmt"
	"github.com/calebtracey/phish-stats-api/internal/models"
	"github.com/calebtracey/phish-stats-api/internal/services/auth"
	"github.com/calebtracey/phish-stats-api/internal/services/phishnet"
	"github.com/calebtracey/phish-stats-api/internal/services/psql"
	"github.com/go-playground/assert/v2"
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
		mockResponse  *models.UserParsedResponse
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
				Email:    "testusername123@email.com",
				Password: "password123",
			},
			query: fmt.Sprintf(psql.FindUserByEmail, "testusername123@email.com"),
			exec:  fmt.Sprintf(psql.UpdateTokens, "39048567301249586", "01938467501934651", updated, "542113"),
			want: models.UserResponse{
				User: &models.UserParsedResponse{
					ID:           "542113",
					FullName:     "Test User",
					Username:     "testusername123",
					Email:        "testusername123@email.com",
					Token:        "39048567301249586",
					RefreshToken: "01938467501934651",
					Shows:        []string{},
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
			mockResponse: &models.UserParsedResponse{
				ID:           "542113",
				FullName:     "Test User",
				Email:        "testusername123@email.com",
				Username:     "testusername123",
				Password:     "password123",
				Token:        "39048567301249586",
				RefreshToken: "01938467501934651",
				CreatedAt:    created,
				UpdatedAt:    updated,
				Shows:        []string{},
			},
			mockErrs: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				PsqlService: tt.fields.PsqlService,
				PNService:   tt.fields.PhishNetService,
				AuthService: tt.fields.AuthService,
				Validator:   tt.fields.Validator,
			}
			mockPsqlSvc.EXPECT().FindUser(tt.ctx, tt.query).Return(tt.mockResponse, tt.mockErrs)
			mockAuthSvc.EXPECT().VerifyPassword(gomock.Any(), gomock.Any()).Return(tt.wantVerify, tt.wantVerifyMsg)
			mockAuthSvc.EXPECT().GenerateAllTokens(gomock.Any()).Return("39048567301249586", "01938467501934651", tt.wantGenErr)
			mockPsqlSvc.EXPECT().InsertOne(tt.ctx, gomock.Any()).Return(tt.wantUpdateErr)
			got := s.LoginUser(tt.ctx, tt.userRequest)
			assert.Equal(t, got.User, tt.want.User)
			assert.Equal(t, got.Message.Status, tt.want.Message.Status)
			//if  !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("LoginUser() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestService_RegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPsqlSvc := psql.NewMockServiceI(ctrl)
	mockPhNetSvc := phishnet.NewMockServiceI(ctrl)
	mockAuthSvc := auth.NewMockServiceI(ctrl)
	mockMapper := psql.NewMockMapperI(ctrl)
	type fields struct {
		PsqlService     psql.ServiceI
		PhishNetService phishnet.ServiceI
		AuthService     auth.ServiceI
		PSQLMapper      psql.MapperI
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
				PSQLMapper:      mockMapper,
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
				User: &models.UserParsedResponse{
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
				PsqlService: tt.fields.PsqlService,
				PNService:   tt.fields.PhishNetService,
				AuthService: tt.fields.AuthService,
				PSQLMapper:  tt.fields.PSQLMapper,
				Validator:   tt.fields.Validator,
			}
			mockMapper.EXPECT().CreatePSQLUserExec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("")
			mockPsqlSvc.EXPECT().InsertNewUser(tt.ctx, gomock.Any()).Return(tt.mockRes, tt.execErrs)
			mockAuthSvc.EXPECT().HashPassword(tt.userRequest.Password).Return(tt.mockHash)
			mockAuthSvc.EXPECT().GenerateAllTokens(gomock.Any()).Return("39048567301249586", "01938467501934651", tt.wantGenErr)
			mockPsqlSvc.EXPECT().InsertOne(tt.ctx, gomock.Any()).Return(tt.wantUpdateErr)

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
	mockMapper := phishnet.NewMockMapperI(ctrl)
	type fields struct {
		PsqlService psql.ServiceI
		PNService   phishnet.ServiceI
		AuthService auth.ServiceI
		PNMapper    phishnet.MapperI
		Validator   *validator.Validate
	}

	tests := []struct {
		name         string
		fields       fields
		ctx          context.Context
		req          models.GetShowRequest
		want         models.ShowResponse
		mockResponse phishnet.PNShowResponse
		mockErr      error
	}{
		{
			name: "Happy Path",
			fields: fields{
				PsqlService: mockPsqlSvc,
				PNService:   mockPhNetSvc,
				AuthService: mockAuthSvc,
				PNMapper:    mockMapper,
				Validator:   validator.New(),
			},
			ctx: context.Background(),
			req: models.GetShowRequest{Date: "11/10/1991"},
			want: models.ShowResponse{
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
			mockResponse: phishnet.PNShowResponse{
				Error:        false,
				ErrorMessage: "",
				Data: []phishnet.PNShowData{
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
				PsqlService: tt.fields.PsqlService,
				PNService:   tt.fields.PNService,
				AuthService: tt.fields.AuthService,
				PNMapper:    tt.fields.PNMapper,
				Validator:   tt.fields.Validator,
			}
			mockMapper.EXPECT().PhishNetResponseToShowResponse(gomock.Any()).Return(tt.want)
			mockPhNetSvc.EXPECT().GetShow(tt.ctx, tt.req.Date).Return(tt.mockResponse, tt.mockErr)

			if got := s.GetShow(tt.ctx, tt.req); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetShow() = %v, want %v", got, tt.want)
			}
		})
	}
}
