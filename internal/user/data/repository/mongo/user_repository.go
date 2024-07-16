package repository

import (
	"context"
	"time"

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
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
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

// UserRepository provides methods to interact with the user collection in the database.
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new UserRepository and ensures the unique email index.
// It performs the following steps:
// 1. Initializes the UserRepository with the given MongoDB collection.
// 2. Ensures a unique index on the email field during initialization.
// 3. Handles any potential errors during index creation.
//
// Parameters:
// - database (*mongo.Database): The MongoDB database instance.
//
// Returns:
// - UserRepository: The initialized UserRepository.
func NewUserRepository(database *mongo.Database) UserRepository {
	repository := UserRepository{collection: database.Collection(constants.UsersTable)}

	// Ensure the unique index on email during initialization.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ensureUniqueEmailIndexError := repository.ensureUniqueEmailIndex(ctx, location+"NewUserRepository")
	if validator.IsError(ensureUniqueEmailIndexError) {
		// Handle index creation error appropriately (e.g., log it, panic, etc.)
		panic(ensureUniqueEmailIndexError)
	}

	return repository
}

// GetAllUsers retrieves a list of users from the database based on pagination parameters.
// It performs the following steps:
// 1. Initializes the query with an empty BSON document.
// 2. Determines the sort order based on the pagination query.
// 3. Counts the total number of users to set up pagination.
// 4. Sets up pagination and sorting options using provided parameters.
// 5. Queries the database to fetch users and processes the results.
// 6. Checks for cursor errors after processing the results.
// 7. Maps the repository model to domain ones and returns the result.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - paginationQuery (commonModel.PaginationQuery): The pagination query parameters.
//
// Returns:
// - commonModel.Result[userModel.Users]: The result containing either the list of users or an error.
func (userRepository UserRepository) GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users] {
	// Initialize the query with an empty BSON document.
	query := bson.M{}

	// Determine the sort order based on the pagination query.
	sortOrder := mongoModel.SetSortOrder(paginationQuery.SortOrder)

	// Count the total number of users to set up pagination.
	totalUsers, countDocumentsError := userRepository.collection.CountDocuments(ctx, query)
	if validator.IsError(countDocumentsError) {
		// If an error occurs while counting documents, log and return an internal error.
		internalError := domainError.NewInternalError(location+"GetAllUsers.collection.CountDocuments", countDocumentsError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.Users](internalError)
	}

	// Set up pagination and sorting options using provided parameters.
	paginationQuery.TotalItems = int(totalUsers)
	paginationQuery = commonModel.SetCorrectPage(paginationQuery)
	option := options.FindOptions{}
	option.SetLimit(int64(paginationQuery.Limit))
	option.SetSkip(int64(paginationQuery.Skip))
	sortOptions := bson.M{paginationQuery.OrderBy: sortOrder}
	option.SetSort(sortOptions)

	// Query the database to fetch users.
	cursor, findError := userRepository.collection.Find(ctx, query, &option)
	if validator.IsError(findError) {
		// If an error occurs while finding documents, log and return an item not found error.
		queryString := commonUtility.ConvertQueryToString(query)
		itemNotFoundError := domainError.NewItemNotFoundError(location+"GetAllUsers.Find", queryString, findError.Error())
		logging.Logger(itemNotFoundError)
		return commonModel.NewResultOnFailure[userModel.Users](itemNotFoundError)
	}
	defer cursor.Close(ctx)

	// Process the results and map them to the repository model.
	fetchedUsers := make([]userRepositoryModel.UserRepository, 0, paginationQuery.Limit)
	for cursor.Next(ctx) {
		user := userRepositoryModel.UserRepository{}
		decodeError := cursor.Decode(&user)
		if validator.IsError(decodeError) {
			// If an error occurs while decoding documents, log and return an internal error.
			internalError := domainError.NewInternalError(location+"GetAllUsers.cursor.decode", decodeError.Error())
			logging.Logger(internalError)
			return commonModel.NewResultOnFailure[userModel.Users](internalError)
		}
		fetchedUsers = append(fetchedUsers, user)
	}

	// Check for cursor errors.
	cursorError := cursor.Err()
	if validator.IsError(cursorError) {
		// If an error occurs with the cursor, log and return an internal error.
		internalError := domainError.NewInternalError(location+"GetAllUsers.cursor.Err", cursorError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.Users](internalError)
	}

	// If no users are fetched, return an empty result.
	if len(fetchedUsers) == 0 {
		return commonModel.NewResultOnSuccess[userModel.Users](userModel.Users{})
	}

	// Map the repository model to domain ones.
	usersRepository := userRepositoryModel.UserRepositoryToUsersRepositoryMapper(fetchedUsers)
	usersRepository.PaginationResponse = commonModel.NewPaginationResponse(paginationQuery)
	return commonModel.NewResultOnSuccess[userModel.Users](userRepositoryModel.UsersRepositoryToUsersMapper(usersRepository))
}

