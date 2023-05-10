package v1

import (
	"context"

	repositoryUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/utility"
	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/model/pb"
	httpUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/utility"

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

	if err := repositoryUtility.VerifyPassword(user.Password, request.GetPassword()); err != nil {

		return nil, status.Errorf(codes.InvalidArgument, "Invalid email or Password")
	}

	// Generate tokens
	accessToken, err := httpUtility.CreateToken(userServer.config.AccessTokenExpiresIn, user.UserID, userServer.config.AccessTokenPrivateKey)

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	refreshToken, err := httpUtility.CreateToken(userServer.config.RefreshTokenExpiresIn, user.UserID, userServer.config.RefreshTokenPrivateKey)

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
