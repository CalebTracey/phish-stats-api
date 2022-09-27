package psql

import (
	"context"
	"database/sql"
	"fmt"
	config "github.com/calebtracey/config-yaml"
	"github.com/calebtracey/phish-stats-api/internal/models"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=mockService.go -package=psql . ServiceI
type ServiceI interface {
	FindUser(ctx context.Context, query string) (*models.UserParsedResponse, []error)
	InsertNewUser(ctx context.Context, exec string) (*models.NewUserResponse, []error)
	InsertOne(ctx context.Context, exec string) error
}

type Service struct {
	db *sql.DB
}

func InitializePsqlService(psqlConfig *config.DatabaseConfig) *Service {
	return &Service{db: psqlConfig.DB}
}

func (s Service) FindUser(ctx context.Context, query string) (*models.UserParsedResponse, []error) {
	var user models.UserParsedResponse
	errs := s.validateDbAction(query)

	if errs != nil || len(errs) > 0 {
		return nil, errs
	}

	rows, err := s.db.QueryContext(ctx, query)
	if err == sql.ErrNoRows {
		return nil, []error{fmt.Errorf("username does not exist in the database")}
	}
	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.Token,
			&user.RefreshToken,
			&user.CreatedAt,
			&user.UpdatedAt,
			pq.Array(&user.Shows),
		)
		if err != nil {
			log.Panicln(err)
		}
	}

	return &user, nil
}

func (s Service) InsertNewUser(ctx context.Context, exec string) (*models.NewUserResponse, []error) {
	var errs []error
	vErrs := s.validateDbAction(exec)

	if vErrs != nil || len(vErrs) > 0 {
		return nil, vErrs
	}

	result, err := s.db.ExecContext(ctx, exec)
	if err != nil {
		return nil, []error{fmt.Errorf("error retrieving data; err: %v", err)}
	}

	rowsAffected, rowErr := result.RowsAffected()
	if rowErr != nil {
		errs = append(errs, rowErr)
	}

	return &models.NewUserResponse{
		RowsAffected: rowsAffected,
	}, errs
}

func (s Service) InsertOne(ctx context.Context, exec string) error {
	_, err := s.db.ExecContext(ctx, exec)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) validateDbAction(query string) []error {
	var errs []error
	if s.db == nil {
		errs = append(errs, fmt.Errorf("no database connection"))
	}
	if query == "" {
		errs = append(errs, fmt.Errorf("missing query/statement"))
	}
	return errs
}
