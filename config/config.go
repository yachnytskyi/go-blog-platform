package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// AppConfig holds the global application configuration.
var AppConfig ApplicationConfig

// Constants for configuration paths and default values.
const (
	version                  = "v1"
	environmentsPath         = "config/environment/.env."
	localDevEnvironment      = "local.dev"
	dockerDevEnvironment     = "docker.dev"
	defaultVersion           = "v1"
	defaultEnvironmentsPath  = "config/environment/.env."
	defaultConfigPath        = "config/yaml/v1/local.dev.application.yaml"
	defaultDockerEnvironment = "local.dev"
	location                 = "config.LoadConfig."
)

// ApplicationConfig defines the structure of the application configuration.
type ApplicationConfig struct {
	Core         Core         `mapstructure:"Core"`
	MongoDB      MongoDB      `mapstructure:"MongoDB"`
	Security     Security     `mapstructure:"Security"`
	Gin          Gin          `mapstructure:"Gin"`
	GRPC         GRPC         `mapstructure:"Grpc"`
	AccessToken  AccessToken  `mapstructure:"Access_Token"`
	RefreshToken RefreshToken `mapstructure:"Refresh_Token"`
	Email        Email        `mapstructure:"Email"`
}

// Core defines the core settings for the application.
type Core struct {
	Logger   string `mapstructure:"Logger"`
	Database string `mapstructure:"Database"`
	UseCase  string `mapstructure:"UseCase"`
	Delivery string `mapstructure:"Delivery"`
}

// Security defines the security settings for the application.
type Security struct {
	CookieSecure                    bool     `mapstructure:"Cookie_Secure"`
	HTTPOnly                        bool     `mapstructure:"HTTP_Only"`
	RateLimit                       float64  `mapstructure:"Rate_Limit"`
	ContentSecurityPolicyHeader     Header   `mapstructure:"Content_Security_Policy_Header"`
	ContentSecurityPolicyHeaderFull Header   `mapstructure:"Content_Security_Policy_Header_Full"`
	StrictTransportSecurityHeader   Header   `mapstructure:"Strict_Transport_Security_Header"`
	XContentTypeOptionsHeader       Header   `mapstructure:"X_Content_Type_Options_Header"`
	AllowedHTTPMethods              []string `mapstructure:"Allowed_HTTP_Methods"`
	AllowedContentTypes             []string `mapstructure:"Allowed_Content_Types"`
}

// Header defines the structure for HTTP headers.
type Header struct {
	Key   string `mapstructure:"Key"`
	Value string `mapstructure:"Value"`
}

// MongoDB defines the MongoDB settings.
type MongoDB struct {
	Name string `mapstructure:"Name"`
	URI  string `mapstructure:"URI"`
}

// Gin defines the settings for the Gin web framework.
type Gin struct {
	Port             string `mapstructure:"Port"`
	AllowOrigins     string `mapstructure:"Allow_Origins"`
	AllowCredentials bool   `mapstructure:"Allow_Credentials"`
	ServerGroup      string `mapstructure:"Server_Group"`
}

// GRPC defines the settings for gRPC.
type GRPC struct {
	ServerUrl string `mapstructure:"Server_Url"`
}

// AccessToken defines the settings for access tokens.
type AccessToken struct {
	PrivateKey string        `mapstructure:"Private_Key"`
	PublicKey  string        `mapstructure:"Public_Key"`
	ExpiredIn  time.Duration `mapstructure:"Expired_In"`
	MaxAge     int           `mapstructure:"Max_Age"`
}

// RefreshToken defines the settings for refresh tokens.
type RefreshToken struct {
	PrivateKey string        `mapstructure:"Private_Key"`
	PublicKey  string        `mapstructure:"Public_Key"`
	ExpiredIn  time.Duration `mapstructure:"Expired_In"`
	MaxAge     int           `mapstructure:"Max_Age"`
}

