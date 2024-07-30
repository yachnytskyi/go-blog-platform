package gin

import (
	"github.com/gin-gonic/gin"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	middleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
)

type PostRouter struct {
	Config         interfaces.Config
	Logger         interfaces.Logger
	PostController interfaces.PostController
}

func NewPostRouter(config interfaces.Config, logger interfaces.Logger, postController interfaces.PostController) interfaces.PostRouter {
	return PostRouter{
		Config:         config,
		Logger:         logger,
		PostController: postController,
	}
}

func (postRouter PostRouter) PostRouter(routerGroup any) {
	ginRouterGroup := routerGroup.(*gin.RouterGroup)
	router := ginRouterGroup.Group(constants.PostsGroupPath)
	router.GET("/", func(ginContext *gin.Context) {
		postRouter.PostController.GetAllPosts(ginContext)
	})
	router.GET("/:postID", func(ginContext *gin.Context) {
		postRouter.PostController.GetPostById(ginContext)
	})

	router.Use(middleware.AuthenticationMiddleware(postRouter.Config, postRouter.Logger))
	router.POST("/", func(ginContext *gin.Context) {
		postRouter.PostController.CreatePost(ginContext)
	})
	router.PUT("/:postID", func(ginContext *gin.Context) {
		postRouter.PostController.UpdatePostById(ginContext)
	})
	router.DELETE("/:postID", func(ginContext *gin.Context) {
		postRouter.PostController.DeletePostByID(ginContext)
	})
}
