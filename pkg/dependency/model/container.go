package model

import "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"

type Container struct {
	RepositoryFactory factory.RepositoryFactory
	DomainFactory     factory.DomainFactory
}

func NewContainer(repositoryFactory factory.RepositoryFactory, domainFactory factory.DomainFactory) *Container {
	return &Container{
		RepositoryFactory: repositoryFactory,
		DomainFactory:     domainFactory,
	}
}
