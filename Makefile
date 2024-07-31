up:
	docker-compose up 

stop:
	docker-compose stop

down:
	docker-compose down

run:
	go run cmd/server/main.go

reflex:
	reflex -s go run cmd/server/main.go 

mongo-local:
	docker-compose up mongodb -d 
	reflex -s go run cmd/server/main.go 

mongo:
	docker compose up mongodb app 

mongo-prod:
	docker compose up mongodb app-production 

reflex:
	reflex -s go run cmd/server/main.go 

unit-tests:
	go test ./test/unit/...

unit-tests-pkg:
	go test ./test/unit/pkg/...

# unit-tests-internal:
# 	go test ./test/unit/internal...








