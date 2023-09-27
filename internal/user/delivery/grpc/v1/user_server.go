package v1

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserGrpcServer struct {
	pb.UnimplementedUserUseCaseServer
	config         config.Config
	userUseCase    user.UserUseCase
	userCollection *mongo.Collection
}

func NewGrpcUserServer(config config.Config, userUseCase user.UserUseCase, userCollection *mongo.Collection) (*UserGrpcServer, error) {

	userGrpcServer := &UserGrpcServer{
		config:         config,
		userUseCase:    userUseCase,
		userCollection: userCollection,
	}

	return userGrpcServer, nil
}
