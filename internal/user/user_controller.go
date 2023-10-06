package user

type UserController interface {
	GetAllUsers(context interface{})
	GetCurrentUser(context interface{})
	GetUserById(context interface{})
	Register(context interface{})
	UpdateUserById(context interface{})
	Delete(context interface{})
	Login(context interface{})
	RefreshAccessToken(context interface{})
	ForgottenPassword(context interface{})
	ResetUserPassword(context interface{})
	Logout(context interface{})
}
