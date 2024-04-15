package repository

import (
	"context"
	"errors"
	"time"

	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postRepositoryModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/data/repository/mongo/model"
	postModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"

	mongoModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	location = "post.data.repository.mongo."
)

type PostRepository struct {
	collection *mongo.Collection
}

func NewPostRepository(db *mongo.Database) post.PostRepository {
	return &PostRepository{collection: db.Collection("posts")}
}

func (postRepository *PostRepository) GetAllPosts(ctx context.Context, page int, limit int) (*postModel.Posts, error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	skip := (page - 1) * limit

	option := options.FindOptions{}
	option.SetLimit(int64(limit))
	option.SetSkip(int64(skip))

	query := bson.M{}
	cursor, err := postRepository.collection.Find(ctx, query, &option)

	if validator.IsError(err) {
		return nil, err
	}

	defer cursor.Close(ctx)

	var fetchedPosts []*postRepositoryModel.PostRepository

	for cursor.Next(ctx) {
		post := &postRepositoryModel.PostRepository{}
		err := cursor.Decode(post)

		if validator.IsError(err) {
			return nil, err
		}

		fetchedPosts = append(fetchedPosts, post)
	}

	err = cursor.Err()
	if validator.IsError(err) {
		return nil, err
	}

	if len(fetchedPosts) == 0 {
		return &postModel.Posts{
			Posts: make([]*postModel.Post, 0),
		}, nil
	}

	return &postModel.Posts{
		Posts: postRepositoryModel.PostsRepositoryToPostsMapper(fetchedPosts),
	}, nil
}

func (postRepository *PostRepository) GetPostById(ctx context.Context, postID string) (*postModel.Post, error) {
	postObjectID, hexToObjectIDMapperError := mongoModel.HexToObjectIDMapper(location+"GetPostById", postID)
	if validator.IsError(hexToObjectIDMapperError) {
		return nil, hexToObjectIDMapperError
	}

	query := bson.M{"_id": postObjectID}
	var fetchedPost *postRepositoryModel.PostRepository
	err := postRepository.collection.FindOne(ctx, query).Decode(&fetchedPost)
	if validator.IsError(err) {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return postRepositoryModel.PostRepositoryToPostMapper(fetchedPost), nil
}

func (postRepository *PostRepository) CreatePost(ctx context.Context, post *postModel.PostCreate) (*postModel.Post, error) {
	postMappedToRepository, postCreateToPostCreateRepositoryMapperError := postRepositoryModel.PostCreateToPostCreateRepositoryMapper(post)
	if validator.IsError(postCreateToPostCreateRepositoryMapperError) {
		return nil, postCreateToPostCreateRepositoryMapperError
	}
	postMappedToRepository.CreatedAt = time.Now()
	postMappedToRepository.UpdatedAt = post.CreatedAt

	result, err := postRepository.collection.InsertOne(ctx, postMappedToRepository)
	if validator.IsError(err) {
		er, ok := err.(mongo.WriteException)
		if ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("post with that title already exists")
		}
		return nil, err
	}

	option := options.Index()
	option.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"title": 1}, Options: option}
	_, err = postRepository.collection.Indexes().CreateOne(ctx, index)
	if validator.IsError(err) {
		return nil, errors.New("could not create an index for a title")
	}

	var createdPost *postRepositoryModel.PostRepository
	query := bson.M{"_id": result.InsertedID}
	err = postRepository.collection.FindOne(ctx, query).Decode(&createdPost)
	if validator.IsError(err) {
		return nil, err
	}

	return postRepositoryModel.PostRepositoryToPostMapper(createdPost), nil
}

func (postRepository *PostRepository) UpdatePostById(ctx context.Context, postID string, post *postModel.PostUpdate) (*postModel.Post, error) {
	postUpdateRepository, postUpdateToPostUpdateRepositoryMapper := postRepositoryModel.PostUpdateToPostUpdateRepositoryMapper(post)
	if validator.IsError(postUpdateToPostUpdateRepositoryMapper) {
		return nil, postUpdateToPostUpdateRepositoryMapper
	}

	postUpdateRepository.UpdatedAt = time.Now()
	postMappedToMongoDB, dataToMongoDocumentMapper := mongoModel.DataToMongoDocumentMapper(location+"UpdatePostById", postUpdateRepository)
	if validator.IsError(dataToMongoDocumentMapper) {
		return nil, dataToMongoDocumentMapper
	}

	postObjectID, hexToObjectIDMapperError := mongoModel.HexToObjectIDMapper(location+"UpdatePostById", postID)
	if validator.IsError(hexToObjectIDMapperError) {
		return nil, hexToObjectIDMapperError
	}

	query := bson.D{{Key: "_id", Value: postObjectID}}
	update := bson.D{{Key: "$set", Value: postMappedToMongoDB}}
	result := postRepository.collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedPost *postModel.Post
	err := result.Decode(&updatedPost)
	if validator.IsError(err) {
		return nil, errors.New("sorry, but this title already exists. Please choose another one")
	}

	return updatedPost, nil
}

func (postRepository *PostRepository) DeletePostByID(ctx context.Context, postID string) error {
	postObjectID, hexToObjectIDMapperError := mongoModel.HexToObjectIDMapper(location+"GetPostById", postID)
	if validator.IsError(hexToObjectIDMapperError) {
		return hexToObjectIDMapperError
	}

	query := bson.M{"_id": postObjectID}
	result, err := postRepository.collection.DeleteOne(ctx, query)
	if validator.IsError(err) {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no document with that Id exists")
	}

	return nil
}
