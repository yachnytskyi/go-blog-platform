package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
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
		if validator.IsValueNotNil(walkError) {
			sendEmailInternalError := domainError.NewInternalError(location+"ParseTemplateDirectory.Walk", walkError.Error())
			// logging.Logger(sendEmailInternalError)
			return sendEmailInternalError
		}
		if validator.IsBooleanNotTrue(info.IsDir()) {
			paths = append(paths, path)
		}
		return nil
	})
	fmt.Println(printlnMessage)
	if validator.IsValueNotNil(filePathWalkError) {
		sendEmailInternalError := domainError.NewInternalError(location+"ParseTemplateDirectory."+printlnMessage, filePathWalkError.Error())
		// logging.Logger(sendEmailInternalError)
		return nil, sendEmailInternalError
	}
	return template.ParseFiles(paths...)
}

func SendEmail(user *userModel.User, data *userModel.EmailData) error {
	loadConfig, loadConfigError := config.LoadConfig(".")
	if validator.IsValueNotNil(loadConfigError) {
		loadConfigInternalError := domainError.NewInternalError(location+"SendEmail.LoadConfig", loadConfigError.Error())
		logging.Logger(loadConfigInternalError)
		return loadConfigInternalError
	}

	// Send data.
	// from := loadConfig.EmailFrom
	smtpPass := loadConfig.SMTPPassword
	smtpUser := loadConfig.SMTPUser
	// to := user.Email
	smtpHost := loadConfig.SMTPHost
	smtpPort := loadConfig.SMTPPort

	// var body bytes.Buffer
	// template, parseTemplateDirectoryError := ParseTemplateDirectory(loadConfig.UserEmailTemplatePath)
	// if validator.IsValueNotNil(parseTemplateDirectoryError) {
	// 	logging.Logger(parseTemplateDirectoryError)
	// 	return parseTemplateDirectoryError
	// }

	// template = template.Lookup(templateName)
	// template.Execute(&body, &data)
	// fmt.Println(template.Name())
	// message := gomail.NewMessage()
	// message.SetHeader("From", from)
	// message.SetHeader("To", to)
	// message.SetHeader("Subject", data.Subject)
	// message.SetBody("text/html", body.String())
	// message.AddAlternative("text/plain", html2text.HTML2Text(body.String()))
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send an email.
	message, prepareSendMessageError := PrepareSendMessage(&loadConfig, user.Email, data)
	if validator.IsValueNotNil(prepareSendMessageError) {
		logging.Logger(prepareSendMessageError)
		return prepareSendMessageError
	}
	dialAndSendError := dialer.DialAndSend(message)
	if validator.IsValueNotNil(dialAndSendError) {
		sendEmailInternalError := domainError.NewInternalError(location+"SendEmail.DialAndSend", dialAndSendError.Error())
		logging.Logger(sendEmailInternalError)
		return sendEmailInternalError
	}
	return nil
}

func PrepareSendMessage(loadConfig *config.Config, userEmail string, data *userModel.EmailData) (*gomail.Message, error) {
	// Prepare data.
	templateName := loadConfig.UserEmailTemplateName
	from := loadConfig.EmailFrom
	to := userEmail

	var body bytes.Buffer
	template, parseTemplateDirectoryError := ParseTemplateDirectory(loadConfig.UserEmailTemplatePath)
	if validator.IsValueNotNil(parseTemplateDirectoryError) {
		return nil, parseTemplateDirectoryError
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
