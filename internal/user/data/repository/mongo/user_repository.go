package repository

import (
	"context"
	"time"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	repository "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo/model"
	repositoryUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/utility"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	useCase "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/data/repository/mongo"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	location       = "user.data.repository.mongo."
	emailKey       = "email"
	passwordKey    = "password"
	resetTokenKey  = "reset_token"
	resetExpiryKey = "reset_expiry"

	invalidEmailOrPassword = "Invalid email or password."
	emailOrPasswordFields  = "email or password"
	passwordsDoNotMatch    = "Passwords do not match."
)

type UserRepository struct {
	Config *config.ApplicationConfig
	Logger interfaces.Logger
	Users  *mongo.Collection
}

func NewUserRepository(config *config.ApplicationConfig, logger interfaces.Logger, database *mongo.Database) UserRepository {
	repository := UserRepository{
		Config: config,
		Logger: logger,
		Users:  database.Collection(constants.UsersTable),
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimer)
	defer cancel()

	// Ensure the unique index on email during initialization.
	ensureUniqueEmailIndexError := repository.ensureUniqueEmailIndex(ctx, location+"NewUserRepository")
	if validator.IsError(ensureUniqueEmailIndexError) {
		logger.Panic(domain.NewInternalError(location+"GetAllUsers.Users.CountDocuments", ensureUniqueEmailIndexError.Error()))
	}

	return repository
}

// GetAllUsers retrieves a list of users from the database based on pagination parameters.
func (userRepository UserRepository) GetAllUsers(ctx context.Context, paginationQuery common.PaginationQuery) common.Result[user.Users] {
	// Count the total number of users to set up pagination.
	query := bson.M{}
	totalUsers, countDocumentsError := userRepository.Users.CountDocuments(ctx, query)
	if validator.IsError(countDocumentsError) {
		internalError := domain.NewInternalError(location+"GetAllUsers.Users.CountDocuments", countDocumentsError.Error())
		userRepository.Logger.Error(internalError)
		return common.NewResultOnFailure[user.Users](internalError)
	}

	// Set up pagination and sorting options using provided parameters.
	paginationQuery.TotalItems = int(totalUsers)
	paginationQuery = common.SetCorrectPage(paginationQuery)
	option := options.FindOptions{}
	option.SetLimit(int64(paginationQuery.Limit))
	option.SetSkip(int64(paginationQuery.Skip))
	sortOptions := bson.M{paginationQuery.OrderBy: utility.SetSortOrder(paginationQuery.SortOrder)}
	option.SetSort(sortOptions)

	// Query the database to fetch users.
	cursor, findError := userRepository.Users.Find(ctx, query, &option)
	if validator.IsError(findError) {
		itemNotFoundError := domain.NewItemNotFoundError(location+"GetAllUsers.Find", utility.BSONToStringMapper(query), findError.Error())
		userRepository.Logger.Error(itemNotFoundError)
		return common.NewResultOnFailure[user.Users](itemNotFoundError)
	}
	defer cursor.Close(ctx)

	// Process the results and map them to the repository model.
	fetchedUsers := make([]repository.UserRepository, 0, paginationQuery.Limit)
	for cursor.Next(ctx) {
		userInstance := repository.UserRepository{}
		decodeError := cursor.Decode(&userInstance)
		if validator.IsError(decodeError) {
			internalError := domain.NewInternalError(location+"GetAllUsers.cursor.decode", decodeError.Error())
			userRepository.Logger.Error(internalError)
			return common.NewResultOnFailure[user.Users](internalError)
		}
		fetchedUsers = append(fetchedUsers, userInstance)
	}

	cursorError := cursor.Err()
	if validator.IsError(cursorError) {
		internalError := domain.NewInternalError(location+"GetAllUsers.cursor.Err", cursorError.Error())
		userRepository.Logger.Error(internalError)
		return common.NewResultOnFailure[user.Users](internalError)
	}

	usersRepository := repository.UserRepositoryToUsersRepositoryMapper(fetchedUsers)
	usersRepository.PaginationResponse = common.NewPaginationResponse(paginationQuery)
	return common.NewResultOnSuccess[user.Users](repository.UsersRepositoryToUsersMapper(usersRepository))
}

