package main

import (
	"log"

	userClient "github.com/yachnytskyi/golang-mongo-grpc/cmd/client/user"
	userProtobufV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "0.0.0.0:8081"
)

func main() {
	connect, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer connect.Close()

	// Register.
	if false {
		registerUserClient := userClient.NewRegisterUserClient(connect)
		createdUser := &userProtobufV1.UserCreate{
			Name:            "Test Test",
			Email:           "test100@gmail.com",
			Password:        "somepassword",
			PasswordConfirm: "somepassword",
		}

		registerUserClient.Register(createdUser)
	}

	// Login.
	if true {
		loginUserClient := userClient.NewLoginUserClient(connect)

		credentials := &userProtobufV1.LoginUser{
			Email:    "test100@gmail.com",
			Password: "somepassword",
		}

		loginUserClient.Login(credentials)
	}

	// Get Me.
	if false {
		getMeClient := userClient.NewGetMeClient(connect)
		id := &userProtobufV1.GetMeRequest{
			Id: "628cffb91e50302d360c1a2c",
		}
		getMeClient.GetMeUser(id)
	}

}
