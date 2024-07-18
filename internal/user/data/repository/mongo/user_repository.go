package repository

import (
	"context"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	userRepositoryMail "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/external/mail"
	userRepositoryModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo/model"
	repositoryUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/utility"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	userValidator "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	mongoModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logger"
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
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) UserRepository {
	repository := UserRepository{collection: database.Collection(constants.UsersTable)}
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimer)
	defer cancel()

	// Ensure the unique index on email during initialization.
	ensureUniqueEmailIndexError := repository.ensureUniqueEmailIndex(ctx, location+"NewUserRepository")
	if validator.IsError(ensureUniqueEmailIndexError) {
		panic(ensureUniqueEmailIndexError)
	}

	return repository
}

// GetAllUsers retrieves a list of users from the database based on pagination parameters.
func (userRepository UserRepository) GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users] {
	// Count the total number of users to set up pagination.
	query := bson.M{}
	totalUsers, countDocumentsError := userRepository.collection.CountDocuments(ctx, query)
	if validator.IsError(countDocumentsError) {
		internalError := domainError.NewInternalError(location+"GetAllUsers.collection.CountDocuments", countDocumentsError.Error())
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.Users](internalError)
	}

	// Set up pagination and sorting options using provided parameters.
	paginationQuery.TotalItems = int(totalUsers)
	paginationQuery = commonModel.SetCorrectPage(paginationQuery)
	option := options.FindOptions{}
	option.SetLimit(int64(paginationQuery.Limit))
	option.SetSkip(int64(paginationQuery.Skip))
	sortOptions := bson.M{paginationQuery.OrderBy: mongoModel.SetSortOrder(paginationQuery.SortOrder)}
	option.SetSort(sortOptions)

	// Query the database to fetch users.
	cursor, findError := userRepository.collection.Find(ctx, query, &option)
	if validator.IsError(findError) {
		queryString := commonUtility.ConvertQueryToString(query)
		itemNotFoundError := domainError.NewItemNotFoundError(location+"GetAllUsers.Find", queryString, findError.Error())
		logger.Logger(itemNotFoundError)
		return commonModel.NewResultOnFailure[userModel.Users](itemNotFoundError)
	}
	defer cursor.Close(ctx)

	// Process the results and map them to the repository model.
	fetchedUsers := make([]userRepositoryModel.UserRepository, 0, paginationQuery.Limit)
	for cursor.Next(ctx) {
		user := userRepositoryModel.UserRepository{}
		decodeError := cursor.Decode(&user)
		if validator.IsError(decodeError) {
			internalError := domainError.NewInternalError(location+"GetAllUsers.cursor.decode", decodeError.Error())
			logger.Logger(internalError)
			return commonModel.NewResultOnFailure[userModel.Users](internalError)
		}
		fetchedUsers = append(fetchedUsers, user)
	}

	cursorError := cursor.Err()
	if validator.IsError(cursorError) {
		internalError := domainError.NewInternalError(location+"GetAllUsers.cursor.Err", cursorError.Error())
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.Users](internalError)
	}

	if len(fetchedUsers) == 0 {
		return commonModel.NewResultOnSuccess[userModel.Users](userModel.Users{})
	}

	usersRepository := userRepositoryModel.UserRepositoryToUsersRepositoryMapper(fetchedUsers)
	usersRepository.PaginationResponse = commonModel.NewPaginationResponse(paginationQuery)
	return commonModel.NewResultOnSuccess[userModel.Users](userRepositoryModel.UsersRepositoryToUsersMapper(usersRepository))
}

// GetUserById retrieves a user by their ID from the database.
func (userRepository UserRepository) GetUserById(ctx context.Context, userID string) commonModel.Result[userModel.User] {
	userObjectID := mongoModel.HexToObjectIDMapper(location+"GetUserById", userID)
	if validator.IsError(userObjectID.Error) {
		return commonModel.NewResultOnFailure[userModel.User](userObjectID.Error)
	}

	query := bson.M{mongoModel.ID: userObjectID.Data}
	return userRepository.getUserByQuery(location+"GetUserById", ctx, query)
}

// GetUserByEmail retrieves a user by their email from the database.
func (userRepository UserRepository) GetUserByEmail(ctx context.Context, email string) commonModel.Result[userModel.User] {
	query := bson.M{emailKey: email}
	return userRepository.getUserByQuery(location+"GetUserByEmail", ctx, query)
}

// CheckEmailDuplicate checks if an email already exists in the database.
func (userRepository UserRepository) CheckEmailDuplicate(ctx context.Context, email string) error {
	fetchedUser := userRepositoryModel.UserRepository{}

	// Find and decode the user.
	// If no user is found, return nil (indicating that the email is unique).
	query := bson.M{emailKey: email}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsError(userFindOneError) {
		if userFindOneError == mongo.ErrNoDocuments {
			return nil
		}

		internalError := domainError.NewInternalError(location+"CheckEmailDuplicate.FindOne.Decode", userFindOneError.Error())
		logger.Logger(internalError)
		return internalError
	}

	// If a user with the given email is found, return a validation error.
	validationError := domainError.NewValidationError(
		location+"CheckEmailDuplicate",
		userValidator.EmailField,
		constants.FieldRequired,
		constants.EmailAlreadyExists,
	)

	logger.Logger(validationError)
	return validationError
}

