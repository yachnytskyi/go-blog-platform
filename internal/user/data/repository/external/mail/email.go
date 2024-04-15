package mail

import (
	"bytes"
	"context"
	"crypto/tls"
	"html/template"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
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
		if validator.IsError(walkError) {
			sendEmailInternalError := domainError.NewInternalError(location+"ParseTemplateDirectory.Walk", walkError.Error())
			return sendEmailInternalError
		}
		if info.IsDir() {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	logging.Logger(loggerMessage)
	if validator.IsError(filePathWalkError) {
		sendEmailInternalError := domainError.NewInternalError(location+"ParseTemplateDirectory."+loggerMessage, filePathWalkError.Error())
		return nil, sendEmailInternalError
	}

	return template.ParseFiles(paths...)
}

func SendEmail(ctx context.Context, user userModel.User, data userModel.EmailData) error {
	emailConfig := config.AppConfig.Email
	smtpPass := emailConfig.SMTPPassword
	smtpUser := emailConfig.SMTPUser
	smtpHost := emailConfig.SMTPHost
	smtpPort := emailConfig.SMTPPort

	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send an email.
	message, prepareSendMessageError := PrepareSendMessage(ctx, user.Email, data)
	if validator.IsError(prepareSendMessageError) {
		return prepareSendMessageError
	}
	dialAndSendError := dialer.DialAndSend(message)
	if validator.IsError(dialAndSendError) {
		sendEmailInternalError := domainError.NewInternalError(location+"SendEmail.DialAndSend", dialAndSendError.Error())
		return sendEmailInternalError
	}

	return nil
}

func PrepareSendMessage(ctx context.Context, userEmail string, data userModel.EmailData) (*gomail.Message, error) {
	emailConfig := config.AppConfig.Email

	// Prepare data.
	from := emailConfig.EmailFrom
	to := userEmail

	var body bytes.Buffer
	template, parseTemplateDirectoryError := ParseTemplateDirectory(data.TemplatePath)
	if validator.IsError(parseTemplateDirectoryError) {
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
