package model

import "time"

type ViperConfig struct {
	Core         ViperCore         `mapstructure:"Core"`
	Security     ViperSecurity     `mapstructure:"Security"`
	MongoDB      ViperMongoDB      `mapstructure:"MongoDB"`
	Gin          ViperGin          `mapstructure:"Gin"`
	GRPC         ViperGRPC         `mapstructure:"Grpc"`
	AccessToken  ViperAccessToken  `mapstructure:"Access_Token"`
	RefreshToken ViperRefreshToken `mapstructure:"Refresh_Token"`
	Email        ViperEmail        `mapstructure:"Email"`
}

type ViperCore struct {
	Logger   string `mapstructure:"Logger"`
	Database string `mapstructure:"Database"`
	UseCase  string `mapstructure:"UseCase"`
	Delivery string `mapstructure:"Delivery"`
}

type ViperSecurity struct {
	CookieSecure                    bool        `mapstructure:"Cookie_Secure"`
	HTTPOnly                        bool        `mapstructure:"HTTP_Only"`
	RateLimit                       float64     `mapstructure:"Rate_Limit"`
	ContentSecurityPolicyHeader     ViperHeader `mapstructure:"Content_Security_Policy_Header"`
	ContentSecurityPolicyHeaderFull ViperHeader `mapstructure:"Content_Security_Policy_Header_Full"`
	StrictTransportSecurityHeader   ViperHeader `mapstructure:"Strict_Transport_Security_Header"`
	XContentTypeOptionsHeader       ViperHeader `mapstructure:"X_Content_Type_Options_Header"`
	AllowedHTTPMethods              []string    `mapstructure:"Allowed_HTTP_Methods"`
	AllowedContentTypes             []string    `mapstructure:"Allowed_Content_Types"`
}

type ViperHeader struct {
	Key   string `mapstructure:"Key"`
	Value string `mapstructure:"Value"`
}

type ViperMongoDB struct {
	Name string `mapstructure:"Name"`
	URI  string `mapstructure:"URI"`
}

type ViperGin struct {
	Port             string `mapstructure:"Port"`
	AllowOrigins     string `mapstructure:"Allow_Origins"`
	AllowCredentials bool   `mapstructure:"Allow_Credentials"`
	ServerGroup      string `mapstructure:"Server_Group"`
}

type ViperGRPC struct {
	ServerUrl string `mapstructure:"Server_Url"`
}

type ViperAccessToken struct {
	PrivateKey string        `mapstructure:"Private_Key"`
	PublicKey  string        `mapstructure:"Public_Key"`
	ExpiredIn  time.Duration `mapstructure:"Expired_In"`
	MaxAge     int           `mapstructure:"Max_Age"`
}

type ViperRefreshToken struct {
	PrivateKey string        `mapstructure:"Private_Key"`
	PublicKey  string        `mapstructure:"Public_Key"`
	ExpiredIn  time.Duration `mapstructure:"Expired_In"`
	MaxAge     int           `mapstructure:"Max_Age"`
}

type ViperEmail struct {
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
