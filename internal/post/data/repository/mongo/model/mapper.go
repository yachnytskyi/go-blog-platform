package model

import (
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "post.data.depository.mongo."
)

func PostsRepositoryToPostsMapper(postsRepository []*PostRepository) []*post.Post {
	posts := make([]*post.Post, len(postsRepository))
	for index, postRepository := range postsRepository {
		posts[index] = PostRepositoryToPostMapper(postRepository)
	}

	return posts
}

func PostRepositoryToPostMapper(postRepository *PostRepository) *post.Post {
	return &post.Post{
		PostID:    postRepository.PostID.Hex(),
		UserID:    postRepository.UserID.Hex(),
		Title:     postRepository.Title,
		Content:   postRepository.Content,
		Image:     postRepository.Image,
		User:      postRepository.User,
		CreatedAt: postRepository.CreatedAt,
		UpdatedAt: postRepository.UpdatedAt,
	}
}

func PostCreateToPostCreateRepositoryMapper(logger interfaces.Logger, postCreate *post.PostCreate) (*PostCreateRepository, error) {
	userObjectID := model.HexToObjectIDMapper(logger, location+"PostCreateToPostCreateRepositoryMapper", postCreate.UserID)
	if validator.IsError(userObjectID.Error) {
		return &PostCreateRepository{}, userObjectID.Error
	}

	return &PostCreateRepository{
		UserID:    userObjectID.Data,
		Title:     postCreate.Title,
		Content:   postCreate.Content,
		Image:     postCreate.Image,
		User:      postCreate.User,
		CreatedAt: postCreate.CreatedAt,
		UpdatedAt: postCreate.UpdatedAt,
	}, nil
}

func PostUpdateToPostUpdateRepositoryMapper(logger interfaces.Logger, postUpdate *post.PostUpdate) (*PostUpdateRepository, error) {
	postObjectID := model.HexToObjectIDMapper(logger, location+"PostUpdateToPostUpdateRepositoryMapper.postObjectID", postUpdate.PostID)
	if validator.IsError(postObjectID.Error) {
		return &PostUpdateRepository{}, postObjectID.Error
	}

	userObjectID := model.HexToObjectIDMapper(logger, location+"PostUpdateToPostUpdateRepositoryMapper.userObjectID", postUpdate.UserID)
	if validator.IsError(userObjectID.Error) {
		return &PostUpdateRepository{}, userObjectID.Error
	}

	return &PostUpdateRepository{
		PostID:    postObjectID.Data,
		UserID:    userObjectID.Data,
		Title:     postUpdate.Title,
		Content:   postUpdate.Content,
		Image:     postUpdate.Image,
		User:      postUpdate.User,
		CreatedAt: postUpdate.CreatedAt,
		UpdatedAt: postUpdate.UpdatedAt,
	}, nil
}
