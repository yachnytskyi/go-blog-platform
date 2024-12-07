initial:
	cp config/yaml/v1/example/local.application.example.yaml config/yaml/v1/local.application.yaml
	cp config/yaml/v1/example/test.application.example.yaml config/yaml/v1/test.application.yaml
	cp config/yaml/v1/example/docker.dev.application.example.yaml config/yaml/v1/docker.dev.application.yaml
	cp config/yaml/v1/example/docker.release.application.example.yaml config/yaml/v1/docker.release.application.yaml
	cp config/yaml/v1/example/docker.prod.application.example.yaml config/yaml/v1/docker.prod.application.yaml

	cp infrastructure/script/data/repository/mongo/example/init-mongo.example.js infrastructure/script/data/repository/mongo/init-mongo.js
	cp infrastructure/script/data/repository/mongo/example/init-test-data-mongo.example.js infrastructure/script/data/repository/mongo/init-test-data-mongo.js

mongo-local:
	docker compose up mongodb -d 
	reflex -s go run cmd/server/main.go 

mongo-local-docker:
	docker compose up mongodb app-local

mongo-dev:
	docker compose up mongodb app-dev

mongo-release:
	docker compose up mongodb app-release

mongo-production:
	docker compose up mongodb app-production 

build-mongo:
	docker compose build mongodb

# build-local-docker:
  	# docker compose build app-local

# build-dev:
# 	docker compose build app-dev

# build-release:
#   	docker compose build app-release

# build-production:
#   	docker compose build app-production

make tests:
	go test ./test/...

unit-tests:
	go test ./test/unit/...

unit-tests-pkg:
	go test ./test/unit/pkg/...

unit-tests-internal:
	go test ./test/unit/internal...

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








