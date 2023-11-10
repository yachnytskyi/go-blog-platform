package repository

import (
	"context"
	"strings"
	"time"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userRepositoryMail "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/external/mail"
	userRepositoryModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo/model"
	repositoryUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/utility"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	userValidator "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	mongoModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	location                = "user.data.repository.mongo."
	updateIsNotSuccessful   = "Update was not successful."
	delitionIsNotSuccessful = "Delition was not successful."
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) user.UserRepository {
	return UserRepository{collection: db.Collection("users")}
}

// GetAllUsers retrieves a list of users from the database based on pagination parameters.
func (userRepository UserRepository) GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users] {
	// Initialize the query with an empty BSON document
	// and determine the sort order based on the pagination query.
	query := bson.M{}
	sortOrder := mongoModel.SetSortOrder(paginationQuery.SortOrder)

	// Count the total number of users to set up pagination.
	totalUsers, countDocumentsError := userRepository.collection.CountDocuments(ctx, query)
	if validator.IsErrorNotNil(countDocumentsError) {
		internalError := domainError.NewInternalError(location+"GetAllUsers.collection.CountDocuments", countDocumentsError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.Users](internalError)
	}

	// Set up pagination and sorting options using provided parameters.
	paginationQuery = commonModel.SetCorrectPage(int(totalUsers), paginationQuery)
	option := options.FindOptions{}
	option.SetLimit(int64(paginationQuery.Limit))
	option.SetSkip(int64(paginationQuery.Skip))
	sortOptions := bson.M{paginationQuery.OrderBy: sortOrder}
	option.SetSort(sortOptions)

	// Query the database to fetch users.
	cursor, cursorFindError := userRepository.collection.Find(ctx, query, &option)
	if validator.IsErrorNotNil(cursorFindError) {
		queryString := commonUtility.ConvertQueryToString(query)
		entityNotFoundError := domainError.NewEntityNotFoundError(location+"GetAllUsers.Find", queryString, cursorFindError.Error())
		logging.Logger(entityNotFoundError)
		return commonModel.NewResultOnFailure[userModel.Users](entityNotFoundError)
	}
	defer cursor.Close(ctx)

	// Process the results and map them to the repository model.
	fetchedUsers := make([]userRepositoryModel.UserRepository, 0, paginationQuery.Limit)
	for cursor.Next(ctx) {
		user := userRepositoryModel.UserRepository{}
		cursorDecodeError := cursor.Decode(&user)
		if validator.IsErrorNotNil(cursorDecodeError) {
			internalError := domainError.NewInternalError(location+"GetAllUsers.cursor.decode", cursorDecodeError.Error())
			logging.Logger(internalError)
			return commonModel.NewResultOnFailure[userModel.Users](internalError)
		}
		fetchedUsers = append(fetchedUsers, user)
	}
	cursorError := cursor.Err()
	if validator.IsErrorNotNil(cursorError) {
		internalError := domainError.NewInternalError(location+"GetAllUsers.cursor.Err", cursorError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.Users](internalError)
	}
	if validator.IsSliceEmpty(fetchedUsers) {
		return commonModel.NewResultOnSuccess[userModel.Users](userModel.Users{})
	}

	// Map the repository model to domain ones.
	usersRepository := userRepositoryModel.UserRepositoryToUsersRepositoryMapper(fetchedUsers)
	paginationResponse := commonModel.NewPaginationResponse(paginationQuery.Page, int(totalUsers), paginationQuery.Limit, paginationQuery.OrderBy)
	usersRepository.PaginationResponse = paginationResponse
	users := userRepositoryModel.UsersRepositoryToUsersMapper(usersRepository)
	return commonModel.NewResultOnSuccess[userModel.Users](users)
}

// GetUserById retrieves a user by their ID from the database.
func (userRepository UserRepository) GetUserById(ctx context.Context, userID string) commonModel.Result[userModel.User] {
	userObjectID, objectIDFromHexError := primitive.ObjectIDFromHex(userID)
	if objectIDFromHexError != nil {
		internalError := domainError.NewInternalError(location+"GetUserById.ObjectIDFromHex", objectIDFromHexError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.User](internalError)
	}

	// Initialize a User object and define the query to find the user by ObjectID.
	fetchedUser := userRepositoryModel.UserRepository{}
	query := bson.M{"_id": userObjectID}

	// Find and decode the user.
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsErrorNotNil(userFindOneError) {
		queryString := commonUtility.ConvertQueryToString(query)
		entityNotFoundError := domainError.NewEntityNotFoundError(location+"GetUserById.FindOne.Decode", queryString, userFindOneError.Error())
		logging.Logger(entityNotFoundError)
		return commonModel.NewResultOnFailure[userModel.User](entityNotFoundError)
	}

	// Map the retrieved User to the UserModel and return a success result.
	user := userRepositoryModel.UserRepositoryToUserMapper(fetchedUser)
	return commonModel.NewResultOnSuccess[userModel.User](user)
}

