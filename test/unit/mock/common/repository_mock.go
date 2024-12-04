package common

import (
	"context"
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

func (m MockRepository) Close(ctx context.Context) {
}
