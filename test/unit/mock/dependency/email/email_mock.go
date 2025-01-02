package common

import (
	"fmt"

	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
)

type MockEmail struct {
	IsError bool
}

func NewMockEmail() MockEmail {
	return MockEmail{}
}

func (mockEmail MockEmail) SendEmail(config *model.ApplicationConfig, logger interfaces.Logger, location string, data any, emailData interfaces.EmailData) error {
	if mockEmail.IsError {
		return fmt.Errorf("")
	}
	return nil
}
