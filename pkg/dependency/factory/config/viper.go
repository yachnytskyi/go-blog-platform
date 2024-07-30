package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

type Viper struct {
	ApplicationConfig *config.ApplicationConfig
}

func NewViper() Viper {
	loadEnvironmentsError := godotenv.Load(constants.EnvironmentsPath + constants.LocalEnvironment)
	if validator.IsError(loadEnvironmentsError) {
		loadEnvironmentsInternalError := domainError.NewInternalError(location+"viper.Load", loadEnvironmentsError.Error())
		log.Println(loadEnvironmentsInternalError)
		loadDefaultEnvironment()
	}

	configPath := os.Getenv(constants.Version)
	viperInstance := viper.New()
	viperInstance.SetConfigFile(configPath)
	viperInstance.AutomaticEnv()

	readInConfigError := viperInstance.ReadInConfig()
	if validator.IsError(readInConfigError) {
		readInInternalError := domainError.NewInternalError(location+"viper.ReadInConfig", readInConfigError.Error())
		log.Println(readInInternalError)
		loadDefaultConfig(viperInstance)
	}

	var viperConfig config.ViperConfig
	unmarshalError := viperInstance.Unmarshal(&viperConfig)
	if validator.IsError(unmarshalError) {
		panic(domainError.NewInternalError(location+"viper.Unmarshal", unmarshalError.Error()))
	}

	applicationConfig := viperConfigToApplicationConfigMapper(&viperConfig)
	return Viper{
		ApplicationConfig: &applicationConfig,
	}
}

func loadDefaultEnvironment() {
	defaultEnvironmentError := godotenv.Load(constants.DefaultEnvironmentsPath + constants.LocalEnvironment)
	if validator.IsError(defaultEnvironmentError) {
		panic(domainError.NewInternalError(location+"viper.loadDefaultEnvironment", defaultEnvironmentError.Error()))
	}
	log.Println(domainError.NewInfoMessage(location+"viper.loadDefaultEnvironment", constants.DefaultConfigPathNotification))
}

func loadDefaultConfig(viper *viper.Viper) {
	viper.SetConfigFile(constants.DefaultConfigPath)
	readInConfigError := viper.ReadInConfig()
	if validator.IsError(readInConfigError) {
		panic(domainError.NewInternalError(location+"viper.loadDefaultConfig", readInConfigError.Error()))
	}
	log.Println(domainError.NewInfoMessage(location+"viper.loadDefaultConfig", constants.DefaultConfigPathNotification))
}

func (viper Viper) GetConfig() *config.ApplicationConfig {
	return viper.ApplicationConfig
}

func viperConfigToApplicationConfigMapper(viperConfig *config.ViperConfig) config.ApplicationConfig {
	return config.ApplicationConfig{
		Core:         convertCore(&viperConfig.Core),
		MongoDB:      convertMongoDB(&viperConfig.MongoDB),
		Security:     convertSecurity(&viperConfig.Security),
		Gin:          convertGin(&viperConfig.Gin),
		GRPC:         convertGRPC(&viperConfig.GRPC),
		AccessToken:  convertAccessToken(&viperConfig.AccessToken),
		RefreshToken: convertRefreshToken(&viperConfig.RefreshToken),
		Email:        convertEmail(&viperConfig.Email),
	}
}

func convertCore(core *config.ViperCore) config.Core {
	return config.Core{
		Logger:   core.Logger,
		Email:    core.Email,
		Database: core.Database,
		UseCase:  core.UseCase,
		Delivery: core.Delivery,
	}
}

func convertMongoDB(mongoDB *config.ViperMongoDB) config.MongoDB {
	return config.MongoDB{
		Name: mongoDB.Name,
		URI:  mongoDB.URI,
	}
}

func convertSecurity(security *config.ViperSecurity) config.Security {
	return config.Security{
		CookieSecure:                    security.CookieSecure,
		HTTPOnly:                        security.HTTPOnly,
		RateLimit:                       security.RateLimit,
		ContentSecurityPolicyHeader:     convertHeader(&security.ContentSecurityPolicyHeader),
		ContentSecurityPolicyHeaderFull: convertHeader(&security.ContentSecurityPolicyHeaderFull),
		StrictTransportSecurityHeader:   convertHeader(&security.StrictTransportSecurityHeader),
		XContentTypeOptionsHeader:       convertHeader(&security.XContentTypeOptionsHeader),
		AllowedHTTPMethods:              security.AllowedHTTPMethods,
		AllowedContentTypes:             security.AllowedContentTypes,
	}
}

func convertHeader(header *config.ViperHeader) config.Header {
	return config.Header{
		Key:   header.Key,
		Value: header.Value,
	}
}

func convertGin(gin *config.ViperGin) config.Gin {
	return config.Gin{
		Port:             gin.Port,
		AllowOrigins:     gin.AllowOrigins,
		AllowCredentials: gin.AllowCredentials,
		ServerGroup:      gin.ServerGroup,
	}
}

func convertGRPC(grpc *config.ViperGRPC) config.GRPC {
	return config.GRPC{
		ServerUrl: grpc.ServerUrl,
	}
}

func convertAccessToken(accessToken *config.ViperAccessToken) config.AccessToken {
	return config.AccessToken{
		PrivateKey: accessToken.PrivateKey,
		PublicKey:  accessToken.PublicKey,
		ExpiredIn:  accessToken.ExpiredIn,
		MaxAge:     accessToken.MaxAge,
	}
}

func convertRefreshToken(refreshToken *config.ViperRefreshToken) config.RefreshToken {
	return config.RefreshToken{
		PrivateKey: refreshToken.PrivateKey,
		PublicKey:  refreshToken.PublicKey,
		ExpiredIn:  refreshToken.ExpiredIn,
		MaxAge:     refreshToken.MaxAge,
	}
}

func convertEmail(email *config.ViperEmail) config.Email {
	return config.Email{
		ClientOriginUrl:               email.ClientOriginUrl,
		EmailFrom:                     email.EmailFrom,
		SMTPHost:                      email.SMTPHost,
		SMTPPassword:                  email.SMTPPassword,
		SMTPPort:                      email.SMTPPort,
		SMTPUser:                      email.SMTPUser,
		UserConfirmationTemplateName:  email.UserConfirmationTemplateName,
		UserConfirmationTemplatePath:  email.UserConfirmationTemplatePath,
		ForgottenPasswordTemplateName: email.ForgottenPasswordTemplateName,
		ForgottenPasswordTemplatePath: email.ForgottenPasswordTemplatePath,
	}
}
