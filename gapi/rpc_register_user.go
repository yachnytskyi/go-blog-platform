package gapi

import (
	"context"
	"strings"

	"github.com/thanhpk/randstr"
	"github.com/yachnytskyi/golang-mongo-grpc/models"
	"github.com/yachnytskyi/golang-mongo-grpc/pb"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) Register(ctx context.Context, req *pb.RegisterUserInput) (*pb.GenericResponse, error) {
	if req.GetPassword() != req.GetPasswordConfirm() {
		return nil, status.Errorf(codes.InvalidArgument, "passwords do not match")
	}

	user := models.UserCreate{
		Name:            req.GetName(),
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}

	newUser, err := server.userService.Register(ctx, &user)

	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())

		}
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	// Generate verification code.
	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	// Update the user in database.
	server.userService.UpdateNewRegisteredUserById(ctx, newUser.UserID.Hex(), "verificationCode", verificationCode)

	var firstName = newUser.Name
	firstName = utils.UserFirstName(firstName)

	// Send an email.
	emailData := utils.EmailData{
		URL:       server.config.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(newUser, &emailData, "verificationCode.html")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "There was an error sending email: %s", err.Error())

	}

	message := "We sent an email with a verification code to " + newUser.Email

	response := &pb.GenericResponse{
		Status:  "success",
		Message: message,
	}
	return response, nil
}
