package model

import (
	"time"

	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	mongoModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersRepository struct {
	Users              []UserRepository
	PaginationResponse common.PaginationResponse
}

type UserRepository struct {
	mongoModel.BaseEntity `bson:",inline"`
	Username              string `bson:"username"`
	Email                 string `bson:"email"`
	Password              string `bson:"password"`
	Role                  string `bson:"role"`
	Verified              bool   `bson:"verified"`
}

type UserCreateRepository struct {
	Username         string    `bson:"username"`
	Email            string    `bson:"email"`
	Password         string    `bson:"password"`
	Role             string    `bson:"role"`
	Verified         bool      `bson:"verified"`
	VerificationCode string    `bson:"verification_code"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}

type UserUpdateRepository struct {
	UserID    primitive.ObjectID `bson:"_id"`
	Username  string             `bson:"username"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type UserForgottenPasswordRepository struct {
	ResetToken  string    `bson:"reset_token"`
	ResetExpiry time.Time `bson:"reset_expiry"`
}

type UserResetPasswordRepository struct {
	Password string `bson:"password"`
}

type UserResetExpiryRepository struct {
	ResetExpiry time.Time `bson:"reset_expiry"`
}

func NewUserCreateRepository(username, email, password, role string, verified bool, verificationCode string) UserCreateRepository {
	return UserCreateRepository{
		Username:         username,
		Email:            email,
		Password:         password,
		Role:             role,
		Verified:         verified,
		VerificationCode: verificationCode,
	}
}

func NewUserUpdateRepository(userID primitive.ObjectID, username string) UserUpdateRepository {
	return UserUpdateRepository{
		UserID:   userID,
		Username: username,
	}
}

func NewUserForgottenPasswordRepository(resetToken string, resetExpiry time.Time) UserForgottenPasswordRepository {
	return UserForgottenPasswordRepository{
		ResetToken:  resetToken,
		ResetExpiry: resetExpiry,
	}
}

func NewUserResetPasswordRepository(password string) UserResetPasswordRepository {
	return UserResetPasswordRepository{
		Password: password,
	}
}

func NewUserResetExpiryRepository(resetExpiry time.Time) UserResetExpiryRepository {
	return UserResetExpiryRepository{
		ResetExpiry: resetExpiry,
	}
}

func NewUsersRepository(users []UserRepository) UsersRepository {
	return UsersRepository{
		Users: users,
	}
}
