package user

import (
	"context"
	"fmt"
	"log"
	"time"

	userProtobufV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	"google.golang.org/grpc"
)

type LoginUserClient struct {
	usecase userProtobufV1.UserUseCaseClient
}

func NewLoginUserClient(connection *grpc.ClientConn) *LoginUserClient {
	usecase := userProtobufV1.NewUserUseCaseClient(connection)

	return &LoginUserClient{usecase}
}

func (loginUserClient *LoginUserClient) Login(credentials *userProtobufV1.LoginUserInput) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := loginUserClient.usecase.Login(ctx, credentials)

	if err != nil {
		log.Fatalf("Login: %v", err)
	}

	fmt.Println(response)
}
