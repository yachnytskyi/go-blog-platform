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
