package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

var (
	AppConfig ApplicationConfig
)

const (
	version                    = "v1"
	environmentsPath           = "config/environment/.env."
	devInvironmentName         = "dev"
	defaultMongoDBName         = "golang-mongodb"
	defaultMongoDBURI          = "mongodb://root:root@localhost:27017/golang_mongodb"
	defaultGinPort             = "8080"
	defaultGinAllowOrirings    = "http://localhost:8080"
	defaultGinAllowCredentuals = true
	defaultGinServerGroup      = "/api"
	location                   = "config.LoadConfig."
)

type ApplicationConfig struct {
	Core         Core         `mapstructure:"Core"`
	MongoDB      MongoDB      `mapstructure:"MongoDB"`
	Gin          Gin          `mapstructure:"Gin"`
	GRPC         GRPC         `mapstructure:"Grpc"`
	Token        Token        `mapstructure:"Token"`
	AccessToken  AccessToken  `mapstructure:"Access_Token"`
	RefreshToken RefreshToken `mapstructure:"Refresh_Token"`
	Email        Email        `mapstructure:"Email"`
}

type Core struct {
	Database string `mapstructure:"Database"`
	Domain   string `mapstructure:"Domain"`
	Delivery string `mapstructure:"Delivery"`
}

type MongoDB struct {
	Name string `mapstructure:"Name"`
	URI  string `mapstructure:"URI"`
}

type Gin struct {
	Port             string `mapstructure:"Port"`
	AllowOrigins     string `mapstructure:"Allow_Origins"`
	AllowCredentials bool   `mapstructure:"Allow_Credentials"`
	ServerGroup      string `mapstructure:"Server_Group"`
}

type GRPC struct {
	ServerUrl string `mapstructure:"Server_Url"`
}

type Token struct {
	CookieSecure bool `mapstructure:"Cookie_Secure"`
	HttpOnly     bool `mapstructure:"Http_Only"`
}

type AccessToken struct {
	PrivateKey string        `mapstructure:"Private_Key"`
	PublicKey  string        `mapstructure:"Public_Key"`
	ExpiredIn  time.Duration `mapstructure:"Expired_In"`
	MaxAge     int           `mapstructure:"Max_Age"`
}

type RefreshToken struct {
	PrivateKey string        `mapstructure:"Private_Key"`
	PublicKey  string        `mapstructure:"Public_Key"`
	ExpiredIn  time.Duration `mapstructure:"Expired_In"`
	MaxAge     int           `mapstructure:"Max_Age"`
}

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

func LoadConfig() (unmarshalError error) {
	loadEnvironmentsError := godotenv.Load(environmentsPath + devInvironmentName)
	if validator.IsErrorNotNil(loadEnvironmentsError) {
		loadEnvironmentsInternalError := domainError.NewInternalError(location+"Load", loadEnvironmentsError.Error())
		logging.Logger(loadEnvironmentsInternalError)
		return loadEnvironmentsInternalError
	}
	configPath := os.Getenv(version)
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	readInConfigError := viper.ReadInConfig()
	if validator.IsErrorNotNil(readInConfigError) {
		readInInternalError := domainError.NewInternalError(location+"ReadInConfig", readInConfigError.Error())
		logging.Logger(readInInternalError)
		return readInInternalError
	}
	viper.SetDefault("Database", constants.MongoDB)
	viper.SetDefault("Domain", constants.UseCase)
	viper.SetDefault("Delivery", constants.Gin)
	viper.SetDefault("MongoDB.Name", defaultMongoDBName)
	viper.SetDefault("MongoDB.URI", defaultMongoDBURI)
	viper.SetDefault("Gin.Port", defaultGinPort)
	viper.SetDefault("Gin.AllowOrigins", defaultGinAllowOrirings)
	viper.SetDefault("Gin.AllowCredentials", defaultGinAllowCredentuals)
	viper.SetDefault("Gin.ServerGroup", defaultGinServerGroup)
	unmarshalError = viper.Unmarshal(&AppConfig)
	return
}
