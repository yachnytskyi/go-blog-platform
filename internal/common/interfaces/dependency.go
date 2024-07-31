package interfaces

import (
	"context"

	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
)

// Config is an interface that defines a method for retrieving the application's configuration.
type Config interface {
	GetConfig() *config.ApplicationConfig
}

// Logger is an interface that defines methods for logging at different levels.
type Logger interface {
	Trace(data error)
	Debug(data error)
	Info(data error)
	Warn(data error)
	Error(data error)
	Fatal(data error)
	Panic(data error)
}

// Email is an interface that defines methods for sending emails.
type Email interface {
	SendEmail(configInstance Config, logger Logger, location string, data any, emailData EmailData) error
}

// Repository is an interface that defines methods for creating and managing repository instances.
type Repository interface {
	CreateRepository(ctx context.Context) any
	NewRepository(createRepository any, repository any) any
	Close
}

// UseCase is an interface that defines methods for creating use case instances.
type UseCase interface {
	NewUseCase(email Email, repository any) any
}

// Delivery is an interface that defines methods for creating delivery components and managing the server.
type Delivery interface {
	CreateDelivery(serverRouters ServerRouters)
	NewController(userUseCase UserUseCase, usecase any) any
	NewRouter(router any) Router
	LaunchServer(ctx context.Context, repository Repository)
	Close
}

// Close is an interface that defines a method for closing resources or services.
type Close interface {
	Close(ctx context.Context) // Closes resources or services.
}

// ServerRouters holds the routers for different modules of the application.
type ServerRouters struct {
	UserRouter Router
	PostRouter Router
	// Add other routers as needed.
}

// NewServerRouters creates a new instance of ServerRouters with the given routers.
func NewServerRouters(userRouter Router, postRouter Router) ServerRouters {
	return ServerRouters{
		UserRouter: userRouter,
		PostRouter: postRouter,
		// Add other routers as needed.
	}
}

// EmailData holds the data required for sending an email.
type EmailData struct {
	Recipient    string // Recipient's email address.
	URL          string // URL to be included in the email.
	TemplateName string // Name of the email template.
	TemplatePath string // Path to the email template.
	FirstName    string // Recipient's first name.
	Subject      string // Subject of the email.
}

func NewEmailData(recipient, url, templateName, templatePath, firstName, subject string) EmailData {
	return EmailData{
		Recipient:    recipient,
		URL:          url,
		TemplateName: templateName,
		TemplatePath: templatePath,
		FirstName:    firstName,
		Subject:      subject,
	}
}
