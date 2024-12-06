package factory

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	factory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"
	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/data/repository"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/delivery"
	email "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/email"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/logger"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/common"
)

const (
	expectedLocation = "pkg.dependency.factory."
)

func TestNewLoggerZerolog(t *testing.T) {
	t.Parallel()
	mockConfig := mock.NewMockConfig()
	mockConfig.Core.Logger = constants.Zerolog

	zerolog := factory.NewLogger(mockConfig)
	assert.IsType(t, zerolog, logger.Zerolog{}, test.EqualMessage)
	assert.Implements(t, (*interfaces.Logger)(nil), zerolog, test.EqualMessage)
}

func TestNewEmailGoMail(t *testing.T) {
	t.Parallel()
	mockConfig := mock.NewMockConfig()
	mockLogger := mock.NewMockLogger()
	mockConfig.Core.Email = constants.GoMail

	goMail := factory.NewEmail(mockConfig, mockLogger)
	assert.IsType(t, goMail, email.GoMail{}, test.EqualMessage)
	assert.Implements(t, (*interfaces.Email)(nil), goMail, test.EqualMessage)
}

func TestNewRepositoryMongoDB(t *testing.T) {
	t.Parallel()
	mockConfig := mock.NewMockConfig()
	mockConfig.Core.Database = constants.MongoDB
	mockLogger := mock.NewMockLogger()

	mongoDBRepository := factory.NewRepositoryFactory(mockConfig, mockLogger)
	assert.IsType(t, &repository.MongoDBRepository{}, mongoDBRepository, test.EqualMessage)
	assert.Implements(t, (*interfaces.Repository)(nil), mongoDBRepository, test.EqualMessage)
}

func TestNewDeliveryGin(t *testing.T) {
	t.Parallel()
	mockConfig := mock.NewMockConfig()
	mockConfig.Core.Delivery = constants.Gin
	mockLogger := mock.NewMockLogger()
	mockRepository := mock.NewMockRepository()

	ctx := context.Background()
	ginDelivery := factory.NewDeliveryFactory(ctx, mockConfig, mockLogger, mockRepository)
	assert.IsType(t, &delivery.GinDelivery{}, ginDelivery, test.EqualMessage)
	assert.Implements(t, (*interfaces.Delivery)(nil), ginDelivery, test.EqualMessage)
}

func TestNewLoggerInvalidType(t *testing.T) {
	t.Parallel()
	mockConfig := mock.NewMockConfig()
	mockConfig.Core.Logger = constants.Zerolog + "1"

	notification := fmt.Sprintf(constants.UnsupportedLogger, mockConfig.Core.Logger)
	expectedError := domain.NewInternalError(expectedLocation+"NewLogger", notification)
	defer func() {
		recover := recover()
		if recover != nil {
			assert.Equal(t, recover, expectedError, test.EqualMessage)
		}
	}()

	factory.NewLogger(mockConfig)
}

func TestNewEmailInvalidType(t *testing.T) {
	t.Parallel()
	mockConfig := mock.NewMockConfig()
	mockConfig.Core.Email = constants.GoMail + "1"

	notification := fmt.Sprintf(constants.UnsupportedEmail, mockConfig.Core.Email)
	expectedError := domain.NewInternalError(expectedLocation+"NewEmail", notification)
	defer func() {
		recover := recover()
		if recover != nil {
			assert.Equal(t, recover, expectedError, test.EqualMessage)
		}
	}()

	factory.NewEmail(mockConfig, nil)
}

func TestNewRepositoryInvalidType(t *testing.T) {
	t.Parallel()
	mockConfig := mock.NewMockConfig()
	mockConfig.Core.Database = constants.MongoDB + "1"
	mockLogger := mock.NewMockLogger()

	notification := fmt.Sprintf(constants.UnsupportedRepository, mockConfig.Core.Database)
	expectedError := domain.NewInternalError(expectedLocation+"NewRepositoryFactory", notification)
	defer func() {
		recover := recover()
		if recover != nil {
			assert.Equal(t, recover, expectedError, test.EqualMessage)
		}
	}()
	
	factory.NewRepositoryFactory(mockConfig, mockLogger)
}

func TestNewDeliveryInvalidType(t *testing.T) {
	t.Parallel()
	mockConfig := mock.NewMockConfig()
	mockConfig.Core.Delivery = constants.Gin + "1"
	mockLogger := mock.NewMockLogger()
	mockRepository := mock.NewMockRepository()

	notification := fmt.Sprintf(constants.UnsupportedDelivery, mockConfig.Core.Delivery)
	expectedError := domain.NewInternalError(expectedLocation+"NewDeliveryFactory", notification)
	defer func() {
		recover := recover()
		if recover != nil {
			assert.Equal(t, recover, expectedError, test.EqualMessage)
		}
	}()

	ctx := context.Background()
	factory.NewDeliveryFactory(ctx, mockConfig, mockLogger, mockRepository)
}