// GetUserById retrieves a user by their ID from the database.
func (userRepository UserRepository) GetUserById(ctx context.Context, userID string) common.Result[user.User] {
	userObjectID := model.HexToObjectIDMapper(userRepository.Logger, location+"GetUserById", userID)
	if validator.IsError(userObjectID.Error) {
		return common.NewResultOnFailure[user.User](userObjectID.Error)
	}

	query := bson.M{model.ID: userObjectID.Data}
	return userRepository.getUserByQuery(location+"GetUserById", ctx, query)
}

// GetUserByEmail retrieves a user by their email from the database.
func (userRepository UserRepository) GetUserByEmail(ctx context.Context, email string) common.Result[user.User] {
	fetchedUser := repository.UserRepository{}

	query := bson.M{emailKey: email}
	userFindOneError := userRepository.Users.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsError(userFindOneError) {
		if repositoryUtility.IsMongoDBError(userFindOneError) {
			internalError := domain.NewInternalError(location+"GetUserByEmail.FindOne.Decode", userFindOneError.Error())
			userRepository.Logger.Error(internalError)
			return common.NewResultOnFailure[user.User](internalError)
		}
		userRepository.Logger.Error(domain.NewItemNotFoundError(location+"GetUserByEmail.FindOne.Decode", utility.BSONToStringMapper(query), userFindOneError.Error()))
		validationError := domain.NewValidationError(location+".checkPasswords.CompareHashAndPassword", emailOrPasswordFields, constants.FieldRequired, passwordsDoNotMatch)
		validationError.Notification = invalidEmailOrPassword
		return common.NewResultOnFailure[user.User](validationError)
	}

	return common.NewResultOnSuccess[user.User](repository.UserRepositoryToUserMapper(fetchedUser))
}

// CheckEmailDuplicate checks if an email already exists in the database.
func (userRepository UserRepository) CheckEmailDuplicate(ctx context.Context, email string) error {
	fetchedUser := repository.UserRepository{}

	// Find and decode the user.
	// If no user is found, return nil (indicating that the email is unique).
	query := bson.M{emailKey: email}

	userFindOneError := userRepository.Users.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsError(userFindOneError) {
		if userFindOneError == mongo.ErrNoDocuments {
			return nil
		}

		internalError := domain.NewInternalError(location+"CheckEmailDuplicate.FindOne.Decode", userFindOneError.Error())
		userRepository.Logger.Error(internalError)
		return internalError
	}

	// If a user with the given email is found, return a validation error.
	validationError := domain.NewValidationError(location+"CheckEmailDuplicate", useCase.EmailField, constants.FieldRequired, constants.EmailAlreadyExists)

	userRepository.Logger.Error(validationError)
	return validationError
}

// Register creates a user in the database based on the provided UserCreate data.
func (userRepository UserRepository) Register(ctx context.Context, userCreate user.UserCreate) common.Result[user.User] {
	userCreateRepository := repository.UserCreateToUserCreateRepositoryMapper(userCreate)
	hashedPassword := repositoryUtility.HashPassword(userRepository.Logger, location+"Register", userCreateRepository.Password)
	if validator.IsError(hashedPassword.Error) {
		return common.NewResultOnFailure[user.User](hashedPassword.Error)
	}

	userCreateRepository.Password = hashedPassword.Data
	userCreateRepository.CreatedAt = time.Now()
	userCreateRepository.UpdatedAt = time.Now()
	insertOneResult, insertOneResultError := userRepository.Users.InsertOne(ctx, &userCreateRepository)
	if validator.IsError(insertOneResultError) {
		internalError := domain.NewInternalError(location+"Register.InsertOne", insertOneResultError.Error())
		userRepository.Logger.Error(internalError)
		return common.NewResultOnFailure[user.User](internalError)
	}

	query := bson.M{model.ID: insertOneResult.InsertedID}
	return userRepository.getUserByQuery(location+"Register", ctx, query)
}

