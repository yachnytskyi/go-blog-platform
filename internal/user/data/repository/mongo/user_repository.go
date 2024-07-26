package repository

import (
	"context"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userRepositoryMail "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/external/mail"
	userRepositoryModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo/model"
	repositoryUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/utility"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	usecase "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	mongoModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
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
	Config     model.Config
	Logger     model.Logger
	Collection *mongo.Collection
}

func NewUserRepository(config model.Config, logger model.Logger, database *mongo.Database) user.UserRepository {
	repository := &UserRepository{
		Config:     config,
		Logger:     logger,
		Collection: database.Collection(constants.UsersTable),
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimer)
	defer cancel()

	// Ensure the unique index on email during initialization.
	ensureUniqueEmailIndexError := repository.ensureUniqueEmailIndex(ctx, location+"NewUserRepository")
	if validator.IsError(ensureUniqueEmailIndexError) {
		logger.Panic(ensureUniqueEmailIndexError)
	}

	return repository
}

// GetAllUsers retrieves a list of users from the database based on pagination parameters.
func (userRepository *UserRepository) GetAllUsers(ctx context.Context, paginationQuery *common.PaginationQuery) common.Result[*userModel.Users] {
	// Count the total number of users to set up pagination.
	query := bson.M{}
	totalUsers, countDocumentsError := userRepository.Collection.CountDocuments(ctx, query)
	if validator.IsError(countDocumentsError) {
		internalError := domainError.NewInternalError(location+"GetAllUsers.Collection.CountDocuments", countDocumentsError.Error())
		userRepository.Logger.Error(internalError)
		return common.NewResultOnFailure[*userModel.Users](internalError)
	}

	// Set up pagination and sorting options using provided parameters.
	paginationQuery.TotalItems = int(totalUsers)
	paginationQuery.SetCorrectPage()
	option := options.FindOptions{}
	option.SetLimit(int64(paginationQuery.Limit))
	option.SetSkip(int64(paginationQuery.Skip))
	sortOptions := bson.M{paginationQuery.OrderBy: mongoModel.SetSortOrder(paginationQuery.SortOrder)}
	option.SetSort(sortOptions)

	// Query the database to fetch users.
	cursor, findError := userRepository.Collection.Find(ctx, query, &option)
	if validator.IsError(findError) {
		queryString := utility.ConvertQueryToString(query)
		itemNotFoundError := domainError.NewItemNotFoundError(location+"GetAllUsers.Find", queryString, findError.Error())
		userRepository.Logger.Error(itemNotFoundError)
		return common.NewResultOnFailure[*userModel.Users](itemNotFoundError)
	}
	defer cursor.Close(ctx)

	// Process the results and map them to the repository model.
	fetchedUsers := make([]*userRepositoryModel.UserRepository, 0, paginationQuery.Limit)
	for cursor.Next(ctx) {
		user := userRepositoryModel.UserRepository{}
		decodeError := cursor.Decode(&user)
		if validator.IsError(decodeError) {
			internalError := domainError.NewInternalError(location+"GetAllUsers.cursor.decode", decodeError.Error())
			userRepository.Logger.Error(internalError)
			return common.NewResultOnFailure[*userModel.Users](internalError)
		}
		fetchedUsers = append(fetchedUsers, &user)
	}

	cursorError := cursor.Err()
	if validator.IsError(cursorError) {
		internalError := domainError.NewInternalError(location+"GetAllUsers.cursor.Err", cursorError.Error())
		userRepository.Logger.Error(internalError)
		return common.NewResultOnFailure[*userModel.Users](internalError)
	}

	if len(fetchedUsers) == 0 {
		return common.NewResultOnSuccess[*userModel.Users](nil)
	}

	usersRepository := userRepositoryModel.UserRepositoryToUsersRepositoryMapper(fetchedUsers)
	usersRepository.PaginationResponse = common.NewPaginationResponse(paginationQuery)
	return common.NewResultOnSuccess[*userModel.Users](userRepositoryModel.UsersRepositoryToUsersMapper(usersRepository))
}

