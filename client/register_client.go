package client

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/yachnytskyi/golang-mongo-grpc/pkg/pb"
	"google.golang.org/grpc"
)

type RegisterUserClient struct {
	service pb.UserServiceClient
}

func NewRegisterUserClient(connection *grpc.ClientConn) *RegisterUserClient {
	service := pb.NewUserServiceClient(connection)

	return &RegisterUserClient{service}
}

func (registerUserClient *RegisterUserClient) Register(credentials *pb.RegisterUserInput) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	response, err := registerUserClient.service.Register(ctx, credentials)

	if err != nil {
		log.Fatalf("Register: %v", err)
	}

	fmt.Println(response)
}
