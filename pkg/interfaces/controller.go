package interfaces

type HealthCheckController interface {
	HealthCheck(controllerContext any)
}

type UserController interface {
	GetAllUsers(controllerContext any)
	GetCurrentUser(controllerContext any)
	GetUserById(controllerContext any)
	Register(controllerContext any)
	UpdateCurrentUser(controllerContext any)
	DeleteCurrentUser(controllerContext any)
	Login(controllerContext any)
	RefreshAccessToken(controllerContext any)
	Logout(controllerContext any)
	ForgottenPassword(controllerContext any)
	ResetUserPassword(controllerContext any)
}

type PostController interface {
	GetAllPosts(controllerContext any)
	GetPostById(controllerContext any)
	CreatePost(controllerContext any)
	UpdatePostById(controllerContext any)
	DeletePostByID(controllerContext any)
}

type Router interface {
	Router(routerGroup any)
}
