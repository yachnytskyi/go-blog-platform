package v1

import (
	"context"
	"strings"

	"github.com/thanhpk/randstr"
	"github.com/yachnytskyi/golang-mongo-grpc/models"
	pb "github.com/yachnytskyi/golang-mongo-grpc/pkg/pb"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (userServer *UserServer) Register(ctx context.Context, request *pb.RegisterUserInput) (*pb.GenericResponse, error) {
	if request.GetPassword() != request.GetPasswordConfirm() {
		return nil, status.Errorf(codes.InvalidArgument, "passwords do not match")
	}

	user := models.UserCreate{
		Name:            request.GetName(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		PasswordConfirm: request.GetPasswordConfirm(),
	}

	createdUser, err := userServer.userService.Register(ctx, &user)

	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())

		}
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	// Generate verification code.
	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	// Update the user in database.
	userServer.userService.UpdateNewRegisteredUserById(ctx, createdUser.UserID.Hex(), "verificationCode", verificationCode)

	var firstName = createdUser.Name
	firstName = utils.UserFirstName(firstName)

	// Send an email.
	emailData := utils.EmailData{
		URL:       userServer.config.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(createdUser, &emailData, "verificationCode.html")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "There was an error sending email: %s", err.Error())

	}

	message := "We sent an email with a verification code to " + createdUser.Email

	response := &pb.GenericResponse{
		Status:  "success",
		Message: message,
	}
	return response, nil
}
