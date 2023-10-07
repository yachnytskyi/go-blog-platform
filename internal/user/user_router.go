package user

type UserRouter interface {
	UserRouter(routerGroup interface{}, userUseCase UserUseCase)
}
