package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	dependency "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

func init() {
	config.LoadConfig()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), constant.DefaultContextTimer)
	defer cancel()
	container := dependency.CreateApplication(ctx)
	go func() {
		container.DeliveryFactory.LaunchServer(ctx, container)
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

<<<<<<< HEAD
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoClient.Disconnect(ctx)

	startGinServer(config)
	// startGrpcServer(config)
}

func startGrpcServer(config config.Config) {
	userGrpcServer, err := userGrpcPackage.NewGrpcUserServer(config, userUseCase, userCollection)

	if err != nil {
		log.Fatal("cannot createt gRPC User Server: ", err)
	}

	postGrpcServer, err := postGrpcPackage.NewGrpcPostServer(postUseCase)

	if err != nil {
		log.Fatal("cannot create gRPC Post Server: ", err)
	}

	grpcServer := grpc.NewServer()

	// Register User gRPC server.
	userProtobufV1.RegisterUserUseCaseServer(grpcServer, userGrpcServer)

	// Register Post gRPC server.
	postProtobufV1.RegisterPostUseCaseServer(grpcServer, postGrpcServer)

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
=======
	// Perform Graceful Shutdown
	applicationModel.GracefulShutdown(ctx, container)
>>>>>>> 70c4f98c2b3734e7dce4c2e16d7ed270a8ba4713
}
