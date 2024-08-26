# Golang Mongo gRPC

Source code for the Golang Mongo gRPC App.

## Introduction

This project is a Golang-based application that uses MongoDB, REST API, gRPC, following hexagonal architecture principles, manual dependency injection, and the abstract factory pattern.

## Prerequisites

Ensure you have the following installed:
- Golang 
- Docker
- Docker Compose

## Initializing

To set up the configuration files, rename the following files by removing "example" from their filenames:

- `config/yaml/v1/local.application.example.yaml` → `config/yaml/v1/local.application.yaml`
- `config/yaml/v1/docker.dev.application.example.yaml` → `config/yaml/v1/docker.dev.application.yaml`
- `config/yaml/v1/test.application.example.yaml` → `config/yaml/v1/test.application.yaml`
- `config/yaml/v1/docker.staging.application.example.yaml` → `config/yaml/v1/docker.staging.application.yaml`
- `config/yaml/v1/docker.prod.application.example.yaml` →  `config/yaml/v1/docker.prod.application.yaml`

## API Endpoints

The API is available at the following URLs:
- `http://localhost:8080/api/posts`
- `http://localhost:8080/api/users`

For a complete list of available API requests/URLs, check the server terminal upon launching the project.

## Build and Run

To build the project, use one of the following commands:
- `make build` (shortcut command from Makefile)
- `docker-compose build` (full command)

To run the project, choose from:
- `make up` (for Docker environment)
- `make local` (for local environment)

Make sure Docker and Docker Compose are installed.

## Stop Docker Compose Services

To stop the services, use:
- `make down` (shortcut command from Makefile)
- `docker-compose down` (full command)

## License

This project is licensed under the [Creative Commons Attribution-NonCommercial 4.0 International License](https://creativecommons.org/licenses/by-nc/4.0/).

You are free to:
- Share — copy and redistribute the material in any medium or format
- Adapt — remix, transform, and build upon the material

The above rights are granted under the following terms:
- **Attribution** — You must give appropriate credit, provide a link to the license, and indicate if changes were made. You may do so in any reasonable manner, but not in any way that suggests the licensor endorses you or your use.
- **NonCommercial** — You may not use the material for commercial purposes.

**Note**: This license does not grant you the rights to use the work for commercial purposes. For more details, visit the [Creative Commons License Deed](https://creativecommons.org/licenses/by-nc/4.0/).

## Ways to Improve

- **Add Unit and Integration Tests**: Enhance testing coverage.
- **Refactor the System**: Improve code structure and efficiency.

Feel free to provide additional ideas or suggestions for further improvements.