// GetUserById retrieves a user by their ID from the database.
func (userRepository *UserRepository) GetUserById(ctx context.Context, userID string) common.Result[userModel.User] {
	userObjectID := mongoModel.HexToObjectIDMapper(userRepository.Logger, location+"GetUserById", userID)
	if validator.IsError(userObjectID.Error) {
		return common.NewResultOnFailure[userModel.User](userObjectID.Error)
	}

	query := bson.M{mongoModel.ID: userObjectID.Data}
	return userRepository.getUserByQuery(location+"GetUserById", ctx, query)
}

// GetUserByEmail retrieves a user by their email from the database.
func (userRepository *UserRepository) GetUserByEmail(ctx context.Context, email string) common.Result[userModel.User] {
	query := bson.M{emailKey: email}
	return userRepository.getUserByQuery(location+"GetUserByEmail", ctx, query)
}

// CheckEmailDuplicate checks if an email already exists in the database.
func (userRepository *UserRepository) CheckEmailDuplicate(ctx context.Context, email string) error {
	fetchedUser := userRepositoryModel.UserRepository{}

	// Find and decode the user.
	// If no user is found, return nil (indicating that the email is unique).
	query := bson.M{emailKey: email}
	userFindOneError := userRepository.Collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsError(userFindOneError) {
		if userFindOneError == mongo.ErrNoDocuments {
			return nil
		}

		internalError := domainError.NewInternalError(location+"CheckEmailDuplicate.FindOne.Decode", userFindOneError.Error())
		userRepository.Logger.Error(internalError)
		return internalError
	}

	// If a user with the given email is found, return a validation error.
	validationError := domainError.NewValidationError(
		location+"CheckEmailDuplicate",
		usecase.EmailField,
		constants.FieldRequired,
		constants.EmailAlreadyExists,
	)

	userRepository.Logger.Error(validationError)
	return validationError
}

// Register creates a user in the database based on the provided UserCreate data.
func (userRepository *UserRepository) Register(ctx context.Context, userCreate userModel.UserCreate) common.Result[userModel.User] {
	userCreateRepository := userRepositoryModel.UserCreateToUserCreateRepositoryMapper(userCreate)
	hashedPassword := repositoryUtility.HashPassword(userRepository.Logger, location+"Register", userCreateRepository.Password)
	if validator.IsError(hashedPassword.Error) {
		return common.NewResultOnFailure[userModel.User](hashedPassword.Error)
	}

	userCreateRepository.Password = hashedPassword.Data
	insertOneResult, insertOneResultError := userRepository.Collection.InsertOne(ctx, &userCreateRepository)
	if validator.IsError(insertOneResultError) {
		internalError := domainError.NewInternalError(location+"Register.InsertOne", insertOneResultError.Error())
		userRepository.Logger.Error(internalError)
		return common.NewResultOnFailure[userModel.User](internalError)
	}

	query := bson.M{mongoModel.ID: insertOneResult.InsertedID}
	return userRepository.getUserByQuery(location+"Register", ctx, query)
}

