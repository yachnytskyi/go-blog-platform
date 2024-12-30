KEY_SIZE ?= 2048

initial:
	go install github.com/cespare/reflex@latest

	# Copy example configuration files to their respective locations.
	cp config/yaml/v1/example/local.application.example.yaml config/yaml/v1/local.application.yaml
	cp config/yaml/v1/example/docker.local.application.example.yaml config/yaml/v1/docker.local.application.yaml
	cp config/yaml/v1/example/docker.develop.application.example.yaml config/yaml/v1/docker.develop.application.yaml
	cp config/yaml/v1/example/docker.release.application.example.yaml config/yaml/v1/docker.release.application.yaml
	cp config/yaml/v1/example/docker.production.application.example.yaml config/yaml/v1/docker.production.application.yaml

	# Copy MongoDB environments.
	cp infrastructure/script/data/repository/mongo/example/.example.env infrastructure/script/data/repository/mongo/.env

	# Encode Docker configuration YAML files to base64.
	# On Windows, you need to use a Linux Subsystem (WSL) or PowerShell equivalent.
	base64 -i config/yaml/v1/docker.develop.application.yaml -o DOCKER_DEVELOP_APPLICATION_CONFIG_YAML.txt
	base64 -i config/yaml/v1/docker.release.application.yaml -o DOCKER_RELEASE_APPLICATION_CONFIG_YAML.txt
	base64 -i config/yaml/v1/docker.production.application.yaml -o DOCKER_PRODUCTION_APPLICATION_CONFIG_YAML.txt

	# Generate Public and Private Keys.
	openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:$(KEY_SIZE)
	openssl rsa -pubout -in private_key.pem -out public_key.pem
	base64 -i private_key.pem > private_key_base64.txt	
	base64 -i public_key.pem > public_key_base64.txt	

mongo-local:
	docker-compose up mongodb -d 
	reflex -s go run cmd/server/main.go 

mongo-local-docker:
	docker-compose up mongodb app-local -d

mongo-develop:
	docker-compose up mongodb app-develop -d

mongo-release:
	docker-compose up mongodb app-release -d

mongo-production:
	docker-compose up mongodb app-production -d

build-mongo:
	docker-compose build mongodb

build-local:
	docker-compose build app-local

build-develop:
	docker-compose build app-develop

build-release:
	docker-compose build app-release

build-production:
	docker-compose build app-production

make tests:
	go test ./test/...

unit-tests:
	go test ./test/unit/...

unit-tests-pkg:
	go test ./test/unit/pkg/...

unit-tests-internal:
	go test ./test/unit/internal...

lint:
	golangci-lint run

update:
	go get -u ./...	
	go mod tidy

up:
	docker-compose up 

build:
	docker-compose build --no-cache

stop:
	docker-compose stop

down:
	docker-compose down

down-v:
	docker-compose down -v

clean:
	rm -f private_key.pem public_key.pem
	rm -f private_key_base64.txt public_key_base64.txt
	rm -f DOCKER_DEVELOP_APPLICATION_CONFIG_YAML.txt
	rm -f DOCKER_RELEASE_APPLICATION_CONFIG_YAML.txt
	rm -f DOCKER_PRODUCTION_APPLICATION_CONFIG_YAML.txt