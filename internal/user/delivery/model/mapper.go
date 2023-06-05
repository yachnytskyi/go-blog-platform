package model

import userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"

func UsersToUsersViewMapper(users *userModel.Users, limit int) UsersView {
	usersView := make([]*UserView, 0, limit)

	for _, user := range users.Users {
		userView := &UserView{}
		userView.UserID = user.UserID
		userView.Name = user.Name
		userView.Email = user.Email
		userView.Role = user.Role
		userView.CreatedAt = user.CreatedAt
		userView.UpdatedAt = user.UpdatedAt
		usersView = append(usersView, userView)
	}

	return UsersView{
		UsersView: usersView,
	}
}

func UserToUserViewMapper(user *userModel.User) UserView {
	return UserView{
		UserID:    user.UserID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserCreateViewToUserCreateMapper(user *UserCreateView) userModel.UserCreate {
	return userModel.UserCreate{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
	}
}
