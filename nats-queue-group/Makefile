.PHONY: build build-worker build-publisher run-worker run-workers run-publisher clean deps docker-nats stop-nats start-demo

# Build all components
build: build-worker build-publisher

# Build the worker
build-worker:
	go build -o bin/worker ./cmd/worker

# Build the publisher
build-publisher:
	go build -o bin/publisher ./cmd/publisher

# Run a worker with specified ID (default: 1)
run-worker: build-worker stop-nats docker-nats
	./bin/worker $(if $(ID),$(ID),1)

# Run the publisher
run-publisher: build-publisher 
	./bin/publisher

# Start multiple workers (default: 3)
run-workers: build-worker
	@for i in $$(seq 1 $(if $(N),$(N),3)); do \
		./bin/worker $$i & \
	done
	@echo "Started $(if $(N),$(N),3) workers"
	@wait

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy

# Start NATS server using Docker
docker-nats:
	docker run -d --name nats -p 4222:4222 -p 8222:8222 nats:latest

# Stop and remove NATS container
stop-nats:
	docker stop nats || true
	docker rm nats || true

# Start demo (NATS + 3 workers)
start-demo: stop-nats docker-nats
	@echo "Starting NATS server..."
	@sleep 2
	@echo "Starting 3 workers..."
	@make run-workers N=3 &
	@sleep 1
	@echo "Starting publisher..."
	@make run-publisher