// Register creates a user in the database based on the provided UserCreate data.
func (userRepository UserRepository) Register(ctx context.Context, userCreate userModel.UserCreate) commonModel.Result[userModel.User] {
	userCreateRepository := userRepositoryModel.UserCreateToUserCreateRepositoryMapper(userCreate)
	hashedPassword := repositoryUtility.HashPassword(location+"Register", userCreateRepository.Password)
	if validator.IsError(hashedPassword.Error) {
		return commonModel.NewResultOnFailure[userModel.User](hashedPassword.Error)
	}

	userCreateRepository.Password = hashedPassword.Data
	insertOneResult, insertOneResultError := userRepository.collection.InsertOne(ctx, &userCreateRepository)
	if validator.IsError(insertOneResultError) {
		internalError := domainError.NewInternalError(location+"Register.InsertOne", insertOneResultError.Error())
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.User](internalError)
	}

	query := bson.M{mongoModel.ID: insertOneResult.InsertedID}
	return userRepository.getUserByQuery(location+"Register", ctx, query)
}

// UpdateCurrentUser updates a user in the database based on the provided UserUpdate data.
func (userRepository UserRepository) UpdateCurrentUser(ctx context.Context, userUpdate userModel.UserUpdate) commonModel.Result[userModel.User] {
	userUpdateRepository := userRepositoryModel.UserUpdateToUserUpdateRepositoryMapper(location+"UpdateCurrentUser", userUpdate)
	if validator.IsError(userUpdateRepository.Error) {
		return commonModel.NewResultOnFailure[userModel.User](userUpdateRepository.Error)
	}

	userUpdateBSON := mongoModel.DataToMongoDocumentMapper(location+"UpdateCurrentUser", userUpdateRepository.Data)
	if validator.IsError(userUpdateBSON.Error) {
		return commonModel.NewResultOnFailure[userModel.User](userUpdateBSON.Error)
	}

	query := bson.D{{Key: mongoModel.ID, Value: userUpdateRepository.Data.UserID}}
	update := bson.D{{Key: mongoModel.Set, Value: userUpdateBSON.Data}}
	result := userRepository.collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))
	updatedUser := userRepositoryModel.UserRepository{}
	decodeError := result.Decode(&updatedUser)
	if validator.IsError(decodeError) {
		internalError := domainError.NewInternalError(location+"UpdateCurrentUser.Decode", decodeError.Error())
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.User](internalError)
	}

	return commonModel.NewResultOnSuccess[userModel.User](userRepositoryModel.UserRepositoryToUserMapper(updatedUser))
}

// DeleteUserById deletes a user in the database based on the provided userID.
func (userRepository UserRepository) DeleteUserById(ctx context.Context, userID string) error {
	userObjectID := mongoModel.HexToObjectIDMapper(location+"GetUserById", userID)
	if validator.IsError(userObjectID.Error) {
		return userObjectID.Error
	}

	query := bson.M{mongoModel.ID: userObjectID.Data}
	result, userDeleteOneError := userRepository.collection.DeleteOne(ctx, query)
	if validator.IsError(userDeleteOneError) {
		internalError := domainError.NewInternalError(location+"Delete.DeleteOne", userDeleteOneError.Error())
		logger.Logger(internalError)
		return internalError
	}

	if result.DeletedCount == 0 {
		internalError := domainError.NewInternalError(location+"Delete.DeleteOne.DeletedCount", mongoModel.DeletionIsNotSuccessful)
		logger.Logger(internalError)
		return internalError
	}

	return nil
}

// GetResetExpiry retrieves a reset token based on the provided reset token from the database.
func (userRepository UserRepository) GetResetExpiry(ctx context.Context, token string) commonModel.Result[userModel.UserResetExpiry] {
	fetchedResetExpiry := userRepositoryModel.UserResetExpiryRepository{}
	query := bson.M{resetTokenKey: token}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedResetExpiry)
	if validator.IsError(userFindOneError) {
		invalidTokenError := domainError.NewInvalidTokenError(location+"GetResetExpiry.Decode", userFindOneError.Error())
		logger.Logger(invalidTokenError)
		invalidTokenError.Notification = constants.InvalidTokenErrorMessage
		return commonModel.NewResultOnFailure[userModel.UserResetExpiry](invalidTokenError)
	}

	return commonModel.NewResultOnSuccess[userModel.UserResetExpiry](userRepositoryModel.UserResetExpiryRepositoryToUserResetExpiryMapper(fetchedResetExpiry))
}

