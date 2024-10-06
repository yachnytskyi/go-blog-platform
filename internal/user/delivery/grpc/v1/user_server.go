package v1

import (
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserGrpcServer struct {
	Config *config.ApplicationConfig
	Logger interfaces.Logger
	pb.UnimplementedUserUseCaseServer
	// applicationConfig config.ApplicationConfig
	userUseCase    user.UserUseCase
	userCollection *mongo.Collection
}

// func NewGrpcUserServer(config *ApplicationConfig, userUseCase user.UserUseCase, userCollection *mongo.Collection) (*UserGrpcServer, error) {
// 	applicationConfig := config.GetGRPC().
// 	userGrpcServer := &UserGrpcServer{
// 		applicationConfig: applicationConfig,
// 		userUseCase:       userUseCase,
// 		userCollection:    userCollection,
// 	}

// 	return userGrpcServer, nil
// }
