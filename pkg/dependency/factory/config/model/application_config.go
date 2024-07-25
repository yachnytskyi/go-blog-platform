package model

import (
	"time"
)

// ApplicationConfig defines the structure of the application configuration.
type ApplicationConfig struct {
	Core         Core
	MongoDB      MongoDB
	Security     Security
	Gin          Gin
	GRPC         GRPC
	AccessToken  AccessToken
	RefreshToken RefreshToken
	Email        Email
}

type Core struct {
	Logger   string
	Database string
	UseCase  string
	Delivery string
}

type Security struct {
	CookieSecure                    bool
	HTTPOnly                        bool
	RateLimit                       float64
	ContentSecurityPolicyHeader     Header
	ContentSecurityPolicyHeaderFull Header
	StrictTransportSecurityHeader   Header
	XContentTypeOptionsHeader       Header
	AllowedHTTPMethods              []string
	AllowedContentTypes             []string
}

type Header struct {
	Key   string
	Value string
}

type MongoDB struct {
	Name string
	URI  string
}

type Gin struct {
	Port             string
	AllowOrigins     string
	AllowCredentials bool
	ServerGroup      string
}

type GRPC struct {
	ServerUrl string
}

type AccessToken struct {
	PrivateKey string
	PublicKey  string
	ExpiredIn  time.Duration
	MaxAge     int
}

type RefreshToken struct {
	PrivateKey string
	PublicKey  string
	ExpiredIn  time.Duration
	MaxAge     int
}

type Email struct {
	ClientOriginUrl               string
	EmailFrom                     string
	SMTPHost                      string
	SMTPPassword                  string
	SMTPPort                      int
	SMTPUser                      string
	UserConfirmationTemplateName  string
	UserConfirmationTemplatePath  string
	ForgottenPasswordTemplateName string
	ForgottenPasswordTemplatePath string
}
