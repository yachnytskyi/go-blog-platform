package v1

import (
	"context"
	"strings"

	postProtobufV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/grpc/v1/model/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (PostGrpcServer *PostGrpcServer) DeletePostById(ctx context.Context, postData *postProtobufV1.PostById) (*postProtobufV1.PostDeleteView, error) {
	postID := postData.GetPostID()
	userID := postData.GetUserID()

	err := PostGrpcServer.postUseCase.DeletePostByID(ctx, postID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			return nil, status.Errorf(codes.NotFound, err.Error(), "error")
		}

		return nil, status.Errorf(codes.Internal, err.Error(), "error")
	}

	postDeleteView := &postProtobufV1.PostDeleteView{
		Success: true,
	}

	return postDeleteView, nil
}
