package interfaces

import (
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
)

// Config is an interface that defines a method for retrieving the application's configuration.
type Config interface {
	GetConfig() *config.ApplicationConfig
}

// Logger is an interface that defines methods for logging at different levels.
type Logger interface {
	Trace(data error)
	Debug(data error)
	Info(data error)
	Warn(data error)
	Error(data error)
	Fatal(data error)
	Panic(data error)
}
