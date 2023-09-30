package main

// Require the packages.
import (
	"context"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"

	postPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postHttpGinPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/http/gin"

	"github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency"

	postUseCasePackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"

	userPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	// dependency "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency"

	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"

	userHttpGinPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"

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

const (
	location = "cmd.server.init."
)

// Init function that will run before the "main" function.
func init() {
	loadConfig, loadConfigError := config.LoadConfig(config.ConfigPath)
	if loadConfigError != nil {
		loadConfigInternalError := domainError.NewInternalError(location+"init.LoadConfig", loadConfigError.Error())
		logging.Logger(loadConfigInternalError)
	}

	// Create a context.
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultContextTimer)
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
	config, err := config.LoadConfig(config.ConfigPath)
	if err != nil {
		log.Fatal("Could not load config", err)
	}
	startGinServer(config)
	// startGrpcServer(config)
}

func startGinServer(config config.Config) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))
	router := server.Group("/api")
	userRouter.UserRouter(router, userUseCase)
	postRouter.PostRouter(router, userUseCase)
	log.Fatal(server.Run(":" + config.GinConfig.Port))
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
