package common

import (
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
)

type MockConfig struct {
	ApplicationConfig *config.ApplicationConfig
}

func NewMockConfig() *config.ApplicationConfig {
	return &config.ApplicationConfig{}
}
