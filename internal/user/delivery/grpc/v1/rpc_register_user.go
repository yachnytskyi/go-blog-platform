package v1

import (
	"context"

	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func (userGrpcServer *UserGrpcServer) Register(ctx context.Context, request *pb.UserCreate) (*pb.GenericResponse, error) {
	user := user.UserCreate{
		Username:        request.GetName(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		PasswordConfirm: request.GetPasswordConfirm(),
	}

	createdUser := userGrpcServer.userUseCase.Register(ctx, user)
	if createdUser.Error != nil {

		switch errorType := createdUser.Error.(type) {
		case *domain.ValidationError:
			return nil, errorType
		case *domain.InternalError:
			return nil, errorType
		default:
			var defaultError *domain.InternalError = new(domain.InternalError)
			defaultError.Notification = "reason:" + " something went wrong, please repeat later"
			return nil, errorType
		}

	}

	message := "We sent an email with a verification code to " + user.Email
	response := &pb.GenericResponse{
		Status:  "success",
		Message: message,
	}

	return response, nil
}
