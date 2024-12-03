# Go Blog Platform

Source code for the Go Blog Platform App.

## Introduction

This project is a Golang-based application that uses MongoDB, REST API, gRPC, following hexagonal architecture principles, manual dependency injection, and the abstract factory pattern.

## Prerequisites

Ensure you have the following installed:
- Golang 
- Docker
- Docker Compose

## Initializing 

To set up the project, run the following command:
- `make initial` 

After this, review and update the configuration settings to match your environment (local, development, release, production, etc.).

## API Endpoints

The API is available at the following URLs:
- `http://localhost:8080/api/posts`
- `http://localhost:8080/api/users`

For a complete list of available API requests/URLs, check the server terminal upon launching the project.

## Build and Run

To run the project, choose from:
- `make mongo-local` (for local environment)
- `make mongo-dev` (for dev Docker environment)
- `make mongo-release` (for release Docker environment)
- `make mongo-production` (for production Docker environment)

Make sure Docker and Docker Compose are installed.

## Stop Docker Compose Services

To stop and clean up the running services, use the following `make` commands:

- **Stop the services**:  
  `make stop` – Stops the services without removing the containers.

- **Remove the containers**:  
  `make down` – Stops the services and removes the containers.

- **Remove containers and volumes**:  
  `make down-v` – Stops the services, removes the containers, and also removes any associated Docker volumes (which may contain persistent data).



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
