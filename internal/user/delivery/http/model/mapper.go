package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	http "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "user.delivery.http.model."
)

func UsersToUsersViewMapper(users userModel.Users) UsersView {
	usersView := make([]UserView, len(users.Users))
	for index, user := range users.Users {
		usersView[index] = UserToUserViewMapper(user)
	}

	return UsersView{
		PaginationResponse: http.PaginationResponse{
			CurrentPage: users.PaginationResponse.Page,
			TotalPages:  users.PaginationResponse.TotalPages,
			PagesLeft:   users.PaginationResponse.PagesLeft,
			TotalItems:  users.PaginationResponse.TotalItems,
			ItemsLeft:   users.PaginationResponse.ItemsLeft,
			Limit:       users.PaginationResponse.Limit,
			OrderBy:     users.PaginationResponse.OrderBy,
		},
		UsersView: usersView,
	}
}

func UserToUserViewMapper(user userModel.User) UserView {
	return UserView{
		UserID:    user.UserID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: common.FormatDate(user.CreatedAt),
		UpdatedAt: common.FormatDate(user.UpdatedAt),
	}
}

func UserViewToUserMapper(user UserView) (userModel.User, error) {
	created_at, parseError := common.ParseDate(location+"UserViewToUserMapper.created_at", user.CreatedAt)
	if validator.IsError(parseError) {
		return userModel.User{}, parseError
	}

	updated_at, parseError := common.ParseDate(location+"UserViewToUserMapper.updated_at", user.UpdatedAt)
	if validator.IsError(parseError) {
		return userModel.User{}, parseError
	}

	return userModel.User{
		UserID:    user.UserID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}, nil
}

func UserCreateViewToUserCreateMapper(user UserCreateView) userModel.UserCreate {
	return userModel.UserCreate{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
	}
}

func UserUpdateViewToUserUpdateMapper(user UserUpdateView) userModel.UserUpdate {
	return userModel.UserUpdate{
		UserID: user.UserID,
		Name:   user.Name,
	}
}

func UserLoginViewToUserLoginMapper(user UserLoginView) userModel.UserLogin {
	return userModel.UserLogin{
		Email:    user.Email,
		Password: user.Password,
	}
}

func UserLoginToUserLoginViewMapper(user userModel.UserLogin) UserLoginView {
	return UserLoginView{
		Email:        user.Email,
		Password:     user.Password,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	}
}

func UserTokenToUserTokenViewMapper(user userModel.UserToken) UserTokenView {
	return UserTokenView{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	}
}

func UserForgottenPasswordViewToUserForgottenPassword(user UserForgottenPasswordView) userModel.UserForgottenPassword {
	return userModel.UserForgottenPassword{
		Email: user.Email,
	}
}

func UserResetPasswordViewToUserResetPassword(user UserResetPasswordView) userModel.UserResetPassword {
	return userModel.UserResetPassword{
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
	}
}
