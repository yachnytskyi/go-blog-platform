package post

import "github.com/yachnytskyi/golang-mongo-grpc/internal/user"

type PostRouter interface {
	PostRouter(routerGroup interface{}, userUseCase user.UserUseCase)
}
