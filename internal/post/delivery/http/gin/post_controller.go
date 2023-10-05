package gin

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/model"
	postModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
)

type PostController struct {
	postUseCase post.PostUseCase
}

func NewPostController(postUseCase post.PostUseCase) PostController {
	return PostController{postUseCase: postUseCase}
}

func (postController PostController) GetAllPosts(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constant.DefaultContextTimer)
	defer cancel()
	page := ginContext.DefaultQuery("page", "1")
	limit := ginContext.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)

	if err != nil {
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intLimit, err := strconv.Atoi(limit)

	if err != nil {
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	fetchedPosts, err := postController.postUseCase.GetAllPosts(ctx, intPage, intLimit)
	if err != nil {
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, postViewModel.PostsToPostsViewMapper(fetchedPosts))
}

func (postController PostController) GetPostById(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constant.DefaultContextTimer)
	defer cancel()
	postID := ginContext.Param("postID")

	fetchedPost, err := postController.postUseCase.GetPostById(ctx, postID)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ginContext.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"status": "success", "data": postViewModel.PostToPostViewMapper(fetchedPost)})
}

func (postController PostController) CreatePost(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constant.DefaultContextTimer)
	defer cancel()
	var createdPostData *postModel.PostCreate = new(postModel.PostCreate)
	currentUser := ginContext.MustGet("user").(userViewModel.UserView)
	createdPostData.User = currentUser.Name
	createdPostData.UserID = currentUser.UserID

	err := ginContext.ShouldBindJSON(&createdPostData)
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, err.Error())
	}

	createdPost, err := postController.postUseCase.CreatePost(ctx, createdPostData)

	if err != nil {
		if strings.Contains(err.Error(), "sorry, but this title already exists. Please choose another one") {
			ginContext.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ginContext.JSON(http.StatusCreated, gin.H{"status": "success", "data": createdPost})
}

func (postController PostController) UpdatePostById(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constant.DefaultContextTimer)
	defer cancel()
	postID := ginContext.Param("postID")
	currentUserID := ginContext.MustGet("userID").(string)

	var updatedPostData *postModel.PostUpdate = new(postModel.PostUpdate)
	updatedPostData.PostID = ginContext.Param("postID")
	updatedPostData.UserID = currentUserID
	err := ginContext.ShouldBindJSON(&updatedPostData)
	if err != nil {
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedPost, err := postController.postUseCase.UpdatePostById(ctx, postID, updatedPostData, currentUserID)
	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ginContext.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "sorry, but you do not have permissions to do that") {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedPost})
}

func (postController PostController) DeletePostByID(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constant.DefaultContextTimer)
	defer cancel()
	postID := ginContext.Param("postID")
	currentUserID := ginContext.MustGet("userID").(string)
	err := postController.postUseCase.DeletePostByID(ctx, postID, currentUserID)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ginContext.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "sorry, but you do not have permissions to do that") {
			ginContext.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ginContext.JSON(http.StatusNoContent, nil)
}