// GetUserById retrieves a user by their ID from the database.
// It performs the following steps:
// 1. Maps the provided userID to an ObjectID.
// 2. Defines the MongoDB query to find the user by ObjectID.
// 3. Uses the getUserByQuery method to retrieve the user from the database.
// 4. Returns the result, which may be either a success with the user data or a failure with an error.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userID (string): The unique identifier of the user to be fetched.
//
// Returns:
// - commonModel.Result[userModel.User]: The result containing either the user data or an error.
func (userRepository UserRepository) GetUserById(ctx context.Context, userID string) commonModel.Result[userModel.User] {
	userObjectID := mongoModel.HexToObjectIDMapper(location+"GetUserById", userID)
	if validator.IsError(userObjectID.Error) {
		return commonModel.NewResultOnFailure[userModel.User](userObjectID.Error)
	}

	// Define the MongoDB query to find the user by ObjectID.
	query := bson.M{mongoModel.ID: userObjectID.Data}
	return userRepository.getUserByQuery(location+"GetUserById", ctx, query)
}

// GetUserByEmail retrieves a user by their email from the repository.
// It performs the following steps:
// 1. Defines the MongoDB query to find the user by email.
// 2. Uses the getUserByQuery method to retrieve the user from the database.
// 3. Returns the result, which may be either a success with the user data or a failure with an error.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - email (string): The email address of the user to be fetched.
//
// Returns:
// - commonModel.Result[userModel.User]: The result containing either the user data or an error.
func (userRepository UserRepository) GetUserByEmail(ctx context.Context, email string) commonModel.Result[userModel.User] {
	// Define the MongoDB query to find the user by email.
	query := bson.M{emailKey: email}

	// Retrieve the user from the database.
	return userRepository.getUserByQuery(location+"GetUserByEmail", ctx, query)
}

// CheckEmailDuplicate checks if an email already exists in the repository.
// It performs the following steps:
// 1. Initializes a User object and defines the MongoDB query to find the user by email.
// 2. Finds and decodes the user.
// 3. If no user is found, returns nil (indicating that the email is unique).
// 4. If an error occurs during the database query, logs it as an internal error and returns the error.
// 5. If a user with the given email is found, returns a validation error indicating that the email already exists.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - email (string): The email address to be checked for duplication.
//
// Returns:
// - error: An error indicating whether the email is unique or already exists in the repository.
func (userRepository UserRepository) CheckEmailDuplicate(ctx context.Context, email string) error {
	// Initialize a User object and define the MongoDB query to find the user by email.
	fetchedUser := userRepositoryModel.UserRepository{}
	query := bson.M{emailKey: email}

	// Find and decode the user.
	// If no user is found, return nil (indicating that the email is unique).
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsError(userFindOneError) {
		if userFindOneError == mongo.ErrNoDocuments {
			// No user found, email is unique.
			return nil
		}

		// If an error occurs during the database query, log it as an internal error.
		internalError := domainError.NewInternalError(location+"CheckEmailDuplicate.FindOne.Decode", userFindOneError.Error())
		logging.Logger(internalError)
		return internalError
	}

	// If a user with the given email is found, return a validation error.
	userFindOneValidationError := domainError.NewValidationError(location+"CheckEmailDuplicate", userValidator.EmailField, constants.FieldRequired, constants.EmailAlreadyExists)
	logging.Logger(userFindOneValidationError)
	return userFindOneValidationError
}

