## Golang Mongo gRPC     
Source code for Golang Mongo gRPC App.

The project uses:  
**Golang**  
**Gin**  
**MongoDB**  
**gRPC**

Hexagonal architecture, manual dependency injection, and abstract factory are implemented in the project.

## Initializing

To set up the configuration files, rename the following files by removing "example" from their filenames:

- Rename `config/yaml/v1/local.application.example.yaml` to `config/yaml/v1/local.application.yaml`.
- Rename `config/yaml/v1/docker.dev.application.example.yaml` to `config/yaml/v1/docker.dev.application.yaml`.

Repeat this for the following files:

- `test.application.yaml`
- `docker.staging.application.yaml`
- `docker.prod.application.yaml`

## API Endpoints

The API will be available at the following URLs:
- `http://localhost:8080/api/posts`
- `http://localhost:8080/api/users`

You can find all possible API requests/URLs when you launch the project in your server terminal.

## Build and Run

To build the project, you can use one of the following commands:
- `make build` (shortcut command from Makefile)
- `docker-compose build` (full command)

After building, run the project using:
- `make up` (for Docker environment)
- `make local` (for local environment)

To run the server, ensure you have Docker and Docker Compose installed.

## Stop Docker Compose Services

To stop the services, use:
- `make down` (shortcut command from Makefile)
- `docker-compose down` (full command)

## Ways to Improve

- **Add Unit and Integration Tests**: Enhance testing coverage.
- **Refactor the System**: Improve code structure and efficiency.

Feel free to provide additional ideas or suggestions for further improvements.
