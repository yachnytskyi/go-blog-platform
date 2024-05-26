up:
	docker-compose up 

down:
	docker-compose down

local:
	docker-compose up mongodb -d 
	reflex -s go run cmd/server/main.go 

run:
	go run cmd/server/main.go

reflex:
	reflex -s go run cmd/server/main.go 

