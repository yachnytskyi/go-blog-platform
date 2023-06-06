up:
	docker-compose up -d

down:
	docker-compose down

run:
	go run cmd/server/main.go

reflex:
	reflex -s go run cmd/server/main.go 