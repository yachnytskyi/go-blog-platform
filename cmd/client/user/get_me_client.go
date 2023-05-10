package user

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/model/pb"
	"google.golang.org/grpc"
)

type GetMeClient struct {
	usecase pb.UserUseCaseClient
}

func NewGetMeClient(connection *grpc.ClientConn) *GetMeClient {
	usecase := pb.NewUserUseCaseClient(connection)

	return &GetMeClient{usecase}
}

func (getMeClient *GetMeClient) GetMeUser(credentials *pb.GetMeRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Microsecond*5000))
	defer cancel()

	response, err := getMeClient.usecase.GetMe(ctx, credentials)

	if err != nil {
		log.Fatalf("GetMe: %v", err)
	}

	fmt.Println(response)

}
