package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	mongoModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "user.data.repository.mongo.model."
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
	userObjectID, hexToObjectIDMapperError := mongoModel.HexToObjectIDMapper(location+"UserUpdateToUserUpdateRepositoryMapper", userUpdate.UserID)
	if validator.IsErrorNotNil(hexToObjectIDMapperError) {
		return UserUpdateRepository{}, hexToObjectIDMapperError
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
