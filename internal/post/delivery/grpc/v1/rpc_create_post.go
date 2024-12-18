package v1

import (
	"context"
	"strings"

	postProtobufV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/grpc/v1/model/pb"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (postGrpcServer *PostGrpcServer) CreatePost(ctx context.Context, createdPostData *postProtobufV1.PostCreate) (*postProtobufV1.PostView, error) {
	post := &post.PostCreate{
		UserID:  createdPostData.UserID,
		Title:   createdPostData.Title,
		Content: createdPostData.Content,
		Image:   createdPostData.Image,
	}

	createdPost, err := postGrpcServer.postUseCase.CreatePost(ctx, post)

	if err != nil {
		if strings.Contains(err.Error(), "sorry, but this title already exists. Please choose another one") {
			return nil, status.Errorf(codes.AlreadyExists, err.Error(), "error")
		}

		return nil, status.Errorf(codes.Internal, err.Error(), "error")
	}

	postView := &postProtobufV1.PostView{
		Post: &postProtobufV1.Post{
			PostID:    createdPost.PostID,
			UserID:    createdPost.UserID,
			Title:     createdPost.Title,
			Content:   createdPost.Content,
			CreatedAt: timestamppb.New(createdPost.CreatedAt),
			UpdatedAt: timestamppb.New(createdPost.UpdatedAt),
		},
	}

	return postView, nil
}