// GetUserByEmail retrieves a user by their email from the repository.
func (userRepository UserRepository) GetUserByEmail(ctx context.Context, email string) commonModel.Result[userModel.User] {
	// Initialize an empty user from the repository model.
	fetchedUser := userRepositoryModel.UserRepository{}
	query := bson.M{"email": email}

	// Find and decode the user.
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsErrorNotNil(userFindOneError) {
		queryString := commonUtility.ConvertQueryToString(query)
		entityNotFoundError := domainError.NewEntityNotFoundError(location+"GetUserByEmail.FindOne.Decode", queryString, userFindOneError.Error())
		logging.Logger(entityNotFoundError)
		return commonModel.NewResultOnFailure[userModel.User](entityNotFoundError)
	}

	// Map the retrieved User to the UserModel and return a success result.
	user := userRepositoryModel.UserRepositoryToUserMapper(fetchedUser)
	return commonModel.NewResultOnSuccess[userModel.User](user)
}

func (userRepository UserRepository) CheckEmailDuplicate(ctx context.Context, email string) error {
	fetchedUser := userRepositoryModel.UserRepository{}
	query := bson.M{"email": email}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsValueNil(fetchedUser) {
		return nil
	}
	if validator.IsErrorNotNil(userFindOneError) {
		internalError := domainError.NewInternalError(location+"CheckEmailDublicate.FindOne.Decode", userFindOneError.Error())
		logging.Logger(internalError)
		return internalError
	}
	userFindOneValidationError := domainError.NewValidationError(location+"CheckEmailDublicate", userValidator.EmailField, constants.FieldRequired, constants.EmailAlreadyExists)
	logging.Logger(userFindOneValidationError)
	return userFindOneValidationError
}

func (userRepository UserRepository) Register(ctx context.Context, userCreate userModel.UserCreate) commonModel.Result[userModel.User] {
	// Map the incoming user data to the repository model.
	// Hash the user's password.
	userCreateRepository := userRepositoryModel.UserCreateToUserCreateRepositoryMapper(userCreate)
	userCreateRepository.Password, _ = repositoryUtility.HashPassword(userCreate.Password)
	userCreateRepository.CreatedAt = time.Now()
	userCreateRepository.UpdatedAt = userCreate.CreatedAt

	// Insert the user data into the database.
	insertOneResult, insertOneResultError := userRepository.collection.InsertOne(ctx, &userCreateRepository)
	if validator.IsErrorNotNil(insertOneResultError) {
		internalError := domainError.NewInternalError(location+"Register.InsertOne", insertOneResultError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.User](internalError)
	}

	// Create a unique index for the email field.
	option := options.Index()
	option.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: option}
	_, userIndexesCreateOneError := userRepository.collection.Indexes().CreateOne(ctx, index)
	if validator.IsErrorNotNil(userIndexesCreateOneError) {
		internalError := domainError.NewInternalError(location+"Register.Indexes.CreateOne", userIndexesCreateOneError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.User](internalError)
	}

	// Retrieve the created user from the database.
	createdUserRepository := userRepositoryModel.UserRepository{}
	query := bson.M{"_id": insertOneResult.InsertedID}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&createdUserRepository)
	if validator.IsErrorNotNil(userFindOneError) {
		queryString := commonUtility.ConvertQueryToString(query)
		entityNotFoundError := domainError.NewEntityNotFoundError(location+"Register.FindOne.Decode", queryString, userFindOneError.Error())
		logging.Logger(entityNotFoundError)
		return commonModel.NewResultOnFailure[userModel.User](entityNotFoundError)
	}

	// Map the retrieved user back to the domain model and return it.
	createdUser := userRepositoryModel.UserRepositoryToUserMapper(createdUserRepository)
	return commonModel.NewResultOnSuccess[userModel.User](createdUser)
}

