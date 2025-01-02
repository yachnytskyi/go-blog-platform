package common

import (
	"context"

	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
)

type MockRepository struct{}

func NewMockRepository() MockRepository {
	return MockRepository{}
}

func (mockRepository MockRepository) CreateRepository(ctx context.Context) any {
	return nil
}

func (mockRepository MockRepository) NewRepository(createRepository any, repository any) any {
	return nil
}

func (mockRepository MockRepository) HealthCheck(interfaces.Delivery) {
}

func (mockRepository MockRepository) DatabasePing() bool {
	return true
}

func (mockRepository MockRepository) Close(ctx context.Context) {
}
