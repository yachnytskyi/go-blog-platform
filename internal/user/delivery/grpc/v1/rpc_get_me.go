package v1

import (
	"context"

	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (userGrpcServer *UserGrpcServer) GetMe(ctx context.Context, request *pb.GetMeRequest) (*pb.UserView, error) {
	userID := request.GetId()
	user := userGrpcServer.userUseCase.GetUserById(ctx, userID)

	if user.Error != nil {
		return nil, status.Errorf(codes.Unimplemented, user.Error.Error())
	}

	response := &pb.UserView{
		User: &pb.User{
			Id:        user.Data.UserID,
			Name:      user.Data.Name,
			Email:     user.Data.Email,
			Role:      user.Data.Role,
			CreatedAt: timestamppb.New(user.Data.CreatedAt),
			UpdatedAt: timestamppb.New(user.Data.UpdatedAt),
		},
	}

	return response, nil
}
