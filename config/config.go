package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database string `mapstructure:"Database"`
	Domain   string `mapstrucrure:"Domain"`

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
	ClientOriginUrl string `mapstructure:"Client_Origin_Url"`
	EmailFrom       string `mapstructure:"Email_From"`
	SMTPHost        string `mapstructure:"SMTP_Host"`
	SMTPPassword    string `mapstructure:"SMTP_Password"`
	SMTPPort        int    `mapstructure:"SMTP_Port"`
	SMTPUser        string `mapstructure:"SMTP_User"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&config)
	return
}
