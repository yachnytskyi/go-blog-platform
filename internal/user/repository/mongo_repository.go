package repository

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

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) user.Repository {
	return &UserRepository{collection: collection}
}

func (userRepository *UserRepository) GetUserById(ctx context.Context, userID string) (*models.UserFullResponse, error) {
	objectUserID, _ := primitive.ObjectIDFromHex(userID)

	var fetchedUser *models.UserFullResponse

	query := bson.M{"_id": objectUserID}
	err := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.UserFullResponse{}, err
		}
		return nil, err
	}

	return fetchedUser, nil
}

func (userRepository *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserFullResponse, error) {
	var fetchedUser *models.UserFullResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.UserFullResponse{}, err
		}
		return nil, err
	}

	return fetchedUser, nil
}

func (userRepository *UserRepository) Register(ctx context.Context, user *models.UserCreate) (*models.UserFullResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = true
	user.Role = "user"
	user.Password, _ = utils.HashPassword(user.Password)

	mappedUser := models.UserCreateMongoDBMapping(user)

	result, err := userRepository.collection.InsertOne(ctx, &mappedUser)

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

	var createdUser *models.UserFullResponse
	query := bson.M{"_id": result.InsertedID}

	err = userRepository.collection.FindOne(ctx, query).Decode(&createdUser)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (userRepository *UserRepository) UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*models.UserFullResponse, error) {
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	query := bson.D{{Key: "_id", Value: userObjectID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: value}}}}
	result, err := userRepository.collection.UpdateOne(ctx, query, update)

	if err != nil {
		return &models.UserFullResponse{}, err
	}

	if result.ModifiedCount == 0 {
		return &models.UserFullResponse{}, err
	}

	return &models.UserFullResponse{}, nil
}

func (userRepository *UserRepository) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, SecondValue time.Time) error {

	query := bson.D{{Key: "email", Value: strings.ToLower(email)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: firstKey, Value: firstValue}, {Key: secondKey, Value: secondKey}}}}
	result, err := userRepository.collection.UpdateOne(ctx, query, update)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return err
	}

	return nil
}

func (userRepository *UserRepository) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	hashedPassword, _ := utils.HashPassword(password)

	query := bson.D{{Key: firstKey, Value: firstValue}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: passwordKey, Value: hashedPassword}}}, {Key: "$unset", Value: bson.D{{Key: firstKey, Value: ""}, {Key: secondKey, Value: ""}}}}

	result, err := userRepository.collection.UpdateOne(ctx, query, update)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return err
	}

	return nil
}

func (userRepository *UserRepository) UpdateUserById(ctx context.Context, userID string, user *models.UserUpdate) (*models.UserResponse, error) {
	user.UpdatedAt = time.Now()

	mappedUser, err := utils.MongoMapping(user)

	if err != nil {
		return &models.UserResponse{}, err
	}

	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	query := bson.D{{Key: "_id", Value: userObjectID}}
	update := bson.D{{Key: "$set", Value: mappedUser}}
	result := userRepository.collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedUser *models.UserResponse

	if err := result.Decode(&updatedUser); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (userRepository *UserRepository) DeleteUserById(ctx context.Context, userID string) error {
	userObjectID, _ := primitive.ObjectIDFromHex(userID)
	query := bson.M{"_id": userObjectID}

	result, err := userRepository.collection.DeleteOne(ctx, query)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return err
	}

	return nil
}
