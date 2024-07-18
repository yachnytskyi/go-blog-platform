package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
)

func UserCreateViewToUserCreateMapper(user UserCreateView) userModel.UserCreate {
	return userModel.NewUserCreate(
		user.Name,
		user.Email,
		user.Password,
		user.PasswordConfirm,
	)
}

func UserUpdateViewToUserUpdateMapper(user UserUpdateView) userModel.UserUpdate {
	return userModel.NewUserUpdate(
		user.ID,
		user.Name,
	)
}

func UserLoginViewToUserLoginMapper(user UserLoginView) userModel.UserLogin {
	return userModel.NewUserLogin(
		user.Email,
		user.Password,
	)
}

func UserForgottenPasswordViewToUserForgottenPassword(userForgottenPasswordView UserForgottenPasswordView) userModel.UserForgottenPassword {
	return userModel.NewUserForgottenPassword(
		userForgottenPasswordView.Email,
	)
}

func UserResetPasswordViewToUserResetPassword(user UserResetPasswordView) userModel.UserResetPassword {
	return userModel.NewUserResetPassword(
		user.ResetToken,
		user.Password,
		user.PasswordConfirm,
	)
}

func UsersToUsersViewMapper(users userModel.Users) UsersView {
	usersView := make([]UserView, len(users.Users))
	for index, user := range users.Users {
		usersView[index] = UserToUserViewMapper(user)
	}

	return NewUsersView(usersView, httpModel.NewHTTPPaginationResponse(
		users.PaginationResponse.Page,
		users.PaginationResponse.TotalPages,
		users.PaginationResponse.PagesLeft,
		users.PaginationResponse.ItemsLeft,
		users.PaginationResponse.TotalItems,
		users.PaginationResponse.Limit,
		users.PaginationResponse.OrderBy,
		users.PaginationResponse.SortOrder,
		users.PaginationResponse.PageLinks,
	))
}

func UserToUserViewMapper(user userModel.User) UserView {
	return NewUserView(
		user.ID,
		user.CreatedAt,
		user.UpdatedAt,
		user.Name,
		user.Email,
		user.Role,
	)
}

func UserLoginToUserLoginViewMapper(user userModel.UserLogin) UserLoginView {
	return NewUserLoginView(
		user.Email,
		user.Password,
	)
}

func UserTokenToUserTokenViewMapper(user userModel.UserToken) UserTokenView {
	return NewUserTokenView(
		user.AccessToken,
		user.RefreshToken,
	)
}

func UserForgottenPasswordToUserForgottenPasswordViewMapper(userForgottenPassword userModel.UserForgottenPassword) UserForgottenPasswordView {
	return NewUserForgottenPasswordView(
		userForgottenPassword.Email,
	)
}
