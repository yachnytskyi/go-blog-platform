package v1

import (
	"context"

	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (userGrpcServer *UserGrpcServer) Login(ctx context.Context, request *pb.LoginUser) (*pb.LoginUserView, error) {
	user := userGrpcServer.userUseCase.GetUserByEmail(ctx, request.GetEmail())

	if user.Error != nil {
		return nil, status.Errorf(codes.Internal, user.Error.Error())
	}

	if !user.Data.Verified {

		return nil, status.Errorf(codes.PermissionDenied, "You are not verified, please verify your email to login")
	}

	// if userUseCase.ArePasswordsEqual(user.Password, request.GetPassword()) {
	// 	return nil, status.Errorf(codes.InvalidArgument, "Invalid email or Password")
	// }

	// Generate tokens
	accessToken, createTokenError := domainUtility.CreateToken(userGrpcServer.applicationConfig.AccessToken.ExpiredIn, user.Data.UserID, userGrpcServer.applicationConfig.AccessToken.PrivateKey)
	if createTokenError != nil {
		return nil, status.Errorf(codes.PermissionDenied, createTokenError.Error())
	}
	refreshToken, createTokenError := domainUtility.CreateToken(userGrpcServer.applicationConfig.RefreshToken.ExpiredIn, user.Data.UserID, userGrpcServer.applicationConfig.RefreshToken.PrivateKey)
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
