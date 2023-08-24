up:
	docker-compose up -d
	reflex -s go run cmd/server/main.go 

down:
	docker-compose down

run:
	go run cmd/server/main.go

reflex:
	reflex -s go run cmd/server/main.go 

build:
	docker-compose build