// Register creates a user in the repository based on the provided UserCreate data.
// It performs the following steps:
// 1. Maps the incoming data to the repository model.
// 2. Hashes the user's password.
// 3. Sets the hashed password in the repository model.
// 4. Inserts the user into the database by executing the MongoDB insert query.
// 5. Retrieves the created user from the database, maps it back to the domain model, and returns the result.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userCreate (userModel.UserCreate): The data for creating the user.
//
// Returns:
// - commonModel.Result[userModel.User]: The result containing either the user data or an error.
func (userRepository UserRepository) Register(ctx context.Context, userCreate userModel.UserCreate) commonModel.Result[userModel.User] {
	// Map the incoming user data to the repository model.
	userCreateRepository := userRepositoryModel.UserCreateToUserCreateRepositoryMapper(userCreate)

	// Hash the user's password.
	hashedPassword := repositoryUtility.HashPassword(location+"Register", userCreateRepository.Password)
	if validator.IsError(hashedPassword.Error) {
		return commonModel.NewResultOnFailure[userModel.User](hashedPassword.Error)
	}

	// Set the hashed password in the repository model.
	userCreateRepository.Password = hashedPassword.Data

	// Insert the user data into the database.
	insertOneResult, insertOneResultError := userRepository.collection.InsertOne(ctx, &userCreateRepository)
	if validator.IsError(insertOneResultError) {
		internalError := domainError.NewInternalError(location+"Register.InsertOne", insertOneResultError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.User](internalError)
	}

	// Define the MongoDB query to find the user by ObjectID.
	// Retrieve the created user from the database.
	query := bson.M{mongoModel.ID: insertOneResult.InsertedID}
	return userRepository.getUserByQuery(location+"Register", ctx, query)
}

