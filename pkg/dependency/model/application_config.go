package model

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
)

var (
	ApplicationConfig config.ApplicationConfig
)

func LoadConfig() {
	ApplicationConfig, _ = config.LoadConfig(constant.ConfigPath)
}
