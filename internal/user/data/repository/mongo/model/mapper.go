package model

import userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"

func UsersRepositoryToUsersMapper(usersRepository []UserRepository) userModel.Users {
	users := make([]userModel.User, 0, len(usersRepository))
	for _, userRepository := range usersRepository {
		user := UserRepositoryToUserMapper(userRepository)
		users = append(users, user)
	}
	return userModel.Users{
		Users: users,
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
