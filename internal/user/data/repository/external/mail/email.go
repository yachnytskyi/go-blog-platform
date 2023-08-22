package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"
	"gopkg.in/gomail.v2"
)

// Email template parser.
func ParseTemplateDirectory(directory string) (*template.Template, error) {
	var paths []string

	walkError := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			var sendEmailInternalError *domainError.InternalError = new(domainError.InternalError)
			sendEmailInternalError.Location = "User.Data.Repository.External.Mail.ParseTemplateDirectory.Walk"
			sendEmailInternalError.Reason = err.Error()
			fmt.Println(sendEmailInternalError)
			return sendEmailInternalError
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	fmt.Println("parsing template...")

	if walkError != nil {
		var sendEmailInternalError *domainError.InternalError = new(domainError.InternalError)
		sendEmailInternalError.Location = "User.Data.Repository.External.Mail.SendEmail.parsing template"
		sendEmailInternalError.Reason = walkError.Error()
		fmt.Println(sendEmailInternalError)
		return nil, sendEmailInternalError
	}

	return template.ParseFiles(paths...)
}

func SendEmail(user *userModel.User, data *userModel.EmailData, templateName string) error {
	loadConfig, loadConfigError := config.LoadConfig(".")

	if loadConfigError != nil {
		var loadConfigInternalError *domainError.InternalError = new(domainError.InternalError)
		loadConfigInternalError.Location = "User.Data.Repository.External.Mail.SendEmail.LoadConfig"
		loadConfigInternalError.Reason = loadConfigError.Error()
		fmt.Println(loadConfigInternalError)
		log.Fatal("could not load config")
		return loadConfigInternalError
	}

	// Send data.
	from := loadConfig.EmailFrom
	smtpPass := loadConfig.SMTPPassword
	smtpUser := loadConfig.SMTPUser
	to := user.Email
	smtpHost := loadConfig.SMTPHost
	smtpPort := loadConfig.SMTPPort

	var body bytes.Buffer
	template, parseTemplateDirectoryError := ParseTemplateDirectory(loadConfig.UserEmailTemplatePath)

	if parseTemplateDirectoryError != nil {
		var sendEmailInternalError *domainError.InternalError = new(domainError.InternalError)
		sendEmailInternalError.Location = "User.Data.Repository.External.Mail.SendEmail.ParseTemplateDirectory"
		sendEmailInternalError.Reason = "could not parse template"
		fmt.Println(sendEmailInternalError)
		return sendEmailInternalError
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
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send an email.
	if dialAndSendError := dialer.DialAndSend(message); dialAndSendError != nil {
		var sendEmailInternalError *domainError.InternalError = new(domainError.InternalError)
		sendEmailInternalError.Location = "User.Data.Repository.External.Mail.SendEmail.DialAndSend"
		sendEmailInternalError.Reason = dialAndSendError.Error()
		fmt.Println(sendEmailInternalError)
		return sendEmailInternalError
	}

	return nil
}
