# NATS JetStream Pub/Sub Example

A complete implementation of a publish-subscribe pattern using NATS JetStream in Go, with separate publisher, subscriber, and monitoring applications.

## Features

- Publisher that sends messages to a NATS JetStream stream
- Subscriber with limited concurrency (max 10 messages at a time)
- Monitor that queries the NATS monitoring interface and logs detailed information
- Docker support for NATS server
- Configuration via environment variables
- Graceful shutdown handling

## Prerequisites

- Go 1.21+
- Docker (for running NATS server)

## Project Structure

```
nats-pubsub/
├── cmd/                 # Application entry points
│   ├── publisher/       # Publisher executable
│   ├── subscriber/      # Subscriber executable
│   └── monitor/        # Monitoring executable
├── internal/           # Internal packages
│   ├── config/         # Configuration handling
│   ├── pubsub/         # Publisher and subscriber implementation
│   ├── monitor/        # Monitoring implementation
│   └── stream/         # JetStream setup and management
├── bin/                # Built executables
├── Makefile           # Build and run commands
└── README.md          # Project documentation
```

## Quick Start

1. Start the NATS server:
   ```bash
   make docker-nats
   ```

2. Build all applications:
   ```bash
   make build
   ```

3. Run applications individually:
   ```bash
   # Run the monitor
   make run-monitor

   # In another terminal, run the publisher
   make run-publisher

   # Optionally, run the subscriber
   make run-subscriber
   ```

## Environment Variables

Configure the applications using these environment variables:

- `APP_NATS_URL`: NATS server URL (default: "nats://localhost:4222")
- `APP_NATS_MONITOR_URL`: NATS monitoring URL (default: "http://localhost:8222")
- `APP_STREAM_NAME`: JetStream stream name (default: "ORDERS")
- `APP_STREAM_SUBJECTS`: JetStream subjects (default: "ORDERS.*")
- `APP_STREAM_SUBJECTNAME`: Subject to publish/subscribe to (default: "ORDERS.received")
- `APP_STREAM_RETENTION`: Stream retention policy (default: "workqueue")
- `APP_STREAM_STORAGE`: Stream storage type (default: "file")
- `APP_STREAM_MAXAGE`: Maximum age of messages in seconds (default: 86400)

## Makefile Commands

- `make build`: Build all applications
- `make run-publisher`: Run the publisher
- `make run-subscriber`: Run the subscriber
- `make run-monitor`: Run the monitor
- `make docker-nats`: Start NATS server in Docker
- `make stop-nats`: Stop NATS server
- `make clean`: Remove built binaries
- `make test`: Run tests
- `make deps`: Install/update dependencies

## Development

To add new features or modify existing ones:

1. Install dependencies:
   ```bash
   make deps
   ```

2. Make your changes

3. Run tests:
   ```bash
   make test
   ```

4. Build and run:
   ```bash
   make build
   ```

## Cleanup

To stop the NATS server and clean up:

```bash
make stop-nats
make clean
```