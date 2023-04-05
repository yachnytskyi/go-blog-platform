package repository

import (
	"context"
	"errors"
	"fmt"
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

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) user.Repository {
	return &UserRepository{collection: collection}
}

func (userRepository *UserRepository) Register(ctx context.Context, user *models.UserCreate) (*models.UserFullResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = true
	user.Role = "user"

	// Encrypt the provided password.
	user.Password, _ = utils.HashPassword(user.Password)
	userResponse, err := userRepository.collection.InsertOne(ctx, &user)

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

	if _, err := userRepository.collection.Indexes().CreateOne(ctx, index); err != nil {
		return nil, errors.New("could not create an index for an email")
	}

	var newUser *models.UserFullResponse
	query := bson.M{"_id": userResponse.InsertedID}

	err = userRepository.collection.FindOne(ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (userRepository *UserRepository) GetUserById(ctx context.Context, userID string) (*models.UserFullResponse, error) {
	objectUserID, _ := primitive.ObjectIDFromHex(userID)

	var user *models.UserFullResponse

	query := bson.M{"_id": objectUserID}
	err := userRepository.collection.FindOne(ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.UserFullResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (userRepository *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserFullResponse, error) {
	var user *models.UserFullResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := userRepository.collection.FindOne(ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.UserFullResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (userRepository *UserRepository) UpdateUserById(ctx context.Context, userID string, key string, value string) (*models.UserFullResponse, error) {
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	query := bson.D{{Key: "_id", Value: userObjectID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: value}}}}
	result, err := userRepository.collection.UpdateOne(ctx, query, update)

	fmt.Println(result.ModifiedCount)

	if err != nil {
		fmt.Println(err)
		return &models.UserFullResponse{}, err
	}

	return &models.UserFullResponse{}, nil
}
