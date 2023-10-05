package gin

import (
	"github.com/gin-gonic/gin"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
)

type PostRouter struct {
	postController PostController
}

func NewPostRouter(postController PostController) PostRouter {
	return PostRouter{postController: postController}
}

func (postRouter *PostRouter) PostRouter(routerGroup *gin.RouterGroup, userUseCase user.UserUseCase) {
	router := routerGroup.Group("/posts")
	router.GET("/", postRouter.postController.GetAllPosts)
	router.GET("/:postID", postRouter.postController.GetPostById)

	router.Use(httpGinMiddleware.DeserializeUser(userUseCase))
	router.POST("/", postRouter.postController.CreatePost)
	router.PUT("/:postID", postRouter.postController.UpdatePostById)
	router.DELETE("/:postID", postRouter.postController.DeletePostByID)
}
