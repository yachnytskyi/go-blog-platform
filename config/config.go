package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

var (
	AppConfig ApplicationConfig
)

const (
	environmentsPath   = "config/environment/.env."
	devInveronmentName = "dev"
	defaultConfigPath  = "config/yaml/dev.application.yaml"
	location           = "config.LoadConfig."
)

type ApplicationConfig struct {
	Database      string        `mapstructure:"Database"`
	Domain        string        `mapstructure:"Domain"`
	MongoDBConfig MongoDBConfig `mapstructure:"MongoDB"`
	GinConfig     GinConfig     `mapstructure:"Gin"`
	GRPCConfig    GRPCConfig    `mapstructure:"Grpc"`
	AccessToken   AccessToken   `mapstructure:"Access_Token"`
	RefreshToken  RefreshToken  `mapstructure:"Refresh_Token"`
	Email         Email         `mapstructure:"Email"`
}

type MongoDBConfig struct {
	Name string `mapstructure:"Name"`
	URI  string `mapstructure:"URI"`
}

type GinConfig struct {
	Port string `mapstructure:"PORT"`
}

type GRPCConfig struct {
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
	loadEnvironmentsError := godotenv.Load(environmentsPath + devInveronmentName)
	if validator.IsErrorNotNil(loadEnvironmentsError) {
		loadEnvironmentsInternalError := domainError.NewInternalError(location+"Load", loadEnvironmentsError.Error())
		logging.Logger(loadEnvironmentsInternalError)
		return loadEnvironmentsInternalError
	}
	configPath := os.Getenv("CONFIG_PATH")
	if validator.IsStringEmpty(configPath) {
		configPath = defaultConfigPath
	}
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	readInConfigError := viper.ReadInConfig()
	if validator.IsErrorNotNil(readInConfigError) {
		readInInternalError := domainError.NewInternalError(location+"ReadInConfig", readInConfigError.Error())
		logging.Logger(readInInternalError)
		return readInInternalError
	}
	viper.SetDefault("Database", "MongoDB")
	viper.SetDefault("Domain", "UseCase")
	viper.SetDefault("MongoDB.Name", "golang_mongodb")
	viper.SetDefault("MongoDB.URI", "mongodb://root:root@localhost:27017")
	viper.SetDefault("Gin.Port", "8080")
	unmarshalError = viper.Unmarshal(&AppConfig)
	return
}
