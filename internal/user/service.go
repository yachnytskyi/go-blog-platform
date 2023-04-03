package user

import "github.com/yachnytskyi/golang-mongo-grpc/models"

type Service interface {
	Register(*models.UserCreate) (*models.UserFullResponse, error)
	Login(*models.UserSignIn) (*models.UserFullResponse, error)
	UserGetById(string) (*models.UserFullResponse, error)
	UserGetByEmail(string) (*models.UserFullResponse, error)
}
