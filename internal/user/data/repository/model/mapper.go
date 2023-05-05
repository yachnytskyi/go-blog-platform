package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
)

func UserCreateToUserCreateRepositoryMapper(user *userModel.UserCreate) UserCreateRepository {
	return UserCreateRepository{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            user.Role,
		Verified:        user.Verified,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

func UserUpdateToUserUpdateRepositoryMapper(user *userModel.UserUpdate) UserUpdateRepository {
	return UserUpdateRepository{
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
	}
}
