package model

import userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"

func UsersRepositoryToUsersMapper(usersRepository []*UserRepository) *userModel.Users {
	users := make([]*userModel.User, 0, len(usersRepository))
	for _, userRepository := range usersRepository {
		user := UserRepositoryToUserMapper(userRepository)
		users = append(users, user)
	}

	return &userModel.Users{
		Users: users,
	}
}

func UserRepositoryToUserMapper(userRepository *UserRepository) *userModel.User {
	return &userModel.User{
		UserID:    userRepository.UserID,
		Name:      userRepository.Name,
		Email:     userRepository.Email,
		Password:  userRepository.Password,
		Role:      userRepository.Role,
		CreatedAt: userRepository.CreatedAt,
		UpdatedAt: userRepository.UpdatedAt,
	}
}

func UserCreateToUserCreateRepositoryMapper(user *userModel.UserCreate) *UserCreateRepository {
	return &UserCreateRepository{
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserUpdateToUserUpdateRepositoryMapper(user *userModel.UserUpdate) *UserUpdateRepository {
	return &UserUpdateRepository{
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
	}
}