// UpdateCurrentUser updates a user in the repository based on the provided UserUpdate data.
// It performs the following steps:
// 1. Maps the incoming data to a repository model.
// 2. Maps the repository model to a MongoDB BSON document.
// 3. Constructs a MongoDB query and update operation.
// 4. Executes the MongoDB update query and retrieves the updated user.
// 5. Decodes the updated user from the result, maps it back to the domain model, and returns the result.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userUpdate (userModel.UserUpdate): The data for updating the user.
//
// Returns:
// - commonModel.Result[userModel.User]: The result containing either the user data or an error.
func (userRepository UserRepository) UpdateCurrentUser(ctx context.Context, userUpdate userModel.UserUpdate) commonModel.Result[userModel.User] {
	// Map user update data to a repository model.
	userUpdateRepository := userRepositoryModel.UserUpdateToUserUpdateRepositoryMapper(userUpdate)
	if validator.IsError(userUpdateRepository.Error) {
		return commonModel.NewResultOnFailure[userModel.User](userUpdateRepository.Error)
	}

	// Map the user update repository to a BSON document for MongoDB update.
	userUpdateBSON := mongoModel.DataToMongoDocumentMapper(location+"UpdateCurrentUser", userUpdateRepository.Data)
	if validator.IsError(userUpdateBSON.Error) {
		return commonModel.NewResultOnFailure[userModel.User](userUpdateBSON.Error)
	}

	// Define the MongoDB query and update operation.
	query := bson.D{{Key: mongoModel.ID, Value: userUpdateRepository.Data.UserID}}
	update := bson.D{{Key: mongoModel.Set, Value: userUpdateBSON.Data}}
	result := userRepository.collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	// Decode the updated user from the result.
	updatedUser := userRepositoryModel.UserRepository{}
	decodeError := result.Decode(&updatedUser)
	if validator.IsError(decodeError) {
		internalError := domainError.NewInternalError(location+"UpdateCurrentUser.Decode", decodeError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[userModel.User](internalError)
	}

	// Map the updated repository model to the user model.
	return commonModel.NewResultOnSuccess[userModel.User](userRepositoryModel.UserRepositoryToUserMapper(updatedUser))
}

// DeleteUserById deletes a user in the repository based on the provided userID.
// It performs the following steps:
// 1. Maps the provided userID string to a MongoDB ObjectID.
// 2. Constructs a MongoDB query to delete the user by their ObjectID.
// 3. Executes the delete operation on the MongoDB collection.
// 4. Handles any errors encountered during the delete operation, logging and returning them if necessary.
// 5. Checks if any documents were deleted and handles cases where no documents were deleted.
// 6. Returns nil to indicate a successful operation.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userID (string): The unique identifier of the user to be deleted.
//
// Returns:
// - error: An error if any occurred during the operation, otherwise nil.
func (userRepository UserRepository) DeleteUserById(ctx context.Context, userID string) error {
	// Maps the provided userID string to a MongoDB ObjectID.
	userObjectID := mongoModel.HexToObjectIDMapper(location+"GetUserById", userID)
	if validator.IsError(userObjectID.Error) {
		return userObjectID.Error
	}

	// Define the MongoDB query to delete the user by ObjectID.
	query := bson.M{mongoModel.ID: userObjectID.Data}

	// Execute the delete operation.
	result, userDeleteOneError := userRepository.collection.DeleteOne(ctx, query)
	if validator.IsError(userDeleteOneError) {
		internalError := domainError.NewInternalError(location+"Delete.DeleteOne", userDeleteOneError.Error())
		logging.Logger(internalError)
		return internalError
	}

	// Check if any user was deleted.
	if result.DeletedCount == 0 {
		internalError := domainError.NewInternalError(location+"Delete.DeleteOne.DeletedCount", mongoModel.DeletionIsNotSuccessful)
		logging.Logger(internalError)
		return internalError
	}

	// Return nil to indicate success.
	return nil
}

// GetUserByResetToken retrieves a user based on the provided reset token from the repository.
// It performs the following steps:
// 1. Defines the MongoDB query to find the user by the reset token.
// 2. Executes the query and decodes the result into a UserResetTokenRepository model.
// 3. Handles any errors that occur during the query or decoding, logging and returning an appropriate error message.
// 4. Maps the repository model to the domain model and returns the result.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - token (string): The reset token of the user to be fetched.
//
// Returns:
// - commonModel.Result[userModel.UserResetToken]: The result containing either the user data or an error.
func (userRepository UserRepository) GetUserByResetToken(ctx context.Context, token string) commonModel.Result[userModel.UserResetToken] {
	// Define the MongoDB query to find the user by the reset token.
	query := bson.M{resetTokenKey: token}

	// Initialize a UserResetTokenRepository object and find the user based on the provided query.
	fetchedUser := userRepositoryModel.UserResetTokenRepository{}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsError(userFindOneError) {
		// Handle errors that occur during the query or decoding.
		invalidTokenError := domainError.NewInvalidTokenError(location+".GetUserByResetToken.Decode", userFindOneError.Error())
		logging.Logger(invalidTokenError)
		invalidTokenError.Notification = constants.InvalidTokenErrorMessage
		return commonModel.NewResultOnFailure[userModel.UserResetToken](invalidTokenError)
	}

	// Map the repository model to the domain model.
	return commonModel.NewResultOnSuccess[userModel.UserResetToken](userRepositoryModel.UserResetTokenRepositoryToUserTokenMapper(fetchedUser))
}

// ForgottenPassword updates a user's record with a reset token and expiration time.
// It performs the following steps:
// 1. Constructs a MongoDB query to find the user by their email address.
// 2. Prepares the update document by mapping the userForgottenPassword to a repository model.
// 3. Executes the update operation on the MongoDB collection.
// 4. Handles any errors encountered during the update operation, logging and returning them if necessary.
// 5. Checks if the update operation modified any documents and handles cases where no documents were modified.
// 6. Returns nil to indicate a successful operation.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userForgottenPassword (userModel.UserForgottenPassword): The data for updating the user's forgotten password.
//
// Returns:
// - error: An error if the operation fails, otherwise nil.
func (userRepository UserRepository) ForgottenPassword(ctx context.Context, userForgottenPassword userModel.UserForgottenPassword) error {
	// Define the MongoDB query.
	query := bson.D{{Key: emailKey, Value: userForgottenPassword.Email}}

	// Map the userForgottenPassword data to a repository model.
	userForgottenPasswordRepository := userRepositoryModel.UserForgottenPasswordToUserForgottenPasswordRepositoryMapper(userForgottenPassword)

	// Map the userForgottenPassword repository model to a BSON document for MongoDB update.
	userForgottenPasswordBSON := mongoModel.DataToMongoDocumentMapper(location+"ForgottenPassword", userForgottenPasswordRepository)
	if validator.IsError(userForgottenPasswordBSON.Error) {
		return domainError.NewInternalError(location+"ForgottenPassword.Mapping", userForgottenPasswordBSON.Error.Error())
	}

	// Define the update operation with the entire model.
	update := bson.D{{Key: mongoModel.Set, Value: userForgottenPasswordBSON.Data}}

	// Execute the update operation.
	result, updateOneError := userRepository.collection.UpdateOne(ctx, query, update)
	if validator.IsError(updateOneError) {
		// Log and return an internal error if the update operation fails.
		internalError := domainError.NewInternalError(location+"ForgottenPassword.UpdateOne", updateOneError.Error())
		logging.Logger(internalError)
		return internalError
	}

	// Check if any document was modified. If not, log and return an error.
	if result.ModifiedCount == 0 {
		internalError := domainError.NewInternalError(location+"ForgottenPassword.UpdateOne.ModifiedCount", mongoModel.UpdateIsNotSuccessful)
		logging.Logger(internalError)
		return internalError
	}

	// Return nil to indicate success.
	return nil
}

// ResetUserPassword updates a user's password based on the provided reset token and new password.
// It performs the following steps:
// 1. Maps the incoming data to a repository model.
// 2. Maps the repository model to a MongoDB BSON document.
// 3. Constructs a MongoDB query and update operation.
// 4. Executes the MongoDB update query and retrieves the updated user.
// 5. Checks if the update operation modified any documents and handles cases where no documents were modified.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userResetPassword (userModel.UserResetPassword): The data for resetting the user's password.
//
// Returns:
// - error: An error if the operation fails, otherwise nil.
func (userRepository UserRepository) ResetUserPassword(ctx context.Context, userResetPassword userModel.UserResetPassword) error {
	// Map the user reset password data to a repository model.
	userResetPasswordRepository := userRepositoryModel.UserResetPasswordToUserResetPasswordRepositoryMapper(userResetPassword)

	// Hash the user's password.
	hashedPassword := repositoryUtility.HashPassword(location+"ResetUserPassword", userResetPassword.Password)
	if validator.IsError(hashedPassword.Error) {
		return hashedPassword.Error
	}

	// Set the hashed password in the repository model.
	userResetPasswordRepository.Password = hashedPassword.Data

	// Map the user reset password repository model to a BSON document for MongoDB update.
	userResetPasswordBSON := mongoModel.DataToMongoDocumentMapper(location+"ResetUserPassword", userResetPasswordRepository)
	if validator.IsError(userResetPasswordBSON.Error) {
		return userResetPasswordBSON.Error
	}

	// Define the MongoDB query.
	query := bson.D{{Key: resetTokenKey, Value: userResetPassword.ResetToken}}

	// Define the update operation with the password update and the fields to unset.
	update := bson.D{
		{Key: mongoModel.Set, Value: userResetPasswordBSON.Data},
		{Key: mongoModel.Unset, Value: bson.D{
			{Key: resetTokenKey, Value: ""},
			{Key: resetExpiryKey, Value: ""},
		}},
	}

	// Execute the update operation.
	result, updateOneError := userRepository.collection.UpdateOne(ctx, query, update)
	if validator.IsError(updateOneError) {
		// Log and return an internal error if the update operation fails.
		internalError := domainError.NewInternalError(location+"ResetUserPassword.UpdateOne", updateOneError.Error())
		logging.Logger(internalError)
		return internalError
	}

	// Check if any document was modified. If not, log and return an error.
	if result.ModifiedCount == 0 {
		internalError := domainError.NewInternalError(location+"ResetUserPassword.UpdateOne.ModifiedCount", mongoModel.UpdateIsNotSuccessful)
		logging.Logger(internalError)
		return internalError
	}

	// Return nil to indicate success.
	return nil
}

// SendEmail sends an email to the specified user with the provided data.
// This method performs the following steps:
// 1. Calls the SendEmail function from the Mail repository, passing the location, user, and data as parameters.
// 2. Checks for any errors returned by the SendEmail function.
// 3. If an error is encountered, it is returned.
// 4. If no error occurs, the method returns nil indicating success.
//
// Parameters:
// - user (userModel.User): The user to whom the email will be sent.
// - data (userModel.EmailData): The data to be included in the email.
//
// Returns:
// - error: An error if the operation fails, otherwise nil.
func (userRepository UserRepository) SendEmail(user userModel.User, data userModel.EmailData) error {
	// Send the email using the Mail repository.
	sendEmailError := userRepositoryMail.SendEmail(location+"SendEmail", user, data)
	if validator.IsError(sendEmailError) {
		// Return the error if encountered.
		return sendEmailError
	}

	return nil
}

// ensureUniqueEmailIndex creates a unique index on the email field to enforce email uniqueness in the repository.
// This method ensures that each user has a unique email address in the database by performing the following steps:
// 1. Creates options for the index and sets it as unique.
// 2. Defines the index model based on the email field.
// 3. Creates the unique index in the collection.
// 4. If an error occurs during index creation, it logs the error and returns an internal error.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - location (string): The location string for logging and error handling.
//
// Returns:
// - error: An error if the operation fails, otherwise nil.
func (userRepository UserRepository) ensureUniqueEmailIndex(ctx context.Context, location string) error {
	// Create options for the index, setting it as unique.
	option := options.Index()
	option.SetUnique(true)

	// Define the index model based on the email field.
	index := mongo.IndexModel{Keys: bson.M{emailKey: 1}, Options: option}

	// Create the unique index in the collection.
	_, userIndexesCreateOneError := userRepository.collection.Indexes().CreateOne(ctx, index)
	if validator.IsError(userIndexesCreateOneError) {
		// Log the error and return an internal error.
		internalError := domainError.NewInternalError(location+".ensureUniqueEmailIndex.Indexes.CreateOne", userIndexesCreateOneError.Error())
		logging.Logger(internalError)
		return internalError
	}

	return nil
}

// getUserByQuery retrieves a user based on the provided query from the repository.
// It performs the following steps:
// 1. Initializes a User object and defines the MongoDB query.
// 2. Executes the query to find the user in the database and decodes the result.
// 3. If an error occurs during the query execution or decoding, it logs the error and returns an item not found error.
// 4. If the user is found, maps the repository model to the domain model and returns the result.
//
// Parameters:
// - location (string): The location string for logging and error handling.
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - query (bson.M): The MongoDB query to find the user.
//
// Returns:
// - commonModel.Result[userModel.User]: The result containing either the user data or an error.
func (userRepository UserRepository) getUserByQuery(location string, ctx context.Context, query bson.M) commonModel.Result[userModel.User] {
	// Initialize a User object and find the user based on the provided query.
	fetchedUser := userRepositoryModel.UserRepository{}
	userFindOneError := userRepository.collection.FindOne(ctx, query).Decode(&fetchedUser)
	if validator.IsError(userFindOneError) {
		// Convert the query to a string for logging.
		queryString := commonUtility.ConvertQueryToString(query)

		// Log the error and return an item not found error.
		itemNotFoundError := domainError.NewItemNotFoundError(location+".getUserByQuery.Decode", queryString, userFindOneError.Error())
		logging.Logger(itemNotFoundError)
		return commonModel.NewResultOnFailure[userModel.User](itemNotFoundError)
	}

	// Map the repository model to the domain model.
	return commonModel.NewResultOnSuccess[userModel.User](userRepositoryModel.UserRepositoryToUserMapper(fetchedUser))
}
