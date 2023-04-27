package main

import (
	"log"

	"github.com/yachnytskyi/golang-mongo-grpc/client"
	pb "github.com/yachnytskyi/golang-mongo-grpc/pkg/proto-generated"
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
		registerUserClient := client.NewRegisterUserClient(connect)
		createdUser := &pb.RegisterUserInput{
			Name:            "Test Test",
			Email:           "test100@gmail.com",
			Password:        "somepassword",
			PasswordConfirm: "somepassword",
		}

		registerUserClient.Register(createdUser)
	}

	// Login.
	if true {
		loginUserClient := client.NewLoginUserClient(connect)

		credentials := &pb.LoginUserInput{
			Email:    "test100@gmail.com",
			Password: "somepassword",
		}

		loginUserClient.Login(credentials)
	}

	// Get Me.
	if false {
		getMeClient := client.NewGetMeClient(connect)
		id := &pb.GetMeRequest{
			Id: "628cffb91e50302d360c1a2c",
		}
		getMeClient.GetMeUser(id)
	}

}
