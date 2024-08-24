package gin

import "github.com/gin-gonic/gin"

func NewMockGinEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}
