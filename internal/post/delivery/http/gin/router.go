package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility"
)

type PostRouter struct {
	postHandler PostHandler
}

func NewPostRouter(postHandler PostHandler) PostRouter {
	return PostRouter{postHandler: postHandler}
}

func (postRouter *PostRouter) PostRouter(routerGroup *gin.RouterGroup, userUseCase user.UserUseCase) {
	router := routerGroup.Group("/posts")

	router.GET("/", postRouter.postHandler.GetAllPosts)
	router.GET("/:postID", postRouter.postHandler.GetPostById)

	router.Use(httpGinUtility.DeserializeUser(userUseCase))

	router.POST("/", postRouter.postHandler.CreatePost)
	router.PUT("/:postID", postRouter.postHandler.UpdatePostById)
	router.DELETE("/:postID", postRouter.postHandler.DeletePostByID)
}
