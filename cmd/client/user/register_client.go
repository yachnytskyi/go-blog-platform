package user

import (
	"context"
	"fmt"
	"log"
	"time"

	userProtobufV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	"google.golang.org/grpc"
)

type RegisterUserClient struct {
	usecase userProtobufV1.UserUseCaseClient
}

func NewRegisterUserClient(connection *grpc.ClientConn) *RegisterUserClient {
	usecase := userProtobufV1.NewUserUseCaseClient(connection)

	return &RegisterUserClient{usecase}
}

func (registerUserClient *RegisterUserClient) Register(credentials *userProtobufV1.UserCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	response, err := registerUserClient.usecase.Register(ctx, credentials)

	if err != nil {
		log.Fatalf("Register: %v", err)
	}

	fmt.Println(response)
}
