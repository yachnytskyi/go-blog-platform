package v1

import (
	"context"

	postProtobufV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/grpc/v1/model/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (PostGrpcServer *PostGrpcServer) GetPosts(postsData *postProtobufV1.Posts, streamOfPosts postProtobufV1.PostUseCase_GetPostsServer) error {
	ctx := context.Background()
	page := postsData.GetPage()
	limit := postsData.GetLimit()

	posts, err := PostGrpcServer.postUseCase.GetAllPosts(ctx, int(page), int(limit))

	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	for _, post := range posts {
		streamOfPosts.Send(&postProtobufV1.Post{
			PostID:    post.PostID,
			Title:     post.Title,
			Content:   post.Content,
			Image:     post.Image,
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		})
	}

	return nil
}
