package user

import (
	"context"
	"fmt"
	"log"
	"time"

	userProtobufV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1/model/pb"
	"google.golang.org/grpc"
)

type GetMeClient struct {
	usecase userProtobufV1.UserUseCaseClient
}

func NewGetMeClient(connection *grpc.ClientConn) *GetMeClient {
	usecase := userProtobufV1.NewUserUseCaseClient(connection)

	return &GetMeClient{usecase}
}

func (getMeClient *GetMeClient) GetMeUser(credentials *userProtobufV1.GetMeRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Microsecond*5000))
	defer cancel()

	response, err := getMeClient.usecase.GetMe(ctx, credentials)

	if err != nil {
		log.Fatalf("GetMe: %v", err)
	}

	fmt.Println(response)

}
