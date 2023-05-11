package v1

import (
	"context"

	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (userServer *UserServer) GetMe(ctx context.Context, request *pb.GetMeRequest) (*pb.UserResponse, error) {
	userID := request.GetId()
	user, err := userServer.userUseCase.GetUserById(ctx, userID)

	if err != nil {
		return nil, status.Errorf(codes.Unimplemented, err.Error())
	}

	response := &pb.UserResponse{
		User: &pb.User{
			Id:        user.UserID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}

	return response, nil
}
