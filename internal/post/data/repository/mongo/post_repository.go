package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	repository "github.com/yachnytskyi/golang-mongo-grpc/internal/post/data/repository/mongo/model"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo/model"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/data/repository/mongo"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	location = "post.data.repository.mongo."
)

type PostRepository struct {
	Logger interfaces.Logger
	posts  *mongo.Collection
	users  *mongo.Collection
}

func NewPostRepository(logger interfaces.Logger, db *mongo.Database) interfaces.PostRepository {
	return &PostRepository{
		Logger: logger,
		posts:  db.Collection(constants.PostsTable),
		users:  db.Collection(constants.UsersTable),
	}
}

func (postRepository *PostRepository) GetAllPosts(ctx context.Context, page int, limit int) (*post.Posts, error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 100
	}

	skip := (page - 1) * limit

	option := options.FindOptions{}
	option.SetLimit(int64(limit))
	option.SetSkip(int64(skip))

	query := bson.M{}
	cursor, err := postRepository.posts.Find(ctx, query, &option)

	if validator.IsError(err) {
		return nil, err
	}

	defer cursor.Close(ctx)

	var fetchedPosts []*repository.PostRepository

	for cursor.Next(ctx) {
		post := &repository.PostRepository{}
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
		return &post.Posts{
			Posts: make([]*post.Post, 0),
		}, nil
	}

	return &post.Posts{
		Posts: repository.PostsRepositoryToPostsMapper(fetchedPosts),
	}, nil
}

func (postRepository *PostRepository) GetPostById(ctx context.Context, postID string) (*post.Post, error) {
	postObjectID := model.HexToObjectIDMapper(postRepository.Logger, location+"GetPostById", postID)
	if validator.IsError(postObjectID.Error) {
		return nil, postObjectID.Error
	}

	var fetchedPost *repository.PostRepository
	query := bson.M{model.ID: postObjectID.Data}
	err := postRepository.posts.FindOne(ctx, query).Decode(&fetchedPost)
	if validator.IsError(err) {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return repository.PostRepositoryToPostMapper(fetchedPost), nil
}

func (postRepository *PostRepository) CreatePost(ctx context.Context, post *post.PostCreate) (*post.Post, error) {
	// Map PostCreate to PostCreateRepository
	postMappedToRepository, postCreateToPostCreateRepositoryMapperError := repository.PostCreateToPostCreateRepositoryMapper(postRepository.Logger, post)
	if validator.IsError(postCreateToPostCreateRepositoryMapperError) {
		return nil, postCreateToPostCreateRepositoryMapperError
	}

	// Fetch the username from the users collection.
	var user user.UserRepository
	query := bson.M{model.ID: postMappedToRepository.UserID}
	err := postRepository.users.FindOne(ctx, query).Decode(&user)
	if validator.IsError(err) {
		fmt.Println(err)
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("error fetching user details")
	}

	postMappedToRepository.Username = user.Username
	postMappedToRepository.CreatedAt = time.Now()
	postMappedToRepository.UpdatedAt = time.Now()

	// Insert the post into the collection
	result, err := postRepository.posts.InsertOne(ctx, postMappedToRepository)
	if validator.IsError(err) {
		er, ok := err.(mongo.WriteException)
		if ok && len(er.WriteErrors) > 0 && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("post with that title already exists")
		}
		return nil, err
	}

	// Create an index for the title field
	option := options.Index()
	option.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"title": 1}, Options: option}
	_, err = postRepository.posts.Indexes().CreateOne(ctx, index)
	if validator.IsError(err) {
		return nil, errors.New("could not create an index for the title")
	}

	// Retrieve the created post
	var createdPost *repository.PostRepository
	query = bson.M{model.ID: result.InsertedID}
	err = postRepository.posts.FindOne(ctx, query).Decode(&createdPost)
	if validator.IsError(err) {
		return nil, err
	}

	return repository.PostRepositoryToPostMapper(createdPost), nil
}

func (postRepository *PostRepository) UpdatePostById(ctx context.Context, postID string, postUpdate *post.PostUpdate) (*post.Post, error) {
	postUpdateRepository, postUpdateToPostUpdateRepositoryMapper := repository.PostUpdateToPostUpdateRepositoryMapper(postRepository.Logger, postUpdate)
	if validator.IsError(postUpdateToPostUpdateRepositoryMapper) {
		return nil, postUpdateToPostUpdateRepositoryMapper
	}

	// Fetch the username from the users collection.
	var user user.UserRepository
	userQuery := bson.M{model.ID: postUpdateRepository.UserID}
	err := postRepository.users.FindOne(ctx, userQuery).Decode(&user)
	if validator.IsError(err) {
		fmt.Println(err)
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("error fetching user details")
	}

	postUpdateRepository.Username = user.Username
	postUpdateRepository.UpdatedAt = time.Now()

	// Map the user update repository to a BSON document for MongoDB update.
	postUpdateBson := model.DataToMongoDocumentMapper(postRepository.Logger, location+"UpdatePostById", postUpdateRepository)
	if validator.IsError(postUpdateBson.Error) {
		return nil, postUpdateBson.Error
	}

	query := bson.D{{Key: model.ID, Value: postUpdateRepository.PostID}}
	update := bson.D{{Key: model.Set, Value: postUpdateBson.Data}}
	result := postRepository.posts.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))
	updatedPost := &repository.PostRepository{}
	decodeError := result.Decode(&updatedPost)
	if validator.IsError(err) {
		fmt.Println(decodeError)
		return nil, decodeError
	}

	return repository.PostRepositoryToPostMapper(updatedPost), nil
}

func (postRepository *PostRepository) DeletePostByID(ctx context.Context, postID string) error {
	postObjectID := model.HexToObjectIDMapper(postRepository.Logger, location+"GetPostById", postID)
	if validator.IsError(postObjectID.Error) {
		return postObjectID.Error
	}

	query := bson.M{model.ID: postObjectID.Data}
	result, err := postRepository.posts.DeleteOne(ctx, query)
	if validator.IsError(err) {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no document with that Id exists")
	}

	return nil
}
