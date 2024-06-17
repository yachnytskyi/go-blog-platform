package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func UsersToUsersViewMapper(users userModel.Users) UsersView {
	usersView := make([]UserView, len(users.Users))
	for index, user := range users.Users {
		usersView[index] = UserToUserViewMapper(user)
	}

	return UsersView{
		PaginationResponse: httpModel.PaginationResponse{
			CurrentPage: users.PaginationResponse.Page,
			TotalPages:  users.PaginationResponse.TotalPages,
			PagesLeft:   users.PaginationResponse.PagesLeft,
			TotalItems:  users.PaginationResponse.TotalItems,
			ItemsLeft:   users.PaginationResponse.ItemsLeft,
			Limit:       users.PaginationResponse.Limit,
			OrderBy:     users.PaginationResponse.OrderBy,
			SortOrder:   users.PaginationResponse.SortOrder,
			PageLinks:   users.PaginationResponse.PageLinks,
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

func UserViewToUserMapper(location string, userView UserView) commonModel.Result[userModel.User] {
	var httpInternalErrorsView httpError.HttpInternalErrorsView
	created_at, parseDateError := common.ParseDate(location+".UserViewToUserMapper.created_at", userView.CreatedAt)
	if validator.IsError(parseDateError) {
		httpInternalErrorsView = append(httpInternalErrorsView, parseDateError)
	}

	updated_at, parseDateError := common.ParseDate(location+"UserViewToUserMapper.updated_at", userView.UpdatedAt)
	if validator.IsError(parseDateError) {
		httpInternalErrorsView = append(httpInternalErrorsView, parseDateError)
	}

	if validator.IsSliceNotEmpty(httpInternalErrorsView) {
		logging.Logger(httpInternalErrorsView)
		return commonModel.NewResultOnFailure[userModel.User](httpInternalErrorsView)
	}

	user := userModel.User{
		UserID:    userView.UserID,
		Name:      userView.Name,
		Email:     userView.Email,
		Role:      userView.Role,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
	}

	return commonModel.NewResultOnSuccess[userModel.User](user)
}
