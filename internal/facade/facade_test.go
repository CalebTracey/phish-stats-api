package facade

import (
	"context"
	"github.com/calebtracey/phish-stats-api/internal/models"
	"github.com/calebtracey/phish-stats-api/internal/services/phishnet"
	"github.com/calebtracey/phish-stats-api/internal/services/psql"
	"github.com/go-playground/validator/v10"
	"reflect"
	"testing"
)

func TestService_LoginUser(t *testing.T) {
	type fields struct {
		PsqlService     psql.ServiceI
		PhishNetService phishnet.ServiceI
		Validator       *validator.Validate
	}
	type args struct {
		ctx         context.Context
		userRequest models.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   models.UserResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				PsqlService:     tt.fields.PsqlService,
				PhishNetService: tt.fields.PhishNetService,
				Validator:       tt.fields.Validator,
			}
			if got := s.LoginUser(tt.args.ctx, tt.args.userRequest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoginUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
