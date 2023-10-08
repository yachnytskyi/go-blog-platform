package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
)

func UserRepositoryToUsersRepositoryMapper(usersRepository []UserRepository) UsersRepository {
	users := make([]UserRepository, 0, len(usersRepository))
	for _, userRepository := range usersRepository {
		users = append(users, userRepository)
	}
	return UsersRepository{
		Users: users,
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

func UserUpdateToUserUpdateRepositoryMapper(userUpdate userModel.UserUpdate) UserUpdateRepository {
	return UserUpdateRepository{
		Name:      userUpdate.Name,
		UpdatedAt: userUpdate.UpdatedAt,
	}
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