// Email defines the settings for email.
type Email struct {
	ClientOriginUrl               string `mapstructure:"Client_Origin_Url"`
	EmailFrom                     string `mapstructure:"Email_From"`
	SMTPHost                      string `mapstructure:"SMTP_Host"`
	SMTPPassword                  string `mapstructure:"SMTP_Password"`
	SMTPPort                      int    `mapstructure:"SMTP_Port"`
	SMTPUser                      string `mapstructure:"SMTP_User"`
	UserConfirmationTemplateName  string `mapstructure:"User_Confirmation_Template_Name"`
	UserConfirmationTemplatePath  string `mapstructure:"User_Confirmation_Template_Path"`
	ForgottenPasswordTemplateName string `mapstructure:"Forgotten_Password_Template_Name"`
	ForgottenPasswordTemplatePath string `mapstructure:"Forgotten_Password_Template_Path"`
}

// LoadConfig loads the application configuration from environment variables and defaults.
// It sets up the global AppConfig variable.
func LoadConfig() {
	loadEnvironmentsError := godotenv.Load(environmentsPath + localDevEnvironment)
	if validator.IsError(loadEnvironmentsError) {
		loadEnvironmentsInternalError := domainError.NewInternalError(location+"Load", loadEnvironmentsError.Error())
		log.Println(loadEnvironmentsInternalError)
		loadDefaultEnvironment()
	}

	// Set the configuration file path from the environment variable.
	configPath := os.Getenv(version)
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	// Read the configuration file.
	readInConfigError := viper.ReadInConfig()
	if validator.IsError(readInConfigError) {
		readInInternalError := domainError.NewInternalError(location+"ReadInConfig", readInConfigError.Error())
		log.Println(readInInternalError)
		loadDefaultConfig()
	}

	// Unmarshal the configuration into the global AppConfig variable.
	unmarshalError := viper.Unmarshal(&AppConfig)
	if validator.IsError(unmarshalError) {
		panic(domainError.NewInternalError(location+"Unmarshal", unmarshalError.Error()))
	}
}

// loadDefaultEnvironment attempts to load the default environment configuration.
func loadDefaultEnvironment() {
	defaultEnvironmentError := godotenv.Load(defaultEnvironmentsPath + defaultDockerEnvironment)
	if validator.IsError(defaultEnvironmentError) {
		panic(domainError.NewInternalError(location+"loadDefaultEnvironment", defaultEnvironmentError.Error()))
	}

	// Log an informational message indicating the use of a default environment path.
	log.Println(domainError.NewInfoMessage(location+"loadDefaultEnvironment", "Using default environment path"))
}

// loadDefaultConfig attempts to load the default application configuration.
func loadDefaultConfig() {
	viper.SetConfigFile(defaultConfigPath)
	defaultConfigError := viper.ReadInConfig()
	if validator.IsError(defaultConfigError) {
		panic(domainError.NewInternalError(location+"loadDefaultConfig", defaultConfigError.Error()))
	}

	// Log an informational message indicating the use of a default configuration path.
	log.Println(domainError.NewInfoMessage(location+"loadDefaultConfig", "Using default configuration path"))
}

// GetCoreConfig retrieves the core configuration.
func GetCoreConfig() Core {
	return AppConfig.Core
}

// GetSecurityConfig retrieves the security configuration.
func GetSecurityConfig() Security {
	return AppConfig.Security
}

// GetMongoDBConfig retrieves the MongoDB configuration.
func GetMongoDBConfig() MongoDB {
	return AppConfig.MongoDB
}

// GetGinConfig retrieves the Gin configuration.
func GetGinConfig() Gin {
	return AppConfig.Gin
}

// GetEmailConfig retrieves the Email configuration.
func GetEmailConfig() Email {
	return AppConfig.Email
}

// GetAccessConfig retrieves the Access Token configuration.
func GetAccessConfig() AccessToken {
	return AppConfig.AccessToken
}

// GetRefreshConfig retrieves the Refresh Token configuration.
func GetRefreshConfig() RefreshToken {
	return AppConfig.RefreshToken
}
