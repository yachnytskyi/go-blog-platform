package model

import userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"

func UsersToUsersViewMapper(users userModel.Users) UsersView {
	usersView := make([]UserView, 0, len(users.Users))
	for _, user := range users.Users {
		userView := UserToUserViewMapper(user)
		usersView = append(usersView, userView)
	}

	return UsersView{
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
