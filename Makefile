dev:
	go run ./cmd/app.go
	
build:
	go build ./cmd/app.go

start:
	docker-compose up --build --force-recreate --no-cache -d

start_interactive:
	docker-compose up --build --force-recreate --no-cache

update:
	git pull
	docker-compose stop
	docker-compose up --build --force-recreate --no-cache -d
	docker image prune -f

stop:
	docker-compose stop

remove:
	docker image prune -f