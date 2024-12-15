package common

import (
	"context"

	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
)

type MockRepository struct{}

func NewMockRepository() MockRepository {
	return MockRepository{}
}

func (m MockRepository) CreateRepository(ctx context.Context) any {
	return nil
}

func (m MockRepository) NewRepository(createRepository any, repository any) any {
	return nil
}

func (m MockRepository) HealthCheck(interfaces.Delivery) {
}

func (m MockRepository) DatabasePing() bool {
	return true
}

func (m MockRepository) Close(ctx context.Context) {
}
