package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpV1Utility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/v1/utility"
)

type PostRouter struct {
	postHandler PostHandler
}

func NewPostRouter(postHandler PostHandler) PostRouter {
	return PostRouter{postHandler: postHandler}
}

func (postRouter *PostRouter) PostRouter(routerGroup *gin.RouterGroup, userUseCase user.UseCase) {
	router := routerGroup.Group("/posts")

	router.GET("/", postRouter.postHandler.GetAllPosts)
	router.GET("/:postID", postRouter.postHandler.GetPostById)

	router.Use(httpV1Utility.DeserializeUser(userUseCase))

	router.POST("/", postRouter.postHandler.CreatePost)
	router.PUT("/:postID", postRouter.postHandler.UpdatePostById)
	router.DELETE("/:postID", postRouter.postHandler.DeletePostByID)
}
