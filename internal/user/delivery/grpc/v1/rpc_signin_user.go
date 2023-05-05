package v1

import (
	"context"

	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/model/pb"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (userServer *UserServer) Login(ctx context.Context, request *pb.LoginUserInput) (*pb.LoginUserResponse, error) {
	user, err := userServer.userUseCase.GetUserByEmail(ctx, request.GetEmail())

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if !user.Verified {

		return nil, status.Errorf(codes.PermissionDenied, "You are not verified, please verify your email to login")
	}

	if err := utils.VerifyPassword(user.Password, request.GetPassword()); err != nil {

		return nil, status.Errorf(codes.InvalidArgument, "Invalid email or Password")
	}

	// Generate tokens
	accessToken, err := utils.CreateToken(userServer.config.AccessTokenExpiresIn, user.UserID, userServer.config.AccessTokenPrivateKey)

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	refreshToken, err := utils.CreateToken(userServer.config.RefreshTokenExpiresIn, user.UserID, userServer.config.RefreshTokenPrivateKey)

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	response := &pb.LoginUserResponse{
		Status:       "success",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}
