package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// UserViewToUserMapper maps a UserView instance to a userModel.User instance.
// Parameters:
// - location: The location for parsing dates.
// - userView: The UserView instance to be mapped.
// Returns:
// - commonModel.Result[userModel.User]: The result containing the mapped User model or an error.
func UserViewToUserMapper(location string, userView UserView) commonModel.Result[userModel.User] {
	createdAt := commonUtility.ParseDate(location, userView.CreatedAt)
	if validator.IsError(createdAt.Error) {
		return commonModel.NewResultOnFailure[userModel.User](createdAt.Error)
	}

	updatedAt := commonUtility.ParseDate(location, userView.UpdatedAt)
	if validator.IsError(updatedAt.Error) {
		return commonModel.NewResultOnFailure[userModel.User](updatedAt.Error)
	}

	return commonModel.NewResultOnSuccess(
		userModel.NewUser(
			userView.ID,
			createdAt.Data,
			updatedAt.Data,
			userView.Name,
			userView.Email,
			userView.Role,
		),
	)
}

// UserCreateViewToUserCreateMapper maps a UserCreateView instance to a userModel.UserCreate instance.
// Parameters:
// - user: The UserCreateView instance to be mapped.
// Returns:
// - userModel.UserCreate: The mapped UserCreate model.
func UserCreateViewToUserCreateMapper(user UserCreateView) userModel.UserCreate {
	return userModel.NewUserCreate(
		user.Name,
		user.Email,
		user.Password,
		user.PasswordConfirm,
	)
}

// UserUpdateViewToUserUpdateMapper maps a UserUpdateView instance to a userModel.UserUpdate instance.
// Parameters:
// - user: The UserUpdateView instance to be mapped.
// Returns:
// - userModel.UserUpdate: The mapped UserUpdate model.
func UserUpdateViewToUserUpdateMapper(user UserUpdateView) userModel.UserUpdate {
	return userModel.NewUserUpdate(
		user.ID,
		user.Name,
	)
}

// UserLoginViewToUserLoginMapper maps a UserLoginView instance to a userModel.UserLogin instance.
// Parameters:
// - user: The UserLoginView instance to be mapped.
// Returns:
// - userModel.UserLogin: The mapped UserLogin model.
func UserLoginViewToUserLoginMapper(user UserLoginView) userModel.UserLogin {
	return userModel.NewUserLogin(
		user.Email,
		user.Password,
	)
}

// UserForgottenPasswordViewToUserForgottenPassword maps a UserForgottenPasswordView instance to a userModel.UserForgottenPassword instance.
// Parameters:
// - userForgottenPasswordView: The UserForgottenPasswordView instance to be mapped.
// Returns:
// - userModel.UserForgottenPassword: The mapped UserForgottenPassword model.
func UserForgottenPasswordViewToUserForgottenPassword(userForgottenPasswordView UserForgottenPasswordView) userModel.UserForgottenPassword {
	return userModel.NewUserForgottenPassword(
		userForgottenPasswordView.Email,
	)
}

// UserResetPasswordViewToUserResetPassword maps a UserResetPasswordView instance to a userModel.UserResetPassword instance.
// Parameters:
// - user: The UserResetPasswordView instance to be mapped.
// Returns:
// - userModel.UserResetPassword: The mapped UserResetPassword model.
func UserResetPasswordViewToUserResetPassword(user UserResetPasswordView) userModel.UserResetPassword {
	return userModel.NewUserResetPassword(
		user.ResetToken,
		user.Password,
		user.PasswordConfirm,
	)
}

// UsersToUsersViewMapper maps a userModel.Users instance to a UsersView instance.
// Parameters:
// - users: The userModel.Users instance to be mapped.
// Returns:
// - UsersView: The mapped UsersView instance.
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

// UserToUserViewMapper maps a userModel.User instance to a UserView instance.
// Parameters:
// - user: The userModel.User instance to be mapped.
// Returns:
// - UserView: The mapped UserView instance.
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

// UserLoginToUserLoginViewMapper maps a userModel.UserLogin instance to a UserLoginView instance.
// Parameters:
// - user: The userModel.UserLogin instance to be mapped.
// Returns:
// - UserLoginView: The mapped UserLoginView instance.
func UserLoginToUserLoginViewMapper(user userModel.UserLogin) UserLoginView {
	return NewUserLoginView(
		user.Email,
		user.Password,
	)
}

// UserTokenToUserTokenViewMapper maps a userModel.UserToken instance to a UserTokenView instance.
// Parameters:
// - user: The userModel.UserToken instance to be mapped.
// Returns:
// - UserTokenView: The mapped UserTokenView instance.
func UserTokenToUserTokenViewMapper(user userModel.UserToken) UserTokenView {
	return NewUserTokenView(
		user.AccessToken,
		user.RefreshToken,
	)
}

// UserForgottenPasswordToUserForgottenPasswordViewMapper maps a userModel.UserForgottenPassword instance to a UserForgottenPasswordView instance.
// Parameters:
// - userForgottenPassword: The userModel.UserForgottenPassword instance to be mapped.
// Returns:
// - UserForgottenPasswordView: The mapped UserForgottenPasswordView instance.
func UserForgottenPasswordToUserForgottenPasswordViewMapper(userForgottenPassword userModel.UserForgottenPassword) UserForgottenPasswordView {
	return NewUserForgottenPasswordView(
		userForgottenPassword.Email,
	)
}
