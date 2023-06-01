package model

import (
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
)

func UserToUserViewMapper(user *model.User) UserView {
	return UserView{
		UserID:    user.UserID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UsersToUsersViewMapper(users *model.Users) UsersView {
	usersView := make([]*UserView, 0, 10)

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
