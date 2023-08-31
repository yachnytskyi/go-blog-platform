package repository

import (
	"context"
	"errors"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postRepositoryModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/data/repository/mongo/model"
	postModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"

	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepository struct {
	collection *mongo.Collection
}

func NewPostRepository(collection *mongo.Collection) post.Repository {
	return &PostRepository{collection: collection}
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

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var fetchedPosts []*postModel.Post

	for cursor.Next(ctx) {
		post := &postModel.Post{}
		err := cursor.Decode(post)

		if err != nil {
			return nil, err
		}

		fetchedPosts = append(fetchedPosts, post)
	}

	err = cursor.Err()
	if err != nil {
		return nil, err
	}

	if len(fetchedPosts) == 0 {
		return &postModel.Posts{
			Posts: make([]*postModel.Post, 0),
		}, nil
	}

	return &postModel.Posts{
		Posts: fetchedPosts,
	}, nil
}

func (postRepository *PostRepository) GetPostById(ctx context.Context, postID string) (*postModel.Post, error) {
	postIDMappedToMongoDB, _ := primitive.ObjectIDFromHex(postID)

	query := bson.M{"_id": postIDMappedToMongoDB}

	var fetchedPost *postModel.Post

	err := postRepository.collection.FindOne(ctx, query).Decode(&fetchedPost)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return fetchedPost, nil
}

func (postRepository *PostRepository) CreatePost(ctx context.Context, post *postModel.PostCreate) (*postModel.Post, error) {
	postMappedToRepository := postRepositoryModel.PostCreateToPostCreateRepositoryMapper(post)
	postMappedToRepository.CreatedAt = time.Now()
	postMappedToRepository.UpdatedAt = post.CreatedAt

	result, err := postRepository.collection.InsertOne(ctx, postMappedToRepository)

	if err != nil {
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
	if err != nil {
		return nil, errors.New("could not create an index for a title")
	}

	var createdPost *postModel.Post
	query := bson.M{"_id": result.InsertedID}

	err = postRepository.collection.FindOne(ctx, query).Decode(&createdPost)
	if err != nil {
		return nil, err
	}

	return createdPost, nil
}

func (postRepository *PostRepository) UpdatePostById(ctx context.Context, postID string, post *postModel.PostUpdate) (*postModel.Post, error) {
	postMappedToRepository := postRepositoryModel.PostUpdateToPostUpdateRepositoryMapper(post)
	postMappedToRepository.UpdatedAt = time.Now()
	postMappedToMongoDB, err := utility.MongoMappper(postMappedToRepository)

	if err != nil {
		return nil, err
	}

	postIDMappedToMongoDB, _ := primitive.ObjectIDFromHex(postID)

	query := bson.D{{Key: "_id", Value: postIDMappedToMongoDB}}
	update := bson.D{{Key: "$set", Value: postMappedToMongoDB}}
	result := postRepository.collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedPost *postModel.Post

	err = result.Decode(&updatedPost)
	if err != nil {
		return nil, errors.New("sorry, but this title already exists. Please choose another one")
	}

	return updatedPost, nil
}

func (postRepository *PostRepository) DeletePostByID(ctx context.Context, postID string) error {
	postIDMappedToMongoDB, _ := primitive.ObjectIDFromHex(postID)

	query := bson.M{"_id": postIDMappedToMongoDB}
	result, err := postRepository.collection.DeleteOne(ctx, query)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no document with that Id exists")
	}

	return nil
}
