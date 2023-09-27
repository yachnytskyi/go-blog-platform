package gin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/model"
	postModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
	ginUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
)

type PostHandler struct {
	postUseCase post.PostUseCase
}

func NewPostHandler(postUseCase post.PostUseCase) PostHandler {
	return PostHandler{postUseCase: postUseCase}
}

func (postHandler *PostHandler) GetAllPosts(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intLimit, err := strconv.Atoi(limit)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	fetchedPosts, err := postHandler.postUseCase.GetAllPosts(ctx.Request.Context(), intPage, intLimit)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, postViewModel.PostsToPostsViewMapper(fetchedPosts))
}

func (postHandler *PostHandler) GetPostById(ctx *gin.Context) {
	postID := ctx.Param("postID")

	fetchedPost, err := postHandler.postUseCase.GetPostById(ctx.Request.Context(), postID)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": postViewModel.PostToPostViewMapper(fetchedPost)})
}

func (postHandler *PostHandler) CreatePost(ctx *gin.Context) {
	var createdPostData *postModel.PostCreate = new(postModel.PostCreate)
	currentUser := ctx.MustGet("currentUser").(*userModel.User)
	createdPostData.User = currentUser.Name
	createdPostData.UserID = currentUser.UserID

	err := ctx.ShouldBindJSON(&createdPostData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}

	createdPost, err := postHandler.postUseCase.CreatePost(ctx.Request.Context(), createdPostData)

	if err != nil {
		if strings.Contains(err.Error(), "sorry, but this title already exists. Please choose another one") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": createdPost})
}

func (postHandler *PostHandler) UpdatePostById(ctx *gin.Context) {
	postID := ctx.Param("postID")
	currentUserID := ginUtility.GetCurrentUserID(ctx)

	var updatedPostData *postModel.PostUpdate = new(postModel.PostUpdate)
	updatedPostData.PostID = ctx.Param("postID")
	updatedPostData.UserID = ginUtility.GetCurrentUserID(ctx)
	err := ctx.ShouldBindJSON(&updatedPostData)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedPost, err := postHandler.postUseCase.UpdatePostById(ctx.Request.Context(), postID, updatedPostData, currentUserID)
	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "sorry, but you do not have permissions to do that") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedPost})
}

func (postHandler *PostHandler) DeletePostByID(ctx *gin.Context) {
	postID := ctx.Param("postID")
	currentUserID := ginUtility.GetCurrentUserID(ctx)

	err := postHandler.postUseCase.DeletePostByID(ctx.Request.Context(), postID, currentUserID)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "sorry, but you do not have permissions to do that") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
