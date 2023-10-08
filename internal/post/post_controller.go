package post

type PostController interface {
	GetAllPosts(controllerContext interface{})
	GetPostById(controllerContext interface{})
	CreatePost(controllerContext interface{})
	UpdatePostById(controllerContext interface{})
	DeletePostByID(controllerContext interface{})
}
