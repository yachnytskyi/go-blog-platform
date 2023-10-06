package gin

import (
	"github.com/gin-gonic/gin"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
)

type PostRouter struct {
	postController post.PostController
}

func NewPostRouter(postController PostController) PostRouter {
	return PostRouter{postController: postController}
}

func (postRouter *PostRouter) PostRouter(routerGroup *gin.RouterGroup, userUseCase user.UserUseCase) {
	router := routerGroup.Group("/posts")
	router.GET("/", func(ginContext *gin.Context) {
		postRouter.postController.GetAllPosts(ginContext)
	})
	router.GET("/:postID", func(ginContext *gin.Context) {
		postRouter.postController.GetPostById(ginContext)
	})

	router.Use(httpGinMiddleware.DeserializeUser(userUseCase))
	router.POST("/", func(ginContext *gin.Context) {
		postRouter.postController.CreatePost(ginContext)
	})
	router.PUT("/:postID", func(ginContext *gin.Context) {
		postRouter.postController.UpdatePostById(ginContext)
	})
	router.DELETE("/:postID", func(ginContext *gin.Context) {
		postRouter.postController.DeletePostByID(ginContext)
	})
}
