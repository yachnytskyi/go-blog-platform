package gapi

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	"github.com/yachnytskyi/golang-mongo-grpc/pb"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	config         config.Config
	userService    user.Service
	userCollection *mongo.Collection
}

func NewGrpcUserServer(config config.Config, userService user.Service, userCollection *mongo.Collection) (*UserServer, error) {

	userServer := &UserServer{
		config:         config,
		userService:    userService,
		userCollection: userCollection,
	}

	return userServer, nil
}
