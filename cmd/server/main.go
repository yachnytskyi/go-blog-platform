package main

// Require the packages.
import (
	"context"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"

	postPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postHttpGinPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/http/gin"

	"github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency"

	postUseCasePackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"

	userPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	// dependency "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency"

	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"

	userHttpGinPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin"
	userUseCasePackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
)

// Create required variables that we'll re-assign later.
var (
	server *gin.Engine

	userRepository userPackage.UserRepository
	userUseCase    userPackage.UserUseCase
	userController userHttpGinPackage.UserController
	userRouter     userHttpGinPackage.UserRouter

	postRepository postPackage.PostRepository
	postUseCase    postPackage.PostUseCase
	postController postHttpGinPackage.PostHandler
	postRouter     postHttpGinPackage.PostRouter
)

// Init function that will run before the "main" function.
func init() {
	loadConfig := commonUtility.LoadConfig()

	// Create a context.
	ctx, cancel := context.WithTimeout(context.Background(), constant.DefaultContextTimer)
	defer cancel()

	container := dependency.CreateApplication(ctx)
	fmt.Println(container)

	repository.InjectRepository(loadConfig, container)

	// Create a DB database instance using the factory.
	db := container.RepositoryFactory.NewRepository(ctx)
	userRepository = container.RepositoryFactory.NewUserRepository(db)
	postRepository = container.RepositoryFactory.NewPostRepository(db)

	// Use Cases.
	userUseCase = userUseCasePackage.NewUserUseCase(userRepository)
	postUseCase = postUseCasePackage.NewPostUseCase(postRepository)

	// Handlers
	userController = userHttpGinPackage.NewUserController(userUseCase)
	postController = postHttpGinPackage.NewPostHandler(postUseCase)

	// Routers.
	userRouter = userHttpGinPackage.NewUserRouter(userController)
	postRouter = postHttpGinPackage.NewPostRouter(postController)

	// Create the Gin Engine instance.
	server = gin.Default()
}

func main() {
	startGinServer()
	// startGrpcServer(config)
}

func startGinServer() {
	loadConfig := commonUtility.LoadConfig()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))
	router := server.Group("/api")
	userRouter.UserRouter(router, userUseCase)
	postRouter.PostRouter(router, userUseCase)
	log.Fatal(server.Run(":" + loadConfig.GinConfig.Port))
}

// func startGrpcServer(config config.Config) {
// 	userGrpcServer, err := userGrpcPackage.NewGrpcUserServer(config, userUseCase, userCollection)

// 	if err != nil {
// 		log.Fatal("cannot createt gRPC User Server: ", err)
// 	}

// 	postGrpcServer, err := postGrpcPackage.NewGrpcPostServer(postUseCase)

// 	if err != nil {
// 		log.Fatal("cannot create gRPC Post Server: ", err)
// 	}

// 	grpcServer := grpc.NewServer()

// 	// Register User gRPC server.
// 	userProtobufV1.RegisterUserUseCaseServer(grpcServer, userGrpcServer)

// 	// Register Post gRPC server.
// 	postProtobufV1.RegisterPostUseCaseServer(grpcServer, postGrpcServer)

// 	reflection.Register(grpcServer)

// 	listener, err := net.Listen("tcp", config.GrpcServerAddress)

// 	if err != nil {
// 		log.Fatal("cannot create grpc server: ", err)
// 	}

// 	log.Printf("start grpc server on %s", listener.Addr().String())
// 	err = grpcServer.Serve(listener)

// 	if err != nil {
// 		log.Fatal("cannot create grpc server: ", err)
// 	}
// }