// UpdateCurrentUser updates a user in the database based on the provided UserUpdate data.
func (userRepository UserRepository) UpdateCurrentUser(ctx context.Context, userUpdate user.UserUpdate) common.Result[user.User] {
	userUpdateRepository := repository.UserUpdateToUserUpdateRepositoryMapper(userRepository.Logger, location+"UpdateCurrentUser", userUpdate)
	if validator.IsError(userUpdateRepository.Error) {
		return common.NewResultOnFailure[user.User](userUpdateRepository.Error)
	}

	userUpdateRepository.Data.UpdatedAt = time.Now()
	userUpdateBSON := model.DataToMongoDocumentMapper(userRepository.Logger, location+"UpdateCurrentUser", userUpdateRepository.Data)
	if validator.IsError(userUpdateBSON.Error) {
		return common.NewResultOnFailure[user.User](userUpdateBSON.Error)
	}

	query := bson.D{{Key: model.ID, Value: userUpdateRepository.Data.UserID}}
	update := bson.D{{Key: model.Set, Value: userUpdateBSON.Data}}
	result := userRepository.Users.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))
	updatedUser := repository.UserRepository{}
	decodeError := result.Decode(&updatedUser)
	if validator.IsError(decodeError) {
		internalError := domain.NewInternalError(location+"UpdateCurrentUser.Decode", decodeError.Error())
		userRepository.Logger.Error(internalError)
		return common.NewResultOnFailure[user.User](internalError)
	}

	return common.NewResultOnSuccess[user.User](repository.UserRepositoryToUserMapper(updatedUser))
}

// DeleteUserById deletes a user in the database based on the provided userID.
func (userRepository UserRepository) DeleteUserById(ctx context.Context, userID string) error {
	userObjectID := model.HexToObjectIDMapper(userRepository.Logger, location+"GetUserById", userID)
	if validator.IsError(userObjectID.Error) {
		return userObjectID.Error
	}

	query := bson.M{model.ID: userObjectID.Data}
	result, userDeleteOneError := userRepository.Users.DeleteOne(ctx, query)
	if validator.IsError(userDeleteOneError) {
		internalError := domain.NewInternalError(location+"Delete.DeleteOne", userDeleteOneError.Error())
		userRepository.Logger.Error(internalError)
		return internalError
	}

	if result.DeletedCount == 0 {
		internalError := domain.NewInternalError(location+"Delete.DeleteOne.DeletedCount", model.DeletionIsNotSuccessful)
		userRepository.Logger.Error(internalError)
		return internalError
	}

	return nil
}

// GetResetExpiry retrieves a reset token based on the provided reset token from the database.
func (userRepository UserRepository) GetResetExpiry(ctx context.Context, token string) common.Result[user.UserResetExpiry] {
	fetchedResetExpiry := repository.UserResetExpiryRepository{}
	query := bson.M{resetTokenKey: token}

	userFindOneError := userRepository.Users.FindOne(ctx, query).Decode(&fetchedResetExpiry)
	if validator.IsError(userFindOneError) {
		if repositoryUtility.IsMongoDBError(userFindOneError) {
			internalError := domain.NewInternalError(location+"GetResetExpiry.FindOne.Decode", userFindOneError.Error())
			userRepository.Logger.Error(internalError)
			return common.NewResultOnFailure[user.UserResetExpiry](internalError)
		}
		invalidTokenError := domain.NewInvalidTokenError(location+"GetResetExpiry.Decode", userFindOneError.Error())
		userRepository.Logger.Error(invalidTokenError)
		invalidTokenError.Notification = constants.InvalidTokenErrorMessage
		return common.NewResultOnFailure[user.UserResetExpiry](invalidTokenError)
	}

	return common.NewResultOnSuccess[user.UserResetExpiry](repository.UserResetExpiryRepositoryToUserResetExpiryMapper(fetchedResetExpiry))
}

// ForgottenPassword updates a user's record with a reset token and expiration time.
func (userRepository UserRepository) ForgottenPassword(ctx context.Context, userForgottenPassword user.UserForgottenPassword) error {
	userForgottenPasswordRepository := repository.UserForgottenPasswordToUserForgottenPasswordRepositoryMapper(userForgottenPassword)
	userForgottenPasswordBSON := model.DataToMongoDocumentMapper(userRepository.Logger, location+"ForgottenPassword", userForgottenPasswordRepository)
	if validator.IsError(userForgottenPasswordBSON.Error) {
		return domain.NewInternalError(location+"ForgottenPassword.Mapping", userForgottenPasswordBSON.Error.Error())
	}

	query := bson.D{{Key: emailKey, Value: userForgottenPassword.Email}}
	update := bson.D{{Key: model.Set, Value: userForgottenPasswordBSON.Data}}
	result, updateOneError := userRepository.Users.UpdateOne(ctx, query, update)
	if validator.IsError(updateOneError) {
		internalError := domain.NewInternalError(location+"ForgottenPassword.UpdateOne", updateOneError.Error())
		userRepository.Logger.Error(internalError)
		return internalError
	}

	if result.ModifiedCount == 0 {
		internalError := domain.NewInternalError(location+"ForgottenPassword.UpdateOne.ModifiedCount", model.UpdateIsNotSuccessful)
		userRepository.Logger.Error(internalError)
		return internalError
	}

	return nil
}

