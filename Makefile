.PHONY: build run docker-build docker-run docker-stop clean

# App settings
APP_NAME=deploy
PORT=8080

# Go building and running
build:
	go build -o app main.go

run: build
	./app

# Docker commands
docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run -d --name $(APP_NAME) -p $(PORT):$(PORT) $(APP_NAME)

docker-stop:
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

docker-logs:
	docker logs -f $(APP_NAME)

# Clean up
clean:
	rm -f app
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true
	docker rmi $(APP_NAME) || true

# Default target
all: build

# Start everything
start: docker-stop docker-build docker-run
	@echo "Application started on http://localhost:$(PORT)"