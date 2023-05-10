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
	postHttpGinPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/http/gin"
	postUseCasePackage "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	userPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userRepositoryPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo"
	userGrpcPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/grpc/v1"
	userHttpGinPackage "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin"
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
	userHandler    userHttpGinPackage.UserHandler
	userRouter     userHttpGinPackage.UserRouter

	postCollection *mongo.Collection
	postRepository postPackage.Repository
	postUseCase    postPackage.UseCase
	postHandler    postHttpGinPackage.PostHandler
	postRouter     postHttpGinPackage.PostRouter

	templateInstance *template.Template
)

// Init function that will run before the "main" function.
func init() {

	// Load the .env variables.
	templateInstance = template.Must(template.ParseGlob("internal/user/delivery/http/utility/templates/*.html"))
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
	userHandler = userHttpGinPackage.NewUserHandler(userUseCase, templateInstance)
	postHandler = postHttpGinPackage.NewPostHandler(postUseCase)

	// Routers.
	userRouter = userHttpGinPackage.NewUserRouter(userHandler)
	postRouter = postHttpGinPackage.NewPostRouter(postHandler)

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
