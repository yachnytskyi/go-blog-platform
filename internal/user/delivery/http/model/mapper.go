package model

import (
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
)

func TokenStringToTokenViewMapper(token string) TokenView {
	return TokenView{
		Token: token,
	}
}

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
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
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
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserLoginViewToUserLoginMapper(user UserLoginView) userModel.UserLogin {
	return userModel.UserLogin{
		Email:    user.Email,
		Password: user.Password,
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
