package gapi

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	"github.com/yachnytskyi/golang-mongo-grpc/pb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	config         config.Config
	userService    user.Service
	userCollection *mongo.Collection
}

func NewGrpcServer(config config.Config, userService user.Service, userCollection *mongo.Collection) (*Server, error) {

	server := &Server{
		config:         config,
		userService:    userService,
		userCollection: userCollection,
	}

	return server, nil
}
