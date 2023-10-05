package gin

import (
	"github.com/gin-gonic/gin"
	httpGinUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility/middleware"
)

type PostRouter struct {
	postHandler PostHandler
}

func NewPostRouter(postHandler PostHandler) PostRouter {
	return PostRouter{postHandler: postHandler}
}

func (postRouter *PostRouter) PostRouter(routerGroup *gin.RouterGroup) {
	router := routerGroup.Group("/posts")

	router.GET("/", postRouter.postHandler.GetAllPosts)
	router.GET("/:postID", postRouter.postHandler.GetPostById)

	router.Use(httpGinUtility.DeserializeUser())

	router.POST("/", postRouter.postHandler.CreatePost)
	router.PUT("/:postID", postRouter.postHandler.UpdatePostById)
	router.DELETE("/:postID", postRouter.postHandler.DeletePostByID)
}
