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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	location                = "user.data.repository.mongo."
	updateIsNotSuccessful   = "Update was not successful."
	deletionIsNotSuccessful = "Deletion was not successful."
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
	if validator.IsError(countDocumentsError) {
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
	if validator.IsError(cursorFindError) {
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
		if validator.IsError(cursorDecodeError) {
			internalError := domainError.NewInternalError(location+"GetAllUsers.cursor.decode", cursorDecodeError.Error())
			logging.Logger(internalError)
			return commonModel.NewResultOnFailure[userModel.Users](internalError)
		}
		fetchedUsers = append(fetchedUsers, user)
	}

	cursorError := cursor.Err()
	if validator.IsError(cursorError) {
		internalError := domainError.NewInternalError(location+"GetAllUsers.cursor.Err", cursorError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.Users](internalError)
	}
	if len(fetchedUsers) == 0 {
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
	// Check context timeout.
	contextError := commonUtility.HandleWithContextError("internal.user.data.repository.mongo.GetUserById", ctx)
	if validator.IsError(contextError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(contextError))
	}

	userObjectID, hexToObjectIDMapperError := mongoModel.HexToObjectIDMapper(location+"GetUserById", userID)
	if validator.IsError(hexToObjectIDMapperError) {
		return commonModel.NewResultOnFailure[userModel.User](hexToObjectIDMapperError)
	}

	// Define the MongoDB query to find the user by ObjectID.
	// Retrieve the user from the database.
	query := bson.M{"_id": userObjectID}
	return userRepository.getUserByQuery(location+"GetUserById", ctx, query)
}

// GetUserByEmail retrieves a user by their email from the repository.
func (userRepository UserRepository) GetUserByEmail(ctx context.Context, email string) commonModel.Result[userModel.User] {
	// Initialize a User object and define the MongoDB query to find the user by Email.
	query := bson.M{"email": email}

	// Retrieve the user from the database.
	return userRepository.getUserByQuery(location+"GetUserByEmail", ctx, query)
}

// CheckEmailDuplicate checks if an email already exists in the UserRepository.
// It returns an error if the email is already associated with a user, or nil if the email is unique.
func (userRepository UserRepository) CheckEmailDuplicate(ctx context.Context, email string) error {
	// Initialize a User object and define th MongoDB query to find the user by Email.
	fetchedUser := userRepositoryModel.UserRepository{}
	query := bson.M{"email": email}

	// Find and decode the user.
	// If no user is found, return nil (indicating that the email is unique).
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsValueEmpty(fetchedUser) {
		return nil
	}

	// If an error occurs during the database query, log it as an internal error.
	if validator.IsError(userFindOneError) {
		userFindOneInternalError := domainError.NewInternalError(location+"CheckEmailDublicate.FindOne.Decode", userFindOneError.Error())
		logging.Logger(userFindOneInternalError)
		return userFindOneInternalError
	}

	// If a user with the given email is found, return a validation error.
	userFindOneValidationError := domainError.NewValidationError(location+"CheckEmailDublicate", userValidator.EmailField, constants.FieldRequired, constants.EmailAlreadyExists)
	logging.Logger(userFindOneValidationError)
	return userFindOneValidationError
}

// Register creates a user in the repository based on the provided UserCreate data.
// It performs the several steps:
// 1. Maps the incoming data.
// 2. Hashes the password.
// 3. Inserts the user into the database by executing the MongoDB insert query.
// 4. Ensures the uniqueness of the email field by creating a unique index.
// 5. Retrieves the created user from the database, maps it back to the domain model, and returns the result.
func (userRepository UserRepository) Register(ctx context.Context, userCreate userModel.UserCreate) commonModel.Result[userModel.User] {
	// Map the incoming user data to the repository model.
	userCreateRepository := userRepositoryModel.UserCreateToUserCreateRepositoryMapper(userCreate)

	// Hash the user's password.
	hashedPassword, hashPasswordError := repositoryUtility.HashPassword(location+"Register", userCreate.Password)
	if validator.IsError(hashPasswordError) {
		return commonModel.NewResultOnFailure[userModel.User](hashPasswordError)
	}

	// Set the hashed password in the repository model.
	userCreateRepository.Password = hashedPassword

	// Insert the user data into the database.
	insertOneResult, insertOneResultError := userRepository.collection.InsertOne(ctx, &userCreateRepository)
	if validator.IsError(insertOneResultError) {
		internalError := domainError.NewInternalError(location+"Register.InsertOne", insertOneResultError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.User](internalError)
	}

	// Ensure uniqueness of the email field by creating a unique index.
	ensureUniqueEmailIndexError := userRepository.ensureUniqueEmailIndex(ctx)
	if validator.IsError(ensureUniqueEmailIndexError) {
		return commonModel.NewResultOnFailure[userModel.User](ensureUniqueEmailIndexError)
	}

	// Retrieve the created user from the database.
	query := bson.M{"_id": insertOneResult.InsertedID}
	createdUser := userRepository.getUserByQuery(location+"Register", ctx, query)
	return createdUser
}

