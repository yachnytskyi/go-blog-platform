package model

import (
	"time"

	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	mongoModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersRepository struct {
	Users              []UserRepository
	PaginationResponse commonModel.PaginationResponse
}

type UserRepository struct {
	mongoModel.BaseEntity `bson:",inline"`
	Name                  string `bson:"name"`
	Email                 string `bson:"email"`
	Password              string `bson:"password"`
	Role                  string `bson:"role"`
	Verified              bool   `bson:"verified"`
}

type UserCreateRepository struct {
	Name             string    `bson:"name"`
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
	Name      string             `bson:"name"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type UserForgottenPasswordRepository struct {
	ResetToken  string    `bson:"reset_token"`
	ResetExpiry time.Time `bson:"reset_expiry"`
}

type UserResetPasswordRepository struct {
	Password string `bson:"password"`
}

type UserResetTokenRepository struct {
	ResetExpiry time.Time `bson:"reset_expiry"`
}
