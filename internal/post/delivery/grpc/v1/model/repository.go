package model

import postModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"

type Repository interface {
	GetPostById(string) (*postModel.Post, error)
	GetAllPosts(page int, limit int) ([]*postModel.Post, error)
	CreatePost(*postModel.PostCreate) (*postModel.Post, error)
	UpdatePostById(string, *postModel.PostUpdate) (*postModel.Post, error)
	DeletePostById(string) error
}
