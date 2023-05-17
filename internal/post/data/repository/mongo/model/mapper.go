package model

import (
	postModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
)

func PostCreateToPostCreateRepositoryMapper(post *postModel.PostCreate) PostCreateRepository {
	return PostCreateRepository{
		Title:     post.Title,
		Content:   post.Content,
		Image:     post.Image,
		UserID:    post.UserID,
		User:      post.User,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

func PostUpdateToPostUpdateRepositoryMapper(post *postModel.PostUpdate) PostUpdateRepository {
	return PostUpdateRepository{
		PostID:    post.PostID,
		Title:     post.Title,
		Content:   post.Content,
		Image:     post.Image,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}
