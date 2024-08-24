package common

import (
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
)

type MockConfig struct {
	ApplicationConfig *config.ApplicationConfig
}

func NewMockConfig() MockConfig {
	return MockConfig{
		ApplicationConfig: &config.ApplicationConfig{},
	}
}

func (mock MockConfig) GetConfig() *config.ApplicationConfig {
	return mock.ApplicationConfig
}
