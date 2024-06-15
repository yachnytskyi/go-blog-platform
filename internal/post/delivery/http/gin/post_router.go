package gin

import (
	"github.com/gin-gonic/gin"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"

	// user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
)

type PostRouter struct {
	postController post.PostController
	// userUseCase    user.UserUseCase
}

func NewPostRouter(postController post.PostController) PostRouter {
	return PostRouter{
		postController: postController,
		// userUseCase:    userUseCase,
	}
}

func (postRouter PostRouter) PostRouter(routerGroup any, userUseCase user.UserUseCase) {
	ginRouterGroup := routerGroup.(*gin.RouterGroup)
	router := ginRouterGroup.Group("/posts")
	router.GET("/", func(ginContext *gin.Context) {
		postRouter.postController.GetAllPosts(ginContext)
	})
	router.GET("/:postID", func(ginContext *gin.Context) {
		postRouter.postController.GetPostById(ginContext)
	})

	router.Use(httpGinMiddleware.AuthenticationMiddleware())
	router.POST("/", httpGinMiddleware.UserContextMiddleware(userUseCase), func(ginContext *gin.Context) {
		postRouter.postController.CreatePost(ginContext)
	})
	router.PUT("/:postID", func(ginContext *gin.Context) {
		postRouter.postController.UpdatePostById(ginContext)
	})
	router.DELETE("/:postID", func(ginContext *gin.Context) {
		postRouter.postController.DeletePostByID(ginContext)
	})
}