// ResetUserPassword updates a user's password based on the provided reset token and new password.
func (userRepository UserRepository) ResetUserPassword(ctx context.Context, userResetPassword user.UserResetPassword) error {
	userResetPasswordRepository := repository.UserResetPasswordToUserResetPasswordRepositoryMapper(userResetPassword)
	hashedPassword := repositoryUtility.HashPassword(userRepository.Logger, location+"ResetUserPassword", userResetPassword.Password)
	if validator.IsError(hashedPassword.Error) {
		return hashedPassword.Error
	}

	userResetPasswordRepository.Password = hashedPassword.Data
	userResetPasswordBSON := model.DataToMongoDocumentMapper(userRepository.Logger, location+"ResetUserPassword", userResetPasswordRepository)
	if validator.IsError(userResetPasswordBSON.Error) {
		return userResetPasswordBSON.Error
	}

	// Define the MongoDB query.
	// Define the update operation with the password update and the fields to unset.
	query := bson.D{{Key: resetTokenKey, Value: userResetPassword.ResetToken}}
	update := bson.D{
		{Key: model.Set, Value: userResetPasswordBSON.Data},
		{Key: model.Unset, Value: bson.D{
			{Key: resetTokenKey, Value: ""},
			{Key: resetExpiryKey, Value: ""},
		}},
	}

	result, updateOneError := userRepository.Users.UpdateOne(ctx, query, update)
	if validator.IsError(updateOneError) {
		internalError := domain.NewInternalError(location+"ResetUserPassword.UpdateOne", updateOneError.Error())
		userRepository.Logger.Error(internalError)
		return internalError
	}

	if result.ModifiedCount == 0 {
		internalError := domain.NewInternalError(location+"ResetUserPassword.UpdateOne.ModifiedCount", model.UpdateIsNotSuccessful)
		userRepository.Logger.Error(internalError)
		return internalError
	}

	return nil
}

// ensureUniqueEmailIndex creates a unique index on the email field to enforce email uniqueness in the database.
func (userRepository UserRepository) ensureUniqueEmailIndex(ctx context.Context, location string) error {
	option := options.Index()
	option.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{emailKey: 1}, Options: option}
	_, userIndexesCreateOneError := userRepository.Users.Indexes().CreateOne(ctx, index)
	if validator.IsError(userIndexesCreateOneError) {
		internalError := domain.NewInternalError(location+".ensureUniqueEmailIndex.Indexes.CreateOne", userIndexesCreateOneError.Error())
		userRepository.Logger.Error(internalError)
		return internalError
	}

	return nil
}

// getUserByQuery retrieves a user based on the provided query from the database.
func (userRepository UserRepository) getUserByQuery(location string, ctx context.Context, query bson.M) common.Result[user.User] {
	fetchedUser := repository.UserRepository{}
	userFindOneError := userRepository.Users.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsError(userFindOneError) {
		if repositoryUtility.IsMongoDBError(userFindOneError) {
			internalError := domain.NewInternalError(location+".getUserByQuery.FindOne.Decode", userFindOneError.Error())
			userRepository.Logger.Error(internalError)
			return common.NewResultOnFailure[user.User](internalError)
		}
		itemNotFoundError := domain.NewItemNotFoundError(location+".getUserByQuery.FindOne.Decode", utility.BSONToStringMapper(query), userFindOneError.Error())
		userRepository.Logger.Error(itemNotFoundError)
		return common.NewResultOnFailure[user.User](itemNotFoundError)
	}

	return common.NewResultOnSuccess[user.User](repository.UserRepositoryToUserMapper(fetchedUser))
}
