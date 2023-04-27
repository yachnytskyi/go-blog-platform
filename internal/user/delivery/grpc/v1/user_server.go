package v1

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	pb "github.com/yachnytskyi/golang-mongo-grpc/pkg/proto-generated"
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
