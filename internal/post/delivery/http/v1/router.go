package v1

import "github.com/gin-gonic/gin"

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
	router.POST("/", postRouter.postHandler.CreatePost)
	router.PUT("/:postID", postRouter.postHandler.UpdatePost)
	router.DELETE("/:postID", postRouter.postHandler.DeletePostByID)
}
