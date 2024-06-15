package post

import "github.com/yachnytskyi/golang-mongo-grpc/internal/user"

type PostRouter interface {
	PostRouter(routerGroup any, userUseCase user.UserUseCase)
}
