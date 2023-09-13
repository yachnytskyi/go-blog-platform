package model

import (
	postModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	createRepositoryMapperObjectIDFromHex = "Post.Data.Repository.MongoDB.PostCreateToPostCreateRepositoryMapper.ObjectIDFromHex"
	updateRepositoryMapperObjectIDFromHex = "Post.Data.Repository.MongoDB.PostUpdateToPostUpdateRepositoryMapper.ObjectIDFromHex"
)

func PostRepositoryToPostMapper(postRepository *PostRepository) *postModel.Post {
	return &postModel.Post{
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

func PostCreateToPostCreateRepositoryMapper(postCreate *postModel.PostCreate) (*PostCreateRepository, error) {
	userObjectID, objectIDFromHexError := primitive.ObjectIDFromHex(postCreate.UserID)
	if validator.IsErrorNotNil(objectIDFromHexError) {
		objectIDFromHexErrorInternalError := domainError.NewInternalError(createRepositoryMapperObjectIDFromHex, objectIDFromHexError.Error())
		logging.Logger(objectIDFromHexErrorInternalError)
		return nil, objectIDFromHexErrorInternalError
	}
	return &PostCreateRepository{
		UserID:    userObjectID,
		Title:     postCreate.Title,
		Content:   postCreate.Content,
		Image:     postCreate.Image,
		User:      postCreate.User,
		CreatedAt: postCreate.CreatedAt,
		UpdatedAt: postCreate.UpdatedAt,
	}, nil
}

func PostUpdateToPostUpdateRepositoryMapper(postUpdate *postModel.PostUpdate) (*PostUpdateRepository, error) {
	postObjectID, objectIDFromHexError := primitive.ObjectIDFromHex(postUpdate.PostID)
	if validator.IsErrorNotNil(objectIDFromHexError) {
		objectIDFromHexErrorInternalError := domainError.NewInternalError(updateRepositoryMapperObjectIDFromHex, objectIDFromHexError.Error())
		logging.Logger(objectIDFromHexErrorInternalError)
		return nil, objectIDFromHexErrorInternalError
	}

	userObjectID, objectIDFromHexError := primitive.ObjectIDFromHex(postUpdate.UserID)
	if validator.IsErrorNotNil(objectIDFromHexError) {
		objectIDFromHexErrorInternalError := domainError.NewInternalError(createRepositoryMapperObjectIDFromHex, objectIDFromHexError.Error())
		logging.Logger(objectIDFromHexErrorInternalError)
		return nil, objectIDFromHexErrorInternalError
	}
	return &PostUpdateRepository{
		PostID:    postObjectID,
		UserID:    userObjectID,
		Title:     postUpdate.Title,
		Content:   postUpdate.Content,
		Image:     postUpdate.Image,
		User:      postUpdate.User,
		CreatedAt: postUpdate.CreatedAt,
		UpdatedAt: postUpdate.UpdatedAt,
	}, nil
}
