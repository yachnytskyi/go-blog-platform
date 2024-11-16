package v1

import (
	"context"
	"strings"
	"time"

	postProtobufV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/grpc/v1/model/pb"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (postGrpcServer *PostGrpcServer) UpdatePostById(ctx context.Context, updatedPostData *postProtobufV1.PostUpdate) (*postProtobufV1.PostView, error) {
	postID := updatedPostData.GetPostID()
	userID := updatedPostData.GetUserID()

	post := &post.PostUpdate{
		Title:     updatedPostData.GetTitle(),
		Content:   updatedPostData.GetContent(),
		Image:     updatedPostData.GetImage(),
		UserID:    updatedPostData.GetUserID(),
		UpdatedAt: time.Now(),
	}

	createdPost, err := postGrpcServer.postUseCase.UpdatePostById(ctx, postID, post, userID)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
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
