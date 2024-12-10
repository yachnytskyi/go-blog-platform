package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "pkg.dependency.factory.config."
)

func NewViper() *config.ApplicationConfig {
	loadEnvironmentsError := godotenv.Load(constants.EnvironmentsPath)
	if validator.IsError(loadEnvironmentsError) {
		loadEnvironmentsInternalError := domain.NewInternalError(location+"viper.Load", loadEnvironmentsError.Error())
		log.Println(loadEnvironmentsInternalError)
		loadDefaultEnvironment()
	}

	configPath := os.Getenv("version")
	viperInstance := viper.New()
	viperInstance.SetConfigFile(configPath)
	viperInstance.AutomaticEnv()

	readInConfigError := viperInstance.ReadInConfig()
	if validator.IsError(readInConfigError) {
		readInInternalError := domain.NewInternalError(location+"viper.ReadInConfig", readInConfigError.Error())
		log.Println(readInInternalError)
		loadDefaultConfig(viperInstance)
	}

	var yamlConfig config.YamlConfig
	unmarshalError := viperInstance.Unmarshal(&yamlConfig)
	if validator.IsError(unmarshalError) {
		panic(domain.NewInternalError(location+"viper.Unmarshal", unmarshalError.Error()))
	}

	applicationConfig := yamlConfigToApplicationConfigMapper(&yamlConfig)
	return &applicationConfig
}

func loadDefaultEnvironment() {
	defaultEnvironmentError := godotenv.Load(constants.DefaultEnvironmentsPath + constants.DefaultEnvironment)
	if validator.IsError(defaultEnvironmentError) {
		panic(domain.NewInternalError(location+"viper.loadDefaultEnvironment", defaultEnvironmentError.Error()))
	}
	log.Println(domain.NewInfoMessage(location+"viper.loadDefaultEnvironment", constants.DefaultConfigPathNotification))
}

func loadDefaultConfig(viper *viper.Viper) {
	viper.SetConfigFile(constants.DefaultConfigPath)
	readInConfigError := viper.ReadInConfig()
	if validator.IsError(readInConfigError) {
		panic(domain.NewInternalError(location+"viper.loadDefaultConfig", readInConfigError.Error()))
	}
	log.Println(domain.NewInfoMessage(location+"viper.loadDefaultConfig", constants.DefaultConfigPathNotification))
}

func yamlConfigToApplicationConfigMapper(yamlConfig *config.YamlConfig) config.ApplicationConfig {
	return config.ApplicationConfig{
		Core:         convertCore(&yamlConfig.Core),
		MongoDB:      convertMongoDB(&yamlConfig.MongoDB),
		Security:     convertSecurity(&yamlConfig.Security),
		Gin:          convertGin(&yamlConfig.Gin),
		GRPC:         convertGRPC(&yamlConfig.GRPC),
		AccessToken:  convertAccessToken(&yamlConfig.AccessToken),
		RefreshToken: convertRefreshToken(&yamlConfig.RefreshToken),
		Email:        convertEmail(&yamlConfig.Email),
	}
}

func convertCore(core *config.YamlCore) config.Core {
	return config.Core{
		Logger:   core.Logger,
		Email:    core.Email,
		Database: core.Database,
		Delivery: core.Delivery,
	}
}

func convertMongoDB(mongoDB *config.YamlMongoDB) config.MongoDB {
	return config.MongoDB{
		Name: mongoDB.Name,
		URI:  mongoDB.URI,
	}
}

func convertSecurity(security *config.YamlSecurity) config.Security {
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

func convertHeader(header *config.YamlHeader) config.Header {
	return config.Header{
		Key:   header.Key,
		Value: header.Value,
	}
}

func convertGin(gin *config.YamlGin) config.Gin {
	return config.Gin{
		Port:             gin.Port,
		AllowOrigins:     gin.AllowOrigins,
		AllowCredentials: gin.AllowCredentials,
		ServerGroup:      gin.ServerGroup,
	}
}

func convertGRPC(grpc *config.YamlGRPC) config.GRPC {
	return config.GRPC{
		ServerUrl: grpc.ServerUrl,
	}
}

func convertAccessToken(accessToken *config.YamlAccessToken) config.AccessToken {
	return config.AccessToken{
		PrivateKey: accessToken.PrivateKey,
		PublicKey:  accessToken.PublicKey,
		ExpiredIn:  accessToken.ExpiredIn,
		MaxAge:     accessToken.MaxAge,
	}
}

func convertRefreshToken(refreshToken *config.YamlRefreshToken) config.RefreshToken {
	return config.RefreshToken{
		PrivateKey: refreshToken.PrivateKey,
		PublicKey:  refreshToken.PublicKey,
		ExpiredIn:  refreshToken.ExpiredIn,
		MaxAge:     refreshToken.MaxAge,
	}
}

func convertEmail(email *config.YamlEmail) config.Email {
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
