package model

import userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"

func UsersRepositoryToUsersMapper(usersRepository []*UserRepository) userModel.Users {
	users := make([]*userModel.User, 0, 10)

	for _, userRepository := range usersRepository {
		user := &userModel.User{}
		user.UserID = userRepository.UserID
		user.Name = userRepository.Name
		user.Email = userRepository.Email
		user.Role = userRepository.Role
		user.CreatedAt = userRepository.CreatedAt
		user.UpdatedAt = userRepository.UpdatedAt
		users = append(users, user)
	}

	return userModel.Users{
		Users: users,
	}
}

func UserRepositoryToUserMapper(userRepository *UserRepository) userModel.User {
	return userModel.User{
		UserID:    userRepository.UserID,
		Name:      userRepository.Name,
		Email:     userRepository.Email,
		Role:      userRepository.Role,
		CreatedAt: userRepository.CreatedAt,
		UpdatedAt: userRepository.UpdatedAt,
	}
}

func UserCreateToUserCreateRepositoryMapper(user *userModel.UserCreate) UserCreateRepository {
	return UserCreateRepository{
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Verified:  user.Verified,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserUpdateToUserUpdateRepositoryMapper(user *userModel.UserUpdate) UserUpdateRepository {
	return UserUpdateRepository{
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
	}
}
