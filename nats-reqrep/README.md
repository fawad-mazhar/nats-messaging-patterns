# NATS Request-Reply Pattern Example

A complete implementation of request-reply pattern using NATS JetStream in Go.

## Features

- Server that processes requests and sends responses
- Client that sends requests and waits for responses
- JetStream persistence for reliability
- Graceful shutdown handling
- Docker support for NATS server

## Prerequisites

- Go 1.21+
- Docker (for running NATS server)

## Project Structure

```
nats-reqrep/
├── cmd/
│   ├── server/        # Server executable
│   └── client/        # Client executable
├── internal/
│   ├── config/        # Configuration
│   └── reqrep/        # Request-reply implementation
├── bin/               # Built executables
├── Makefile          # Build and run commands
└── README.md         # Project documentation
```

## Quick Start

1. Start the NATS server:
   ```bash
   make docker-nats
   ```

2. Build the applications:
   ```bash
   make build
   ```

3. In one terminal, start the server:
   ```bash
   make run-server
   ```

4. In another terminal, run the client:
   ```bash
   make run-client
   ```

## Implementation Details

- Server creates a JetStream stream for persistence
- Uses durable consumers for reliable message delivery
- Client generates unique subjects for each request
- Implements timeout handling for requests
- Proper message acknowledgment

## Makefile Commands

- `make build`: Build server and client
- `make run-server`: Run the server
- `make run-client`: Run the client
- `make docker-nats`: Start NATS server in Docker
- `make stop-nats`: Stop NATS server
- `make clean`: Remove built binaries
- `make deps`: Install/update dependencies

## Cleanup

To stop the NATS server and clean up:

```bash
make stop-nats
make clean
```