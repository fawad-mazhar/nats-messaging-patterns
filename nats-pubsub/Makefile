.PHONY: build run clean test deps docker-nats stop-nats

# Build all applications
build:
	go build -o bin/publisher ./cmd/publisher
	go build -o bin/subscriber ./cmd/subscriber
	go build -o bin/monitor ./cmd/monitor

# Run the publisher
run-publisher: build
	./bin/publisher

# Run the subscriber
run-subscriber: build
	./bin/subscriber

# Run the monitor
run-monitor: build
	./bin/monitor

# Start NATS server using Docker
docker-nats:
	docker run -d --name nats -p 4222:4222 -p 8222:8222 nats:latest -js -m 8222

# Stop and remove NATS container
stop-nats:
	docker stop nats || true
	docker rm nats || true

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod tidy

# Start everything locally (NATS + Monitor + Publisher)
start-local: stop-nats docker-nats
	@echo "Starting NATS server..."
	@sleep 2
	@echo "Starting monitor..."
	./bin/monitor & \
	@echo "Starting publisher..."
	./bin/publisher

