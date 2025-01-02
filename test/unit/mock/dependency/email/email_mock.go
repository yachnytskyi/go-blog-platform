package common

import (
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
)

type MockEmail struct{}

func NewMockEmail() MockEmail {
	return MockEmail{}
}

// Mock SendEmail function to always return nil (no error).
func (m MockEmail) SendEmail(config *model.ApplicationConfig, logger interfaces.Logger, location string, data any, emailData interfaces.EmailData) error {
	return nil
}
