package model

import "time"

type YamlConfig struct {
	Core         YamlCore         `mapstructure:"Core"`
	Security     YamlSecurity     `mapstructure:"Security"`
	MongoDB      YamlMongoDB      `mapstructure:"MongoDB"`
	Gin          YamlGin          `mapstructure:"Gin"`
	GRPC         YamlGRPC         `mapstructure:"Grpc"`
	AccessToken  YamlAccessToken  `mapstructure:"Access_Token"`
	RefreshToken YamlRefreshToken `mapstructure:"Refresh_Token"`
	Email        YamlEmail        `mapstructure:"Email"`
}

type YamlCore struct {
	Logger   string `mapstructure:"Logger"`
	Email    string `mapstructure:"Email"`
	Database string `mapstructure:"Database"`
	UseCase  string `mapstructure:"UseCase"`
	Delivery string `mapstructure:"Delivery"`
}

type YamlSecurity struct {
	CookieSecure                    bool       `mapstructure:"Cookie_Secure"`
	HTTPOnly                        bool       `mapstructure:"HTTP_Only"`
	RateLimit                       float64    `mapstructure:"Rate_Limit"`
	ContentSecurityPolicyHeader     YamlHeader `mapstructure:"Content_Security_Policy_Header"`
	ContentSecurityPolicyHeaderFull YamlHeader `mapstructure:"Content_Security_Policy_Header_Full"`
	StrictTransportSecurityHeader   YamlHeader `mapstructure:"Strict_Transport_Security_Header"`
	XContentTypeOptionsHeader       YamlHeader `mapstructure:"X_Content_Type_Options_Header"`
	AllowedHTTPMethods              []string   `mapstructure:"Allowed_HTTP_Methods"`
	AllowedContentTypes             []string   `mapstructure:"Allowed_Content_Types"`
}

type YamlHeader struct {
	Key   string `mapstructure:"Key"`
	Value string `mapstructure:"Value"`
}

type YamlMongoDB struct {
	Name string `mapstructure:"Name"`
	URI  string `mapstructure:"URI"`
}

type YamlGin struct {
	Port             string `mapstructure:"Port"`
	AllowOrigins     string `mapstructure:"Allow_Origins"`
	AllowCredentials bool   `mapstructure:"Allow_Credentials"`
	ServerGroup      string `mapstructure:"Server_Group"`
}

type YamlGRPC struct {
	ServerUrl string `mapstructure:"Server_Url"`
}

type YamlAccessToken struct {
	PrivateKey string        `mapstructure:"Private_Key"`
	PublicKey  string        `mapstructure:"Public_Key"`
	ExpiredIn  time.Duration `mapstructure:"Expired_In"`
	MaxAge     int           `mapstructure:"Max_Age"`
}

type YamlRefreshToken struct {
	PrivateKey string        `mapstructure:"Private_Key"`
	PublicKey  string        `mapstructure:"Public_Key"`
	ExpiredIn  time.Duration `mapstructure:"Expired_In"`
	MaxAge     int           `mapstructure:"Max_Age"`
}

type YamlEmail struct {
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
