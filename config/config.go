package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

var (
	AppConfig ApplicationConfig
)

const (
	version            = "V1"
	environmentsPath   = "config/environment/.env."
	devInvironmentName = "dev"
	defaultMongoDBName = "golang-mongodb"
	defaultMongoDBURI  = "mongodb://root:root@localhost:27017"
	defaultServerPort  = "8080"
	location           = "config.LoadConfig."
)

type ApplicationConfig struct {
	Core         Core         `mapstructure:"Core"`
	MongoDB      MongoDB      `mapstructure:"MongoDB"`
	Gin          Gin          `mapstructure:"Gin"`
	GRPC         GRPC         `mapstructure:"Grpc"`
	AccessToken  AccessToken  `mapstructure:"Access_Token"`
	RefreshToken RefreshToken `mapstructure:"Refresh_Token"`
	Email        Email        `mapstructure:"Email"`
}

type Core struct {
	Database string `mapstructure:"Database"`
	Domain   string `mapstructure:"Domain"`
}

type MongoDB struct {
	Name string `mapstructure:"Name"`
	URI  string `mapstructure:"URI"`
}

type Gin struct {
	Port string `mapstructure:"PORT"`
}

type GRPC struct {
	ServerUrl string `mapstructure:"Server_Url"`
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
	viper.SetDefault("Database", constant.MongoDB)
	viper.SetDefault("Domain", constant.UseCase)
	viper.SetDefault("MongoDB.Name", defaultMongoDBName)
	viper.SetDefault("MongoDB.URI", defaultMongoDBURI)
	viper.SetDefault("Gin.Port", defaultServerPort)
	unmarshalError = viper.Unmarshal(&AppConfig)
	return
}
