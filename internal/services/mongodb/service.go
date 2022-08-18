package mongodb

import (
	"context"
	"fmt"
	config "github.com/calebtracey/config-yaml"
	"github.com/calebtracey/phish-stats-api/internal/models"
	"github.com/calebtracey/phish-stats-api/internal/services/auth"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//go:generate mockgen -destination=mockService.go -package=mongodb . ServiceI
type ServiceI interface {
	AllUsersByEmail(ctx context.Context, value *string) ([]*models.User, error)
	AllUsersByPhone(ctx context.Context, value *string) ([]*models.User, error)
	AllUsersByUsername(ctx context.Context, value *string) ([]*models.User, error)
	UpdateTokensById(ctx context.Context, userId string, updateObj primitive.D) error
	AddNewUser(ctx context.Context, user *models.User) (*models.User, error)
	FindUserByUsername(ctx context.Context, user *models.User) (models.User, error)
	UpdateAllTokens(ctx context.Context, token string, signedRefreshToken string, userId string) error
}

type Service struct {
	Database string
	Client   *mongo.Client
	//Mapper   Mapper
}

func InitializeMongoService(appConfig *config.Config) (*Service, error) {
	mongoConfig, err := appConfig.GetDatabaseConfig("MONGO")
	if err != nil {
		return nil, err
	}
	return &Service{
		Database: mongoConfig.Database.Value,
		Client:   mongoConfig.MongoClient,
	}, nil
}

func (s *Service) UpdateTokensById(ctx context.Context, userId string, updateObj primitive.D) error {
	dbName := s.Database
	database := s.Client.Database(dbName)
	var err error
	upsert := true
	filter := bson.M{"userId": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err = database.Collection("users").UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)
	return err
}

func (s *Service) AddNewUser(ctx context.Context, user *models.User) (*models.User, error) {
	var response *models.User

	dbName := s.Database
	database := s.Client.Database(dbName)

	emails, err := database.Collection("users").Find(ctx, bson.M{"email": user.Email})
	if err != nil {
		return response, err
	}
	if emails.RemainingBatchLength() > 0 {
		return response, fmt.Errorf("email: %v already exists in the database", *user.Email)
	}
	//
	//phoneNumbers, err := database.Collection("users").Find(ctx, bson.M{"phone": user.Phone})
	//if err != nil {
	//	return response, err
	//}
	//if phoneNumbers.RemainingBatchLength() > 0 {
	//	return response, fmt.Errorf("phone number: %v already exists in the database", *user.Phone)
	//}

	usernames, err := database.Collection("users").Find(ctx, bson.M{"username": user.Username})
	if err != nil {
		return response, err
	}
	if usernames.RemainingBatchLength() > 0 {
		return response, fmt.Errorf("username: %v already exists in the database", *user.Username)
	}

	password := auth.HashPassword(*user.Password)
	user.Password = &password

	_, err = database.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return response, err
	}
	log.Infof("inserted new user: %v into database", *user.Username)

	response = user

	return response, nil
}

func (s *Service) FindUserByUsername(ctx context.Context, user *models.User) (models.User, error) {
	var foundUser models.User

	dbName := s.Database
	database := s.Client.Database(dbName)

	err := database.Collection("users").FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
	if err != nil {
		return foundUser, fmt.Errorf("login or passowrd is incorrect; err: %v", err.Error())
	}

	return foundUser, nil
}

func (s *Service) AllUsersByEmail(ctx context.Context, value *string) ([]*models.User, error) {
	dbName := s.Database
	database := s.Client.Database(dbName)
	var results []*models.User
	var err error

	cursor, err := database.Collection("users").Find(ctx, bson.M{"email": value})
	if err != nil {
		return results, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			err = fmt.Errorf("failed to close mongodb cursor; err: %v", err.Error())
		}
	}(cursor, ctx)

	if curErr := cursor.All(ctx, &results); curErr != nil {
		return nil, curErr
	}

	return results, nil
}

func (s *Service) AllUsersByPhone(ctx context.Context, value *string) ([]*models.User, error) {
	dbName := s.Database
	database := s.Client.Database(dbName)
	var results []*models.User
	var err error

	cursor, err := database.Collection("users").Find(ctx, bson.M{"phone": value})
	if err != nil {
		return results, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			err = fmt.Errorf("failed to close mongodb cursor; err: %v", err.Error())
		}
	}(cursor, ctx)

	if curErr := cursor.All(ctx, &results); curErr != nil {
		return nil, curErr
	}

	return results, nil
}

func (s *Service) AllUsersByUsername(ctx context.Context, value *string) ([]*models.User, error) {
	dbName := s.Database
	database := s.Client.Database(dbName)
	var results []*models.User
	var err error

	cursor, err := database.Collection("users").Find(ctx, bson.M{"username": value})
	if err != nil {
		return results, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			err = fmt.Errorf("failed to close mongodb cursor; err: %v", err.Error())
		}
	}(cursor, ctx)

	if curErr := cursor.All(ctx, &results); curErr != nil {
		return nil, curErr
	}

	return results, nil
}

func (s *Service) UpdateAllTokens(ctx context.Context, token string, signedRefreshToken string, userId string) error {

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: token})
	updateObj = append(updateObj, bson.E{Key: "refreshToken", Value: signedRefreshToken})

	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: UpdatedAt})

	err := s.UpdateTokensById(ctx, userId, updateObj)

	if err != nil {
		return err
	}

	return nil
}
