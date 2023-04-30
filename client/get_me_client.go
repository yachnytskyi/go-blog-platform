package client

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/yachnytskyi/golang-mongo-grpc/pkg/pb"
	"google.golang.org/grpc"
)

type GetMeClient struct {
	service pb.UserServiceClient
}

func NewGetMeClient(connection *grpc.ClientConn) *GetMeClient {
	service := pb.NewUserServiceClient(connection)

	return &GetMeClient{service}
}

func (getMeClient *GetMeClient) GetMeUser(credentials *pb.GetMeRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Microsecond*5000))
	defer cancel()

	response, err := getMeClient.service.GetMe(ctx, credentials)

	if err != nil {
		log.Fatalf("GetMe: %v", err)
	}

	fmt.Println(response)

}
