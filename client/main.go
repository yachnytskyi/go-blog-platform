package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/pb"
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

	client := pb.NewUserServiceClient(connect)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	newUser := &pb.RegisterUserInput{
		Name:            "Test Test",
		Email:           "test100@gmail.com",
		Password:        "somepassword",
		PasswordConfirm: "somepassword",
	}

	response, err := client.Register(ctx, newUser)
	if err != nil {
		log.Fatalf("Register: %v", err)
	}

	fmt.Println(response)

}
