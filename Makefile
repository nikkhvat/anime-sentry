
dev:
	go run ./cmd/app.go
	
build:
	go build ./cmd/app.go

start:
	docker-compose up --build -d

start_with_out_backgroud:
	docker-compose up --build

stop_anime_schedule_bot:
	docker-compose stop anime_schedule_bot
	docker-compose rm -f anime_schedule_bot

stop_db:
	docker-compose stop db
	docker-compose rm -f db

stop:
	make stop_anime_schedule_bot
	make stop_db