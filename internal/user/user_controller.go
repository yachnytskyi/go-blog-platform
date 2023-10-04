package user

type UserDelivery interface {
	GetAllUsers(context interface{})
	GetMe(context interface{})
	GetUserById(context interface{})
	GetUserByEmail(context interface{})
	Register(context interface{})
	UpdateUserById(context interface{})
	Delete(context interface{})
	Login(context interface{})
	RefreshAccessToken(context interface{})
	ForgottenPassword(context interface{})
	ResetUserPassword(context interface{})
	Logout(context interface{})
}
