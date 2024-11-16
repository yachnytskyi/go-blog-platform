package model

import (
	"github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
)

func PostsToPostsViewMapper(posts *model.Posts) PostsView {
	postsView := make([]*PostView, 0, 10)
	for _, post := range posts.Posts {
		postView := &PostView{}
		postView.PostID = post.PostID
		postView.Title = post.Title
		postView.Content = post.Content
		postView.Image = post.Image
		postView.CreatedAt = post.CreatedAt
		postView.UpdatedAt = post.UpdatedAt
		postsView = append(postsView, postView)
	}

	return PostsView{
		PostsView: postsView,
	}
}

func PostToPostViewMapper(post *model.Post) PostView {
	return PostView{
		PostID:    post.PostID,
		UserID:    post.UserID,
		Title:     post.Title,
		Content:   post.Content,
		Image:     post.Image,
		Username:  post.Username,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}
