package user

type UserRouter interface {
	UserRouter(routerGroup any, userUseCase UserUseCase)
}