// UpdateUserById updates a user in the repository based on the provided UserUpdate data.
// It performs the following steps:
// 1. Maps the incoming data.
// 2. Maps the repository model to a MongoDB model.
// 3. Updates the user in the database by executing the MongoDB update query.
// 4. Retrieves the updated user from the database, maps it back to the domain model, and returns the result.
func (userRepository UserRepository) UpdateCurrentUser(ctx context.Context, user userModel.UserUpdate) commonModel.Result[userModel.User] {
	// Map user update data to a repository model.
	userUpdateRepository, userUpdateError := userRepositoryModel.UserUpdateToUserUpdateRepositoryMapper(user)
	if validator.IsError(userUpdateError) {
		return commonModel.NewResultOnFailure[userModel.User](userUpdateError)
	}

	// Map repository model to a MongoDB model.
	userUpdateMongo, dataToMongoDocumentMapper := mongoModel.DataToMongoDocumentMapper(location+"UpdateCurrentUser", userUpdateRepository)
	if validator.IsError(dataToMongoDocumentMapper) {
		return commonModel.NewResultOnFailure[userModel.User](dataToMongoDocumentMapper)
	}

	// Define the MongoDB query.
	// Define the update operation.
	// Execute the update query and retrieve the updated user.
	query := bson.D{{Key: "_id", Value: userUpdateRepository.UserID}}
	update := bson.D{{Key: "$set", Value: userUpdateMongo}}
	result := userRepository.collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	// Decode the updated user from the result.
	updatedUserRepository := userRepositoryModel.UserRepository{}
	decodeError := result.Decode(&updatedUserRepository)
	if validator.IsError(decodeError) {
		userUpdateError := domainError.NewInternalError(location+"UpdateUserById.Decode", decodeError.Error())
		logging.Logger(userUpdateError)
		return commonModel.NewResultOnFailure[userModel.User](userUpdateError)
	}

	// Map the updated repository model to the user model.
	updatedUser := userRepositoryModel.UserRepositoryToUserMapper(updatedUserRepository)
	return commonModel.NewResultOnSuccess[userModel.User](updatedUser)
}

// DeleteUserById deletes a user in the repository based on the provided userID.
func (userRepository UserRepository) DeleteUserById(ctx context.Context, userID string) error {
	userObjectID, hexToObjectIDMapperError := mongoModel.HexToObjectIDMapper(location+"GetUserById", userID)
	if validator.IsError(hexToObjectIDMapperError) {
		return hexToObjectIDMapperError
	}

	// Define the MongoDB query to delete the user by ObjectID.
	// Execute the delete query.
	query := bson.M{"_id": userObjectID}
	result, userDeleteOneError := userRepository.collection.DeleteOne(ctx, query)
	if validator.IsError(userDeleteOneError) {
		deletedUserError := domainError.NewInternalError(location+"Delete.DeleteOne", userDeleteOneError.Error())
		logging.Logger(deletedUserError)
		return deletedUserError
	}

	// Check if any user was deleted
	if result.DeletedCount == 0 {
		deletedUserError := domainError.NewInternalError(location+"Delete.DeleteOne.DeletedCount", deletionIsNotSuccessful)
		logging.Logger(deletedUserError)
		return deletedUserError
	}
	return nil
}

func (userRepository UserRepository) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	// Hash the user's password.
	hashedPassword, hashPasswordError := repositoryUtility.HashPassword(location+"ResetUserPassword", password)
	if validator.IsError(hashPasswordError) {
		return hashPasswordError
	}

	query := bson.D{{Key: firstKey, Value: firstValue}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: passwordKey, Value: hashedPassword}}},
		{Key: "$unset", Value: bson.D{{Key: firstKey, Value: constants.EmptyString}, {Key: secondKey, Value: constants.EmptyString}}}}
	result, updateUserPasswordUpdateOneError := userRepository.collection.UpdateOne(ctx, query, update)
	if validator.IsError(updateUserPasswordUpdateOneError) {
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
	if validator.IsError(updateUserResetTokenUpdateOneError) {
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
	if validator.IsError(sendEmailError) {
		logging.Logger(sendEmailError)
		return sendEmailError
	}
	return nil
}

func (userRepository UserRepository) SendEmailForgottenPasswordMessage(ctx context.Context, user userModel.User, data userModel.EmailData) error {
	sendEmailError := userRepositoryMail.SendEmail(ctx, user, data)
	if validator.IsError(sendEmailError) {
		logging.Logger(sendEmailError)
		return sendEmailError
	}
	return nil
}

// ensureUniqueEmailIndex creates a unique index on the email field to enforce email uniqueness in the repository.
// This ensures that each user has a unique email address in the database.
func (userRepository UserRepository) ensureUniqueEmailIndex(ctx context.Context) error {
	// Create options for the index, setting it as unique.
	option := options.Index()
	option.SetUnique(true)

	// Define the index model based on the email field.
	// Create the unique index in the collection.
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: option}
	_, userIndexesCreateOneError := userRepository.collection.Indexes().CreateOne(ctx, index)
	if validator.IsError(userIndexesCreateOneError) {
		internalError := domainError.NewInternalError(location+"ensureUniqueEmailIndex.Indexes.CreateOne", userIndexesCreateOneError.Error())
		logging.Logger(internalError)
		return internalError
	}
	return nil
}

// getUserByQuery retrieves a user based on the provided query from the repository.
func (userRepository UserRepository) getUserByQuery(location string, ctx context.Context, query bson.M) commonModel.Result[userModel.User] {
	// Check context timeout.
	contextError := commonUtility.HandleWithContextError("internal.user.data.repository.mongo.GetUserById", ctx)
	if validator.IsError(contextError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(contextError))
	}

	// Initialize a User object and find the user based on the provided query.
	fetchedUser := userRepositoryModel.UserRepository{}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsError(userFindOneError) {
		queryString := commonUtility.ConvertQueryToString(query)
		entityNotFoundError := domainError.NewEntityNotFoundError(location+".getUserByQuery.Decode", queryString, userFindOneError.Error())
		logging.Logger(entityNotFoundError)
		return commonModel.NewResultOnFailure[userModel.User](entityNotFoundError)
	}

	// Map the repository model to domain ones.
	user := userRepositoryModel.UserRepositoryToUserMapper(fetchedUser)
	return commonModel.NewResultOnSuccess[userModel.User](user)
}
