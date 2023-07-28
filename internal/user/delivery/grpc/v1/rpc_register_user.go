package v1

import (
	"context"

	"github.com/thanhpk/randstr"
	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"

	httpUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/utility"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"
	"google.golang.org/grpc/codes"
)

func (userGrpcServer *UserGrpcServer) Register(ctx context.Context, request *pb.UserCreate) (*pb.GenericResponse, error) {
	user := userModel.UserCreate{
		Name:            request.GetName(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		PasswordConfirm: request.GetPasswordConfirm(),
	}

	createdUser, createdUserErrors := userGrpcServer.userUseCase.Register(ctx, &user)

	if len(createdUserErrors) != 0 {
		for _, userCreateErrorType := range createdUserErrors {
			if userCreateErrorView, ok := userCreateErrorType.(*domainError.ValidationError); ok {
				return nil, userCreateErrorView

			} else {
				var userCreateErrorView *domainError.InternalError = new(domainError.InternalError)

				userCreateErrorView.Location = "UserCreate.Delivery.Grpc.V1.Register.createdUserErrors"
				userCreateErrorView.Code = codes.Internal.String()
				userCreateErrorView.Reason = "reason:" + " something went wrong, please repeat later"

				return nil, userCreateErrorView
			}
		}
	}

	// Generate verification code.
	code := randstr.String(20)

	verificationCode := utility.Encode(code)

	// Update the user in database.
	userGrpcServer.userUseCase.UpdateNewRegisteredUserById(ctx, createdUser.UserID, "verificationCode", verificationCode)

	var firstName = createdUser.Name
	firstName = httpUtility.UserFirstName(firstName)

	// Send an email.
	emailData := httpUtility.EmailData{
		URL:       userGrpcServer.config.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err := httpUtility.SendEmail(createdUser, &emailData, "verificationCode.html")

	if err != nil {
		var userCreateErrorView *domainError.InternalError = new(domainError.InternalError)

		userCreateErrorView.Location = "UserCreate.Delivery.Grpc.V1.Register.createdUserErrors"
		userCreateErrorView.Code = codes.Internal.String()
		userCreateErrorView.Reason = err.Error()

		return nil, userCreateErrorView

	}

	message := "We sent an email with a verification code to " + createdUser.Email

	response := &pb.GenericResponse{
		Status:  "success",
		Message: message,
	}
	return response, nil
}
