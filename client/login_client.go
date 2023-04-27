package client

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/yachnytskyi/golang-mongo-grpc/pkg/proto-generated"
	"google.golang.org/grpc"
)

type LoginUserClient struct {
	service pb.UserServiceClient
}

func NewLoginUserClient(connection *grpc.ClientConn) *LoginUserClient {
	service := pb.NewUserServiceClient(connection)

	return &LoginUserClient{service}
}

func (loginUserClient *LoginUserClient) Login(credentials *pb.LoginUserInput) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := loginUserClient.service.Login(ctx, credentials)

	if err != nil {
		log.Fatalf("Login: %v", err)
	}

	fmt.Println(response)
}