// UpdateCurrentUser updates a user in the database based on the provided UserUpdate data.
func (userRepository *UserRepository) UpdateCurrentUser(ctx context.Context, userUpdate userModel.UserUpdate) common.Result[userModel.User] {
	userUpdateRepository := userRepositoryModel.UserUpdateToUserUpdateRepositoryMapper(userRepository.Logger, location+"UpdateCurrentUser", userUpdate)
	if validator.IsError(userUpdateRepository.Error) {
		return common.NewResultOnFailure[userModel.User](userUpdateRepository.Error)
	}

	userUpdateBSON := mongoModel.DataToMongoDocumentMapper(userRepository.Logger, location+"UpdateCurrentUser", userUpdateRepository.Data)
	if validator.IsError(userUpdateBSON.Error) {
		return common.NewResultOnFailure[userModel.User](userUpdateBSON.Error)
	}

	query := bson.D{{Key: mongoModel.ID, Value: userUpdateRepository.Data.UserID}}
	update := bson.D{{Key: mongoModel.Set, Value: userUpdateBSON.Data}}
	result := userRepository.Collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))
	updatedUser := userRepositoryModel.UserRepository{}
	decodeError := result.Decode(&updatedUser)
	if validator.IsError(decodeError) {
		internalError := domainError.NewInternalError(location+"UpdateCurrentUser.Decode", decodeError.Error())
		userRepository.Logger.Error(internalError)
		return common.NewResultOnFailure[userModel.User](internalError)
	}

	return common.NewResultOnSuccess[userModel.User](userRepositoryModel.UserRepositoryToUserMapper(updatedUser))
}

// DeleteUserById deletes a user in the database based on the provided userID.
func (userRepository *UserRepository) DeleteUserById(ctx context.Context, userID string) error {
	userObjectID := mongoModel.HexToObjectIDMapper(userRepository.Logger, location+"GetUserById", userID)
	if validator.IsError(userObjectID.Error) {
		return userObjectID.Error
	}

	query := bson.M{mongoModel.ID: userObjectID.Data}
	result, userDeleteOneError := userRepository.Collection.DeleteOne(ctx, query)
	if validator.IsError(userDeleteOneError) {
		internalError := domainError.NewInternalError(location+"Delete.DeleteOne", userDeleteOneError.Error())
		userRepository.Logger.Error(internalError)
		return internalError
	}

	if result.DeletedCount == 0 {
		internalError := domainError.NewInternalError(location+"Delete.DeleteOne.DeletedCount", mongoModel.DeletionIsNotSuccessful)
		userRepository.Logger.Error(internalError)
		return internalError
	}

	return nil
}

// GetResetExpiry retrieves a reset token based on the provided reset token from the database.
func (userRepository *UserRepository) GetResetExpiry(ctx context.Context, token string) common.Result[userModel.UserResetExpiry] {
	fetchedResetExpiry := userRepositoryModel.UserResetExpiryRepository{}
	query := bson.M{resetTokenKey: token}
	userFindOneError := userRepository.Collection.FindOne(ctx, query).Decode(&fetchedResetExpiry)
	if validator.IsError(userFindOneError) {
		invalidTokenError := domainError.NewInvalidTokenError(location+"GetResetExpiry.Decode", userFindOneError.Error())
		userRepository.Logger.Error(invalidTokenError)
		invalidTokenError.Notification = constants.InvalidTokenErrorMessage
		return common.NewResultOnFailure[userModel.UserResetExpiry](invalidTokenError)
	}

	return common.NewResultOnSuccess[userModel.UserResetExpiry](userRepositoryModel.UserResetExpiryRepositoryToUserResetExpiryMapper(fetchedResetExpiry))
}

// ForgottenPassword updates a user's record with a reset token and expiration time.
func (userRepository *UserRepository) ForgottenPassword(ctx context.Context, userForgottenPassword userModel.UserForgottenPassword) error {
	userForgottenPasswordRepository := userRepositoryModel.UserForgottenPasswordToUserForgottenPasswordRepositoryMapper(userForgottenPassword)
	userForgottenPasswordBSON := mongoModel.DataToMongoDocumentMapper(userRepository.Logger, location+"ForgottenPassword", userForgottenPasswordRepository)
	if validator.IsError(userForgottenPasswordBSON.Error) {
		return domainError.NewInternalError(location+"ForgottenPassword.Mapping", userForgottenPasswordBSON.Error.Error())
	}

	query := bson.D{{Key: emailKey, Value: userForgottenPassword.Email}}
	update := bson.D{{Key: mongoModel.Set, Value: userForgottenPasswordBSON.Data}}
	result, updateOneError := userRepository.Collection.UpdateOne(ctx, query, update)
	if validator.IsError(updateOneError) {
		internalError := domainError.NewInternalError(location+"ForgottenPassword.UpdateOne", updateOneError.Error())
		userRepository.Logger.Error(internalError)
		return internalError
	}

	if result.ModifiedCount == 0 {
		internalError := domainError.NewInternalError(location+"ForgottenPassword.UpdateOne.ModifiedCount", mongoModel.UpdateIsNotSuccessful)
		userRepository.Logger.Error(internalError)
		return internalError
	}

	return nil
}

