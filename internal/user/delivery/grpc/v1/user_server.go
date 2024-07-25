package v1

import (
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserGrpcServer struct {
	Config model.Config
	Logger model.Logger
	pb.UnimplementedUserUseCaseServer
	// applicationConfig config.ApplicationConfig
	userUseCase    user.UserUseCase
	userCollection *mongo.Collection
}

// func NewGrpcUserServer(config model.Config, userUseCase user.UserUseCase, userCollection *mongo.Collection) (*UserGrpcServer, error) {
// 	applicationConfig := config.GetGRPC().
// 	userGrpcServer := &UserGrpcServer{
// 		applicationConfig: applicationConfig,
// 		userUseCase:       userUseCase,
// 		userCollection:    userCollection,
// 	}

// 	return userGrpcServer, nil
// }
