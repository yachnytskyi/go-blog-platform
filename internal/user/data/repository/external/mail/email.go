package mail

import (
	"bytes"
	"context"
	"crypto/tls"
	"html/template"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"gopkg.in/gomail.v2"
)

const (
	loggerMessage = "parsing template..."
	location      = "User.Data.Repository.External.Mail."
)

// Email template parser.
func ParseTemplateDirectory(templatePath string) (*template.Template, error) {
	var paths []string
	filePathWalkError := filepath.Walk(templatePath, func(path string, info os.FileInfo, walkError error) error {
		if validator.IsErrorNotNil(walkError) {
			sendEmailInternalError := domainError.NewInternalError(location+"ParseTemplateDirectory.Walk", walkError.Error())
			return sendEmailInternalError
		}
		if validator.IsBooleanNotTrue(info.IsDir()) {
			paths = append(paths, path)
		}
		return nil
	})
	logging.Logger(loggerMessage)
	if validator.IsErrorNotNil(filePathWalkError) {
		sendEmailInternalError := domainError.NewInternalError(location+"ParseTemplateDirectory."+loggerMessage, filePathWalkError.Error())
		return nil, sendEmailInternalError
	}
	return template.ParseFiles(paths...)
}

func SendEmail(ctx context.Context, user userModel.User, data userModel.EmailData) error {
	loadConfig, loadConfigError := config.LoadConfig(config.ConfigPath)
	if validator.IsErrorNotNil(loadConfigError) {
		loadConfigInternalError := domainError.NewInternalError(location+"SendEmail.LoadConfig", loadConfigError.Error())
		return loadConfigInternalError
	}
	smtpPass := loadConfig.Email.SMTPPassword
	smtpUser := loadConfig.Email.SMTPUser
	smtpHost := loadConfig.Email.SMTPHost
	smtpPort := loadConfig.Email.SMTPPort

	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send an email.
	message, prepareSendMessageError := PrepareSendMessage(ctx, user.Email, data)
	if validator.IsErrorNotNil(prepareSendMessageError) {
		return prepareSendMessageError
	}
	dialAndSendError := dialer.DialAndSend(message)
	if validator.IsErrorNotNil(dialAndSendError) {
		sendEmailInternalError := domainError.NewInternalError(location+"SendEmail.DialAndSend", dialAndSendError.Error())
		return sendEmailInternalError
	}
	return nil
}

func PrepareSendMessage(ctx context.Context, userEmail string, data userModel.EmailData) (*gomail.Message, error) {
	loadConfig, loadConfigError := config.LoadConfig(config.ConfigPath)
	if validator.IsErrorNotNil(loadConfigError) {
		loadConfigInternalError := domainError.NewInternalError(location+"SendEmail.PrepareSendMessage.LoadConfig", loadConfigError.Error())
		return nil, loadConfigInternalError
	}
	// Prepare data.
	from := loadConfig.Email.EmailFrom
	to := userEmail

	var body bytes.Buffer
	template, parseTemplateDirectoryError := ParseTemplateDirectory(data.TemplatePath)
	if validator.IsErrorNotNil(parseTemplateDirectoryError) {
		parseTemplateDirectoryInternalError := domainError.NewInternalError(location+"SendEmail.PrepareSendMessage.ParseTemplateDirectory", parseTemplateDirectoryError.Error())
		return nil, parseTemplateDirectoryInternalError
	}

	template = template.Lookup(data.TemplateName)
	template.Execute(&body, &data)
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", data.Subject)
	message.SetBody("text/html", body.String())
	message.AddAlternative("text/plain", html2text.HTML2Text(body.String()))
	return message, nil
}
