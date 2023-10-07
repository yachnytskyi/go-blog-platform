package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postController "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/http/gin"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userController "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin"
)

type GinFactory struct {
	Gin    config.Gin
	Server *gin.Engine
}

func (ginFactory GinFactory) CloseServer() {
}

func (ginFactory GinFactory) NewUserController(domain interface{}) user.UserController {
	userUseCase := domain.(user.UserUseCase)
	return userController.NewUserController(userUseCase)
}

func (ginFactory GinFactory) NewPostController(domain interface{}) post.PostController {
	postUseCase := domain.(post.PostUseCase)
	return postController.NewPostController(postUseCase)
}