func (userRepository UserRepository) UpdateUserById(ctx context.Context, user userModel.UserUpdate) commonModel.Result[userModel.User] {
	userUpdateRepository, userUpdateError := userRepositoryModel.UserUpdateToUserUpdateRepositoryMapper(user)
	if validator.IsErrorNotNil(userUpdateError) {
		return commonModel.NewResultOnFailure[userModel.User](userUpdateError)
	}
	userUpdateRepository.UpdatedAt = time.Now()

	userUpdateRepositoryMappedToMongoDB, mongoMapperError := mongoModel.MongoMapper(userUpdateRepository)
	if validator.IsErrorNotNil(mongoMapperError) {
		return commonModel.NewResultOnFailure[userModel.User](mongoMapperError)
	}

	query := bson.D{{Key: "_id", Value: userUpdateRepository.UserID}}
	update := bson.D{{Key: "$set", Value: userUpdateRepositoryMappedToMongoDB}}
	result := userRepository.collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))
	updatedUserRepository := userRepositoryModel.UserRepository{}
	updateUserRepositoryDecodeError := result.Decode(&updatedUserRepository)
	if validator.IsErrorNotNil(updateUserRepositoryDecodeError) {
		userUpdateError := domainError.NewInternalError(location+"UpdateUserById.Decode", updateUserRepositoryDecodeError.Error())
		logging.Logger(userUpdateError)
		return commonModel.NewResultOnFailure[userModel.User](userUpdateError)
	}

	updatedUser := userRepositoryModel.UserRepositoryToUserMapper(updatedUserRepository)
	return commonModel.NewResultOnSuccess[userModel.User](updatedUser)
}

func (userRepository UserRepository) DeleteUser(ctx context.Context, userID string) error {
	userIDMappedToMongoDB, _ := primitive.ObjectIDFromHex(userID)
	query := bson.M{"_id": userIDMappedToMongoDB}
	result, userDeleteOneError := userRepository.collection.DeleteOne(ctx, query)
	if validator.IsErrorNotNil(userDeleteOneError) {
		deletedUserError := domainError.NewInternalError(location+"Delete.DeleteOne", userDeleteOneError.Error())
		logging.Logger(deletedUserError)
		return deletedUserError
	}
	if result.DeletedCount == 0 {
		deletedUserError := domainError.NewInternalError(location+"Delete.DeleteOne.DeletedCount", delitionIsNotSuccessful)
		logging.Logger(deletedUserError)
		return deletedUserError
	}
	return nil
}

func (userRepository UserRepository) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	hashedPassword, _ := repositoryUtility.HashPassword(password)
	query := bson.D{{Key: firstKey, Value: firstValue}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: passwordKey, Value: hashedPassword}}}, {Key: "$unset", Value: bson.D{{Key: firstKey, Value: ""}, {Key: secondKey, Value: ""}}}}
	result, updateUserPasswordUpdateOneError := userRepository.collection.UpdateOne(ctx, query, update)
	if validator.IsErrorNotNil(updateUserPasswordUpdateOneError) {
		updatedUserPasswordError := domainError.NewInternalError(location+"ResetUserPassword.UpdateOne", updateUserPasswordUpdateOneError.Error())
		logging.Logger(updatedUserPasswordError)
		return updatedUserPasswordError
	}
	if result.ModifiedCount == 0 {
		updatedUserPasswordError := domainError.NewInternalError(location+"ResetUserPassword.UpdateOne.ModifiedCount", updateUserPasswordUpdateOneError.Error())
		logging.Logger(updatedUserPasswordError)
		return updatedUserPasswordError
	}
	return nil
}

func (userRepository UserRepository) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, SecondValue time.Time) error {

	query := bson.D{{Key: "email", Value: strings.ToLower(email)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: firstKey, Value: firstValue}, {Key: secondKey, Value: secondKey}}}}
	result, updateUserResetTokenUpdateOneError := userRepository.collection.UpdateOne(ctx, query, update)
	if validator.IsErrorNotNil(updateUserResetTokenUpdateOneError) {
		updatedUserResetTokenError := domainError.NewInternalError(location+"UpdatePasswordResetTokenUserByEmail.UpdateOne", updateUserResetTokenUpdateOneError.Error())
		logging.Logger(updatedUserResetTokenError)
		return updatedUserResetTokenError
	}
	if result.ModifiedCount == 0 {
		updatedUserResetTokenError := domainError.NewInternalError(location+"UpdatePasswordResetTokenUserByEmail.UpdateOne.ModifiedCount", updateIsNotSuccessful)
		logging.Logger(updatedUserResetTokenError)
		return updatedUserResetTokenError
	}
	return nil
}

func (userRepository UserRepository) SendEmailVerificationMessage(ctx context.Context, user userModel.User, data userModel.EmailData) error {
	sendEmailError := userRepositoryMail.SendEmail(ctx, user, data)
	if validator.IsErrorNotNil(sendEmailError) {
		logging.Logger(sendEmailError)
		return sendEmailError
	}
	return nil
}

func (userRepository UserRepository) SendEmailForgottenPasswordMessage(ctx context.Context, user userModel.User, data userModel.EmailData) error {
	sendEmailError := userRepositoryMail.SendEmail(ctx, user, data)
	if validator.IsErrorNotNil(sendEmailError) {
		logging.Logger(sendEmailError)
		return sendEmailError
	}
	return nil
}
