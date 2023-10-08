package user

type UserController interface {
	GetAllUsers(controllerContext interface{})
	GetCurrentUser(controllerContext interface{})
	GetUserById(controllerContext interface{})
	Register(controllerContext interface{})
	UpdateUserById(controllerContext interface{})
	Delete(controllerContext interface{})
	Login(controllerContext interface{})
	RefreshAccessToken(controllerContext interface{})
	ForgottenPassword(controllerContext interface{})
	ResetUserPassword(controllerContext interface{})
	Logout(controllerContext interface{})
}
