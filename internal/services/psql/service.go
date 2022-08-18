package psql

import (
	"context"
	"database/sql"
	"fmt"
	config "github.com/calebtracey/config-yaml"
	"github.com/calebtracey/phish-stats-api/internal/models"
	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=mockService.go -package=psql . ServiceI
type ServiceI interface {
	FindUserByUsername(ctx context.Context, query string) (*models.UserPsqlResponse, []error)
	InsertNewUser(ctx context.Context, exec string) (*models.NewUserResponse, []error)
	UpdateAllTokens(ctx context.Context, exec string) error
}

type Service struct {
	db *sql.DB
}

func InitializePsqlService(psqlConfig *config.DatabaseConfig) *Service {
	return &Service{db: psqlConfig.DB}
}

func (s *Service) FindUserByUsername(ctx context.Context, query string) (*models.UserPsqlResponse, []error) {
	var usr models.UserPsqlResponse
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
			&usr.ID,
			&usr.FullName,
			&usr.Email,
			&usr.Username,
			&usr.Password,
			&usr.Token,
			&usr.CreatedAt,
			&usr.UpdatedAt,
			&usr.RefreshToken,
		)
		if err != nil {
			log.Panicln(err)
		}
	}

	return &usr, nil
}

func (s *Service) InsertNewUser(ctx context.Context, exec string) (*models.NewUserResponse, []error) {
	var errs []error
	vErrs := s.validateDbAction(exec)

	if vErrs != nil || len(vErrs) > 0 {
		return nil, vErrs
	}

	result, err := s.db.ExecContext(ctx, exec)
	if err != nil {
		return nil, []error{fmt.Errorf("error retrieving data; err: %v", err)}
	}

	lastInsertedId, idErr := result.LastInsertId()
	if idErr != nil {
		errs = append(errs, idErr)
	}
	rowsAffected, rowErr := result.RowsAffected()
	if rowErr != nil {
		errs = append(errs, rowErr)
	}

	return &models.NewUserResponse{
		LastInsertedId: lastInsertedId,
		RowsAffected:   rowsAffected,
	}, errs
}

func (s *Service) UpdateAllTokens(ctx context.Context, exec string) error {
	_, err := s.db.ExecContext(ctx, exec)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) validateDbAction(query string) []error {
	var errs []error
	if s.db == nil {
		errs = append(errs, fmt.Errorf("no database connection"))
	}
	if query == "" {
		errs = append(errs, fmt.Errorf("missing query/statement"))
	}
	return errs
}
