package v1

import (
	"context"
	"strings"

	postProtobufV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/grpc/v1/model/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (PostGrpcServer *PostGrpcServer) GetPostById(ctx context.Context, postData *postProtobufV1.PostById) (*postProtobufV1.PostView, error) {
	postID := postData.GetPostID()

	post, err := PostGrpcServer.postUseCase.GetPostById(ctx, postID)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			return nil, status.Errorf(codes.NotFound, err.Error(), "error")
		}

		return nil, status.Errorf(codes.Internal, err.Error(), "error")
	}

	postView := &postProtobufV1.PostView{
		Post: &postProtobufV1.Post{
			PostID:    post.PostID,
			UserID:    post.UserID,
			Title:     post.Title,
			Content:   post.Content,
			Image:     post.Image,
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		},
	}

	return postView, nil
}
