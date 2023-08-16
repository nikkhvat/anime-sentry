
dev:
	go run ./cmd/app.go
	
build:
	go build ./cmd/app.go

start:
	docker-compose up --build -d

start_with_out_backgroud:
	docker-compose up --build

stop:
	docker compose down