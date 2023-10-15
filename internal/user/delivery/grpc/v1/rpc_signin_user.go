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
	accessToken, createTokenError := httpUtility.CreateToken(userGrpcServer.applicationConfig.AccessToken.ExpiredIn, user.UserID, userGrpcServer.applicationConfig.AccessToken.PrivateKey)
	if createTokenError != nil {
		return nil, status.Errorf(codes.PermissionDenied, createTokenError.Error())
	}
	refreshToken, createTokenError := httpUtility.CreateToken(userGrpcServer.applicationConfig.RefreshToken.ExpiredIn, user.UserID, userGrpcServer.applicationConfig.RefreshToken.PrivateKey)
	if createTokenError != nil {
		return nil, status.Errorf(codes.PermissionDenied, createTokenError.Error())
	}

	response := &pb.LoginUserView{
		Status:       "success",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return response, nil
}
