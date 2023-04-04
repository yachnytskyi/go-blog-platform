package user

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	UserRouter(routerGroup *gin.RouterGroup, userService Service)
}
