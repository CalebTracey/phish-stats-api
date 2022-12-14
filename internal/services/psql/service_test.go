package psql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/calebtracey/phish-stats-api/internal/models"
	"github.com/lib/pq"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestService_InsertNewUser(t *testing.T) {
	db, mock, _ := sqlmock.New()

	tests := []struct {
		name        string
		db          *sql.DB
		ctx         context.Context
		exec        string
		want        *models.NewUserResponse
		mockResErrs []error
	}{
		{
			name: "Happy Path",
			db:   db,
			ctx:  context.Background(),
			exec: fmt.Sprintf(AddUser, "13sdubf94", "Test Name", "test@email.com", "Test Username", "1208931bnd08128dn1", "1908wbhn190cb10cb1b0c", "19bc10cb10w8cb10w8cb", "11/10/2021", "11/10/2021", []string{}),
			want: &models.NewUserResponse{
				RowsAffected: int64(10),
			},
			mockResErrs: nil,
		},
		{
			name:        "Sad Path: validation error, missing exec",
			db:          db,
			ctx:         context.Background(),
			exec:        "",
			want:        nil,
			mockResErrs: []error{fmt.Errorf("missing query/statement")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				db: tt.db,
			}
			if tt.want != nil {
				mock.ExpectExec(regexp.QuoteMeta(tt.exec)).WillReturnResult(sqlmock.NewResult(tt.want.LastInsertedId, tt.want.RowsAffected))
			}
			got, got1 := s.InsertNewUser(tt.ctx, tt.exec)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertNewUser() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.mockResErrs) {
				t.Errorf("InsertNewUser() got1 = %v, want %v", got1, tt.mockResErrs)
			}
		})
	}
}

func TestService_UpdateAllTokens(t *testing.T) {
	db, mock, _ := sqlmock.New()
	tests := []struct {
		name               string
		db                 *sql.DB
		ctx                context.Context
		exec               string
		wantRowsAffected   int64
		wantLastInsertedId int64
		wantErr            bool
	}{
		{
			name:               "Happy Path",
			db:                 db,
			ctx:                context.Background(),
			exec:               fmt.Sprintf(UpdateTokens, "1208931bnd08128dn1", "1908wbhn190cb10cb1b0c", "11/10/2021", "12032302030asa"),
			wantRowsAffected:   int64(4),
			wantLastInsertedId: int64(12312123123),
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				db: tt.db,
			}
			mock.ExpectExec(regexp.QuoteMeta(tt.exec)).WillReturnResult(sqlmock.NewResult(tt.wantLastInsertedId, tt.wantRowsAffected))
			if err := s.InsertOne(tt.ctx, tt.exec); (err != nil) != tt.wantErr {
				t.Errorf("InsertOne() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_FindUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	tests := []struct {
		name        string
		db          *sql.DB
		ctx         context.Context
		query       string
		want        *models.UserParsedResponse
		mockResRows []string
		mockResErrs []error
		mockShows   []uint8
	}{
		{
			name:  "Happy Path",
			db:    db,
			ctx:   context.Background(),
			query: fmt.Sprintf(FindUserByEmail, "TestUsername"),
			want: &models.UserParsedResponse{
				ID:           "542113",
				FullName:     "Test User",
				Email:        "test@email.com",
				Username:     "testusername123",
				Password:     "password123",
				Token:        "39048567301249586",
				RefreshToken: "01938467501934651",
				Shows:        []string{},
				CreatedAt:    time.Now().Format(time.RFC3339),
				UpdatedAt:    time.Now().Format(time.RFC3339),
			},
			mockShows:   []uint8{},
			mockResRows: []string{"id", "fullname", "email", "username", "password", "token", "refreshtoken", "created", "updated", "shows"},
			mockResErrs: nil,
		},
		{
			name:        "Sad Path: validation error, missing query",
			db:          db,
			ctx:         context.Background(),
			query:       "",
			want:        nil,
			mockResRows: []string{},
			mockResErrs: []error{fmt.Errorf("missing query/statement")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				db: tt.db,
			}

			if tt.want != nil {
				mock.ExpectQuery(regexp.QuoteMeta(tt.query)).WillReturnRows(sqlmock.NewRows([]string{"id", "fullname", "email", "username", "password", "token", "refreshtoken", "created", "updated", "shows"}).AddRow(tt.want.ID, tt.want.FullName, tt.want.Email, tt.want.Username, tt.want.Password, tt.want.Token, tt.want.RefreshToken, tt.want.CreatedAt, tt.want.UpdatedAt, pq.Array(tt.want.Shows)))
			}
			got, got1 := s.FindUser(tt.ctx, tt.query)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindUser() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.mockResErrs) {
				t.Errorf("FindUser() got1 = %v, want %v", got1, tt.mockResErrs)
			}
		})
	}
}
