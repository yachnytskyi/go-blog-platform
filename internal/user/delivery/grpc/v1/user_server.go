package v1

import (
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserGrpcServer struct {
	Config interfaces.Config
	Logger interfaces.Logger
	pb.UnimplementedUserUseCaseServer
	// applicationConfig config.ApplicationConfig
	userUseCase    interfaces.UserUseCase
	userCollection *mongo.Collection
}

// func NewGrpcUserServer(config interfaces.Config, userUseCase user.UserUseCase, userCollection *mongo.Collection) (*UserGrpcServer, error) {
// 	applicationConfig := config.GetGRPC().
// 	userGrpcServer := &UserGrpcServer{
// 		applicationConfig: applicationConfig,
// 		userUseCase:       userUseCase,
// 		userCollection:    userCollection,
// 	}

// 	return userGrpcServer, nil
// }
