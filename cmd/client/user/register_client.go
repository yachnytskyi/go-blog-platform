package user

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/model/pb"
	"google.golang.org/grpc"
)

type RegisterUserClient struct {
	usecase pb.UserUseCaseClient
}

func NewRegisterUserClient(connection *grpc.ClientConn) *RegisterUserClient {
	usecase := pb.NewUserUseCaseClient(connection)

	return &RegisterUserClient{usecase}
}

func (registerUserClient *RegisterUserClient) Register(credentials *pb.RegisterUserInput) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*5000))
	defer cancel()

	response, err := registerUserClient.usecase.Register(ctx, credentials)

	if err != nil {
		log.Fatalf("Register: %v", err)
	}

	fmt.Println(response)
}