// ForgottenPassword updates a user's record with a reset token and expiration time.
func (userRepository UserRepository) ForgottenPassword(ctx context.Context, userForgottenPassword userModel.UserForgottenPassword) error {
	userForgottenPasswordRepository := userRepositoryModel.UserForgottenPasswordToUserForgottenPasswordRepositoryMapper(userForgottenPassword)
	userForgottenPasswordBSON := mongoModel.DataToMongoDocumentMapper(location+"ForgottenPassword", userForgottenPasswordRepository)
	if validator.IsError(userForgottenPasswordBSON.Error) {
		return domainError.NewInternalError(location+"ForgottenPassword.Mapping", userForgottenPasswordBSON.Error.Error())
	}

	query := bson.D{{Key: emailKey, Value: userForgottenPassword.Email}}
	update := bson.D{{Key: mongoModel.Set, Value: userForgottenPasswordBSON.Data}}
	result, updateOneError := userRepository.collection.UpdateOne(ctx, query, update)
	if validator.IsError(updateOneError) {
		internalError := domainError.NewInternalError(location+"ForgottenPassword.UpdateOne", updateOneError.Error())
		logger.Logger(internalError)
		return internalError
	}

	if result.ModifiedCount == 0 {
		internalError := domainError.NewInternalError(location+"ForgottenPassword.UpdateOne.ModifiedCount", mongoModel.UpdateIsNotSuccessful)
		logger.Logger(internalError)
		return internalError
	}

	return nil
}

// ResetUserPassword updates a user's password based on the provided reset token and new password.
func (userRepository UserRepository) ResetUserPassword(ctx context.Context, userResetPassword userModel.UserResetPassword) error {
	userResetPasswordRepository := userRepositoryModel.UserResetPasswordToUserResetPasswordRepositoryMapper(userResetPassword)
	hashedPassword := repositoryUtility.HashPassword(location+"ResetUserPassword", userResetPassword.Password)
	if validator.IsError(hashedPassword.Error) {
		return hashedPassword.Error
	}

	userResetPasswordRepository.Password = hashedPassword.Data
	userResetPasswordBSON := mongoModel.DataToMongoDocumentMapper(location+"ResetUserPassword", userResetPasswordRepository)
	if validator.IsError(userResetPasswordBSON.Error) {
		return userResetPasswordBSON.Error
	}

	// Define the MongoDB query.
	// Define the update operation with the password update and the fields to unset.
	query := bson.D{{Key: resetTokenKey, Value: userResetPassword.ResetToken}}
	update := bson.D{
		{Key: mongoModel.Set, Value: userResetPasswordBSON.Data},
		{Key: mongoModel.Unset, Value: bson.D{
			{Key: resetTokenKey, Value: ""},
			{Key: resetExpiryKey, Value: ""},
		}},
	}

	result, updateOneError := userRepository.collection.UpdateOne(ctx, query, update)
	if validator.IsError(updateOneError) {
		internalError := domainError.NewInternalError(location+"ResetUserPassword.UpdateOne", updateOneError.Error())
		logger.Logger(internalError)
		return internalError
	}

	if result.ModifiedCount == 0 {
		internalError := domainError.NewInternalError(location+"ResetUserPassword.UpdateOne.ModifiedCount", mongoModel.UpdateIsNotSuccessful)
		logger.Logger(internalError)
		return internalError
	}

	return nil
}

// SendEmail sends an email to the specified user with the provided data.
func (userRepository UserRepository) SendEmail(user userModel.User, data userModel.EmailData) error {
	sendEmailError := userRepositoryMail.SendEmail(location+"SendEmail", user, data)
	if validator.IsError(sendEmailError) {
		return sendEmailError
	}

	return nil
}

// ensureUniqueEmailIndex creates a unique index on the email field to enforce email uniqueness in the database.
func (userRepository UserRepository) ensureUniqueEmailIndex(ctx context.Context, location string) error {
	option := options.Index()
	option.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{emailKey: 1}, Options: option}
	_, userIndexesCreateOneError := userRepository.collection.Indexes().CreateOne(ctx, index)
	if validator.IsError(userIndexesCreateOneError) {
		internalError := domainError.NewInternalError(location+".ensureUniqueEmailIndex.Indexes.CreateOne", userIndexesCreateOneError.Error())
		logger.Logger(internalError)
		return internalError
	}

	return nil
}

// getUserByQuery retrieves a user based on the provided query from the database.
func (userRepository UserRepository) getUserByQuery(location string, ctx context.Context, query bson.M) commonModel.Result[userModel.User] {
	fetchedUser := userRepositoryModel.UserRepository{}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsError(userFindOneError) {
		queryString := commonUtility.ConvertQueryToString(query)
		itemNotFoundError := domainError.NewItemNotFoundError(location+".getUserByQuery.Decode", queryString, userFindOneError.Error())
		logger.Logger(itemNotFoundError)
		return commonModel.NewResultOnFailure[userModel.User](itemNotFoundError)
	}

	return commonModel.NewResultOnSuccess[userModel.User](userRepositoryModel.UserRepositoryToUserMapper(fetchedUser))
}
