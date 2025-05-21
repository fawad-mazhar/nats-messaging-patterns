.PHONY: build run-server run-client clean deps docker-nats stop-nats

# Build all applications
build:
	go build -o bin/server ./cmd/server
	go build -o bin/client ./cmd/client

# Run the server
run-server: build
	./bin/server

# Run the client
run-client: build
	./bin/client

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy

# Start NATS server using Docker
docker-nats:
	docker run -d --name nats -p 4222:4222 -p 8222:8222 nats:latest -js -m 8222

# Stop and remove NATS container
stop-nats:
	docker stop nats || true
	docker rm nats || true

# Start everything (NATS + Server)
start: stop-nats docker-nats
	@echo "Starting NATS server..."
	@sleep 2
	@echo "Starting request-reply server..."
	./bin/server
