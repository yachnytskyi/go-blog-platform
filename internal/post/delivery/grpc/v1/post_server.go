package v1

import (
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/grpc/v1/model/pb"
)

type PostGrpcServer struct {
	pb.UnimplementedPostUseCaseServer
	postUseCase interfaces.PostUseCase
}

func NewGrpcPostServer(postUseCase interfaces.PostUseCase) (*PostGrpcServer, error) {
	postGrpcServer := &PostGrpcServer{
		postUseCase: postUseCase,
	}

	return postGrpcServer, nil
}
