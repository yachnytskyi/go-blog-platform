initial:
	# Copy example configuration files to their respective locations.
	cp config/yaml/v1/example/local.application.example.yaml config/yaml/v1/local.application.yaml
	cp config/yaml/v1/example/docker.local.application.example.yaml config/yaml/v1/docker.local.application.yaml
	cp config/yaml/v1/example/docker.develop.application.example.yaml config/yaml/v1/docker.develop.application.yaml
	cp config/yaml/v1/example/docker.release.application.example.yaml config/yaml/v1/docker.release.application.yaml
	cp config/yaml/v1/example/docker.production.application.example.yaml config/yaml/v1/docker.production.application.yaml

	# Copy MongoDB initialization scripts.
	cp infrastructure/script/data/repository/mongo/example/init-mongo.example.js infrastructure/script/data/repository/mongo/init-mongo.js
	cp infrastructure/script/data/repository/mongo/example/init-test-data-mongo.example.js infrastructure/script/data/repository/mongo/init-test-data-mongo.js

	# Encode Docker configuration YAML files to base64.
	# On Windows, you need to use a Linux Subsystem (WSL) or PowerShell equivalent.
	base64 -i config/yaml/v1/docker.develop.application.yaml -o config/yaml/v1/docker_config_develop.txt
	base64 -i config/yaml/v1/docker.develop.application.yaml -o config/yaml/v1/docker_config_production.txt
	base64 -i config/yaml/v1/docker.develop.application.yaml -o config/yaml/v1/docker_config_release.txt

mongo-local:
	docker compose up mongodb -d 
	reflex -s go run cmd/server/main.go 

mongo-local-docker:
	docker compose up mongodb app-local -d

mongo-develop:
	docker compose up mongodb app-develop -d

mongo-release:
	docker compose up mongodb app-release -d

mongo-production:
	docker compose up mongodb app-production -d

build-mongo:
	docker compose build mongodb

build-local:
	docker compose build app-local

build-develop:
	docker compose build app-develop

build-release:
	docker compose build app-release

build-production:
	docker compose build app-production

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
	docker compose up 

build:
	docker compose build --no-cache

stop:
	docker compose stop

down:
	docker compose down

down-v:
	docker compose down -v
