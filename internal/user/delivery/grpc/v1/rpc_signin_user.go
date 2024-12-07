package v1

import (
	"context"
	"time"

	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	model "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	utility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	location = "internal.user.delivery.grpc.v1."
)

func (userGrpcServer *UserGrpcServer) Login(ctx context.Context, request *pb.LoginUser) (*pb.LoginUserView, error) {
	user := userGrpcServer.userUseCase.GetUserByEmail(ctx, request.GetEmail())

	if user.Error != nil {
		return nil, status.Errorf(codes.Internal, user.Error.Error(), "error")
	}

	if !user.Data.Verified {

		return nil, status.Errorf(codes.PermissionDenied, "You are not verified, please verify your email to login")
	}

	// if userUseCase.ArePasswordsEqual(user.Password, request.GetPassword()) {
	// 	return nil, status.Errorf(codes.InvalidArgument, "Invalid email or Password")
	// }

	// Generate the UserTokenPayload.
	userTokenPayload := model.NewUserTokenPayload(user.Data.ID, user.Data.Role)

	// Generate tokens.
	accessToken := utility.GenerateJWTToken(
		userGrpcServer.Logger,
		location+"Login",
		"",
		time.Duration(2),
		// userGrpcServer.applicationConfig.AccessToken.PrivateKey,
		// userGrpcServer.applicationConfig.AccessToken.ExpiredIn,
		userTokenPayload,
	)
	if validator.IsError(accessToken.Error) {
		return nil, status.Errorf(codes.PermissionDenied, accessToken.Error.Error(), "error")
	}

	refreshToken := utility.GenerateJWTToken(
		userGrpcServer.Logger,
		location+"Login",
		"",
		time.Duration(2),
		// userGrpcServer.applicationConfig.RefreshToken.PrivateKey,
		// userGrpcServer.applicationConfig.RefreshToken.ExpiredIn,
		userTokenPayload,
	)
	if validator.IsError(refreshToken.Error) {
		return nil, status.Errorf(codes.PermissionDenied, refreshToken.Error.Error(), "error")
	}

	response := &pb.LoginUserView{
		Status:       "success",
		AccessToken:  accessToken.Data,
		RefreshToken: refreshToken.Data,
	}

	return response, nil
}
