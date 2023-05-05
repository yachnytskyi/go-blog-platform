package main

// Require the packages.
import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	pb "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/model/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	postPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postRepositoryPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post/data/repository/mongo"
	postHttpPackageV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/http/v1"
	postUseCasePackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	userPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userRepositoryPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo"
	userGrpcPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1"
	userHttpPackageV1 "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/v1"
	userUseCasePackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Create required variables that we'll re-assign later.
var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client
	redisClient *redis.Client

	userCollection *mongo.Collection
	userRepository userPackage.Repository
	userUseCase    userPackage.UseCase
	userHandler    userHttpPackageV1.UserHandler
	userRouter     userHttpPackageV1.UserRouter

	postCollection *mongo.Collection
	postRepository postPackage.Repository
	postUseCase    postPackage.UseCase
	postHandler    postHttpPackageV1.PostHandler
	postRouter     postHttpPackageV1.PostRouter

	templateInstance *template.Template
)

// Init function that will run before the "main" function.
func init() {

	// Load the .env variables.
	templateInstance = template.Must(template.ParseGlob("pkg/templates/*.html"))
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	// Create a context.
	ctx = context.TODO()

	// Connect to MongoDB.
	mongoconn := options.Client().ApplyURI(config.MongoURI)
	mongoClient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// Connect to Redis.
	redisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisURI,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisClient.Set(ctx, "test", "Redis has been launched", 0).Err()

	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")

	// Collections.
	userCollection = mongoClient.Database("golang_mongodb").Collection("users")
	postCollection = mongoClient.Database("golang_mongodb").Collection("posts")

	// Repositories.
	userRepository = userRepositoryPackage.NewUserRepository(userCollection)
	postRepository = postRepositoryPackage.NewPostRepository(postCollection)

	// Use Cases.
	userUseCase = userUseCasePackage.NewUserUseCase(userRepository)
	postUseCase = postUseCasePackage.NewUseCase(postRepository)

	// Handlers
	userHandler = userHttpPackageV1.NewUserHandler(userUseCase, templateInstance)
	postHandler = postHttpPackageV1.NewPostHandler(postUseCase)

	// Routers.
	userRouter = userHttpPackageV1.NewUserRouter(userHandler)
	postRouter = postHttpPackageV1.NewPostRouter(postHandler)

	// Create the Gin Engine instance.
	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoClient.Disconnect(ctx)

	startGinServer(config)
	// startGrpcServer(config)
}

func startGrpcServer(config config.Config) {
	userServer, err := userGrpcPackage.NewGrpcUserServer(config, userUseCase, userCollection)

	if err != nil {
		log.Fatal("cannot createt grpc server: ", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserUseCaseServer(grpcServer, userServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)

	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

	log.Printf("start grpc server on %s", listener.Addr().String())
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot create grpc server: ", err)
	}

}

func startGinServer(config config.Config) {
	value, err := redisClient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})

	userRouter.UserRouter(router, userUseCase)
	postRouter.PostRouter(router, userUseCase)

	log.Fatal(server.Run(":" + config.Port))
}
