package post

type PostController interface {
	GetAllPosts(controllerContext any)
	GetPostById(controllerContext any)
	CreatePost(controllerContext any)
	UpdatePostById(controllerContext any)
	DeletePostByID(controllerContext any)
}
