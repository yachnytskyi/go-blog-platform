package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	mongoModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
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

func UserUpdateToUserUpdateRepositoryMapper(userUpdate userModel.UserUpdate) commonModel.Result[UserUpdateRepository] {
	userObjectID := mongoModel.HexToObjectIDMapper(location+"UserUpdateToUserUpdateRepositoryMapper", userUpdate.ID)
	if validator.IsError(userObjectID.Error) {
		return commonModel.NewResultOnFailure[UserUpdateRepository](userObjectID.Error)
	}

	userUpdateRepository := UserUpdateRepository{
		UserID:    userObjectID.Data,
		Name:      userUpdate.Name,
		UpdatedAt: userUpdate.UpdatedAt,
	}

	return commonModel.NewResultOnSuccess[UserUpdateRepository](userUpdateRepository)
}

func UserRepositoryToUserMapper(userRepository UserRepository) userModel.User {
	return userModel.User{
		BaseEntity: domain.NewBaseEntity(userRepository.ID.Hex(), userRepository.CreatedAt, userRepository.UpdatedAt),
		Name:       userRepository.Name,
		Email:      userRepository.Email,
		Password:   userRepository.Password,
		Role:       userRepository.Role,
		Verified:   userRepository.Verified,
	}
}
