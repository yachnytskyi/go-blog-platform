package post

import user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"

type PostRouter interface {
	PostRouter(routerGroup any, userUseCase user.UserUseCase)
}
