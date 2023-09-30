package common

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func LoadConfig() config.Config {
	loadConfig, loadConfigError := config.LoadConfig(constant.ConfigPath)
	if validator.IsErrorNotNil(loadConfigError) {
		var loadConfigInternalError domainError.InternalError
		loadConfigInternalError.Location = "pkg/utility/common/LoadConfig"
		loadConfigInternalError.Reason = loadConfigError.Error()
		logging.Logger(loadConfigInternalError)
	}
	return loadConfig
}
