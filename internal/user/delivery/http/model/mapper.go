package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	domainModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func UserViewToUserMapper(location string, userView UserView) commonModel.Result[userModel.User] {
	createdAt := commonUtility.ParseDate(location, userView.CreatedAt)
	if validator.IsError(createdAt.Error) {
		return commonModel.NewResultOnFailure[userModel.User](createdAt.Error)
	}

	updatedAt := commonUtility.ParseDate(location, userView.UpdatedAt)
	if validator.IsError(updatedAt.Error) {
		return commonModel.NewResultOnFailure[userModel.User](updatedAt.Error)
	}

	return commonModel.NewResultOnSuccess(userModel.User{
		BaseEntity: domainModel.NewBaseEntity(userView.ID, createdAt.Data, updatedAt.Data),
		Name:       userView.Name,
		Email:      userView.Email,
		Role:       userView.Role,
	})
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
		ID:   user.ID,
		Name: user.Name,
	}
}

func UserLoginViewToUserLoginMapper(user UserLoginView) userModel.UserLogin {
	return userModel.UserLogin{
		Email:    user.Email,
		Password: user.Password,
	}
}

func UserForgottenPasswordViewToUserForgottenPassword(userForgottenPasswordView UserForgottenPasswordView) userModel.UserForgottenPassword {
	return userModel.UserForgottenPassword{
		Email: userForgottenPasswordView.Email,
	}
}

func UserResetPasswordViewToUserResetPassword(user UserResetPasswordView) userModel.UserResetPassword {
	return userModel.UserResetPassword{
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
	}
}

func UsersToUsersViewMapper(users userModel.Users) UsersView {
	usersView := make([]UserView, len(users.Users))
	for index, user := range users.Users {
		usersView[index] = UserToUserViewMapper(user)
	}

	return UsersView{
		HTTPPaginationResponse: httpModel.NewHTTPPaginationResponse(
			users.PaginationResponse.Page,
			users.PaginationResponse.TotalPages,
			users.PaginationResponse.PagesLeft,
			users.PaginationResponse.ItemsLeft,
			users.PaginationResponse.TotalItems,
			users.PaginationResponse.Limit,
			users.PaginationResponse.OrderBy,
			users.PaginationResponse.SortOrder,
			users.PaginationResponse.PageLinks,
		),
		UsersView: usersView,
	}
}

func UserToUserViewMapper(user userModel.User) UserView {
	return UserView{
		BaseEntity: httpModel.NewBaseEntity(user.ID, user.CreatedAt, user.UpdatedAt),
		Name:       user.Name,
		Email:      user.Email,
		Role:       user.Role,
	}
}

func UserLoginToUserLoginViewMapper(user userModel.UserLogin) UserLoginView {
	return UserLoginView{
		Email:    user.Email,
		Password: user.Password,
	}
}

func UserTokenToUserTokenViewMapper(user userModel.UserToken) UserTokenView {
	return UserTokenView{
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	}
}

func UserForgottenPasswordToUserForgottenPasswordViewMapper(userForgottenPassword userModel.UserForgottenPassword) UserForgottenPasswordView {
	return UserForgottenPasswordView{
		Email: userForgottenPassword.Email,
	}
}
