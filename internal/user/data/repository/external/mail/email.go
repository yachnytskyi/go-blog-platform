package mail

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"gopkg.in/gomail.v2"
)

const (
	printlnMessage = "parsing template..."
	location       = "User.Data.Repository.External.Mail."
)

// Email template parser.
func ParseTemplateDirectory(directory string) (*template.Template, error) {
	var paths []string
	filePathWalkError := filepath.Walk(directory, func(path string, info os.FileInfo, walkError error) error {
		if validator.IsErrorNotNil(walkError) {
			sendEmailInternalError := domainError.NewInternalError(location+"ParseTemplateDirectory.Walk", walkError.Error())
			return sendEmailInternalError
		}
		if validator.IsBooleanNotTrue(info.IsDir()) {
			paths = append(paths, path)
		}
		return nil
	})
	fmt.Println(printlnMessage)
	if validator.IsErrorNotNil(filePathWalkError) {
		sendEmailInternalError := domainError.NewInternalError(location+"ParseTemplateDirectory."+printlnMessage, filePathWalkError.Error())
		return nil, sendEmailInternalError
	}
	return template.ParseFiles(paths...)
}

func SendEmail(ctx context.Context, user *userModel.User, data *userModel.EmailData) error {
	loadConfig, loadConfigError := config.LoadConfig(".")
	if validator.IsErrorNotNil(loadConfigError) {
		loadConfigInternalError := domainError.NewInternalError(location+"SendEmail.LoadConfig", loadConfigError.Error())
		return loadConfigInternalError
	}

	smtpPass := loadConfig.SMTPPassword
	smtpUser := loadConfig.SMTPUser
	smtpHost := loadConfig.SMTPHost
	smtpPort := loadConfig.SMTPPort

	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send an email.
	message, prepareSendMessageError := PrepareSendMessage(ctx, &loadConfig, user.Email, data)
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

func PrepareSendMessage(ctx context.Context, loadConfig *config.Config, userEmail string, data *userModel.EmailData) (*gomail.Message, error) {
	// Prepare data.
	templateName := loadConfig.UserEmailTemplateName
	from := loadConfig.EmailFrom
	to := userEmail

	var body bytes.Buffer
	template, parseTemplateDirectoryError := ParseTemplateDirectory(loadConfig.UserEmailTemplatePath)
	if validator.IsErrorNotNil(parseTemplateDirectoryError) {
		parseTemplateDirectoryInternalError := domainError.NewInternalError(location+"SendEmail.PrepareSendMessage.ParseTemplateDirectory", parseTemplateDirectoryError.Error())
		return nil, parseTemplateDirectoryInternalError
	}

	template = template.Lookup(templateName)
	template.Execute(&body, &data)
	fmt.Println(template.Name())
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", data.Subject)
	message.SetBody("text/html", body.String())
	message.AddAlternative("text/plain", html2text.HTML2Text(body.String()))
	return message, nil
}
