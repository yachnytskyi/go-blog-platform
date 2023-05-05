package client

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/model/pb"
	"google.golang.org/grpc"
)

type LoginUserClient struct {
	usecase pb.UserUseCaseClient
}

func NewLoginUserClient(connection *grpc.ClientConn) *LoginUserClient {
	usecase := pb.NewUserUseCaseClient(connection)

	return &LoginUserClient{usecase}
}

func (loginUserClient *LoginUserClient) Login(credentials *pb.LoginUserInput) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := loginUserClient.usecase.Login(ctx, credentials)

	if err != nil {
		log.Fatalf("Login: %v", err)
	}

	fmt.Println(response)
}
