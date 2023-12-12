package user

type UserController interface {
	GetAllUsers(controllerContext any)
	GetCurrentUser(controllerContext any)
	GetUserById(controllerContext any)
	Register(controllerContext any)
	UpdateCurrentUser(controllerContext any)
	DeleteCurrentUser(controllerContext any)
	Login(controllerContext any)
	RefreshAccessToken(controllerContext any)
	ForgottenPassword(controllerContext any)
	ResetUserPassword(controllerContext any)
	Logout(controllerContext any)
}
