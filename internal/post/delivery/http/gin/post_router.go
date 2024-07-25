package gin

import (
	"github.com/gin-gonic/gin"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"

	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
)

type PostRouter struct {
	Config         model.Config
	Logger         model.Logger
	PostController post.PostController
}

func NewPostRouter(config model.Config, logger model.Logger, postController post.PostController) PostRouter {
	return PostRouter{
		Config:         config,
		Logger:         logger,
		PostController: postController,
	}
}

func (postRouter PostRouter) PostRouter(routerGroup any, userUseCase user.UserUseCase) {
	ginRouterGroup := routerGroup.(*gin.RouterGroup)
	router := ginRouterGroup.Group("/posts")
	router.GET("/", func(ginContext *gin.Context) {
		postRouter.PostController.GetAllPosts(ginContext)
	})
	router.GET("/:postID", func(ginContext *gin.Context) {
		postRouter.PostController.GetPostById(ginContext)
	})

	router.Use(httpGinMiddleware.AuthenticationMiddleware(postRouter.Config, postRouter.Logger))
	router.POST("/", httpGinMiddleware.UserContextMiddleware(postRouter.Logger, userUseCase), func(ginContext *gin.Context) {
		postRouter.PostController.CreatePost(ginContext)
	})
	router.PUT("/:postID", func(ginContext *gin.Context) {
		postRouter.PostController.UpdatePostById(ginContext)
	})
	router.DELETE("/:postID", func(ginContext *gin.Context) {
		postRouter.PostController.DeletePostByID(ginContext)
	})
}
