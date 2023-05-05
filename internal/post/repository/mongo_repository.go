package repository

import (
	"context"
	"errors"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	"github.com/yachnytskyi/golang-mongo-grpc/models"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utils"
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

func (postRepository *PostRepository) GetPostById(ctx context.Context, postID string) (*models.Post, error) {
	postObjectID, _ := primitive.ObjectIDFromHex(postID)

	query := bson.M{"_id": postObjectID}

	var fetchedPost *models.Post

	if err := postRepository.collection.FindOne(ctx, query).Decode(&fetchedPost); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return fetchedPost, nil
}

func (postRepository *PostRepository) GetAllPosts(ctx context.Context, page int, limit int) ([]*models.Post, error) {
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

	var fetchedPosts []*models.Post

	for cursor.Next(ctx) {
		post := &models.Post{}
		err := cursor.Decode(post)

		if err != nil {
			return nil, err
		}

		fetchedPosts = append(fetchedPosts, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(fetchedPosts) == 0 {
		return []*models.Post{}, nil
	}

	return fetchedPosts, nil
}

func (postRepository *PostRepository) CreatePost(ctx context.Context, post *models.PostCreate) (*models.Post, error) {
	post.CreatedAt = time.Now()
	post.UpdatedAt = post.CreatedAt

	postMappedToRepository := models.PostCreateToPostCreateRepositoryMapper(post)
	result, err := postRepository.collection.InsertOne(ctx, postMappedToRepository)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("post with that title already exists")
		}
		return nil, err
	}

	option := options.Index()
	option.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"title": 1}, Options: option}

	if _, err := postRepository.collection.Indexes().CreateOne(ctx, index); err != nil {
		return nil, errors.New("could not create an index for a title")
	}

	var createdPost *models.Post
	query := bson.M{"_id": result.InsertedID}

	if err = postRepository.collection.FindOne(ctx, query).Decode(&createdPost); err != nil {
		return nil, err
	}

	return createdPost, nil
}

func (postRepository *PostRepository) UpdatePostById(ctx context.Context, postID string, post *models.PostUpdate) (*models.Post, error) {
	postMappedToRepository := models.PostUpdateToPostUpdateRepositoryMapper(post)
	postMappedToMongoDB, err := utils.MongoMapping(postMappedToRepository)

	if err != nil {
		return nil, err
	}

	postObjectID, _ := primitive.ObjectIDFromHex(postID)
	query := bson.D{{Key: "_id", Value: postObjectID}}
	update := bson.D{{Key: "$set", Value: postMappedToMongoDB}}
	result := postRepository.collection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedPost *models.Post

	if err := result.Decode(&updatedPost); err != nil {
		return nil, errors.New("sorry, but this title already exists. Please choose another one")
	}

	return updatedPost, nil
}

func (postRepository *PostRepository) DeletePostByID(ctx context.Context, postID string) error {
	postObjectID, _ := primitive.ObjectIDFromHex(postID)
	query := bson.M{"_id": postObjectID}

	result, err := postRepository.collection.DeleteOne(ctx, query)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no document with that Id exists")
	}

	return nil
}
