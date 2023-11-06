package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	location = "User.Data.Repository.Mongo.Model."
)

func UserRepositoryToUsersRepositoryMapper(usersRepository []UserRepository) UsersRepository {
	return UsersRepository{
		Users: append([]UserRepository{}, usersRepository...),
	}
}

func UsersRepositoryToUsersMapper(usersRepository UsersRepository) userModel.Users {
	users := make([]userModel.User, len(usersRepository.Users))
	for index, userRepository := range usersRepository.Users {
		users[index] = UserRepositoryToUserMapper(userRepository)
	}
	return userModel.Users{
		PaginationResponse: usersRepository.PaginationResponse,
		Users:              users,
	}
}

func UserCreateToUserCreateRepositoryMapper(userCreate userModel.UserCreate) UserCreateRepository {
	return UserCreateRepository{
		Name:             userCreate.Name,
		Email:            userCreate.Email,
		Password:         userCreate.Password,
		Role:             userCreate.Role,
		Verified:         userCreate.Verified,
		VerificationCode: userCreate.VerificationCode,
		CreatedAt:        userCreate.CreatedAt,
		UpdatedAt:        userCreate.UpdatedAt,
	}
}

func UserUpdateToUserUpdateRepositoryMapper(userUpdate userModel.UserUpdate) (UserUpdateRepository, error) {
	userObjectID, objectIDFromHexError := primitive.ObjectIDFromHex(userUpdate.UserID)
	if objectIDFromHexError != nil {
		objectIDFromHexError := domainError.NewInternalError(location+"UserUpdateToUserUpdateRepositoryMapper.ObjectIDFromHex", objectIDFromHexError.Error())
		logging.Logger(objectIDFromHexError)
		return UserUpdateRepository{}, objectIDFromHexError
	}
	return UserUpdateRepository{
		UserID:    userObjectID,
		Name:      userUpdate.Name,
		UpdatedAt: userUpdate.UpdatedAt,
	}, nil
}

func UserRepositoryToUserMapper(userRepository UserRepository) userModel.User {
	return userModel.User{
		UserID:    userRepository.UserID.Hex(),
		Name:      userRepository.Name,
		Email:     userRepository.Email,
		Password:  userRepository.Password,
		Role:      userRepository.Role,
		Verified:  userRepository.Verified,
		CreatedAt: userRepository.CreatedAt,
		UpdatedAt: userRepository.UpdatedAt,
	}
}

func objectIDFromHex(userID string) primitive.ObjectID {
	userObjectID, objectIDFromHexError := primitive.ObjectIDFromHex(userID)
	if objectIDFromHexError != nil {
		objectIDFromHexError := domainError.NewInternalError(location+"GetUserById.ObjectIDFromHex", objectIDFromHexError.Error())
		logging.Logger(objectIDFromHexError)
	}
	return userObjectID
}
