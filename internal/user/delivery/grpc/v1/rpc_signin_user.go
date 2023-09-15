package v1

import (
	"context"

	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	httpUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/utility"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (userGrpcServer *UserGrpcServer) Login(ctx context.Context, request *pb.LoginUser) (*pb.LoginUserView, error) {
	user, err := userGrpcServer.userUseCase.GetUserByEmail(ctx, request.GetEmail())

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if !user.Verified {

		return nil, status.Errorf(codes.PermissionDenied, "You are not verified, please verify your email to login")
	}

	// if userUseCase.ArePasswordsEqual(user.Password, request.GetPassword()) {
	// 	return nil, status.Errorf(codes.InvalidArgument, "Invalid email or Password")
	// }

	// Generate tokens
	accessToken, err := httpUtility.CreateToken(userGrpcServer.config.AccessTokenExpiresIn, user.UserID, userGrpcServer.config.AccessTokenPrivateKey)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}
	refreshToken, err := httpUtility.CreateToken(userGrpcServer.config.RefreshTokenExpiresIn, user.UserID, userGrpcServer.config.RefreshTokenPrivateKey)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	response := &pb.LoginUserView{
		Status:       "success",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response, nil
}
