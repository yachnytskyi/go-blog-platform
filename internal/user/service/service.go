package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	"github.com/yachnytskyi/golang-mongo-grpc/models"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userService struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserService(collection *mongo.Collection, ctx context.Context) user.Service {
	return &userService{collection: collection, ctx: ctx}
}

func (userService *userService) Register(user *models.UserCreate) (*models.UserFullResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = true
	user.Role = "user"

	// Encrypt the provided password.
	user.Password, _ = utils.HashPassword(user.Password)
	userResponse, err := userService.collection.InsertOne(userService.ctx, &user)

	if err != nil {
		if err, ok := err.(mongo.WriteException); ok && err.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with this email already exists")
		}
		return nil, err
	}

	// Create a unique index for the email field.
	option := options.Index()
	option.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: option}

	if _, err := userService.collection.Indexes().CreateOne(userService.ctx, index); err != nil {
		return nil, errors.New("could not create an index for an email")
	}

	var newUser *models.UserFullResponse
	query := bson.M{"_id": userResponse.InsertedID}

	err = userService.collection.FindOne(userService.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (userService *userService) Login(*models.UserSignIn) (*models.UserFullResponse, error) {
	return nil, nil
}

func (userService *userService) UserGetById(userID string) (*models.UserFullResponse, error) {
	objectUserID, _ := primitive.ObjectIDFromHex(userID)

	var user *models.UserFullResponse

	query := bson.M{"_id": objectUserID}
	err := userService.collection.FindOne(userService.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.UserFullResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (userService *userService) UserGetByEmail(email string) (*models.UserFullResponse, error) {
	var user *models.UserFullResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := userService.collection.FindOne(userService.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.UserFullResponse{}, err
		}
		return nil, err
	}

	return user, nil
}