// ResetUserPassword updates a user's password based on the provided reset token and new password.
func (userRepository *UserRepository) ResetUserPassword(ctx context.Context, userResetPassword userModel.UserResetPassword) error {
	userResetPasswordRepository := userRepositoryModel.UserResetPasswordToUserResetPasswordRepositoryMapper(userResetPassword)
	hashedPassword := repositoryUtility.HashPassword(userRepository.Logger, location+"ResetUserPassword", userResetPassword.Password)
	if validator.IsError(hashedPassword.Error) {
		return hashedPassword.Error
	}

	userResetPasswordRepository.Password = hashedPassword.Data
	userResetPasswordBSON := mongoModel.DataToMongoDocumentMapper(userRepository.Logger, location+"ResetUserPassword", userResetPasswordRepository)
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

	result, updateOneError := userRepository.Collection.UpdateOne(ctx, query, update)
	if validator.IsError(updateOneError) {
		internalError := domainError.NewInternalError(location+"ResetUserPassword.UpdateOne", updateOneError.Error())
		userRepository.Logger.Error(internalError)
		return internalError
	}

	if result.ModifiedCount == 0 {
		internalError := domainError.NewInternalError(location+"ResetUserPassword.UpdateOne.ModifiedCount", mongoModel.UpdateIsNotSuccessful)
		userRepository.Logger.Error(internalError)
		return internalError
	}

	return nil
}

// SendEmail sends an email to the specified user with the provided data.
func (userRepository *UserRepository) SendEmail(user userModel.User, data userModel.EmailData) error {
	sendEmailError := userRepositoryMail.SendEmail(userRepository.Config, userRepository.Logger, location+"SendEmail", user, data)
	if validator.IsError(sendEmailError) {
		return sendEmailError
	}

	return nil
}

// ensureUniqueEmailIndex creates a unique index on the email field to enforce email uniqueness in the database.
func (userRepository *UserRepository) ensureUniqueEmailIndex(ctx context.Context, location string) error {
	option := options.Index()
	option.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{emailKey: 1}, Options: option}
	_, userIndexesCreateOneError := userRepository.Collection.Indexes().CreateOne(ctx, index)
	if validator.IsError(userIndexesCreateOneError) {
		internalError := domainError.NewInternalError(location+".ensureUniqueEmailIndex.Indexes.CreateOne", userIndexesCreateOneError.Error())
		userRepository.Logger.Error(internalError)
		return internalError
	}

	return nil
}

// getUserByQuery retrieves a user based on the provided query from the database.
func (userRepository *UserRepository) getUserByQuery(location string, ctx context.Context, query bson.M) common.Result[userModel.User] {
	fetchedUser := userRepositoryModel.UserRepository{}
	userFindOneError := userRepository.Collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsError(userFindOneError) {
		queryString := mongoModel.BSONToString(query)
		itemNotFoundError := domainError.NewItemNotFoundError(location+".getUserByQuery.Decode", queryString, userFindOneError.Error())
		userRepository.Logger.Error(itemNotFoundError)
		return common.NewResultOnFailure[userModel.User](itemNotFoundError)
	}

	return common.NewResultOnSuccess[userModel.User](userRepositoryModel.UserRepositoryToUserMapper(fetchedUser))
}
