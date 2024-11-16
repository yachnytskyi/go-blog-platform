package model

import (
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
)

func UserCreateViewToUserCreateMapper(userCreateView UserCreateView) user.UserCreate {
	return user.NewUserCreate(
		userCreateView.Username,
		userCreateView.Email,
		userCreateView.Password,
		userCreateView.PasswordConfirm,
	)
}

func UserUpdateViewToUserUpdateMapper(userUpdateView UserUpdateView) user.UserUpdate {
	return user.NewUserUpdate(
		userUpdateView.ID,
		userUpdateView.Username,
	)
}

func UserLoginViewToUserLoginMapper(userLoginView UserLoginView) user.UserLogin {
	return user.NewUserLogin(
		userLoginView.Email,
		userLoginView.Password,
	)
}

func UserForgottenPasswordViewToUserForgottenPassword(userForgottenPasswordView UserForgottenPasswordView) user.UserForgottenPassword {
	return user.NewUserForgottenPassword(
		userForgottenPasswordView.Email,
	)
}

func UserResetPasswordViewToUserResetPassword(userResetPasswordView UserResetPasswordView) user.UserResetPassword {
	return user.NewUserResetPassword(
		userResetPasswordView.ResetToken,
		userResetPasswordView.Password,
		userResetPasswordView.PasswordConfirm,
	)
}

func UsersToUsersViewMapper(users user.Users) UsersView {
	usersView := make([]UserView, len(users.Users))
	for index, user := range users.Users {
		usersView[index] = UserToUserViewMapper(user)
	}

	return NewUsersView(usersView, model.NewHTTPPaginationResponse(
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

func UserToUserViewMapper(user user.User) UserView {
	return NewUserView(
		user.ID,
		user.Username,
		user.Email,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)
}

func UserLoginToUserLoginViewMapper(userLogin user.UserLogin) UserLoginView {
	return NewUserLoginView(
		userLogin.Email,
		userLogin.Password,
	)
}

func UserTokenToUserTokenViewMapper(userToken user.UserToken) UserTokenView {
	return NewUserTokenView(
		userToken.AccessToken,
		userToken.RefreshToken,
	)
}

func UserForgottenPasswordToUserForgottenPasswordViewMapper(userForgottenPassword user.UserForgottenPassword) UserForgottenPasswordView {
	return NewUserForgottenPasswordView(
		userForgottenPassword.Email,
	)
}
